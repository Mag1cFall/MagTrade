// AI 相關資料模型
//
// 包含 AI 對話歷史、推薦建議和異常警報
// ChatHistory：儲存使用者與 AI 助手的對話記錄
// AIRecommendation：儲存 AI 針對秒殺活動的分析建議
package model

import (
	"time"
)

// ChatRole 對話角色
type ChatRole string

const (
	ChatRoleUser      ChatRole = "user"      // 使用者
	ChatRoleAssistant ChatRole = "assistant" // AI 助手
	ChatRoleSystem    ChatRole = "system"    // 系統提示詞
)

// ChatHistory AI 對話歷史模型
type ChatHistory struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    int64     `gorm:"index;not null" json:"user_id"`
	SessionID string    `gorm:"type:varchar(64);index;not null" json:"session_id"` // 對話 Session
	Role      ChatRole  `gorm:"type:varchar(20);not null" json:"role"`
	Content   string    `gorm:"type:text;not null" json:"content"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}

func (ChatHistory) TableName() string {
	return "chat_histories"
}

// RecommendationType AI 推薦類型
type RecommendationType string

const (
	RecommendationTypePriceSuggestion RecommendationType = "price_suggestion" // 價格建議
	RecommendationTypeTimingAdvice    RecommendationType = "timing_advice"    // 時機建議
	RecommendationTypeRiskAlert       RecommendationType = "risk_alert"       // 風險警報
)

// AIRecommendation AI 推薦建議模型
type AIRecommendation struct {
	ID                 int64              `gorm:"primaryKey;autoIncrement" json:"id"`
	FlashSaleID        int64              `gorm:"index;not null" json:"flash_sale_id"`
	RecommendationType RecommendationType `gorm:"type:varchar(50);not null" json:"recommendation_type"`
	Content            string             `gorm:"type:jsonb;not null" json:"content"` // JSON 格式建議內容
	ConfidenceScore    float64            `gorm:"type:float" json:"confidence_score"` // 信心分數 0-1
	CreatedAt          time.Time          `gorm:"autoCreateTime" json:"created_at"`
}

func (AIRecommendation) TableName() string {
	return "ai_recommendations"
}

// AnomalyAlert 異常警報（記憶體結構，不持久化）
type AnomalyAlert struct {
	AlertType  string                 `json:"alert_type"` // ip_frequency/multi_device/bot_pattern
	Severity   string                 `json:"severity"`   // low/medium/high
	Details    map[string]interface{} `json:"details"`
	AutoAction string                 `json:"auto_action,omitempty"` // 自動處理動作
	Timestamp  time.Time              `json:"timestamp"`
}
