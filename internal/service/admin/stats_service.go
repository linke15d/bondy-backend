// Package admin 后台管理业务逻辑层
package admin

import (
	"time"

	"gorm.io/gorm"
)

// AdminStatsService 运营数据统计业务逻辑
type AdminStatsService struct {
	db *gorm.DB
}

// NewAdminStatsService 创建 AdminStatsService 实例
func NewAdminStatsService(db *gorm.DB) *AdminStatsService {
	return &AdminStatsService{db: db}
}

// DashboardStats 运营数据总览
type DashboardStats struct {
	// TotalUsers 总注册用户数
	TotalUsers int64 `json:"total_users"`

	// TodayNewUsers 今日新增用户数
	TodayNewUsers int64 `json:"today_new_users"`

	// TotalCouples 总伴侣对数
	TotalCouples int64 `json:"total_couples"`

	// TotalRecords 总亲密记录数
	TotalRecords int64 `json:"total_records"`

	// TodayNewRecords 今日新增记录数
	TodayNewRecords int64 `json:"today_new_records"`

	// TotalPremiumUsers 总付费用户数
	TotalPremiumUsers int64 `json:"total_premium_users"`

	// ActiveUsersLast7Days 近7天活跃用户数（有新记录）
	ActiveUsersLast7Days int64 `json:"active_users_last_7_days"`
}

// GetDashboard 获取运营数据总览
func (s *AdminStatsService) GetDashboard() (*DashboardStats, error) {
	today := time.Now().Truncate(24 * time.Hour)
	sevenDaysAgo := time.Now().AddDate(0, 0, -7)

	var stats DashboardStats

	s.db.Table("users").Where("deleted_at IS NULL").Count(&stats.TotalUsers)
	s.db.Table("users").Where("deleted_at IS NULL AND created_at >= ?", today).Count(&stats.TodayNewUsers)
	s.db.Table("couples").Where("unlinked_at IS NULL").Count(&stats.TotalCouples)
	s.db.Table("records").Where("is_deleted = false").Count(&stats.TotalRecords)
	s.db.Table("records").Where("is_deleted = false AND created_at >= ?", today).Count(&stats.TodayNewRecords)
	s.db.Table("subscriptions").Where("status = 'ACTIVE'").Count(&stats.TotalPremiumUsers)
	s.db.Table("users").
		Joins("JOIN records ON records.created_by_id = users.id").
		Where("records.created_at >= ? AND records.is_deleted = false", sevenDaysAgo).
		Distinct("users.id").
		Count(&stats.ActiveUsersLast7Days)

	return &stats, nil
}
