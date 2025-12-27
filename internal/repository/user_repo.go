// 使用者資料存取層
//
// 本檔案封裝使用者表的 CRUD 操作
// 使用 GORM 自動處理軟刪除
package repository

import (
	"context"
	"errors"

	"github.com/Mag1cFall/magtrade/internal/database"
	"github.com/Mag1cFall/magtrade/internal/model"
	"gorm.io/gorm"
)

// 錯誤定義
var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
)

// UserRepository 使用者資料存取
type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository() *UserRepository {
	return &UserRepository{db: database.Get()}
}

// Create 建立使用者
func (r *UserRepository) Create(ctx context.Context, user *model.User) error {
	result := r.db.WithContext(ctx).Create(user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return ErrUserAlreadyExists
		}
		return result.Error
	}
	return nil
}

// GetByID 根據 ID 查詢使用者
func (r *UserRepository) GetByID(ctx context.Context, id int64) (*model.User, error) {
	var user model.User
	result := r.db.WithContext(ctx).First(&user, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, result.Error
	}
	return &user, nil
}

// GetByUsername 根據使用者名稱查詢
func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User
	result := r.db.WithContext(ctx).Where("username = ?", username).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, result.Error
	}
	return &user, nil
}

// GetByEmail 根據郵件查詢
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	result := r.db.WithContext(ctx).Where("email = ?", email).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, result.Error
	}
	return &user, nil
}

// ExistsByUsername 檢查使用者名稱是否存在
func (r *UserRepository) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	var count int64
	result := r.db.WithContext(ctx).Model(&model.User{}).Where("username = ?", username).Count(&count)
	return count > 0, result.Error
}

// ExistsByEmail 檢查郵件是否存在
func (r *UserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var count int64
	result := r.db.WithContext(ctx).Model(&model.User{}).Where("email = ?", email).Count(&count)
	return count > 0, result.Error
}

// Update 更新使用者
func (r *UserRepository) Update(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

// UpdateEmailVerified 更新郵件驗證狀態
func (r *UserRepository) UpdateEmailVerified(ctx context.Context, userID int64, verified bool) error {
	return r.db.WithContext(ctx).Model(&model.User{}).Where("id = ?", userID).Update("email_verified", verified).Error
}
