// Package handler HTTP 请求处理层
package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/linke15d/bondy-backend/internal/service"
	"github.com/linke15d/bondy-backend/pkg/response"
)

// UserHandler 用户信息相关请求处理器
type UserHandler struct {
	userService *service.UserService
}

// NewUserHandler 创建 UserHandler 实例
func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// GetProfile 获取当前登录用户信息
//
//	@Summary		获取个人信息
//	@Description	获取当前登录用户的详细信息，包括昵称、头像、生日等
//	@Tags			用户
//	@Produce		json
//	@Param			Authorization	header		string	true	"Bearer {access_token}"
//	@Success		200				{object}	response.Response{data=model.User}	"用户信息"
//	@Failure		401				{object}	response.Response					"未登录或 token 已过期"
//	@Failure		404				{object}	response.Response					"用户不存在"
//	@Security		BearerAuth
//	@Router			/api/v1/user/profile [get]
func (h *UserHandler) GetProfile(c *gin.Context) {
	// 从 JWT 中间件注入的 context 中获取当前用户 ID
	userID := c.GetString("userID")

	user, err := h.userService.GetProfile(userID)
	if err != nil {
		response.NotFound(c, "用户不存在")
		return
	}

	response.Success(c, user)
}

// UpdateProfile 更新个人信息
//
//	@Summary		更新个人信息
//	@Description	更新当前登录用户的昵称、头像、生日等基本信息，只传需要修改的字段即可
//	@Tags			用户
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string								true	"Bearer {access_token}"
//	@Param			body			body		service.UpdateProfileInput			true	"要更新的字段"
//	@Success		200				{object}	response.Response{data=model.User}	"更新后的用户信息"
//	@Failure		400				{object}	response.Response					"参数格式错误"
//	@Failure		401				{object}	response.Response					"未登录或 token 已过期"
//	@Security		BearerAuth
//	@Router			/api/v1/user/profile [put]
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID := c.GetString("userID")

	var input service.UpdateProfileInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	user, err := h.userService.UpdateProfile(userID, input)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, user)
}

// ChangePassword 修改密码
//
//	@Summary		修改密码
//	@Description	修改当前登录用户的密码，需要提供旧密码验证身份。修改成功后所有设备的 token 将失效，需要重新登录
//	@Tags			用户
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string						true	"Bearer {access_token}"
//	@Param			body			body		service.ChangePasswordInput	true	"旧密码和新密码"
//	@Success		200				{object}	response.Response			"修改成功"
//	@Failure		400				{object}	response.Response			"参数错误或旧密码不正确"
//	@Failure		401				{object}	response.Response			"未登录或 token 已过期"
//	@Security		BearerAuth
//	@Router			/api/v1/user/password [put]
func (h *UserHandler) ChangePassword(c *gin.Context) {
	userID := c.GetString("userID")

	var input service.ChangePasswordInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.userService.ChangePassword(userID, input); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, nil)
}
