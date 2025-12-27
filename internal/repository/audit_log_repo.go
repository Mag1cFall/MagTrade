// 審計日誌資料存取層
//
// 本檔案封裝審計日誌表的讀寫操作
// 記錄使用者敏感操作，如登入、支付、刪除等
package repository

import (
	"context"

	"github.com/Mag1cFall/magtrade/internal/database"
	"github.com/Mag1cFall/magtrade/internal/model"
	"gorm.io/gorm"
)

// AuditLogRepository 審計日誌資料存取
type AuditLogRepository struct {
	db *gorm.DB
}

func NewAuditLogRepository() *AuditLogRepository {
	return &AuditLogRepository{db: database.Get()}
}

// Create 建立審計日誌
func (r *AuditLogRepository) Create(ctx context.Context, log *model.AuditLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}

// ListByUserID 查詢指定使用者的審計日誌
func (r *AuditLogRepository) ListByUserID(ctx context.Context, userID int64, page, pageSize int) ([]model.AuditLog, int64, error) {
	var logs []model.AuditLog
	var total int64

	offset := (page - 1) * pageSize

	if err := r.db.WithContext(ctx).Model(&model.AuditLog{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&logs).Error; err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}
