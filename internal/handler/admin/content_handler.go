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

// CreateCategoryRequest Swagger 展示用
type CreateCategoryRequest struct {
	// Names 各语言名称列表，至少传一种语言
	Names []CategoryNameRequest `json:"names" example:"[{\"language_code\":\"zh-CN\",\"name\":\"经典\"}]"`
	// SortOrder 排序值，数字越小越靠前
	SortOrder int `json:"sort_order" example:"1"`
	// IsActive 修改启用状态
	IsActive *bool `json:"is_active" example:"true"`
}

// CategoryNameRequest 单个语言名称
type CategoryNameRequest struct {
	// LanguageCode 语言代码，如 zh-CN / en / ja / ko
	LanguageCode string `json:"language_code" example:"zh-CN"`
	// Name 该语言下的分类名称
	Name string `json:"name" example:"经典"`
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

// CreatePositionCategory 创建姿势分类
//
//	@Summary		创建姿势分类
//	@Description	创建一个新的姿势分类，支持后台配置多语言名称
//	@Tags			后台管理-内容管理
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string													true	"Bearer {access_token}"
//	@Param			body			body		CreateCategoryRequest						true	"分类信息"
//	@Success		201				{object}	response.Response{data=model.PositionCategory}			"创建成功"
//	@Failure		400				{object}	response.Response										"参数错误或分类代码已存在"
//	@Failure		401				{object}	response.Response										"未登录"
//	@Security		BearerAuth
//	@Router			/admin/v1/content/categories/create [post]
func (h *AdminContentHandler) CreatePositionCategory(c *gin.Context) {
	var input adminService.CreateCategoryInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	category, err := h.contentService.CreatePositionCategory(input)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Created(c, category)
}

// ListPositionCategories 获取姿势分类列表
//
//	@Summary		姿势分类列表
//	@Description	获取所有姿势分类，包含每个分类的多语言翻译
//	@Tags			后台管理-内容管理
//	@Produce		json
//	@Param			Authorization	header		string														true	"Bearer {access_token}"
//	@Success		200				{object}	response.Response{data=[]model.PositionCategory}			"分类列表"
//	@Failure		401				{object}	response.Response											"未登录"
//	@Security		BearerAuth
//	@Router			/admin/v1/content/categories/list [post]
func (h *AdminContentHandler) ListPositionCategories(c *gin.Context) {
	categories, err := h.contentService.ListPositionCategories()
	if err != nil {
		response.ServerError(c)
		return
	}
	response.Success(c, categories)
}

// UpdateCategoryRequest Swagger 展示用
type UpdateCategoryRequest struct {
	// ID 分类 ID
	ID string `json:"id" binding:"required"`
	// Names 更新各语言名称，传入的语言会覆盖，未传的语言保持不变
	Names []CategoryNameRequest `json:"names"`
	// SortOrder 修改排序
	SortOrder *int `json:"sort_order" example:"1"`
	// IsActive 修改启用状态
	IsActive *bool `json:"is_active" example:"true"`
}

// UpdatePositionCategory 更新姿势分类
//
//	@Summary		更新姿势分类
//	@Description	修改分类的多语言名称、排序或启用状态，传入的语言名称会覆盖原有内容，未传的语言保持不变
//	@Tags			内容管理
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string											true	"Bearer {access_token}"
//	@Param			body			body		UpdateCategoryRequest							true	"更新内容（需包含 id）"
//	@Success		200				{object}	response.Response{data=model.PositionCategory}	"更新后的分类"
//	@Failure		400				{object}	response.Response								"参数错误"
//	@Failure		401				{object}	response.Response								"未登录"
//	@Security		BearerAuth
//	@Router			/admin/v1/content/categories/update [post]
func (h *AdminContentHandler) UpdatePositionCategory(c *gin.Context) {
	var req struct {
		ID string `json:"id" binding:"required"`
		adminService.UpdateCategoryInput
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	category, err := h.contentService.UpdatePositionCategory(req.ID, req.UpdateCategoryInput)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, category)
}

// DeletePositionCategory 删除姿势分类
//
//	@Summary		删除姿势分类
//	@Description	删除指定分类，如果该分类下有姿势则无法删除
//	@Tags			后台管理-内容管理
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string				true	"Bearer {access_token}"
//	@Param			body			body		ContentIDInput		true	"分类 ID"
//	@Success		200				{object}	response.Response	"删除成功"
//	@Failure		400				{object}	response.Response	"分类下有姿势，无法删除"
//	@Failure		401				{object}	response.Response	"未登录"
//	@Security		BearerAuth
//	@Router			/admin/v1/content/categories/delete [post]
func (h *AdminContentHandler) DeletePositionCategory(c *gin.Context) {
	var input ContentIDInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.contentService.DeletePositionCategory(input.ID); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, nil)
}
