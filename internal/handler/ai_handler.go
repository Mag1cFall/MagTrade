package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Mag1cFall/magtrade/internal/config"
	"github.com/Mag1cFall/magtrade/internal/middleware"
	"github.com/Mag1cFall/magtrade/internal/pkg/response"
	"github.com/Mag1cFall/magtrade/internal/service/ai"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AIHandler struct {
	customerService *ai.CustomerService
	strategyAdvisor *ai.StrategyAdvisor
	log             *zap.Logger
}

func NewAIHandler(cfg *config.AIConfig, log *zap.Logger) *AIHandler {
	llmClient := ai.NewLLMClient(cfg, log)

	return &AIHandler{
		customerService: ai.NewCustomerService(llmClient, log),
		strategyAdvisor: ai.NewStrategyAdvisor(llmClient, log),
		log:             log,
	}
}

func (h *AIHandler) Chat(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "authentication required")
		return
	}

	var req ai.ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	result, err := h.customerService.Chat(c.Request.Context(), userID, &req)
	if err != nil {
		h.log.Error("chat error", zap.Error(err))
		response.InternalError(c, "AI服务暂时不可用，请稍后重试")
		return
	}

	response.Success(c, result)
}

func (h *AIHandler) ChatStream(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		c.JSON(401, gin.H{"error": "authentication required"})
		return
	}

	var req ai.ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("X-Accel-Buffering", "no")

	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		c.JSON(500, gin.H{"error": "streaming not supported"})
		return
	}

	err := h.customerService.ChatStream(c.Request.Context(), userID, &req, func(chunk ai.StreamChunk) error {
		data, _ := json.Marshal(chunk)
		_, err := fmt.Fprintf(c.Writer, "data: %s\n\n", data)
		if err != nil {
			return err
		}
		flusher.Flush()
		return nil
	})

	if err != nil {
		h.log.Error("chat stream error", zap.Error(err))
		fmt.Fprintf(c.Writer, "data: {\"type\":\"error\",\"content\":\"%s\"}\n\n", "AI服务暂时不可用")
		flusher.Flush()
	}
}

func (h *AIHandler) GetChatHistory(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "authentication required")
		return
	}

	sessionID := c.Query("session_id")
	if sessionID == "" {
		response.BadRequest(c, "session_id is required")
		return
	}

	history, err := h.customerService.GetHistory(c.Request.Context(), userID, sessionID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, gin.H{"history": history})
}

func (h *AIHandler) ClearChatHistory(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "authentication required")
		return
	}

	sessionID := c.Query("session_id")
	if sessionID == "" {
		response.BadRequest(c, "session_id is required")
		return
	}

	if err := h.customerService.ClearHistory(c.Request.Context(), userID, sessionID); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, nil)
}

func (h *AIHandler) GetRecommendation(c *gin.Context) {
	flashSaleID, err := strconv.ParseInt(c.Param("flash_sale_id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid flash_sale_id")
		return
	}

	result, err := h.strategyAdvisor.AnalyzeFlashSale(c.Request.Context(), flashSaleID)
	if err != nil {
		h.log.Error("strategy analysis error", zap.Error(err))
		response.InternalError(c, "分析服务暂时不可用")
		return
	}

	response.Success(c, result)
}

func (h *AIHandler) TriggerAnalysis(c *gin.Context) {
	flashSaleID, err := strconv.ParseInt(c.Param("flash_sale_id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid flash_sale_id")
		return
	}

	result, err := h.strategyAdvisor.AnalyzeFlashSale(c.Request.Context(), flashSaleID)
	if err != nil {
		h.log.Error("admin analysis error", zap.Error(err))
		response.InternalError(c, "分析失败")
		return
	}

	response.Success(c, result)
}
