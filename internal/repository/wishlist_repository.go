// Package repository 数据访问层
package repository

import (
	"github.com/linke15d/bondy-backend/internal/model"
	"gorm.io/gorm"
)

// WishlistRepository 心愿清单数据访问对象
type WishlistRepository struct {
	db *gorm.DB
}

// NewWishlistRepository 创建 WishlistRepository 实例
func NewWishlistRepository(db *gorm.DB) *WishlistRepository {
	return &WishlistRepository{db: db}
}

// WishlistFilter 心愿列表过滤条件
type WishlistFilter struct {
	// Scope 按范围过滤：COUPLE / PERSONAL，空字符串表示全部
	Scope string
	// IsCompleted 按完成状态过滤：nil 表示全部，true 已完成，false 未完成
	IsCompleted *bool
	// Page 页码
	Page int
	// PageSize 每页数量
	PageSize int
}

// Create 创建心愿
func (r *WishlistRepository) Create(wishlist *model.Wishlist) error {
	return r.db.Create(wishlist).Error
}

// FindByID 通过 ID 查找心愿
func (r *WishlistRepository) FindByID(id string, coupleID string) (*model.Wishlist, error) {
	var wishlist model.Wishlist
	err := r.db.
		Where("id = ? AND couple_id = ?", id, coupleID).
		First(&wishlist).Error
	if err != nil {
		return nil, err
	}
	return &wishlist, nil
}

// List 获取心愿列表
// 按热度倒序，热度相同按创建时间倒序
func (r *WishlistRepository) List(coupleID string, filter WishlistFilter) ([]model.Wishlist, int64, error) {
	var wishlists []model.Wishlist
	var total int64

	query := r.db.Model(&model.Wishlist{}).
		Where("couple_id = ?", coupleID)

	// 按范围过滤
	if filter.Scope != "" {
		query = query.Where("scope = ?", filter.Scope)
	}

	// 按完成状态过滤
	if filter.IsCompleted != nil {
		query = query.Where("is_completed = ?", *filter.IsCompleted)
	}

	query.Count(&total)

	offset := (filter.Page - 1) * filter.PageSize
	err := query.
		Order("is_completed ASC, heat DESC, created_at DESC").
		Offset(offset).
		Limit(filter.PageSize).
		Find(&wishlists).Error

	return wishlists, total, err
}

// Update 更新心愿
func (r *WishlistRepository) Update(wishlist *model.Wishlist) error {
	return r.db.Save(wishlist).Error
}

// Delete 删除心愿（物理删除，心愿清单允许彻底删除）
func (r *WishlistRepository) Delete(id string, coupleID string) error {
	return r.db.
		Where("id = ? AND couple_id = ?", id, coupleID).
		Delete(&model.Wishlist{}).Error
}

// IncrHeat 增加热度值
func (r *WishlistRepository) IncrHeat(id string, coupleID string) error {
	return r.db.Model(&model.Wishlist{}).
		Where("id = ? AND couple_id = ?", id, coupleID).
		UpdateColumn("heat", gorm.Expr("heat + 1")).Error
}

// SetCompleted 标记完成状态
func (r *WishlistRepository) SetCompleted(id string, coupleID string, completed bool) error {
	updates := map[string]interface{}{
		"is_completed": completed,
	}
	if completed {
		now := "NOW()"
		updates["completed_at"] = gorm.Expr(now)
	} else {
		updates["completed_at"] = nil
	}
	return r.db.Model(&model.Wishlist{}).
		Where("id = ? AND couple_id = ?", id, coupleID).
		Updates(updates).Error
}
