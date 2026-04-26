// Package service 业务逻辑层
package service

import (
	"errors"

	"github.com/linke15d/bondy-backend/internal/model"
	"github.com/linke15d/bondy-backend/internal/repository"
)

// WishlistService 心愿清单业务逻辑
type WishlistService struct {
	wishlistRepo *repository.WishlistRepository
	coupleRepo   *repository.CoupleRepository
}

// NewWishlistService 创建 WishlistService 实例
func NewWishlistService(wishlistRepo *repository.WishlistRepository, coupleRepo *repository.CoupleRepository) *WishlistService {
	return &WishlistService{
		wishlistRepo: wishlistRepo,
		coupleRepo:   coupleRepo,
	}
}

// CreateWishlistInput 创建心愿请求参数
type CreateWishlistInput struct {
	// Title 心愿标题，最多100个字符
	Title string `json:"title" binding:"required,max=100" example:"去海边看日出"`

	// Description 心愿详细描述，可选
	Description *string `json:"description" binding:"omitempty,max=500" example:"找一个好天气，一起去看日出"`

	// IsAnonymous 是否匿名提案，匿名时对方看不到提案人
	IsAnonymous bool `json:"is_anonymous" example:"false"`

	// Scope 心愿范围：COUPLE（共同心愿）或 PERSONAL（个人心愿）
	Scope string `json:"scope" binding:"required,oneof=COUPLE PERSONAL" example:"COUPLE"`
}

// UpdateWishlistInput 更新心愿请求参数
type UpdateWishlistInput struct {
	// Title 修改标题
	Title *string `json:"title" binding:"omitempty,max=100" example:"去山里露营"`

	// Description 修改描述
	Description *string `json:"description" binding:"omitempty,max=500"`

	// IsAnonymous 修改是否匿名
	IsAnonymous *bool `json:"is_anonymous"`

	// Scope 修改范围
	Scope *string `json:"scope" binding:"omitempty,oneof=COUPLE PERSONAL"`
}

// WishlistListInput 获取心愿列表请求参数
type WishlistListInput struct {
	// Scope 按范围过滤：COUPLE / PERSONAL，不传返回全部
	Scope string `json:"scope" binding:"omitempty,oneof=COUPLE PERSONAL" example:"COUPLE"`

	// IsCompleted 按完成状态过滤，不传返回全部
	IsCompleted *bool `json:"is_completed" example:"false"`

	// Page 页码，默认 1
	Page int `json:"page" binding:"omitempty,min=1" example:"1"`

	// PageSize 每页数量，默认 20
	PageSize int `json:"page_size" binding:"omitempty,min=1,max=50" example:"20"`
}

// WishlistListResult 心愿列表返回结构
type WishlistListResult struct {
	// List 心愿列表
	List []model.Wishlist `json:"list"`

	// Total 总数量
	Total int64 `json:"total"`

	// Page 当前页码
	Page int `json:"page"`

	// PageSize 每页数量
	PageSize int `json:"page_size"`
}

// SetCompletedInput 设置完成状态请求参数
type SetCompletedInput struct {
	// ID 心愿 ID
	ID string `json:"id" binding:"required"`

	// IsCompleted 是否完成
	IsCompleted bool `json:"is_completed" example:"true"`
}

// CreateWishlist 创建心愿
func (s *WishlistService) CreateWishlist(userID string, input CreateWishlistInput) (*model.Wishlist, error) {
	couple, err := s.coupleRepo.FindByUserID(userID)
	if err != nil {
		return nil, errors.New("请先绑定伴侣后再添加心愿")
	}

	wishlist := &model.Wishlist{
		CoupleID:    couple.ID,
		CreatedByID: userID,
		Title:       input.Title,
		Description: input.Description,
		IsAnonymous: input.IsAnonymous,
		Scope:       input.Scope,
	}

	if err := s.wishlistRepo.Create(wishlist); err != nil {
		return nil, errors.New("创建心愿失败，请重试")
	}

	return wishlist, nil
}

