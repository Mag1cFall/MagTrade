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

type FlashSaleHandler struct {
	flashSaleService *service.FlashSaleService
	anomalyDetector  *ai.AnomalyDetector
	log              *zap.Logger
}

func NewFlashSaleHandler(producer *mq.Producer, anomalyDetector *ai.AnomalyDetector, log *zap.Logger) *FlashSaleHandler {
	return &FlashSaleHandler{
		flashSaleService: service.NewFlashSaleService(producer, log),
		anomalyDetector:  anomalyDetector,
		log:              log,
	}
}

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

func (h *FlashSaleHandler) Rush(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid flash sale id")
		return
	}

	userID := middleware.GetUserID(c)
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

			if alert.Severity == "high" {
				response.TooManyRequests(c, "请求异常，请稍后重试")
				return
			}
		}
	}

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
