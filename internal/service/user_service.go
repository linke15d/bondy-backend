package service

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/linke15d/bondy-backend/internal/model"
	"github.com/linke15d/bondy-backend/internal/repository"
)

// UserService 用户信息业务逻辑
type UserService struct {
	userRepo *repository.UserRepository
}

// NewUserService 创建 UserService 实例
func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

// UpdateProfileInput 更新个人信息请求参数
// 所有字段均为可选，只传需要修改的字段
type UpdateProfileInput struct {
	// Nickname 昵称，2到20个字符
	Nickname *string `json:"nickname" binding:"omitempty,min=2,max=20" example:"小红"`

	// AvatarURL 头像图片地址，需要先通过上传接口获得 URL 再传入
	AvatarURL *string `json:"avatar_url" binding:"omitempty,url" example:"https://cdn.example.com/avatar.jpg"`

	// Birthday 生日，格式 RFC3339
	Birthday *time.Time `json:"birthday" binding:"omitempty" example:"1995-06-15T00:00:00Z"`
}

// ChangePasswordInput 修改密码请求参数
type ChangePasswordInput struct {
	// OldPassword 当前使用的旧密码，用于验证身份
	OldPassword string `json:"old_password" binding:"required" example:"oldpassword123"`

	// NewPassword 新密码，最少8位
	NewPassword string `json:"new_password" binding:"required,min=8" example:"newpassword456"`
}

// GetProfile 获取用户个人信息
func (s *UserService) GetProfile(userID string) (*model.User, error) {
	return s.userRepo.FindByID(userID)
}

// UpdateProfile 更新用户个人信息
// 只更新传入的非 nil 字段，未传的字段保持不变
func (s *UserService) UpdateProfile(userID string, input UpdateProfileInput) (*model.User, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New("用户不存在")
	}

	// 只更新传入的字段
	if input.Nickname != nil {
		user.Nickname = input.Nickname
	}
	if input.AvatarURL != nil {
		user.AvatarURL = input.AvatarURL
	}
	if input.Birthday != nil {
		user.Birthday = input.Birthday
	}

	if err := s.userRepo.Update(user); err != nil {
		return nil, errors.New("更新失败，请重试")
	}

	return user, nil
}

// ChangePassword 修改密码
// 验证旧密码后更新为新密码，并强制清除所有 refresh token（踢掉所有设备）
func (s *UserService) ChangePassword(userID string, input ChangePasswordInput) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return errors.New("用户不存在")
	}

	// 验证旧密码
	if user.PasswordHash == nil {
		return errors.New("当前账号未设置密码")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(*user.PasswordHash), []byte(input.OldPassword)); err != nil {
		return errors.New("旧密码不正确")
	}

	// 哈希新密码
	hash, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), 12)
	if err != nil {
		return errors.New("密码处理失败")
	}

	hashStr := string(hash)
	user.PasswordHash = &hashStr

	if err := s.userRepo.Update(user); err != nil {
		return errors.New("密码修改失败，请重试")
	}

	// 强制所有设备下线（清除所有 refresh token）
	_ = s.userRepo.DeleteAllUserRefreshTokens(userID)

	return nil
}
