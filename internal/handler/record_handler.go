// Package handler HTTP 请求处理层
package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/linke15d/bondy-backend/internal/service"
	"github.com/linke15d/bondy-backend/pkg/response"
)

// RecordHandler 亲密记录相关请求处理器
type RecordHandler struct {
	recordService *service.RecordService
}

// NewRecordHandler 创建 RecordHandler 实例
func NewRecordHandler(recordService *service.RecordService) *RecordHandler {
	return &RecordHandler{recordService: recordService}
}

// RecordIDInput 通过 ID 操作记录的请求参数
type RecordIDInput struct {
	// ID 记录唯一标识
	ID string `json:"id" binding:"required" example:"550e8400-e29b-41d4-a716-446655440000"`
}

// CreateRecord 创建亲密记录
//
//	@Summary		创建记录
//	@Description	创建一条新的亲密记录。必须先绑定伴侣才能创建记录。备注内容请在客户端加密后传入密文，后端不会解密存储。
//	@Tags			亲密记录
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string											true	"Bearer {access_token}"
//	@Param			body			body		service.CreateRecordInput						true	"记录内容"
//	@Success		201				{object}	response.Response{data=model.Record}			"创建成功"
//	@Failure		400				{object}	response.Response								"参数错误或未绑定伴侣"
//	@Failure		401				{object}	response.Response								"未登录"
//	@Security		BearerAuth
//	@Router			/api/v1/records/create [post]
func (h *RecordHandler) CreateRecord(c *gin.Context) {
	userID := c.GetString("userID")

	var input service.CreateRecordInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	record, err := h.recordService.CreateRecord(userID, input)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Created(c, record)
}

// GetRecord 获取记录详情
//
//	@Summary		获取记录详情
//	@Description	获取单条亲密记录的详细信息，包括关联的标签和姿势。只能查看自己伴侣关系下的记录。
//	@Tags			亲密记录
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string									true	"Bearer {access_token}"
//	@Param			body			body		RecordIDInput							true	"记录 ID"
//	@Success		200				{object}	response.Response{data=model.Record}	"记录详情"
//	@Failure		400				{object}	response.Response						"参数错误"
//	@Failure		401				{object}	response.Response						"未登录"
//	@Failure		404				{object}	response.Response						"记录不存在"
//	@Security		BearerAuth
//	@Router			/api/v1/records/detail [post]
func (h *RecordHandler) GetRecord(c *gin.Context) {
	userID := c.GetString("userID")

	var input RecordIDInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	record, err := h.recordService.GetRecord(userID, input.ID)
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}

	response.Success(c, record)
}

// ListRecords 获取记录列表
//
//	@Summary		获取记录列表
//	@Description	分页获取当前伴侣的亲密记录列表，按发生时间倒序排列。支持按年、月过滤。
//	@Tags			亲密记录
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string													true	"Bearer {access_token}"
//	@Param			body			body		service.RecordListInput									true	"查询条件"
//	@Success		200				{object}	response.Response{data=service.RecordListResult}		"记录列表"
//	@Failure		400				{object}	response.Response										"参数错误"
//	@Failure		401				{object}	response.Response										"未登录"
//	@Security		BearerAuth
//	@Router			/api/v1/records/list [post]
func (h *RecordHandler) ListRecords(c *gin.Context) {
	userID := c.GetString("userID")

	var input service.RecordListInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	result, err := h.recordService.ListRecords(userID, input)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, result)
}

// UpdateRecord 更新记录
//
//	@Summary		更新记录
//	@Description	更新一条亲密记录，只传需要修改的字段。标签和姿势列表为全量替换，传空数组表示清除所有关联。
//	@Tags			亲密记录
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string									true	"Bearer {access_token}"
//	@Param			body			body		service.UpdateRecordInput				true	"要更新的内容（需包含 id 字段）"
//	@Success		200				{object}	response.Response{data=model.Record}	"更新后的记录"
//	@Failure		400				{object}	response.Response						"参数错误或记录不存在"
//	@Failure		401				{object}	response.Response						"未登录"
//	@Security		BearerAuth
//	@Router			/api/v1/records/update [post]
func (h *RecordHandler) UpdateRecord(c *gin.Context) {
	userID := c.GetString("userID")

	// 更新请求需要包含 ID
	var req struct {
		ID string `json:"id" binding:"required"`
		service.UpdateRecordInput
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	record, err := h.recordService.UpdateRecord(userID, req.ID, req.UpdateRecordInput)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, record)
}

