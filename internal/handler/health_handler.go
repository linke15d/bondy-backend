// Package handler HTTP 请求处理层
package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/linke15d/bondy-backend/internal/service"
	"github.com/linke15d/bondy-backend/pkg/response"
)

// HealthHandler 健康记录相关请求处理器
type HealthHandler struct {
	healthService *service.HealthService
}

// NewHealthHandler 创建 HealthHandler 实例
func NewHealthHandler(healthService *service.HealthService) *HealthHandler {
	return &HealthHandler{healthService: healthService}
}

// HealthIDInput 通过 ID 操作健康记录的请求参数
type HealthIDInput struct {
	// ID 健康记录唯一标识
	ID string `json:"id" binding:"required" example:"550e8400-e29b-41d4-a716-446655440000"`
}

// CreateHealthRecord 创建健康记录
//
//	@Summary		创建健康记录
//	@Description	添加一条 STI 检测或疫苗接种记录。健康记录属于个人私密数据，伴侣无法查看。备注请在客户端加密后传入。
//	@Tags			健康记录
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string												true	"Bearer {access_token}"
//	@Param			body			body		service.CreateHealthRecordInput						true	"健康记录信息"
//	@Success		201				{object}	response.Response{data=model.HealthRecord}			"创建成功"
//	@Failure		400				{object}	response.Response									"参数错误"
//	@Failure		401				{object}	response.Response									"未登录"
//	@Security		BearerAuth
//	@Router			/api/v1/health/create [post]
func (h *HealthHandler) CreateHealthRecord(c *gin.Context) {
	userID := c.GetString("userID")

	var input service.CreateHealthRecordInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	record, err := h.healthService.CreateHealthRecord(userID, input)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Created(c, record)
}

// ListHealthRecords 获取健康记录列表
//
//	@Summary		获取健康记录列表
//	@Description	获取当前用户的健康记录，只返回本人数据。可按类型过滤，按检测日期倒序排列。
//	@Tags			健康记录
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string												true	"Bearer {access_token}"
//	@Param			body			body		service.HealthListInput								true	"过滤条件"
//	@Success		200				{object}	response.Response{data=service.HealthListResult}	"健康记录列表"
//	@Failure		401				{object}	response.Response									"未登录"
//	@Security		BearerAuth
//	@Router			/api/v1/health/list [post]
func (h *HealthHandler) ListHealthRecords(c *gin.Context) {
	userID := c.GetString("userID")

	var input service.HealthListInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	result, err := h.healthService.ListHealthRecords(userID, input)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, result)
}

// GetHealthRecord 获取健康记录详情
//
//	@Summary		获取健康记录详情
//	@Description	获取单条健康记录详情，只能查看本人的记录
//	@Tags			健康记录
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string										true	"Bearer {access_token}"
//	@Param			body			body		HealthIDInput								true	"记录 ID"
//	@Success		200				{object}	response.Response{data=model.HealthRecord}	"记录详情"
//	@Failure		400				{object}	response.Response							"参数错误"
//	@Failure		401				{object}	response.Response							"未登录"
//	@Failure		404				{object}	response.Response							"记录不存在"
//	@Security		BearerAuth
//	@Router			/api/v1/health/detail [post]
func (h *HealthHandler) GetHealthRecord(c *gin.Context) {
	userID := c.GetString("userID")

	var input HealthIDInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	record, err := h.healthService.GetHealthRecord(userID, input.ID)
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}

	response.Success(c, record)
}

// UpdateHealthRecord 更新健康记录
//
//	@Summary		更新健康记录
//	@Description	更新健康记录内容，只传需要修改的字段。将 next_remind_at 传 null 可取消提醒。
//	@Tags			健康记录
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string										true	"Bearer {access_token}"
//	@Param			body			body		service.UpdateHealthRecordInput				true	"要更新的内容（需包含 id 字段）"
//	@Success		200				{object}	response.Response{data=model.HealthRecord}	"更新后的记录"
//	@Failure		400				{object}	response.Response							"参数错误"
//	@Failure		401				{object}	response.Response							"未登录"
//	@Security		BearerAuth
//	@Router			/api/v1/health/update [post]
func (h *HealthHandler) UpdateHealthRecord(c *gin.Context) {
	userID := c.GetString("userID")

	var req struct {
		ID string `json:"id" binding:"required"`
		service.UpdateHealthRecordInput
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	record, err := h.healthService.UpdateHealthRecord(userID, req.ID, req.UpdateHealthRecordInput)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, record)
}

// DeleteHealthRecord 删除健康记录
//
//	@Summary		删除健康记录
//	@Description	删除一条健康记录，删除后不可恢复
//	@Tags			健康记录
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string				true	"Bearer {access_token}"
//	@Param			body			body		HealthIDInput		true	"记录 ID"
//	@Success		200				{object}	response.Response	"删除成功"
//	@Failure		400				{object}	response.Response	"参数错误"
//	@Failure		401				{object}	response.Response	"未登录"
//	@Security		BearerAuth
//	@Router			/api/v1/health/delete [post]
func (h *HealthHandler) DeleteHealthRecord(c *gin.Context) {
	userID := c.GetString("userID")

	var input HealthIDInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.healthService.DeleteHealthRecord(userID, input.ID); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, nil)
}

// CancelReminder 取消提醒
//
//	@Summary		取消健康提醒
//	@Description	取消某条健康记录的下次提醒，取消后不会再收到该记录的推送通知
//	@Tags			健康记录
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string				true	"Bearer {access_token}"
//	@Param			body			body		HealthIDInput		true	"记录 ID"
//	@Success		200				{object}	response.Response	"取消成功"
//	@Failure		400				{object}	response.Response	"参数错误"
//	@Failure		401				{object}	response.Response	"未登录"
//	@Security		BearerAuth
//	@Router			/api/v1/health/reminder/cancel [post]
func (h *HealthHandler) CancelReminder(c *gin.Context) {
	userID := c.GetString("userID")

	var input HealthIDInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.healthService.CancelReminder(userID, input.ID); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, nil)
}
