// Package admin 后台管理 Handler 层
package admin

import (
	"github.com/gin-gonic/gin"
	_ "github.com/linke15d/bondy-backend/internal/model"
	adminService "github.com/linke15d/bondy-backend/internal/service/admin"
	"github.com/linke15d/bondy-backend/pkg/response"
)

// AdminUserHandler 用户管理请求处理器
type AdminUserHandler struct {
	userService *adminService.AdminUserService
}

// NewAdminUserHandler 创建 AdminUserHandler 实例
func NewAdminUserHandler(userService *adminService.AdminUserService) *AdminUserHandler {
	return &AdminUserHandler{userService: userService}
}

// UserIDInput 用户 ID 请求参数
type UserIDInput struct {
	// ID 用户唯一标识
	ID string `json:"id" binding:"required"`
}

// ListUsers 获取用户列表
//
//	@Summary		用户列表
//	@Description	获取所有注册用户列表，支持关键词搜索和封禁状态过滤
//	@Tags			后台管理-用户管理
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string													true	"Bearer {access_token}"
//	@Param			body			body		adminService.AdminUserListInput							true	"查询条件"
//	@Success		200				{object}	response.Response{data=adminService.AdminUserListResult}	"用户列表"
//	@Failure		401				{object}	response.Response										"未登录"
//	@Security		BearerAuth
//	@Router			/admin/v1/users/list [post]
func (h *AdminUserHandler) ListUsers(c *gin.Context) {
	var input adminService.AdminUserListInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	result, err := h.userService.ListUsers(input)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, result)
}

// GetUserDetail 获取用户详情
//
//	@Summary		用户详情
//	@Description	获取单个用户的详细信息，包括伴侣关系
//	@Tags			后台管理-用户管理
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string												true	"Bearer {access_token}"
//	@Param			body			body		UserIDInput											true	"用户 ID"
//	@Success		200				{object}	response.Response{data=adminService.AdminUserDetail}	"用户详情"
//	@Failure		401				{object}	response.Response									"未登录"
//	@Failure		404				{object}	response.Response									"用户不存在"
//	@Security		BearerAuth
//	@Router			/admin/v1/users/detail [post]
func (h *AdminUserHandler) GetUserDetail(c *gin.Context) {
	var input UserIDInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	result, err := h.userService.GetUserDetail(input.ID)
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}

	response.Success(c, result)
}

// BlockUser 封禁用户
//
//	@Summary		封禁用户
//	@Description	封禁指定用户，封禁后该用户无法登录
//	@Tags			后台管理-用户管理
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string				true	"Bearer {access_token}"
//	@Param			body			body		UserIDInput			true	"用户 ID"
//	@Success		200				{object}	response.Response	"封禁成功"
//	@Failure		400				{object}	response.Response	"用户不存在或已封禁"
//	@Failure		401				{object}	response.Response	"未登录"
//	@Security		BearerAuth
//	@Router			/admin/v1/users/block [post]
func (h *AdminUserHandler) BlockUser(c *gin.Context) {
	var input UserIDInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.userService.BlockUser(input.ID); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, nil)
}

// UnblockUser 解封用户
//
//	@Summary		解封用户
//	@Description	解除对指定用户的封禁
//	@Tags			后台管理-用户管理
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string				true	"Bearer {access_token}"
//	@Param			body			body		UserIDInput			true	"用户 ID"
//	@Success		200				{object}	response.Response	"解封成功"
//	@Failure		400				{object}	response.Response	"用户不存在或未封禁"
//	@Failure		401				{object}	response.Response	"未登录"
//	@Security		BearerAuth
//	@Router			/admin/v1/users/unblock [post]
func (h *AdminUserHandler) UnblockUser(c *gin.Context) {
	var input UserIDInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.userService.UnblockUser(input.ID); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, nil)
}

// DeleteUser 注销用户
//
//	@Summary		注销用户
//	@Description	软删除指定用户账号，注销后用户无法登录，数据保留
//	@Tags			后台管理-用户管理
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string				true	"Bearer {access_token}"
//	@Param			body			body		UserIDInput			true	"用户 ID"
//	@Success		200				{object}	response.Response	"注销成功"
//	@Failure		400				{object}	response.Response	"用户不存在"
//	@Failure		401				{object}	response.Response	"未登录"
//	@Security		BearerAuth
//	@Router			/admin/v1/users/delete [post]
func (h *AdminUserHandler) DeleteUser(c *gin.Context) {
	var input UserIDInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.userService.DeleteUser(input.ID); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, nil)
}
