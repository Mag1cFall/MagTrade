package repository

import (
	"context"

	"github.com/Mag1cFall/magtrade/internal/database"
	"github.com/Mag1cFall/magtrade/internal/model"
	"gorm.io/gorm"
)

type ChatHistoryRepository struct {
	db *gorm.DB
}

func NewChatHistoryRepository() *ChatHistoryRepository {
	return &ChatHistoryRepository{db: database.Get()}
}

func (r *ChatHistoryRepository) Create(ctx context.Context, chat *model.ChatHistory) error {
	return r.db.WithContext(ctx).Create(chat).Error
}

func (r *ChatHistoryRepository) GetBySession(ctx context.Context, userID int64, sessionID string, limit int) ([]model.ChatHistory, error) {
	var chats []model.ChatHistory
	result := r.db.WithContext(ctx).
		Where("user_id = ? AND session_id = ?", userID, sessionID).
		Order("created_at ASC").
		Limit(limit).
		Find(&chats)

	return chats, result.Error
}

func (r *ChatHistoryRepository) DeleteBySession(ctx context.Context, userID int64, sessionID string) error {
	return r.db.WithContext(ctx).
		Where("user_id = ? AND session_id = ?", userID, sessionID).
		Delete(&model.ChatHistory{}).Error
}

type AIRecommendationRepository struct {
	db *gorm.DB
}

func NewAIRecommendationRepository() *AIRecommendationRepository {
	return &AIRecommendationRepository{db: database.Get()}
}

func (r *AIRecommendationRepository) Create(ctx context.Context, rec *model.AIRecommendation) error {
	return r.db.WithContext(ctx).Create(rec).Error
}

func (r *AIRecommendationRepository) GetByFlashSale(ctx context.Context, flashSaleID int64) ([]model.AIRecommendation, error) {
	var recs []model.AIRecommendation
	result := r.db.WithContext(ctx).
		Where("flash_sale_id = ?", flashSaleID).
		Order("created_at DESC").
		Find(&recs)

	return recs, result.Error
}

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
