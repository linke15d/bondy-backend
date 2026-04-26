// Package handler HTTP 请求处理层
package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/linke15d/bondy-backend/internal/service"
	"github.com/linke15d/bondy-backend/pkg/response"
)

// SubscriptionHandler 订阅会员相关请求处理器
type SubscriptionHandler struct {
	subService *service.SubscriptionService
}

// NewSubscriptionHandler 创建 SubscriptionHandler 实例
func NewSubscriptionHandler(subService *service.SubscriptionService) *SubscriptionHandler {
	return &SubscriptionHandler{subService: subService}
}

// GetStatus 获取会员状态
//
//	@Summary		获取会员状态
//	@Description	获取当前用户的会员状态，包括是否是有效会员、套餐类型、过期时间等。非会员返回 is_premium: false。
//	@Tags			会员
//	@Produce		json
//	@Param			Authorization	header		string													true	"Bearer {access_token}"
//	@Success		200				{object}	response.Response{data=service.SubscriptionStatus}		"会员状态"
//	@Failure		401				{object}	response.Response										"未登录"
//	@Security		BearerAuth
//	@Router			/api/v1/subscription/status [post]
func (h *SubscriptionHandler) GetStatus(c *gin.Context) {
	userID := c.GetString("userID")

	status, err := h.subService.GetStatus(userID)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, status)
}

// Purchase 购买会员
//
//	@Summary		购买会员
//	@Description	提交支付凭证激活会员。客户端完成 Apple/Google/Stripe 支付后，将收据数据传入此接口验证并激活会员。MONTHLY 套餐有效期30天，LIFETIME 套餐永久有效。
//	@Tags			会员
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string													true	"Bearer {access_token}"
//	@Param			body			body		service.PurchaseInput									true	"购买信息和支付凭证"
//	@Success		200				{object}	response.Response{data=service.SubscriptionStatus}		"激活成功，返回最新会员状态"
//	@Failure		400				{object}	response.Response										"参数错误或收据验证失败"
//	@Failure		401				{object}	response.Response										"未登录"
//	@Security		BearerAuth
//	@Router			/api/v1/subscription/purchase [post]
func (h *SubscriptionHandler) Purchase(c *gin.Context) {
	userID := c.GetString("userID")

	var input service.PurchaseInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	status, err := h.subService.Purchase(userID, input)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, status)
}

// Cancel 取消订阅
//
//	@Summary		取消订阅
//	@Description	取消当前的订阅续费。取消后当前周期内仍然可以使用会员功能，到期后不再自动续费。买断套餐无法取消。
//	@Tags			会员
//	@Produce		json
//	@Param			Authorization	header		string				true	"Bearer {access_token}"
//	@Success		200				{object}	response.Response	"取消成功"
//	@Failure		400				{object}	response.Response	"无订阅记录或已取消"
//	@Failure		401				{object}	response.Response	"未登录"
//	@Security		BearerAuth
//	@Router			/api/v1/subscription/cancel [post]
func (h *SubscriptionHandler) Cancel(c *gin.Context) {
	userID := c.GetString("userID")

	if err := h.subService.Cancel(userID); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, nil)
}
