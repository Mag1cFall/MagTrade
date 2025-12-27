// 統一 HTTP 回應格式
//
// 本檔案定義標準化 API 回應結構和快捷方法
// 包含：成功回應、錯誤回應、業務錯誤碼
// 所有 Handler 使用此模組返回一致的 JSON 格式
package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 標準回應結構
type Response struct {
	Code    int         `json:"code"` // 業務碼：0=成功，其他=錯誤
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"` // 資料載體
}

// 業務錯誤碼
const (
	CodeSuccess            = 0    // 成功
	CodeBadRequest         = 400  // 請求參數錯誤
	CodeUnauthorized       = 401  // 未認證
	CodeForbidden          = 403  // 無權限
	CodeNotFound           = 404  // 資源不存在
	CodeConflict           = 409  // 資源衝突
	CodeTooManyRequests    = 429  // 請求過於頻繁
	CodeInternalError      = 500  // 伺服器內部錯誤
	CodeStockInsufficient  = 1001 // 庫存不足
	CodeLimitExceeded      = 1002 // 超出限購數量
	CodeFlashSaleNotActive = 1003 // 秒殺活動未開始或已結束
	CodeOrderNotFound      = 1004 // 訂單不存在
	CodeOrderStatusInvalid = 1005 // 訂單狀態不允許此操作
)

// 錯誤碼對應訊息
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

// Success 成功回應
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    CodeSuccess,
		Message: "success",
		Data:    data,
	})
}

// SuccessWithMessage 成功回應（自訂訊息）
func SuccessWithMessage(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    CodeSuccess,
		Message: message,
		Data:    data,
	})
}

// Error 錯誤回應
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

// BadRequest 400 錯誤
func BadRequest(c *gin.Context, message string) {
	Error(c, http.StatusBadRequest, CodeBadRequest, message)
}

// Unauthorized 401 錯誤
func Unauthorized(c *gin.Context, message string, data ...interface{}) {
	resp := Response{
		Code:    CodeUnauthorized,
		Message: message,
	}
	if len(data) > 0 {
		resp.Data = data[0]
	}
	c.JSON(http.StatusUnauthorized, resp)
}

// Forbidden 403 錯誤
func Forbidden(c *gin.Context, message string) {
	Error(c, http.StatusForbidden, CodeForbidden, message)
}

// NotFound 404 錯誤
func NotFound(c *gin.Context, message string) {
	Error(c, http.StatusNotFound, CodeNotFound, message)
}

// Conflict 409 錯誤
func Conflict(c *gin.Context, message string) {
	Error(c, http.StatusConflict, CodeConflict, message)
}

// TooManyRequests 429 錯誤
func TooManyRequests(c *gin.Context, message string) {
	Error(c, http.StatusTooManyRequests, CodeTooManyRequests, message)
}

// InternalError 500 錯誤
func InternalError(c *gin.Context, message string) {
	Error(c, http.StatusInternalServerError, CodeInternalError, message)
}

// StockInsufficient 庫存不足
func StockInsufficient(c *gin.Context) {
	Error(c, http.StatusOK, CodeStockInsufficient, "")
}

// LimitExceeded 超出限購
func LimitExceeded(c *gin.Context) {
	Error(c, http.StatusOK, CodeLimitExceeded, "")
}

// FlashSaleNotActive 活動未開始或已結束
func FlashSaleNotActive(c *gin.Context) {
	Error(c, http.StatusOK, CodeFlashSaleNotActive, "")
}
