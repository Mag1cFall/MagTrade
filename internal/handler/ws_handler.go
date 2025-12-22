package handler

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/Mag1cFall/magtrade/internal/config"
	"github.com/Mag1cFall/magtrade/internal/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WSHub struct {
	clients    map[int64]*WSClient
	clientsMux sync.RWMutex
	broadcast  chan *WSMessage
	register   chan *WSClient
	unregister chan *WSClient
	log        *zap.Logger
}

type WSClient struct {
	userID int64
	conn   *websocket.Conn
	send   chan []byte
}

type WSMessage struct {
	Type   string      `json:"type"`
	Data   interface{} `json:"data"`
	UserID int64       `json:"-"`
}

func NewWSHub(log *zap.Logger) *WSHub {
	return &WSHub{
		clients:    make(map[int64]*WSClient),
		broadcast:  make(chan *WSMessage, 256),
		register:   make(chan *WSClient),
		unregister: make(chan *WSClient),
		log:        log,
	}
}

func (h *WSHub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clientsMux.Lock()
			h.clients[client.userID] = client
			h.clientsMux.Unlock()
			h.log.Debug("websocket client registered", zap.Int64("user_id", client.userID))

		case client := <-h.unregister:
			h.clientsMux.Lock()
			if _, ok := h.clients[client.userID]; ok {
				delete(h.clients, client.userID)
				close(client.send)
			}
			h.clientsMux.Unlock()
			h.log.Debug("websocket client unregistered", zap.Int64("user_id", client.userID))

		case message := <-h.broadcast:
			h.clientsMux.RLock()
			if message.UserID > 0 {
				if client, ok := h.clients[message.UserID]; ok {
					data, _ := json.Marshal(message)
					select {
					case client.send <- data:
					default:
						close(client.send)
						delete(h.clients, client.userID)
					}
				}
			} else {
				for _, client := range h.clients {
					data, _ := json.Marshal(message)
					select {
					case client.send <- data:
					default:
						close(client.send)
					}
				}
			}
			h.clientsMux.RUnlock()
		}
	}
}

func (h *WSHub) SendToUser(userID int64, msgType string, data interface{}) {
	h.broadcast <- &WSMessage{
		Type:   msgType,
		Data:   data,
		UserID: userID,
	}
}

func (h *WSHub) Broadcast(msgType string, data interface{}) {
	h.broadcast <- &WSMessage{
		Type: msgType,
		Data: data,
	}
}

type WSHandler struct {
	hub    *WSHub
	jwtCfg *config.JWTConfig
	log    *zap.Logger
}

func NewWSHandler(hub *WSHub, jwtCfg *config.JWTConfig, log *zap.Logger) *WSHandler {
	return &WSHandler{
		hub:    hub,
		jwtCfg: jwtCfg,
		log:    log,
	}
}

func (h *WSHandler) HandleConnection(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
		return
	}

	claims, err := utils.ValidateAccessToken(token, h.jwtCfg.Secret)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		h.log.Error("websocket upgrade error", zap.Error(err))
		return
	}

	client := &WSClient{
		userID: claims.UserID,
		conn:   conn,
		send:   make(chan []byte, 256),
	}

	h.hub.register <- client

	go h.writePump(client)
	go h.readPump(client)
}

func (h *WSHandler) readPump(client *WSClient) {
	defer func() {
		h.hub.unregister <- client
		client.conn.Close()
	}()

	client.conn.SetReadLimit(512)
	_ = client.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	client.conn.SetPongHandler(func(string) error {
		_ = client.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, _, err := client.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				h.log.Error("websocket read error", zap.Error(err))
			}
			break
		}
	}
}

func (h *WSHandler) writePump(client *WSClient) {
	ticker := time.NewTicker(30 * time.Second)
	defer func() {
		ticker.Stop()
		client.conn.Close()
	}()

	for {
		select {
		case message, ok := <-client.send:
			_ = client.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				_ = client.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := client.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			_, _ = w.Write(message)

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			_ = client.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := client.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
