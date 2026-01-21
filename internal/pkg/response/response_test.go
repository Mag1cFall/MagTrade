// 統一 HTTP 回應格式單元測試
//
// 測試覆蓋：
// - Response 結構體: Code, Message, Data 欄位
// - 業務錯誤碼常量: 0-500 標準碼 + 1001-1005 業務碼
// - Success/SuccessWithMessage: 成功回應
// - BadRequest/Unauthorized/Forbidden/NotFound/Conflict: HTTP 錯誤回應
// - TooManyRequests/InternalError: 限流與伺服器錯誤
// - StockInsufficient/LimitExceeded/FlashSaleNotActive: 秒殺業務錯誤
// - Error 空消息時使用預設訊息
// - Error 未知錯誤碼時返回 "unknown error"
package response

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestResponse_Structure(t *testing.T) {
	resp := Response{
		Code:    CodeSuccess,
		Message: "success",
		Data:    map[string]string{"key": "value"},
	}

	if resp.Code != 0 {
		t.Errorf("Response.Code = %v, want %v", resp.Code, 0)
	}
	if resp.Message != "success" {
		t.Errorf("Response.Message = %v, want %v", resp.Message, "success")
	}
	if resp.Data == nil {
		t.Error("Response.Data should not be nil")
	}
}

func TestBusinessCodes(t *testing.T) {
	tests := []struct {
		name string
		code int
		want int
	}{
		{"success", CodeSuccess, 0},
		{"bad request", CodeBadRequest, 400},
		{"unauthorized", CodeUnauthorized, 401},
		{"forbidden", CodeForbidden, 403},
		{"not found", CodeNotFound, 404},
		{"conflict", CodeConflict, 409},
		{"too many requests", CodeTooManyRequests, 429},
		{"internal error", CodeInternalError, 500},
		{"stock insufficient", CodeStockInsufficient, 1001},
		{"limit exceeded", CodeLimitExceeded, 1002},
		{"flash sale not active", CodeFlashSaleNotActive, 1003},
		{"order not found", CodeOrderNotFound, 1004},
		{"order status invalid", CodeOrderStatusInvalid, 1005},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.code != tt.want {
				t.Errorf("Code %s = %v, want %v", tt.name, tt.code, tt.want)
			}
		})
	}
}

func TestSuccess(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	testData := map[string]string{"id": "123"}
	Success(c, testData)

	if w.Code != http.StatusOK {
		t.Errorf("HTTP status = %v, want %v", w.Code, http.StatusOK)
	}

	var resp Response
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if resp.Code != CodeSuccess {
		t.Errorf("Response.Code = %v, want %v", resp.Code, CodeSuccess)
	}
	if resp.Message != "success" {
		t.Errorf("Response.Message = %v, want %v", resp.Message, "success")
	}
}

func TestSuccessWithMessage(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	testData := map[string]string{"id": "123"}
	SuccessWithMessage(c, "操作成功", testData)

	var resp Response
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if resp.Message != "操作成功" {
		t.Errorf("Response.Message = %v, want %v", resp.Message, "操作成功")
	}
}

func TestBadRequest(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	BadRequest(c, "参数错误")

	if w.Code != http.StatusBadRequest {
		t.Errorf("HTTP status = %v, want %v", w.Code, http.StatusBadRequest)
	}

	var resp Response
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if resp.Code != CodeBadRequest {
		t.Errorf("Response.Code = %v, want %v", resp.Code, CodeBadRequest)
	}
	if resp.Message != "参数错误" {
		t.Errorf("Response.Message = %v, want %v", resp.Message, "参数错误")
	}
}

func TestUnauthorized(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	Unauthorized(c, "请先登录")

	if w.Code != http.StatusUnauthorized {
		t.Errorf("HTTP status = %v, want %v", w.Code, http.StatusUnauthorized)
	}

	var resp Response
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if resp.Code != CodeUnauthorized {
		t.Errorf("Response.Code = %v, want %v", resp.Code, CodeUnauthorized)
	}
}

func TestUnauthorized_WithData(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	Unauthorized(c, "token expired", map[string]bool{"need_refresh": true})

	var resp Response
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if resp.Data == nil {
		t.Error("Response.Data should not be nil when data is provided")
	}
}

func TestForbidden(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	Forbidden(c, "无权限访问")

	if w.Code != http.StatusForbidden {
		t.Errorf("HTTP status = %v, want %v", w.Code, http.StatusForbidden)
	}

	var resp Response
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if resp.Code != CodeForbidden {
		t.Errorf("Response.Code = %v, want %v", resp.Code, CodeForbidden)
	}
}

func TestNotFound(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	NotFound(c, "资源不存在")

	if w.Code != http.StatusNotFound {
		t.Errorf("HTTP status = %v, want %v", w.Code, http.StatusNotFound)
	}

	var resp Response
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if resp.Code != CodeNotFound {
		t.Errorf("Response.Code = %v, want %v", resp.Code, CodeNotFound)
	}
}

func TestConflict(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	Conflict(c, "资源已存在")

	if w.Code != http.StatusConflict {
		t.Errorf("HTTP status = %v, want %v", w.Code, http.StatusConflict)
	}
}

func TestTooManyRequests(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	TooManyRequests(c, "请求过于频繁")

	if w.Code != http.StatusTooManyRequests {
		t.Errorf("HTTP status = %v, want %v", w.Code, http.StatusTooManyRequests)
	}
}

func TestInternalError(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	InternalError(c, "服务器内部错误")

	if w.Code != http.StatusInternalServerError {
		t.Errorf("HTTP status = %v, want %v", w.Code, http.StatusInternalServerError)
	}
}

func TestStockInsufficient(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	StockInsufficient(c)

	if w.Code != http.StatusOK {
		t.Errorf("HTTP status = %v, want %v", w.Code, http.StatusOK)
	}

	var resp Response
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if resp.Code != CodeStockInsufficient {
		t.Errorf("Response.Code = %v, want %v", resp.Code, CodeStockInsufficient)
	}
	if resp.Message != "库存不足" {
		t.Errorf("Response.Message = %v, want %v", resp.Message, "库存不足")
	}
}

func TestLimitExceeded(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	LimitExceeded(c)

	var resp Response
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if resp.Code != CodeLimitExceeded {
		t.Errorf("Response.Code = %v, want %v", resp.Code, CodeLimitExceeded)
	}
	if resp.Message != "超出限购数量" {
		t.Errorf("Response.Message = %v, want %v", resp.Message, "超出限购数量")
	}
}

func TestFlashSaleNotActive(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	FlashSaleNotActive(c)

	var resp Response
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if resp.Code != CodeFlashSaleNotActive {
		t.Errorf("Response.Code = %v, want %v", resp.Code, CodeFlashSaleNotActive)
	}
	if resp.Message != "秒杀活动未开始或已结束" {
		t.Errorf("Response.Message = %v, want %v", resp.Message, "秒杀活动未开始或已结束")
	}
}

func TestError_EmptyMessage_UsesDefault(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	Error(c, http.StatusBadRequest, CodeBadRequest, "")

	var resp Response
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if resp.Message != "bad request" {
		t.Errorf("Response.Message = %v, want %v", resp.Message, "bad request")
	}
}

func TestError_UnknownCode_UsesUnknownError(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	Error(c, http.StatusBadRequest, 99999, "")

	var resp Response
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if resp.Message != "unknown error" {
		t.Errorf("Response.Message = %v, want %v", resp.Message, "unknown error")
	}
}
