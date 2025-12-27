// 秒殺活動 HTTP 處理器
//
// 本檔案處理所有秒殺相關的 HTTP 請求
// 包含：活動列表、詳情、庫存查詢、建立活動、秒殺搶購
// Rush 方法會先進行 AI 異常檢測，再呼叫業務層處理
package handler

import (
	"strconv"

	"github.com/Mag1cFall/magtrade/internal/middleware"
	"github.com/Mag1cFall/magtrade/internal/model"
	"github.com/Mag1cFall/magtrade/internal/mq"
	"github.com/Mag1cFall/magtrade/internal/pkg/response"
	"github.com/Mag1cFall/magtrade/internal/service"
	"github.com/Mag1cFall/magtrade/internal/service/ai"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// FlashSaleHandler 秒殺 HTTP 處理器
type FlashSaleHandler struct {
	flashSaleService *service.FlashSaleService
	anomalyDetector  *ai.AnomalyDetector // AI 異常偵測（可選，若為 nil 則跳過）
	log              *zap.Logger
}

func NewFlashSaleHandler(producer *mq.Producer, anomalyDetector *ai.AnomalyDetector, log *zap.Logger) *FlashSaleHandler {
	return &FlashSaleHandler{
		flashSaleService: service.NewFlashSaleService(producer, log),
		anomalyDetector:  anomalyDetector,
		log:              log,
	}
}

// List 查詢秒殺活動列表
// GET /api/v1/flash-sales?page=1&page_size=20&status=1
func (h *FlashSaleHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	var status *model.FlashSaleStatus
	if statusStr := c.Query("status"); statusStr != "" {
		s, err := strconv.Atoi(statusStr)
		if err == nil {
			st := model.FlashSaleStatus(s)
			status = &st
		}
	}

	flashSales, total, err := h.flashSaleService.List(c.Request.Context(), page, pageSize, status)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"flash_sales": flashSales,
		"total":       total,
		"page":        page,
		"page_size":   pageSize,
	})
}

// GetByID 查詢活動詳情
// GET /api/v1/flash-sales/:id
func (h *FlashSaleHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid flash sale id")
		return
	}

	result, err := h.flashSaleService.GetByID(c.Request.Context(), id)
	if err != nil {
		response.NotFound(c, "flash sale not found")
		return
	}

	response.Success(c, result)
}

// GetStock 查詢即時庫存（高頻輪詢介面）
// GET /api/v1/flash-sales/:id/stock
func (h *FlashSaleHandler) GetStock(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid flash sale id")
		return
	}

	stock, err := h.flashSaleService.GetStock(c.Request.Context(), id)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, gin.H{"stock": stock})
}

// Create 建立秒殺活動（管理員專用）
// POST /api/v1/flash-sales
func (h *FlashSaleHandler) Create(c *gin.Context) {
	var req service.CreateFlashSaleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	flashSale, err := h.flashSaleService.Create(c.Request.Context(), &req)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, flashSale)
}

// Rush 秒殺搶購（核心介面）
// POST /api/v1/flash-sales/:id/rush
// 需要 JWT 認證，返回排隊憑證或錯誤訊息
func (h *FlashSaleHandler) Rush(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid flash sale id")
		return
	}

	userID := middleware.GetUserID(c) // 從 JWT 中取得使用者 ID
	if userID == 0 {
		response.Unauthorized(c, "authentication required")
		return
	}

	var req service.RushRequest
	if err := c.ShouldBindJSON(&req); err != nil && err.Error() != "EOF" {
		response.BadRequest(c, err.Error())
		return
	}

	if req.Quantity <= 0 {
		req.Quantity = 1
	}

	// AI 異常偵測（防機器人、頻率異常）
	if h.anomalyDetector != nil {
		event := &ai.RequestEvent{
			UserID:      userID,
			IP:          c.ClientIP(),
			UserAgent:   c.GetHeader("User-Agent"),
			FlashSaleID: id,
		}

		if alert := h.anomalyDetector.Detect(c.Request.Context(), event); alert != nil {
			h.log.Warn("anomaly detected during rush",
				zap.String("alert_type", alert.AlertType),
				zap.Int64("user_id", userID),
			)

			if alert.Severity == "high" { // 高風險直接拒絕
				response.TooManyRequests(c, "请求异常，请稍后重试")
				return
			}
		}
	}

	// 呼叫業務層執行秒殺邏輯
	result, err := h.flashSaleService.Rush(c.Request.Context(), userID, id, req.Quantity)
	if err != nil {
		switch err {
		case service.ErrFlashSaleNotActive:
			response.FlashSaleNotActive(c)
		case service.ErrStockInsufficient:
			response.StockInsufficient(c)
		case service.ErrLimitExceeded:
			response.LimitExceeded(c)
		case service.ErrAlreadyPurchased:
			response.Conflict(c, "您已参与过本次秒杀")
		case service.ErrFlashSaleNotStarted:
			response.BadRequest(c, "秒杀活动尚未开始")
		case service.ErrFlashSaleEnded:
			response.BadRequest(c, "秒杀活动已结束")
		default:
			response.InternalError(c, err.Error())
		}
		return
	}

	response.Success(c, result)
}
