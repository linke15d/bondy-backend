// Package repository 数据访问层
package repository

import (
	"github.com/linke15d/bondy-backend/internal/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// SubscriptionRepository 订阅会员数据访问对象
type SubscriptionRepository struct {
	db *gorm.DB
}

// NewSubscriptionRepository 创建 SubscriptionRepository 实例
func NewSubscriptionRepository(db *gorm.DB) *SubscriptionRepository {
	return &SubscriptionRepository{db: db}
}

// FindByUserID 通过用户 ID 查找订阅记录
func (r *SubscriptionRepository) FindByUserID(userID string) (*model.Subscription, error) {
	var sub model.Subscription
	err := r.db.Where("user_id = ?", userID).First(&sub).Error
	if err != nil {
		return nil, err
	}
	return &sub, nil
}

// Upsert 创建或更新订阅记录
// 使用 ON CONFLICT 保证幂等，同一用户多次购买只保留最新记录
func (r *SubscriptionRepository) Upsert(sub *model.Subscription) error {
	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"plan", "status", "start_at", "expires_at", "provider", "provider_sub_id", "updated_at"}),
	}).Create(sub).Error
}

// UpdateStatus 更新订阅状态
func (r *SubscriptionRepository) UpdateStatus(userID string, status string) error {
	return r.db.Model(&model.Subscription{}).
		Where("user_id = ?", userID).
		Update("status", status).Error
}