// GetWishlist 获取单个心愿详情
// 匿名心愿对非提案人隐藏 created_by_id
func (s *WishlistService) GetWishlist(userID string, wishlistID string) (*model.Wishlist, error) {
	couple, err := s.coupleRepo.FindByUserID(userID)
	if err != nil {
		return nil, errors.New("暂无伴侣关系")
	}

	wishlist, err := s.wishlistRepo.FindByID(wishlistID, couple.ID)
	if err != nil {
		return nil, errors.New("心愿不存在")
	}

	// 匿名心愿对非提案人隐藏提案人 ID
	if wishlist.IsAnonymous && wishlist.CreatedByID != userID {
		wishlist.CreatedByID = ""
	}

	return wishlist, nil
}

// ListWishlists 获取心愿列表
func (s *WishlistService) ListWishlists(userID string, input WishlistListInput) (*WishlistListResult, error) {
	couple, err := s.coupleRepo.FindByUserID(userID)
	if err != nil {
		return nil, errors.New("暂无伴侣关系")
	}

	if input.Page == 0 {
		input.Page = 1
	}
	if input.PageSize == 0 {
		input.PageSize = 20
	}

	filter := repository.WishlistFilter{
		Scope:       input.Scope,
		IsCompleted: input.IsCompleted,
		Page:        input.Page,
		PageSize:    input.PageSize,
	}

	wishlists, total, err := s.wishlistRepo.List(couple.ID, filter)
	if err != nil {
		return nil, errors.New("获取心愿列表失败")
	}

	// 对匿名心愿隐藏非本人的提案人 ID
	for i := range wishlists {
		if wishlists[i].IsAnonymous && wishlists[i].CreatedByID != userID {
			wishlists[i].CreatedByID = ""
		}
	}

	return &WishlistListResult{
		List:     wishlists,
		Total:    total,
		Page:     input.Page,
		PageSize: input.PageSize,
	}, nil
}

// UpdateWishlist 更新心愿
// 只有提案人才能修改
func (s *WishlistService) UpdateWishlist(userID string, wishlistID string, input UpdateWishlistInput) (*model.Wishlist, error) {
	couple, err := s.coupleRepo.FindByUserID(userID)
	if err != nil {
		return nil, errors.New("暂无伴侣关系")
	}

	wishlist, err := s.wishlistRepo.FindByID(wishlistID, couple.ID)
	if err != nil {
		return nil, errors.New("心愿不存在")
	}

	// 只有提案人才能修改
	if wishlist.CreatedByID != userID {
		return nil, errors.New("只有提案人才能修改心愿")
	}

	if input.Title != nil {
		wishlist.Title = *input.Title
	}
	if input.Description != nil {
		wishlist.Description = input.Description
	}
	if input.IsAnonymous != nil {
		wishlist.IsAnonymous = *input.IsAnonymous
	}
	if input.Scope != nil {
		wishlist.Scope = *input.Scope
	}

	if err := s.wishlistRepo.Update(wishlist); err != nil {
		return nil, errors.New("更新失败，请重试")
	}

	return wishlist, nil
}

// DeleteWishlist 删除心愿
// 只有提案人才能删除
func (s *WishlistService) DeleteWishlist(userID string, wishlistID string) error {
	couple, err := s.coupleRepo.FindByUserID(userID)
	if err != nil {
		return errors.New("暂无伴侣关系")
	}

	wishlist, err := s.wishlistRepo.FindByID(wishlistID, couple.ID)
	if err != nil {
		return errors.New("心愿不存在")
	}

	if wishlist.CreatedByID != userID {
		return errors.New("只有提案人才能删除心愿")
	}

	return s.wishlistRepo.Delete(wishlistID, couple.ID)
}

// LikeWishlist 为心愿点赞（提升热度）
// 双方都可以点赞，不限制次数
func (s *WishlistService) LikeWishlist(userID string, wishlistID string) error {
	couple, err := s.coupleRepo.FindByUserID(userID)
	if err != nil {
		return errors.New("暂无伴侣关系")
	}

	if err := s.wishlistRepo.IncrHeat(wishlistID, couple.ID); err != nil {
		return errors.New("操作失败，请重试")
	}

	return nil
}

// SetCompleted 标记心愿完成或未完成
// 双方都可以操作
func (s *WishlistService) SetCompleted(userID string, input SetCompletedInput) error {
	couple, err := s.coupleRepo.FindByUserID(userID)
	if err != nil {
		return errors.New("暂无伴侣关系")
	}

	if err := s.wishlistRepo.SetCompleted(input.ID, couple.ID, input.IsCompleted); err != nil {
		return errors.New("操作失败，请重试")
	}

	return nil
}
