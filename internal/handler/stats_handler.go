// Package handler HTTP 请求处理层
package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/linke15d/bondy-backend/internal/service"
	"github.com/linke15d/bondy-backend/pkg/response"
)

// StatsHandler 数据统计相关请求处理器
type StatsHandler struct {
	statsService *service.StatsService
}

// NewStatsHandler 创建 StatsHandler 实例
func NewStatsHandler(statsService *service.StatsService) *StatsHandler {
	return &StatsHandler{statsService: statsService}
}

// GetOverview 获取总览统计
//
//	@Summary		总览统计
//	@Description	获取伴侣关系建立至今的累计数据，包括总次数、总时长、平均评分、连续记录天数、有记录的年份列表
//	@Tags			数据统计
//	@Produce		json
//	@Param			Authorization	header		string											true	"Bearer {access_token}"
//	@Success		200				{object}	response.Response{data=service.OverviewResult}	"总览数据"
//	@Failure		401				{object}	response.Response								"未登录"
//	@Failure		400				{object}	response.Response								"暂无伴侣关系"
//	@Security		BearerAuth
//	@Router			/api/v1/stats/overview [post]
func (h *StatsHandler) GetOverview(c *gin.Context) {
	userID := c.GetString("userID")

	result, err := h.statsService.GetOverview(userID)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, result)
}

// GetYearlyStats 获取年度统计
//
//	@Summary		年度统计
//	@Description	获取指定年份的完整统计数据，包括每月数据、全年热力图、标签/姿势使用排行、时段分布、星期分布
//	@Tags			数据统计
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string											true	"Bearer {access_token}"
//	@Param			body			body		service.YearInput								true	"年份"
//	@Success		200				{object}	response.Response{data=service.YearlyStatsResult}	"年度统计数据"
//	@Failure		400				{object}	response.Response								"参数错误"
//	@Failure		401				{object}	response.Response								"未登录"
//	@Security		BearerAuth
//	@Router			/api/v1/stats/yearly [post]
func (h *StatsHandler) GetYearlyStats(c *gin.Context) {
	userID := c.GetString("userID")

	var input service.YearInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	result, err := h.statsService.GetYearlyStats(userID, input)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, result)
}

// GetMonthlyStats 获取月度统计
//
//	@Summary		月度统计
//	@Description	获取指定月份的统计数据和日历视图数据，包含当月每天的记录列表，用于日历展示
//	@Tags			数据统计
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string										true	"Bearer {access_token}"
//	@Param			body			body		service.MonthInput							true	"年份和月份"
//	@Success		200				{object}	response.Response							"月度统计数据"
//	@Failure		400				{object}	response.Response							"参数错误"
//	@Failure		401				{object}	response.Response							"未登录"
//	@Security		BearerAuth
//	@Router			/api/v1/stats/monthly [post]
func (h *StatsHandler) GetMonthlyStats(c *gin.Context) {
	userID := c.GetString("userID")

	var input service.MonthInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	result, err := h.statsService.GetMonthlyStats(userID, input)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, result)
}
