// Package admin 后台管理业务逻辑层
package admin

import (
	"github.com/linke15d/bondy-backend/internal/model"
	"github.com/linke15d/bondy-backend/internal/repository"
	"gorm.io/gorm"
)

// AdminSubService 订阅管理业务逻辑
type AdminSubService struct {
	subRepo *repository.SubscriptionRepository
	db      *gorm.DB
}

// NewAdminSubService 创建 AdminSubService 实例
func NewAdminSubService(subRepo *repository.SubscriptionRepository) *AdminSubService {
	return &AdminSubService{subRepo: subRepo}
}

// AdminSubListResult 订阅列表返回结构
type AdminSubListResult struct {
	// List 订阅列表
	List []model.Subscription `json:"list"`

	// Total 总数量
	Total int64 `json:"total"`

	// Page 当前页码
	Page int `json:"page"`

	// PageSize 每页数量
	PageSize int `json:"page_size"`
}
