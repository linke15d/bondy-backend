// Package admin 后台管理业务逻辑层
package admin

import (
	"errors"
	"time"

	"github.com/linke15d/bondy-backend/internal/model"
	"github.com/linke15d/bondy-backend/internal/repository"
)

// AdminUserService 用户管理业务逻辑
type AdminUserService struct {
	userRepo   *repository.UserRepository
	coupleRepo *repository.CoupleRepository
}

// NewAdminUserService 创建 AdminUserService 实例
func NewAdminUserService(userRepo *repository.UserRepository, coupleRepo *repository.CoupleRepository) *AdminUserService {
	return &AdminUserService{
		userRepo:   userRepo,
		coupleRepo: coupleRepo,
	}
}

// AdminUserListInput 用户列表查询参数
type AdminUserListInput struct {
	// Keyword 搜索关键词，匹配邮箱或昵称
	Keyword string `json:"keyword" example:"test"`

	// IsBlocked 按封禁状态过滤，不传返回全部
	IsBlocked *bool `json:"is_blocked"`

	// Page 页码
	Page int `json:"page" binding:"omitempty,min=1" example:"1"`

	// PageSize 每页数量
	PageSize int `json:"page_size" binding:"omitempty,min=1,max=100" example:"20"`
}

// AdminUserListResult 用户列表返回结构
type AdminUserListResult struct {
	// List 用户列表
	List []model.User `json:"list"`

	// Total 总用户数
	Total int64 `json:"total"`

	// Page 当前页码
	Page int `json:"page"`

	// PageSize 每页数量
	PageSize int `json:"page_size"`
}

// AdminUserDetail 用户详情（含伴侣信息）
type AdminUserDetail struct {
	// User 用户基本信息
	User model.User `json:"user"`

	// Couple 伴侣关系信息，没有伴侣时为 null
	Couple *model.Couple `json:"couple,omitempty"`
}

// ListUsers 获取用户列表
func (s *AdminUserService) ListUsers(input AdminUserListInput) (*AdminUserListResult, error) {
	if input.Page == 0 {
		input.Page = 1
	}
	if input.PageSize == 0 {
		input.PageSize = 20
	}

	users, total, err := s.userRepo.AdminList(input.Keyword, input.IsBlocked, input.Page, input.PageSize)
	if err != nil {
		return nil, errors.New("获取用户列表失败")
	}

	return &AdminUserListResult{
		List:     users,
		Total:    total,
		Page:     input.Page,
		PageSize: input.PageSize,
	}, nil
}

// GetUserDetail 获取用户详情
func (s *AdminUserService) GetUserDetail(userID string) (*AdminUserDetail, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New("用户不存在")
	}

	couple, _ := s.coupleRepo.FindByUserID(userID)

	return &AdminUserDetail{
		User:   *user,
		Couple: couple,
	}, nil
}

// BlockUser 封禁用户
func (s *AdminUserService) BlockUser(userID string) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return errors.New("用户不存在")
	}

	if user.IsBlocked {
		return errors.New("该用户已被封禁")
	}

	user.IsBlocked = true
	return s.userRepo.Update(user)
}

// UnblockUser 解封用户
func (s *AdminUserService) UnblockUser(userID string) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return errors.New("用户不存在")
	}

	if !user.IsBlocked {
		return errors.New("该用户未被封禁")
	}

	user.IsBlocked = false
	return s.userRepo.Update(user)
}

// DeleteUser 注销用户（软删除）
func (s *AdminUserService) DeleteUser(userID string) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return errors.New("用户不存在")
	}

	now := time.Now()
	user.DeletedAt = &now
	return s.userRepo.Update(user)
}
