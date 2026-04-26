// Package admin 后台管理 Handler 层
package admin

import (
	_ "github.com/linke15d/bondy-backend/internal/model"
	adminService "github.com/linke15d/bondy-backend/internal/service/admin"
)

// AdminSubHandler 订阅管理请求处理器
type AdminSubHandler struct {
	subService *adminService.AdminSubService
}

// NewAdminSubHandler 创建 AdminSubHandler 实例
func NewAdminSubHandler(subService *adminService.AdminSubService) *AdminSubHandler {
	return &AdminSubHandler{subService: subService}
}
