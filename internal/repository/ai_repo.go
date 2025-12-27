// AI 相關資料存取層
//
// 本檔案包含對話歷史和 AI 建議的資料存取
// ChatHistoryRepository：對話歷史記錄
// AIRecommendationRepository：AI 分析建議記錄
package repository

import (
	"context"

	"github.com/Mag1cFall/magtrade/internal/database"
	"github.com/Mag1cFall/magtrade/internal/model"
	"gorm.io/gorm"
)

// ChatHistoryRepository 對話歷史資料存取
type ChatHistoryRepository struct {
	db *gorm.DB
}

func NewChatHistoryRepository() *ChatHistoryRepository {
	return &ChatHistoryRepository{db: database.Get()}
}

// Create 建立對話記錄
func (r *ChatHistoryRepository) Create(ctx context.Context, chat *model.ChatHistory) error {
	return r.db.WithContext(ctx).Create(chat).Error
}

// GetBySession 取得指定 Session 的對話歷史
func (r *ChatHistoryRepository) GetBySession(ctx context.Context, userID int64, sessionID string, limit int) ([]model.ChatHistory, error) {
	var chats []model.ChatHistory
	result := r.db.WithContext(ctx).
		Where("user_id = ? AND session_id = ?", userID, sessionID).
		Order("created_at ASC"). // 按時間正序
		Limit(limit).
		Find(&chats)

	return chats, result.Error
}

// DeleteBySession 刪除指定 Session 的對話歷史
func (r *ChatHistoryRepository) DeleteBySession(ctx context.Context, userID int64, sessionID string) error {
	return r.db.WithContext(ctx).
		Where("user_id = ? AND session_id = ?", userID, sessionID).
		Delete(&model.ChatHistory{}).Error
}

// AIRecommendationRepository AI 建議資料存取
type AIRecommendationRepository struct {
	db *gorm.DB
}

func NewAIRecommendationRepository() *AIRecommendationRepository {
	return &AIRecommendationRepository{db: database.Get()}
}

// Create 建立 AI 建議
func (r *AIRecommendationRepository) Create(ctx context.Context, rec *model.AIRecommendation) error {
	return r.db.WithContext(ctx).Create(rec).Error
}

// GetByFlashSale 取得指定秒殺活動的所有建議
func (r *AIRecommendationRepository) GetByFlashSale(ctx context.Context, flashSaleID int64) ([]model.AIRecommendation, error) {
	var recs []model.AIRecommendation
	result := r.db.WithContext(ctx).
		Where("flash_sale_id = ?", flashSaleID).
		Order("created_at DESC").
		Find(&recs)

	return recs, result.Error
}

// GetLatestByType 取得指定類型的最新建議
func (r *AIRecommendationRepository) GetLatestByType(ctx context.Context, flashSaleID int64, recType model.RecommendationType) (*model.AIRecommendation, error) {
	var rec model.AIRecommendation
	result := r.db.WithContext(ctx).
		Where("flash_sale_id = ? AND recommendation_type = ?", flashSaleID, recType).
		Order("created_at DESC").
		First(&rec)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}

	return &rec, nil
}
