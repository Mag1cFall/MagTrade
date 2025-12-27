// 資料驗證工具
//
// 本檔案提供各類業務資料驗證方法
// 包含：使用者名稱、郵件、密碼、手機號、訂單號等
// 返回結構化驗證錯誤
package validator

import (
	"regexp"
	"unicode"
)

// 正則表達式
var (
	usernameRegex = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_]{2,29}$`) // 字母開頭，3-30 位
	emailRegex    = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	phoneRegex    = regexp.MustCompile(`^1[3-9]\d{9}$`) // 中國大陸手機號
)

// ValidationError 驗證錯誤
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (e *ValidationError) Error() string {
	return e.Message
}

// ValidateUsername 驗證使用者名稱
func ValidateUsername(username string) *ValidationError {
	if len(username) < 3 {
		return &ValidationError{Field: "username", Message: "用户名至少3个字符"}
	}
	if len(username) > 30 {
		return &ValidationError{Field: "username", Message: "用户名最多30个字符"}
	}
	if !usernameRegex.MatchString(username) {
		return &ValidationError{Field: "username", Message: "用户名必须以字母开头，只能包含字母、数字和下划线"}
	}
	return nil
}

// ValidateEmail 驗證郵件格式
func ValidateEmail(email string) *ValidationError {
	if len(email) == 0 {
		return &ValidationError{Field: "email", Message: "邮箱不能为空"}
	}
	if len(email) > 100 {
		return &ValidationError{Field: "email", Message: "邮箱最多100个字符"}
	}
	if !emailRegex.MatchString(email) {
		return &ValidationError{Field: "email", Message: "邮箱格式不正确"}
	}
	return nil
}

// ValidatePassword 驗證密碼強度
func ValidatePassword(password string) *ValidationError {
	if len(password) < 6 {
		return &ValidationError{Field: "password", Message: "密码至少6个字符"}
	}
	if len(password) > 50 {
		return &ValidationError{Field: "password", Message: "密码最多50个字符"}
	}

	// 檢查密碼複雜度
	var hasUpper, hasLower, hasNumber bool
	for _, c := range password {
		switch {
		case unicode.IsUpper(c):
			hasUpper = true
		case unicode.IsLower(c):
			hasLower = true
		case unicode.IsNumber(c):
			hasNumber = true
		}
	}

	if len(password) >= 8 && !(hasUpper || hasLower) {
		return &ValidationError{Field: "password", Message: "密码建议包含大小写字母和数字"}
	}

	_ = hasUpper && hasLower && hasNumber

	return nil
}

// ValidatePhone 驗證手機號（可選欄位）
func ValidatePhone(phone string) *ValidationError {
	if len(phone) == 0 {
		return nil
	}
	if !phoneRegex.MatchString(phone) {
		return &ValidationError{Field: "phone", Message: "手机号格式不正确"}
	}
	return nil
}

// ValidateQuantity 驗證購買數量
func ValidateQuantity(quantity int) *ValidationError {
	if quantity < 1 {
		return &ValidationError{Field: "quantity", Message: "数量至少为1"}
	}
	if quantity > 10 {
		return &ValidationError{Field: "quantity", Message: "单次购买数量不能超过10"}
	}
	return nil
}

// ValidateFlashSaleID 驗證秒殺活動 ID
func ValidateFlashSaleID(id int64) *ValidationError {
	if id <= 0 {
		return &ValidationError{Field: "flash_sale_id", Message: "无效的秒杀活动ID"}
	}
	return nil
}

// ValidateOrderNo 驗證訂單號
func ValidateOrderNo(orderNo string) *ValidationError {
	if len(orderNo) == 0 {
		return &ValidationError{Field: "order_no", Message: "订单号不能为空"}
	}
	if len(orderNo) < 10 || len(orderNo) > 50 {
		return &ValidationError{Field: "order_no", Message: "订单号格式不正确"}
	}
	return nil
}

// ValidateSessionID 驗證 Session ID
func ValidateSessionID(sessionID string) *ValidationError {
	if len(sessionID) == 0 {
		return &ValidationError{Field: "session_id", Message: "会话ID不能为空"}
	}
	if len(sessionID) > 64 {
		return &ValidationError{Field: "session_id", Message: "会话ID过长"}
	}
	return nil
}

// ValidateMessage 驗證消息內容
func ValidateMessage(message string) *ValidationError {
	if len(message) == 0 {
		return &ValidationError{Field: "message", Message: "消息不能为空"}
	}
	if len(message) > 2000 {
		return &ValidationError{Field: "message", Message: "消息最多2000个字符"}
	}
	return nil
}
