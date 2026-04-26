// Package repository 数据访问层
package repository

import (
	"github.com/linke15d/bondy-backend/internal/model"
	"gorm.io/gorm"
)

// HealthRepository 健康记录数据访问对象
type HealthRepository struct {
	db *gorm.DB
}

// NewHealthRepository 创建 HealthRepository 实例
func NewHealthRepository(db *gorm.DB) *HealthRepository {
	return &HealthRepository{db: db}
}

// HealthListFilter 健康记录列表过滤条件
type HealthListFilter struct {
	// Type 按类型过滤：STI_TEST / VACCINE，空字符串表示全部
	Type string
	// Page 页码
	Page int
	// PageSize 每页数量
	PageSize int
}

// Create 创建健康记录
func (r *HealthRepository) Create(record *model.HealthRecord) error {
	return r.db.Create(record).Error
}

// FindByID 通过 ID 查找健康记录
// 同时验证所属用户，防止越权访问
func (r *HealthRepository) FindByID(id string, userID string) (*model.HealthRecord, error) {
	var record model.HealthRecord
	err := r.db.
		Where("id = ? AND user_id = ?", id, userID).
		First(&record).Error
	if err != nil {
		return nil, err
	}
	return &record, nil
}

// List 获取用户的健康记录列表
// 按检测日期倒序排列
func (r *HealthRepository) List(userID string, filter HealthListFilter) ([]model.HealthRecord, int64, error) {
	var records []model.HealthRecord
	var total int64

	query := r.db.Model(&model.HealthRecord{}).
		Where("user_id = ?", userID)

	if filter.Type != "" {
		query = query.Where("type = ?", filter.Type)
	}

	query.Count(&total)

	offset := (filter.Page - 1) * filter.PageSize
	err := query.
		Order("tested_at DESC").
		Offset(offset).
		Limit(filter.PageSize).
		Find(&records).Error

	return records, total, err
}

// Update 更新健康记录
func (r *HealthRepository) Update(record *model.HealthRecord) error {
	return r.db.Save(record).Error
}

// Delete 删除健康记录
func (r *HealthRepository) Delete(id string, userID string) error {
	return r.db.
		Where("id = ? AND user_id = ?", id, userID).
		Delete(&model.HealthRecord{}).Error
}

// FindUpcomingReminders 查找即将到期的提醒记录
// 用于定时任务推送提醒通知
// daysAhead: 提前几天提醒
func (r *HealthRepository) FindUpcomingReminders(daysAhead int) ([]model.HealthRecord, error) {
	var records []model.HealthRecord
	err := r.db.
		Preload("User").
		Where("next_remind_at IS NOT NULL AND next_remind_at <= NOW() + INTERVAL '? days'", daysAhead).
		Find(&records).Error
	return records, err
}
