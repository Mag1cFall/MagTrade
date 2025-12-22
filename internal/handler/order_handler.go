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

type OrderHandler struct {
	orderService *service.OrderService
}

func NewOrderHandler(producer *mq.Producer, log *zap.Logger) *OrderHandler {
	return &OrderHandler{
		orderService: service.NewOrderService(producer, log),
	}
}

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
