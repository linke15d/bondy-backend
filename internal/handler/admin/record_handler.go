// Package admin 后台管理 Handler 层
package admin

import (
	"github.com/gin-gonic/gin"
	_ "github.com/linke15d/bondy-backend/internal/model"
	adminService "github.com/linke15d/bondy-backend/internal/service/admin"
	"github.com/linke15d/bondy-backend/pkg/response"
)

// AdminRecordHandler 记录管理请求处理器
type AdminRecordHandler struct {
	recordService *adminService.AdminRecordService
}

// NewAdminRecordHandler 创建 AdminRecordHandler 实例
func NewAdminRecordHandler(recordService *adminService.AdminRecordService) *AdminRecordHandler {
	return &AdminRecordHandler{recordService: recordService}
}

// ListRecords 获取记录列表
//
//	@Summary		记录列表
//	@Description	获取所有亲密记录，可按伴侣关系过滤
//	@Tags			后台管理-记录管理
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string														true	"Bearer {access_token}"
//	@Param			body			body		adminService.AdminRecordListInput							true	"查询条件"
//	@Success		200				{object}	response.Response{data=adminService.AdminRecordListResult}	"记录列表"
//	@Failure		401				{object}	response.Response											"未登录"
//	@Security		BearerAuth
//	@Router			/admin/v1/records/list [post]
func (h *AdminRecordHandler) ListRecords(c *gin.Context) {
	var input adminService.AdminRecordListInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	result, err := h.recordService.ListRecords(input)
	if err != nil {
		response.ServerError(c)
		return
	}

	response.Success(c, result)
}
