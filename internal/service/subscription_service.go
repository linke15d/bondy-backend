// Package service 业务逻辑层
package service

import (
	"errors"
	"time"

	"github.com/linke15d/bondy-backend/internal/model"
	"github.com/linke15d/bondy-backend/internal/repository"
	"gorm.io/gorm"
)

// SubscriptionService 订阅会员业务逻辑
type SubscriptionService struct {
	subRepo  *repository.SubscriptionRepository
	userRepo *repository.UserRepository
}

// NewSubscriptionService 创建 SubscriptionService 实例
func NewSubscriptionService(subRepo *repository.SubscriptionRepository, userRepo *repository.UserRepository) *SubscriptionService {
	return &SubscriptionService{
		subRepo:  subRepo,
		userRepo: userRepo,
	}
}

// SubscriptionStatus 订阅状态返回结构
type SubscriptionStatus struct {
	// IsPremium 是否是有效会员
	IsPremium bool `json:"is_premium"`

	// Plan 当前套餐，非会员时为空
	Plan string `json:"plan,omitempty" example:"MONTHLY"`

	// Status 订阅状态，非会员时为空
	Status string `json:"status,omitempty" example:"ACTIVE"`

	// ExpiresAt 过期时间，买断套餐和非会员时为空
	ExpiresAt *time.Time `json:"expires_at,omitempty"`

	// Provider 支付渠道
	Provider string `json:"provider,omitempty" example:"apple"`
}

// PurchaseInput 购买/激活会员请求参数
type PurchaseInput struct {
	// Plan 套餐类型：MONTHLY（月付）或 LIFETIME（买断永久）
	Plan string `json:"plan" binding:"required,oneof=MONTHLY LIFETIME" example:"MONTHLY"`

	// Provider 支付渠道：apple / google / stripe
	Provider string `json:"provider" binding:"required,oneof=apple google stripe" example:"apple"`

	// ReceiptData 支付凭证，由客户端完成支付后传入
	// Apple: base64 编码的收据数据
	// Google: purchase token
	// Stripe: payment intent id
	ReceiptData string `json:"receipt_data" binding:"required" example:"MIIT..."`
}

// GetStatus 获取当前用户的会员状态
func (s *SubscriptionService) GetStatus(userID string) (*SubscriptionStatus, error) {
	sub, err := s.subRepo.FindByUserID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 没有订阅记录，返回非会员状态
			return &SubscriptionStatus{IsPremium: false}, nil
		}
		return nil, errors.New("获取会员状态失败")
	}

	// 检查是否有效
	isPremium := sub.Status == "ACTIVE" &&
		(sub.ExpiresAt == nil || sub.ExpiresAt.After(time.Now()))

	// 如果月付已过期，自动更新状态
	if sub.Status == "ACTIVE" && sub.ExpiresAt != nil && sub.ExpiresAt.Before(time.Now()) {
		_ = s.subRepo.UpdateStatus(userID, "EXPIRED")
		isPremium = false
		sub.Status = "EXPIRED"
	}

	return &SubscriptionStatus{
		IsPremium: isPremium,
		Plan:      sub.Plan,
		Status:    sub.Status,
		ExpiresAt: sub.ExpiresAt,
		Provider:  sub.Provider,
	}, nil
}

// Purchase 购买/激活会员
// 实际项目中需要调用 Apple/Google/Stripe 服务端验证收据
// 这里先做基础流程，收据验证后续接入
func (s *SubscriptionService) Purchase(userID string, input PurchaseInput) (*SubscriptionStatus, error) {
	// TODO: 根据 provider 调用对应的收据验证服务
	// Apple:  验证 App Store 收据
	// Google: 验证 Google Play purchase token
	// Stripe: 验证 payment intent
	// 验证失败直接返回错误，不激活会员

	now := time.Now()
	sub := &model.Subscription{
		UserID:        userID,
		Plan:          input.Plan,
		Status:        "ACTIVE",
		StartAt:       now,
		Provider:      input.Provider,
		ProviderSubID: &input.ReceiptData,
	}

	// 月付套餐设置30天有效期，买断套餐不设过期时间
	if input.Plan == "MONTHLY" {
		expiresAt := now.AddDate(0, 1, 0)
		sub.ExpiresAt = &expiresAt
	}

	if err := s.subRepo.Upsert(sub); err != nil {
		return nil, errors.New("会员激活失败，请重试")
	}

	return s.GetStatus(userID)
}

// Cancel 取消订阅
// 取消后当前周期仍然有效，到期后不再续费
func (s *SubscriptionService) Cancel(userID string) error {
	sub, err := s.subRepo.FindByUserID(userID)
	if err != nil {
		return errors.New("暂无订阅记录")
	}

	if sub.Status != "ACTIVE" {
		return errors.New("当前订阅已不是有效状态")
	}

	if err := s.subRepo.UpdateStatus(userID, "CANCELLED"); err != nil {
		return errors.New("取消失败，请重试")
	}

	return nil
}
