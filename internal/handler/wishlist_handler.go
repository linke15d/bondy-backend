// Package handler HTTP 请求处理层
package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/linke15d/bondy-backend/internal/service"
	"github.com/linke15d/bondy-backend/pkg/response"
)

// WishlistHandler 心愿清单相关请求处理器
type WishlistHandler struct {
	wishlistService *service.WishlistService
}

// NewWishlistHandler 创建 WishlistHandler 实例
func NewWishlistHandler(wishlistService *service.WishlistService) *WishlistHandler {
	return &WishlistHandler{wishlistService: wishlistService}
}

// WishlistIDInput 通过 ID 操作心愿的请求参数
type WishlistIDInput struct {
	// ID 心愿唯一标识
	ID string `json:"id" binding:"required" example:"550e8400-e29b-41d4-a716-446655440000"`
}

// CreateWishlist 创建心愿
//
//	@Summary		创建心愿
//	@Description	添加一条新的心愿到清单。支持匿名提案，匿名时对方看不到提案人是谁。必须先绑定伴侣才能添加。
//	@Tags			心愿清单
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string											true	"Bearer {access_token}"
//	@Param			body			body		service.CreateWishlistInput						true	"心愿信息"
//	@Success		201				{object}	response.Response{data=model.Wishlist}			"创建成功"
//	@Failure		400				{object}	response.Response								"参数错误或未绑定伴侣"
//	@Failure		401				{object}	response.Response								"未登录"
//	@Security		BearerAuth
//	@Router			/api/v1/wishlist/create [post]
func (h *WishlistHandler) CreateWishlist(c *gin.Context) {
	userID := c.GetString("userID")

	var input service.CreateWishlistInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	wishlist, err := h.wishlistService.CreateWishlist(userID, input)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Created(c, wishlist)
}

// ListWishlists 获取心愿列表
//
//	@Summary		获取心愿列表
//	@Description	获取当前伴侣的心愿清单，按热度倒序排列，未完成的排在完成的前面。匿名心愿会隐藏对方的提案人信息。
//	@Tags			心愿清单
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string												true	"Bearer {access_token}"
//	@Param			body			body		service.WishlistListInput							true	"过滤条件"
//	@Success		200				{object}	response.Response{data=service.WishlistListResult}	"心愿列表"
//	@Failure		401				{object}	response.Response									"未登录"
//	@Security		BearerAuth
//	@Router			/api/v1/wishlist/list [post]
func (h *WishlistHandler) ListWishlists(c *gin.Context) {
	userID := c.GetString("userID")

	var input service.WishlistListInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	result, err := h.wishlistService.ListWishlists(userID, input)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, result)
}

// GetWishlist 获取心愿详情
//
//	@Summary		获取心愿详情
//	@Description	获取单条心愿的详细信息。匿名心愿对非提案人隐藏提案人信息。
//	@Tags			心愿清单
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string									true	"Bearer {access_token}"
//	@Param			body			body		WishlistIDInput							true	"心愿 ID"
//	@Success		200				{object}	response.Response{data=model.Wishlist}	"心愿详情"
//	@Failure		400				{object}	response.Response						"参数错误"
//	@Failure		401				{object}	response.Response						"未登录"
//	@Failure		404				{object}	response.Response						"心愿不存在"
//	@Security		BearerAuth
//	@Router			/api/v1/wishlist/detail [post]
func (h *WishlistHandler) GetWishlist(c *gin.Context) {
	userID := c.GetString("userID")

	var input WishlistIDInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	wishlist, err := h.wishlistService.GetWishlist(userID, input.ID)
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}

	response.Success(c, wishlist)
}

// UpdateWishlist 更新心愿
//
//	@Summary		更新心愿
//	@Description	修改心愿内容，只有提案人才能修改。只传需要修改的字段。
//	@Tags			心愿清单
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string									true	"Bearer {access_token}"
//	@Param			body			body		service.UpdateWishlistInput				true	"要更新的内容（需包含 id 字段）"
//	@Success		200				{object}	response.Response{data=model.Wishlist}	"更新后的心愿"
//	@Failure		400				{object}	response.Response						"参数错误或无权限"
//	@Failure		401				{object}	response.Response						"未登录"
//	@Security		BearerAuth
//	@Router			/api/v1/wishlist/update [post]
func (h *WishlistHandler) UpdateWishlist(c *gin.Context) {
	userID := c.GetString("userID")

	var req struct {
		ID string `json:"id" binding:"required"`
		service.UpdateWishlistInput
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	wishlist, err := h.wishlistService.UpdateWishlist(userID, req.ID, req.UpdateWishlistInput)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, wishlist)
}

// DeleteWishlist 删除心愿
//
//	@Summary		删除心愿
//	@Description	删除一条心愿，只有提案人才能删除，删除后不可恢复。
//	@Tags			心愿清单
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string				true	"Bearer {access_token}"
//	@Param			body			body		WishlistIDInput		true	"心愿 ID"
//	@Success		200				{object}	response.Response	"删除成功"
//	@Failure		400				{object}	response.Response	"参数错误或无权限"
//	@Failure		401				{object}	response.Response	"未登录"
//	@Security		BearerAuth
//	@Router			/api/v1/wishlist/delete [post]
func (h *WishlistHandler) DeleteWishlist(c *gin.Context) {
	userID := c.GetString("userID")

	var input WishlistIDInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.wishlistService.DeleteWishlist(userID, input.ID); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, nil)
}

// LikeWishlist 为心愿点赞
//
//	@Summary		点赞心愿
//	@Description	为心愿点赞，提升热度值。热度越高在列表中排名越靠前。双方都可以点赞，不限次数。
//	@Tags			心愿清单
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string				true	"Bearer {access_token}"
//	@Param			body			body		WishlistIDInput		true	"心愿 ID"
//	@Success		200				{object}	response.Response	"点赞成功"
//	@Failure		400				{object}	response.Response	"参数错误"
//	@Failure		401				{object}	response.Response	"未登录"
//	@Security		BearerAuth
//	@Router			/api/v1/wishlist/like [post]
func (h *WishlistHandler) LikeWishlist(c *gin.Context) {
	userID := c.GetString("userID")

	var input WishlistIDInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.wishlistService.LikeWishlist(userID, input.ID); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, nil)
}

// SetCompleted 设置完成状态
//
//	@Summary		标记完成/未完成
//	@Description	将心愿标记为已完成或未完成。双方都可以操作，完成后记录完成时间。
//	@Tags			心愿清单
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string							true	"Bearer {access_token}"
//	@Param			body			body		service.SetCompletedInput		true	"心愿 ID 和完成状态"
//	@Success		200				{object}	response.Response				"操作成功"
//	@Failure		400				{object}	response.Response				"参数错误"
//	@Failure		401				{object}	response.Response				"未登录"
//	@Security		BearerAuth
//	@Router			/api/v1/wishlist/complete [post]
func (h *WishlistHandler) SetCompleted(c *gin.Context) {
	userID := c.GetString("userID")

	var input service.SetCompletedInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.wishlistService.SetCompleted(userID, input); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, nil)
}
