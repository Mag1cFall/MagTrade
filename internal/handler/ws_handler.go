// WebSocket Hub 與連線處理器
//
// 本檔案實現 WebSocket 即時通訊功能
// WSHub：管理所有連線，支援單播和廣播
// 用於推送秒殺結果、訂單狀態變更等即時通知
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

// WebSocket 升級器配置
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // 允許所有來源，生產環境應限制
	},
}

// WSHub WebSocket 連線中心
// 管理所有客戶端連線，處理註冊、登出、消息分發
type WSHub struct {
	clients    map[int64]*WSClient // UserID → Client 映射
	clientsMux sync.RWMutex        // 讀寫鎖
	broadcast  chan *WSMessage     // 廣播頻道
	register   chan *WSClient      // 註冊頻道
	unregister chan *WSClient      // 登出頻道
	log        *zap.Logger
}

// WSClient 單個 WebSocket 客戶端
type WSClient struct {
	userID int64
	conn   *websocket.Conn
	send   chan []byte // 發送緩衝區
}

// WSMessage WebSocket 消息
type WSMessage struct {
	Type   string      `json:"type"` // 消息類型：order_result/stock_update/notification
	Data   interface{} `json:"data"`
	UserID int64       `json:"-"` // 目標使用者 ID（0 表示廣播）
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

// Run 啟動 Hub 事件迴圈（需在獨立 goroutine 執行）
func (h *WSHub) Run() {
	for {
		select {
		case client := <-h.register: // 新連線註冊
			h.clientsMux.Lock()
			h.clients[client.userID] = client
			h.clientsMux.Unlock()
			h.log.Debug("websocket client registered", zap.Int64("user_id", client.userID))

		case client := <-h.unregister: // 連線登出
			h.clientsMux.Lock()
			if _, ok := h.clients[client.userID]; ok {
				delete(h.clients, client.userID)
				close(client.send)
			}
			h.clientsMux.Unlock()
			h.log.Debug("websocket client unregistered", zap.Int64("user_id", client.userID))

		case message := <-h.broadcast: // 消息分發
			h.clientsMux.RLock()
			if message.UserID > 0 { // 單播：發送給指定使用者
				if client, ok := h.clients[message.UserID]; ok {
					data, _ := json.Marshal(message)
					select {
					case client.send <- data:
					default: // 發送緩衝區滿，關閉連線
						close(client.send)
						delete(h.clients, client.userID)
					}
				}
			} else { // 廣播：發送給所有使用者
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

// SendToUser 發送消息給指定使用者
func (h *WSHub) SendToUser(userID int64, msgType string, data interface{}) {
	h.broadcast <- &WSMessage{
		Type:   msgType,
		Data:   data,
		UserID: userID,
	}
}

// Broadcast 廣播消息給所有使用者
func (h *WSHub) Broadcast(msgType string, data interface{}) {
	h.broadcast <- &WSMessage{
		Type: msgType,
		Data: data,
	}
}

// WSHandler WebSocket HTTP 處理器
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

// HandleConnection 處理 WebSocket 連線請求
// GET /ws?token=xxx
// Token 通過 Query 傳遞（因 WebSocket 不支援自訂 Header）
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

	// HTTP 升級為 WebSocket
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

	h.hub.register <- client // 註冊到 Hub

	go h.writePump(client) // 啟動寫 goroutine
	go h.readPump(client)  // 啟動讀 goroutine
}

// readPump 讀取客戶端消息（主要處理 Pong 心跳）
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

// writePump 發送消息給客戶端
func (h *WSHandler) writePump(client *WSClient) {
	ticker := time.NewTicker(30 * time.Second) // Ping 心跳間隔
	defer func() {
		ticker.Stop()
		client.conn.Close()
	}()

	for {
		select {
		case message, ok := <-client.send:
			_ = client.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok { // 頻道已關閉
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

		case <-ticker.C: // 定時發送 Ping
			_ = client.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := client.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
