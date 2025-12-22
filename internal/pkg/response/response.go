package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

const (
	CodeSuccess            = 0
	CodeBadRequest         = 400
	CodeUnauthorized       = 401
	CodeForbidden          = 403
	CodeNotFound           = 404
	CodeConflict           = 409
	CodeTooManyRequests    = 429
	CodeInternalError      = 500
	CodeStockInsufficient  = 1001
	CodeLimitExceeded      = 1002
	CodeFlashSaleNotActive = 1003
	CodeOrderNotFound      = 1004
	CodeOrderStatusInvalid = 1005
)

var codeMessages = map[int]string{
	CodeSuccess:            "success",
	CodeBadRequest:         "bad request",
	CodeUnauthorized:       "unauthorized",
	CodeForbidden:          "forbidden",
	CodeNotFound:           "not found",
	CodeConflict:           "conflict",
	CodeTooManyRequests:    "too many requests",
	CodeInternalError:      "internal server error",
	CodeStockInsufficient:  "库存不足",
	CodeLimitExceeded:      "超出限购数量",
	CodeFlashSaleNotActive: "秒杀活动未开始或已结束",
	CodeOrderNotFound:      "订单不存在",
	CodeOrderStatusInvalid: "订单状态不允许此操作",
}

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    CodeSuccess,
		Message: "success",
		Data:    data,
	})
}

func SuccessWithMessage(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    CodeSuccess,
		Message: message,
		Data:    data,
	})
}

func Error(c *gin.Context, httpStatus int, code int, message string) {
	if message == "" {
		if msg, ok := codeMessages[code]; ok {
			message = msg
		} else {
			message = "unknown error"
		}
	}
	c.JSON(httpStatus, Response{
		Code:    code,
		Message: message,
	})
}

func BadRequest(c *gin.Context, message string) {
	Error(c, http.StatusBadRequest, CodeBadRequest, message)
}

func Unauthorized(c *gin.Context, message string) {
	Error(c, http.StatusUnauthorized, CodeUnauthorized, message)
}

func Forbidden(c *gin.Context, message string) {
	Error(c, http.StatusForbidden, CodeForbidden, message)
}

func NotFound(c *gin.Context, message string) {
	Error(c, http.StatusNotFound, CodeNotFound, message)
}

func Conflict(c *gin.Context, message string) {
	Error(c, http.StatusConflict, CodeConflict, message)
}

func TooManyRequests(c *gin.Context, message string) {
	Error(c, http.StatusTooManyRequests, CodeTooManyRequests, message)
}

func InternalError(c *gin.Context, message string) {
	Error(c, http.StatusInternalServerError, CodeInternalError, message)
}

func StockInsufficient(c *gin.Context) {
	Error(c, http.StatusOK, CodeStockInsufficient, "")
}

func LimitExceeded(c *gin.Context) {
	Error(c, http.StatusOK, CodeLimitExceeded, "")
}

func FlashSaleNotActive(c *gin.Context) {
	Error(c, http.StatusOK, CodeFlashSaleNotActive, "")
}
