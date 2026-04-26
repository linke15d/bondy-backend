// Package admin 后台管理 Handler 层
package admin

import (
	"github.com/gin-gonic/gin"
	_ "github.com/linke15d/bondy-backend/internal/model"
	adminService "github.com/linke15d/bondy-backend/internal/service/admin"
	"github.com/linke15d/bondy-backend/pkg/response"
)

// AdminStatsHandler 运营数据统计请求处理器
type AdminStatsHandler struct {
	statsService *adminService.AdminStatsService
}

// NewAdminStatsHandler 创建 AdminStatsHandler 实例
func NewAdminStatsHandler(statsService *adminService.AdminStatsService) *AdminStatsHandler {
	return &AdminStatsHandler{statsService: statsService}
}

// GetDashboard 获取运营数据总览
//
//	@Summary		运营数据总览
//	@Description	获取系统核心运营指标，包括总用户数、今日新增、总记录数、付费用户数等
//	@Tags			后台管理-数据统计
//	@Produce		json
//	@Param			Authorization	header		string												true	"Bearer {access_token}"
//	@Success		200				{object}	response.Response{data=adminService.DashboardStats}	"运营数据"
//	@Failure		401				{object}	response.Response									"未登录"
//	@Security		BearerAuth
//	@Router			/admin/v1/stats/dashboard [post]
func (h *AdminStatsHandler) GetDashboard(c *gin.Context) {
	stats, err := h.statsService.GetDashboard()
	if err != nil {
		response.ServerError(c)
		return
	}

	response.Success(c, stats)
}
