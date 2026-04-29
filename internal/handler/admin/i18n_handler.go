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
//	@Tags			多语言管理
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
//	@Tags			多语言管理
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

// SaveTranslations 批量保存翻译
//
//	@Summary		保存翻译
//	@Description	批量保存某条记录的翻译内容，已存在的语言会更新，不存在的会新建。一次可以提交多种语言的翻译。
//	@Tags			多语言管理
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string								true	"Bearer {access_token}"
//	@Param			body			body		adminService.SaveTranslationsInput	true	"翻译内容"
//	@Success		200				{object}	response.Response					"保存成功"
//	@Failure		400				{object}	response.Response					"参数错误"
//	@Failure		401				{object}	response.Response					"未登录"
//	@Security		BearerAuth
//	@Router			/admin/v1/i18n/translations/save [post]
func (h *I18nHandler) SaveTranslations(c *gin.Context) {
	var input adminService.SaveTranslationsInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.i18nService.SaveTranslations(input); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, nil)
}

// GetTranslations 获取翻译列表
//
//	@Summary		获取翻译
//	@Description	获取某条记录某个字段的所有语言翻译，用于编辑页面回显
//	@Tags			多语言管理
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string											true	"Bearer {access_token}"
//	@Param			body			body		adminService.GetTranslationsInput				true	"查询条件"
//	@Success		200				{object}	response.Response{data=adminService.TranslationResult}	"翻译列表"
//	@Failure		401				{object}	response.Response								"未登录"
//	@Security		BearerAuth
//	@Router			/admin/v1/i18n/translations/get [post]
func (h *I18nHandler) GetTranslations(c *gin.Context) {
	var input adminService.GetTranslationsInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	result, err := h.i18nService.GetTranslations(input)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, result)
}

// GetAllTranslationsByRef 获取某条记录的全部翻译
//
//	@Summary		获取记录全部翻译
//	@Description	获取某条记录所有字段的所有语言翻译，用于编辑姿势/标签时一次性加载全部翻译内容
//	@Tags			多语言管理
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string										true	"Bearer {access_token}"
//	@Param			body			body		GetAllTranslationsRequest					true	"查询条件"
//	@Success		200				{object}	response.Response{data=[]model.Translation}	"翻译列表"
//	@Failure		401				{object}	response.Response							"未登录"
//	@Security		BearerAuth
//	@Router			/admin/v1/i18n/translations/all [post]
func (h *I18nHandler) GetAllTranslationsByRef(c *gin.Context) {
	var input GetAllTranslationsRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	result, err := h.i18nService.GetAllTranslationsByRef(input.Module, input.RefID)
	if err != nil {
		response.ServerError(c)
		return
	}

	response.Success(c, result)
}

// GetAllTranslationsRequest 获取全部翻译请求参数
type GetAllTranslationsRequest struct {
	// Module 模块名：position / tag
	Module string `json:"module" binding:"required"`
	// RefID 关联记录 ID
	RefID string `json:"ref_id" binding:"required"`
}
