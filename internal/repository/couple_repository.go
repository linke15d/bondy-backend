// Package repository 数据访问层
package repository

import (
	"time"

	"github.com/linke15d/bondy-backend/internal/model"
	"gorm.io/gorm"
)

// CoupleRepository 伴侣关系数据访问对象
type CoupleRepository struct {
	db *gorm.DB
}

// NewCoupleRepository 创建 CoupleRepository 实例
func NewCoupleRepository(db *gorm.DB) *CoupleRepository {
	return &CoupleRepository{db: db}
}

// Create 创建伴侣关系记录
func (r *CoupleRepository) Create(couple *model.Couple) error {
	return r.db.Create(couple).Error
}

// FindByUserID 通过用户 ID 查找其当前有效的伴侣关系
// 用户可能是 user1 也可能是 user2，两个字段都需要查
// 只返回未解绑的关系（unlinked_at IS NULL）
func (r *CoupleRepository) FindByUserID(userID string) (*model.Couple, error) {
	var couple model.Couple
	err := r.db.
		Preload("User1").
		Preload("User2").
		Where("(user1_id = ? OR user2_id = ?) AND unlinked_at IS NULL", userID, userID).
		First(&couple).Error
	if err != nil {
		return nil, err
	}
	return &couple, nil
}

// FindByInviteCode 通过邀请码查找伴侣关系记录
// 用于验证邀请码是否存在且未过期
func (r *CoupleRepository) FindByInviteCode(code string) (*model.Couple, error) {
	var couple model.Couple
	err := r.db.
		Where("invite_code = ? AND invite_expires_at > ?", code, time.Now()).
		First(&couple).Error
	if err != nil {
		return nil, err
	}
	return &couple, nil
}

// UpdateInviteCode 更新邀请码和过期时间
// 每次调用生成邀请码接口都会刷新邀请码
func (r *CoupleRepository) UpdateInviteCode(coupleID string, code string, expiresAt time.Time) error {
	return r.db.Model(&model.Couple{}).
		Where("id = ?", coupleID).
		Updates(map[string]interface{}{
			"invite_code":       code,
			"invite_expires_at": expiresAt,
		}).Error
}

// BindCouple 完成绑定：填入 User2ID 并清除邀请码
func (r *CoupleRepository) BindCouple(coupleID string, user2ID string) error {
	return r.db.Model(&model.Couple{}).
		Where("id = ?", coupleID).
		Updates(map[string]interface{}{
			"user2_id":          user2ID,
			"invite_code":       nil,
			"invite_expires_at": nil,
		}).Error
}

// Unlink 解除绑定：记录解绑时间，保留历史数据
func (r *CoupleRepository) Unlink(coupleID string) error {
	now := time.Now()
	return r.db.Model(&model.Couple{}).
		Where("id = ?", coupleID).
		Update("unlinked_at", now).Error
}

// ClearInviteCode 清除邀请码（邀请码使用后立即失效）
func (r *CoupleRepository) ClearInviteCode(coupleID string) error {
	return r.db.Model(&model.Couple{}).
		Where("id = ?", coupleID).
		Updates(map[string]interface{}{
			"invite_code":       nil,
			"invite_expires_at": nil,
		}).Error
}
