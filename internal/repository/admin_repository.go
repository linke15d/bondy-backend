// Package repository 数据访问层
package repository

import (
	"time"

	"github.com/linke15d/bondy-backend/internal/model"
	"gorm.io/gorm"
)

// AdminRepository 管理员数据访问对象
type AdminRepository struct {
	db *gorm.DB
}

// NewAdminRepository 创建 AdminRepository 实例
func NewAdminRepository(db *gorm.DB) *AdminRepository {
	return &AdminRepository{db: db}
}

// FindByUsername 通过用户名查找管理员
func (r *AdminRepository) FindByUsername(username string) (*model.Admin, error) {
	var admin model.Admin
	err := r.db.Where("username = ? AND is_active = true", username).First(&admin).Error
	if err != nil {
		return nil, err
	}
	return &admin, nil
}

// FindByID 通过 ID 查找管理员
func (r *AdminRepository) FindByID(id string) (*model.Admin, error) {
	var admin model.Admin
	err := r.db.Where("id = ?", id).First(&admin).Error
	if err != nil {
		return nil, err
	}
	return &admin, nil
}

// UpdateLastLogin 更新最后登录时间
func (r *AdminRepository) UpdateLastLogin(id string) error {
	now := time.Now()
	return r.db.Model(&model.Admin{}).
		Where("id = ?", id).
		Update("last_login_at", now).Error
}

// Create 创建管理员
func (r *AdminRepository) Create(admin *model.Admin) error {
	return r.db.Create(admin).Error
}

// List 获取管理员列表
func (r *AdminRepository) List(page, pageSize int) ([]model.Admin, int64, error) {
	var admins []model.Admin
	var total int64

	r.db.Model(&model.Admin{}).Count(&total)

	offset := (page - 1) * pageSize
	err := r.db.Order("created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&admins).Error

	return admins, total, err
}