// DeleteRecord 删除记录
//
//	@Summary		删除记录
//	@Description	软删除一条亲密记录，删除后不可恢复。只能删除自己伴侣关系下的记录。
//	@Tags			亲密记录
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string				true	"Bearer {access_token}"
//	@Param			body			body		RecordIDInput		true	"记录 ID"
//	@Success		200				{object}	response.Response	"删除成功"
//	@Failure		400				{object}	response.Response	"参数错误或记录不存在"
//	@Failure		401				{object}	response.Response	"未登录"
//	@Security		BearerAuth
//	@Router			/api/v1/records/delete [post]
func (h *RecordHandler) DeleteRecord(c *gin.Context) {
	userID := c.GetString("userID")

	var input RecordIDInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.recordService.DeleteRecord(userID, input.ID); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, nil)
}

// GetTags 获取标签列表
//
// @Summary      获取标签列表
// @Tags         亲密记录
// @Produce      json
// @Param        Authorization  header    string                                    true   "Bearer {access_token}"
// @Param        lang           query     string                                    false  "语言代码，默认 zh-CN"
// @Success      200            {object}  response.Response{data=[]model.Tag}
// @Router       /api/v1/records/tags [post]
func (h *RecordHandler) GetTags(c *gin.Context) {
	lang := c.Query("lang")
	if lang == "" {
		lang = "zh-CN"
	}
	tags, err := h.recordService.GetTags(lang)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, tags)
}

// GetPositions 获取姿势列表
//
//	@Summary		获取姿势列表
//	@Description	获取可用的姿势列表，包括系统预设姿势和当前伴侣自定义姿势
//	@Tags			亲密记录
//	@Produce		json
//	@Param			Authorization	header		string									true	"Bearer {access_token}"
//	@Success		200				{object}	response.Response{data=[]model.Position}	"姿势列表"
//	@Failure		401				{object}	response.Response						"未登录"
//	@Security		BearerAuth
//	@Router			/api/v1/records/positions [post]
func (h *RecordHandler) GetPositions(c *gin.Context) {
	userID := c.GetString("userID")

	// 从 context 取语言，默认中文
	lang := c.GetString("lang")
	if lang == "" {
		lang = "zh-CN"
	}

	positions, err := h.recordService.GetPositions(userID, lang)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, positions)
}

// GetTagsInput 获取标签列表请求参数
type GetTagsInput struct {
	// Type 标签类型：LOCATION（地点标签）或 ACTIVITY（活动标签）
	Type string `json:"type" binding:"required,oneof=LOCATION ACTIVITY" example:"LOCATION"`
}

// GetPositionCategories 获取姿势分类列表
//
//	@Summary		获取姿势分类列表
//	@Description	获取所有启用的姿势分类，根据 Accept-Language Header 返回对应语言的分类名称
//	@Tags			亲密记录
//	@Produce		json
//	@Param			Authorization	header		string													true	"Bearer {access_token}"
//	@Param			Accept-Language	header		string													false	"语言代码，如 zh-CN / en，默认 zh-CN"
//	@Success		200				{object}	response.Response{data=[]model.PositionCategory}		"分类列表"
//	@Failure		401				{object}	response.Response										"未登录"
//	@Security		BearerAuth
//	@Router			/api/v1/records/categories [post]
func (h *RecordHandler) GetPositionCategories(c *gin.Context) {
	lang := c.GetString("lang")
	if lang == "" {
		lang = "zh-CN"
	}

	categories, err := h.recordService.GetPositionCategories(lang)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, categories)
}

// GetLocations 获取地点列表
//
//	@Summary		获取地点列表
//	@Description	获取所有启用的系统预设地点，根据 Accept-Language Header 返回对应语言名称
//	@Tags			亲密记录
//	@Produce		json
//	@Param			Authorization	header		string									true	"Bearer {access_token}"
//	@Param			Accept-Language	header		string									false	"语言代码，默认 zh-CN"
//	@Success		200				{object}	response.Response{data=[]model.Location}	"地点列表"
//	@Failure		401				{object}	response.Response						"未登录"
//	@Security		BearerAuth
//	@Router			/api/v1/records/locations [post]
func (h *RecordHandler) GetLocations(c *gin.Context) {
	lang := c.GetString("lang")
	if lang == "" {
		lang = "zh-CN"
	}

	locations, err := h.recordService.GetLocations(lang)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, locations)
}
