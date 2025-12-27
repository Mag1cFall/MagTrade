// 訂單相關 HTTP 處理器
//
// 本檔案處理訂單相關請求
// 包含：訂單列表、訂單詳情、支付、取消
package handler

import (
	"strconv"

	"github.com/Mag1cFall/magtrade/internal/middleware"
	"github.com/Mag1cFall/magtrade/internal/mq"
	"github.com/Mag1cFall/magtrade/internal/pkg/response"
	"github.com/Mag1cFall/magtrade/internal/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// OrderHandler 訂單處理器
type OrderHandler struct {
	orderService *service.OrderService
}

func NewOrderHandler(producer *mq.Producer, log *zap.Logger) *OrderHandler {
	return &OrderHandler{
		orderService: service.NewOrderService(producer, log),
	}
}

// List 查詢當前使用者的訂單列表
// GET /api/v1/orders?page=1&page_size=20
func (h *OrderHandler) List(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "authentication required")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	result, err := h.orderService.ListByUser(c.Request.Context(), userID, page, pageSize)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, result)
}

// GetByOrderNo 根據訂單號查詢訂單詳情
// GET /api/v1/orders/:order_no
func (h *OrderHandler) GetByOrderNo(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "authentication required")
		return
	}

	orderNo := c.Param("order_no")
	if orderNo == "" {
		response.BadRequest(c, "order_no is required")
		return
	}

	order, err := h.orderService.GetByOrderNo(c.Request.Context(), userID, orderNo)
	if err != nil {
		response.NotFound(c, "order not found")
		return
	}

	response.Success(c, order)
}

// Pay 支付訂單（模擬支付）
// POST /api/v1/orders/:order_no/pay
func (h *OrderHandler) Pay(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "authentication required")
		return
	}

	orderNo := c.Param("order_no")
	if orderNo == "" {
		response.BadRequest(c, "order_no is required")
		return
	}

	order, err := h.orderService.Pay(c.Request.Context(), userID, orderNo)
	if err != nil {
		if err == service.ErrOrderStatusInvalid {
			response.BadRequest(c, "订单状态不允许支付")
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.SuccessWithMessage(c, "支付成功", order)
}

// Cancel 取消訂單
// POST /api/v1/orders/:order_no/cancel
func (h *OrderHandler) Cancel(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "authentication required")
		return
	}

	orderNo := c.Param("order_no")
	if orderNo == "" {
		response.BadRequest(c, "order_no is required")
		return
	}

	order, err := h.orderService.Cancel(c.Request.Context(), userID, orderNo)
	if err != nil {
		if err == service.ErrOrderStatusInvalid {
			response.BadRequest(c, "订单状态不允许取消")
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.SuccessWithMessage(c, "订单已取消", order)
}
