// Package handler HTTP 请求处理层
package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/linke15d/bondy-backend/internal/service"
	"github.com/linke15d/bondy-backend/pkg/response"
)

// CoupleHandler 伴侣配对相关请求处理器
type CoupleHandler struct {
	coupleService *service.CoupleService
}

// NewCoupleHandler 创建 CoupleHandler 实例
func NewCoupleHandler(coupleService *service.CoupleService) *CoupleHandler {
	return &CoupleHandler{coupleService: coupleService}
}

// BindInput 绑定伴侣请求参数
type BindInput struct {
	// InviteCode 对方生成的6位邀请码
	InviteCode string `json:"invite_code" binding:"required,len=6" example:"A3KP7X"`
}

// GenerateInviteCode 生成邀请码
//
//	@Summary		生成邀请码
//	@Description	生成一个6位邀请码，有效期15分钟。将邀请码分享给伴侣，对方通过邀请码完成绑定。如果已有未使用的邀请码，调用此接口会刷新邀请码并重置有效期。
//	@Tags			伴侣
//	@Produce		json
//	@Param			Authorization	header		string											true	"Bearer {access_token}"
//	@Success		200				{object}	response.Response{data=service.InviteCodeResult}	"邀请码信息"
//	@Failure		400				{object}	response.Response								"已有伴侣，无法生成邀请码"
//	@Failure		401				{object}	response.Response								"未登录"
//	@Security		BearerAuth
//	@Router			/api/v1/couple/invite [post]
func (h *CoupleHandler) GenerateInviteCode(c *gin.Context) {
	userID := c.GetString("userID")

	result, err := h.coupleService.GenerateInviteCode(userID)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, result)
}

// BindPartner 绑定伴侣
//
//	@Summary		绑定伴侣
//	@Description	输入对方生成的邀请码完成伴侣绑定。绑定成功后，双方可以共享亲密记录、心愿清单等数据。每个用户同时只能有一段有效的伴侣关系。
//	@Tags			伴侣
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string										true	"Bearer {access_token}"
//	@Param			body			body		BindInput									true	"邀请码"
//	@Success		200				{object}	response.Response{data=service.CoupleInfo}	"绑定成功，返回伴侣信息"
//	@Failure		400				{object}	response.Response							"邀请码无效/已过期/不能绑定自己/已有伴侣"
//	@Failure		401				{object}	response.Response							"未登录"
//	@Security		BearerAuth
//	@Router			/api/v1/couple/bind [post]
func (h *CoupleHandler) BindPartner(c *gin.Context) {
	userID := c.GetString("userID")

	var input BindInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, "邀请码格式错误，必须是6位字符")
		return
	}

	result, err := h.coupleService.BindPartner(userID, input.InviteCode)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, result)
}

// GetCoupleInfo 获取伴侣信息
//
//	@Summary		获取伴侣信息
//	@Description	获取当前用户的伴侣关系信息，包括对方的用户资料和绑定时间
//	@Tags			伴侣
//	@Produce		json
//	@Param			Authorization	header		string									true	"Bearer {access_token}"
//	@Success		200				{object}	response.Response{data=service.CoupleInfo}	"伴侣信息"
//	@Failure		401				{object}	response.Response						"未登录"
//	@Failure		404				{object}	response.Response						"暂无伴侣关系"
//	@Security		BearerAuth
//	@Router			/api/v1/couple/info [post]
func (h *CoupleHandler) GetCoupleInfo(c *gin.Context) {
	userID := c.GetString("userID")

	info, err := h.coupleService.GetCoupleInfo(userID)
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}

	response.Success(c, info)
}

// Unlink 解除伴侣关系
//
//	@Summary		解除伴侣关系
//	@Description	解除与当前伴侣的绑定关系。解除后历史记录仍然保留，但不再共享新数据。此操作不可撤销。
//	@Tags			伴侣
//	@Produce		json
//	@Param			Authorization	header		string				true	"Bearer {access_token}"
//	@Success		200				{object}	response.Response	"解除成功"
//	@Failure		400				{object}	response.Response	"暂无伴侣关系"
//	@Failure		401				{object}	response.Response	"未登录"
//	@Security		BearerAuth
//	@Router			/api/v1/couple/unbind [post]
func (h *CoupleHandler) Unlink(c *gin.Context) {
	userID := c.GetString("userID")

	if err := h.coupleService.Unlink(userID); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, nil)
}
