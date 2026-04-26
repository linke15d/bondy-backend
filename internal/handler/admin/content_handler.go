// Package admin 后台管理 Handler 层
package admin

import (
	"github.com/gin-gonic/gin"
	_ "github.com/linke15d/bondy-backend/internal/model"
	adminService "github.com/linke15d/bondy-backend/internal/service/admin"
	"github.com/linke15d/bondy-backend/pkg/response"
)

// AdminContentHandler 内容管理请求处理器
type AdminContentHandler struct {
	contentService *adminService.AdminContentService
}

// NewAdminContentHandler 创建 AdminContentHandler 实例
func NewAdminContentHandler(contentService *adminService.AdminContentService) *AdminContentHandler {
	return &AdminContentHandler{contentService: contentService}
}

// ContentIDInput 内容 ID 请求参数
type ContentIDInput struct {
	ID string `json:"id" binding:"required"`
}

// GetTagTypeInput 获取标签列表请求参数
type GetTagTypeInput struct {
	// Type 标签类型：LOCATION / ACTIVITY，不传返回全部
	Type string `json:"type" binding:"omitempty,oneof=LOCATION ACTIVITY"`
}

// ListSystemTags 获取系统标签列表
//
//	@Summary		系统标签列表
//	@Description	获取所有系统预设标签，可按类型过滤
//	@Tags			后台管理-内容管理
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string								true	"Bearer {access_token}"
//	@Param			body			body		GetTagTypeInput						true	"标签类型"
//	@Success		200				{object}	response.Response{data=[]model.Tag}	"标签列表"
//	@Failure		401				{object}	response.Response					"未登录"
//	@Security		BearerAuth
//	@Router			/admin/v1/content/tags/list [post]
func (h *AdminContentHandler) ListSystemTags(c *gin.Context) {
	var input GetTagTypeInput
	c.ShouldBindJSON(&input)

	tags, err := h.contentService.ListSystemTags(input.Type)
	if err != nil {
		response.ServerError(c)
		return
	}

	response.Success(c, tags)
}

// CreateSystemTag 创建系统标签
//
//	@Summary		创建系统标签
//	@Description	创建一个新的系统预设标签，创建后所有用户都可以使用
//	@Tags			后台管理-内容管理
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string								true	"Bearer {access_token}"
//	@Param			body			body		adminService.CreateTagInput			true	"标签信息"
//	@Success		201				{object}	response.Response{data=model.Tag}	"创建成功"
//	@Failure		400				{object}	response.Response					"参数错误"
//	@Failure		401				{object}	response.Response					"未登录"
//	@Security		BearerAuth
//	@Router			/admin/v1/content/tags/create [post]
func (h *AdminContentHandler) CreateSystemTag(c *gin.Context) {
	var input adminService.CreateTagInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	tag, err := h.contentService.CreateSystemTag(input)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Created(c, tag)
}

// DeleteSystemTag 删除系统标签
//
//	@Summary		删除系统标签
//	@Description	删除指定的系统预设标签
//	@Tags			后台管理-内容管理
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string				true	"Bearer {access_token}"
//	@Param			body			body		ContentIDInput		true	"标签 ID"
//	@Success		200				{object}	response.Response	"删除成功"
//	@Failure		401				{object}	response.Response	"未登录"
//	@Security		BearerAuth
//	@Router			/admin/v1/content/tags/delete [post]
func (h *AdminContentHandler) DeleteSystemTag(c *gin.Context) {
	var input ContentIDInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.contentService.DeleteSystemTag(input.ID); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, nil)
}

// ListSystemPositions 获取系统姿势列表
//
//	@Summary		系统姿势列表
//	@Description	获取所有系统预设姿势
//	@Tags			后台管理-内容管理
//	@Produce		json
//	@Param			Authorization	header		string									true	"Bearer {access_token}"
//	@Success		200				{object}	response.Response{data=[]model.Position}	"姿势列表"
//	@Failure		401				{object}	response.Response						"未登录"
//	@Security		BearerAuth
//	@Router			/admin/v1/content/positions/list [post]
func (h *AdminContentHandler) ListSystemPositions(c *gin.Context) {
	positions, err := h.contentService.ListSystemPositions()
	if err != nil {
		response.ServerError(c)
		return
	}

	response.Success(c, positions)
}

// CreateSystemPosition 创建系统姿势
//
//	@Summary		创建系统姿势
//	@Description	创建一个新的系统预设姿势
//	@Tags			后台管理-内容管理
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string									true	"Bearer {access_token}"
//	@Param			body			body		adminService.CreatePositionInput		true	"姿势信息"
//	@Success		201				{object}	response.Response{data=model.Position}	"创建成功"
//	@Failure		400				{object}	response.Response						"参数错误"
//	@Failure		401				{object}	response.Response						"未登录"
//	@Security		BearerAuth
//	@Router			/admin/v1/content/positions/create [post]
func (h *AdminContentHandler) CreateSystemPosition(c *gin.Context) {
	var input adminService.CreatePositionInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	position, err := h.contentService.CreateSystemPosition(input)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Created(c, position)
}

// DeleteSystemPosition 删除系统姿势
//
//	@Summary		删除系统姿势
//	@Description	删除指定的系统预设姿势
//	@Tags			后台管理-内容管理
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string				true	"Bearer {access_token}"
//	@Param			body			body		ContentIDInput		true	"姿势 ID"
//	@Success		200				{object}	response.Response	"删除成功"
//	@Failure		401				{object}	response.Response	"未登录"
//	@Security		BearerAuth
//	@Router			/admin/v1/content/positions/delete [post]
func (h *AdminContentHandler) DeleteSystemPosition(c *gin.Context) {
	var input ContentIDInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.contentService.DeleteSystemPosition(input.ID); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, nil)
}
