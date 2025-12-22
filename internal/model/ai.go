package model

import (
	"time"
)

type ChatRole string

const (
	ChatRoleUser      ChatRole = "user"
	ChatRoleAssistant ChatRole = "assistant"
	ChatRoleSystem    ChatRole = "system"
)

type ChatHistory struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    int64     `gorm:"index;not null" json:"user_id"`
	SessionID string    `gorm:"type:varchar(64);index;not null" json:"session_id"`
	Role      ChatRole  `gorm:"type:varchar(20);not null" json:"role"`
	Content   string    `gorm:"type:text;not null" json:"content"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}

func (ChatHistory) TableName() string {
	return "chat_histories"
}

type RecommendationType string

const (
	RecommendationTypePriceSuggestion RecommendationType = "price_suggestion"
	RecommendationTypeTimingAdvice    RecommendationType = "timing_advice"
	RecommendationTypeRiskAlert       RecommendationType = "risk_alert"
)

type AIRecommendation struct {
	ID                 int64              `gorm:"primaryKey;autoIncrement" json:"id"`
	FlashSaleID        int64              `gorm:"index;not null" json:"flash_sale_id"`
	RecommendationType RecommendationType `gorm:"type:varchar(50);not null" json:"recommendation_type"`
	Content            string             `gorm:"type:jsonb;not null" json:"content"`
	ConfidenceScore    float64            `gorm:"type:float" json:"confidence_score"`
	CreatedAt          time.Time          `gorm:"autoCreateTime" json:"created_at"`
}

func (AIRecommendation) TableName() string {
	return "ai_recommendations"
}

type AnomalyAlert struct {
	AlertType  string                 `json:"alert_type"`
	Severity   string                 `json:"severity"`
	Details    map[string]interface{} `json:"details"`
	AutoAction string                 `json:"auto_action,omitempty"`
	Timestamp  time.Time              `json:"timestamp"`
}
