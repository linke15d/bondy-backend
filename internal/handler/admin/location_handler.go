// Package admin 后台管理 Handler 层
package admin

import (
	"github.com/gin-gonic/gin"
	adminService "github.com/linke15d/bondy-backend/internal/service/admin"
	"github.com/linke15d/bondy-backend/pkg/response"
)

// LocationHandler 地点管理请求处理器
type LocationHandler struct {
	locationService *adminService.LocationService
}

// NewLocationHandler 创建 LocationHandler 实例
func NewLocationHandler(locationService *adminService.LocationService) *LocationHandler {
	return &LocationHandler{locationService: locationService}
}

// LocationNameRequest 单个语言地点名称（Swagger 用）
type LocationNameRequest struct {
	LanguageCode string `json:"language_code" example:"zh-CN"`
	Name         string `json:"name" example:"家里"`
}

// CreateLocationRequest 创建地点请求（Swagger 用）
type CreateLocationRequest struct {
	Names      []LocationNameRequest `json:"names"`
	IconBase64 *string               `json:"icon_base64"`
	SortOrder  int                   `json:"sort_order" example:"1"`
}

// UpdateLocationRequest 更新地点请求（Swagger 用）
type UpdateLocationRequest struct {
	ID         string                `json:"id" binding:"required"`
	Names      []LocationNameRequest `json:"names"`
	IconBase64 *string               `json:"icon_base64"`
	SortOrder  *int                  `json:"sort_order"`
	IsActive   *bool                 `json:"is_active"`
}

// CreateLocation 创建地点
//
//	@Summary		创建地点
//	@Description	创建系统预设地点，支持多语言名称
//	@Tags			后台管理-内容管理
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string									true	"Bearer {access_token}"
//	@Param			body			body		CreateLocationRequest					true	"地点信息"
//	@Success		201				{object}	response.Response{data=model.Location}	"创建成功"
//	@Failure		400				{object}	response.Response						"参数错误或名称重复"
//	@Failure		401				{object}	response.Response						"未登录"
//	@Security		BearerAuth
//	@Router			/admin/v1/content/locations/create [post]
func (h *LocationHandler) CreateLocation(c *gin.Context) {
	var input adminService.CreateLocationInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	location, err := h.locationService.CreateLocation(input)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Created(c, location)
}

// ListLocations 获取地点列表
//
//	@Summary		地点列表
//	@Description	获取系统预设地点列表，支持关键词搜索、启用状态过滤和分页，返回默认中文名称及所有语言名称
//	@Tags			后台管理-内容管理
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string													true	"Bearer {access_token}"
//	@Param			body			body		adminService.LocationListInput							true	"查询条件"
//	@Success		200				{object}	response.Response{data=adminService.LocationListResult}	"地点列表"
//	@Failure		401				{object}	response.Response										"未登录"
//	@Security		BearerAuth
//	@Router			/admin/v1/content/locations/list [post]
func (h *LocationHandler) ListLocations(c *gin.Context) {
	var input adminService.LocationListInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	result, err := h.locationService.ListLocations(input)
	if err != nil {
		response.ServerError(c)
		return
	}

	response.Success(c, result)
}

// UpdateLocation 更新地点
//
//	@Summary		更新地点
//	@Description	修改地点的多语言名称、图标、排序或启用状态
//	@Tags			后台管理-内容管理
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string									true	"Bearer {access_token}"
//	@Param			body			body		UpdateLocationRequest					true	"更新内容（需包含 id）"
//	@Success		200				{object}	response.Response{data=model.Location}	"更新后的地点"
//	@Failure		400				{object}	response.Response						"参数错误"
//	@Failure		401				{object}	response.Response						"未登录"
//	@Security		BearerAuth
//	@Router			/admin/v1/content/locations/update [post]
func (h *LocationHandler) UpdateLocation(c *gin.Context) {
	var req struct {
		ID string `json:"id" binding:"required"`
		adminService.UpdateLocationInput
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	location, err := h.locationService.UpdateLocation(req.ID, req.UpdateLocationInput)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, location)
}

// DeleteLocation 删除地点
//
//	@Summary		删除地点
//	@Description	删除指定地点，同时删除所有语言名称
//	@Tags			后台管理-内容管理
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string				true	"Bearer {access_token}"
//	@Param			body			body		ContentIDInput		true	"地点 ID"
//	@Success		200				{object}	response.Response	"删除成功"
//	@Failure		400				{object}	response.Response	"参数错误"
//	@Failure		401				{object}	response.Response	"未登录"
//	@Security		BearerAuth
//	@Router			/admin/v1/content/locations/delete [post]
func (h *LocationHandler) DeleteLocation(c *gin.Context) {
	var input ContentIDInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.locationService.DeleteLocation(input.ID); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, nil)
}
