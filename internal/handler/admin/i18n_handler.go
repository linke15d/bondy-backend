// Package admin 后台管理 Handler 层
package admin

import (
	"github.com/gin-gonic/gin"
	adminService "github.com/linke15d/bondy-backend/internal/service/admin"
	"github.com/linke15d/bondy-backend/pkg/response"
)

// I18nHandler 多语言管理请求处理器
type I18nHandler struct {
	i18nService *adminService.I18nService
}

// NewI18nHandler 创建 I18nHandler 实例
func NewI18nHandler(i18nService *adminService.I18nService) *I18nHandler {
	return &I18nHandler{i18nService: i18nService}
}

// CreateLanguage 创建语言
//
//	@Summary		创建语言
//	@Description	添加一种新的支持语言，如英语、日语等
//	@Tags			后台管理-多语言管理
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string												true	"Bearer {access_token}"
//	@Param			body			body		adminService.CreateLanguageInput					true	"语言信息"
//	@Success		201				{object}	response.Response{data=model.SupportedLanguage}		"创建成功"
//	@Failure		400				{object}	response.Response									"参数错误或语言已存在"
//	@Failure		401				{object}	response.Response									"未登录"
//	@Security		BearerAuth
//	@Router			/admin/v1/i18n/languages/create [post]
func (h *I18nHandler) CreateLanguage(c *gin.Context) {
	var input adminService.CreateLanguageInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	lang, err := h.i18nService.CreateLanguage(input)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Created(c, lang)
}

// ListLanguages 获取语言列表
//
//	@Summary		语言列表
//	@Description	获取所有支持的语言列表，按排序值升序返回
//	@Tags			后台管理-多语言管理
//	@Produce		json
//	@Param			Authorization	header		string													true	"Bearer {access_token}"
//	@Success		200				{object}	response.Response{data=[]model.SupportedLanguage}		"语言列表"
//	@Failure		401				{object}	response.Response										"未登录"
//	@Security		BearerAuth
//	@Router			/admin/v1/i18n/languages/list [post]
func (h *I18nHandler) ListLanguages(c *gin.Context) {
	languages, err := h.i18nService.ListLanguages()
	if err != nil {
		response.ServerError(c)
		return
	}
	response.Success(c, languages)
}

// UpdateLanguage 更新语言
//
//	@Summary		更新语言
//	@Description	修改语言名称、排序、启用状态或设为默认语言
//	@Tags			后台管理-多语言管理
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string												true	"Bearer {access_token}"
//	@Param			body			body		UpdateLanguageRequest								true	"更新内容（需包含 id）"
//	@Success		200				{object}	response.Response{data=model.SupportedLanguage}		"更新后的语言"
//	@Failure		400				{object}	response.Response									"参数错误"
//	@Failure		401				{object}	response.Response									"未登录"
//	@Security		BearerAuth
//	@Router			/admin/v1/i18n/languages/update [post]
func (h *I18nHandler) UpdateLanguage(c *gin.Context) {
	var req struct {
		ID string `json:"id" binding:"required"`
		adminService.UpdateLanguageInput
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	lang, err := h.i18nService.UpdateLanguage(req.ID, req.UpdateLanguageInput)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, lang)
}

// DeleteLanguage 删除语言
//
//	@Summary		删除语言
//	@Description	删除指定语言，默认语言和有翻译内容的语言不允许删除
//	@Tags			后台管理-多语言管理
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string				true	"Bearer {access_token}"
//	@Param			body			body		ContentIDInput		true	"语言 ID"
//	@Success		200				{object}	response.Response	"删除成功"
//	@Failure		400				{object}	response.Response	"默认语言不允许删除或有翻译内容"
//	@Failure		401				{object}	response.Response	"未登录"
//	@Security		BearerAuth
//	@Router			/admin/v1/i18n/languages/delete [post]
func (h *I18nHandler) DeleteLanguage(c *gin.Context) {
	var input ContentIDInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.i18nService.DeleteLanguage(input.ID); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, nil)
}

// UpdateLanguageRequest 更新语言请求（Swagger 用）
type UpdateLanguageRequest struct {
	ID        string  `json:"id" binding:"required"`
	Name      *string `json:"name"`
	IsDefault *bool   `json:"is_default"`
	IsActive  *bool   `json:"is_active"`
	SortOrder *int    `json:"sort_order"`
}
