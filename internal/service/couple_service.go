// Package service 业务逻辑层
package service

import (
	"crypto/rand"
	"errors"
	"math/big"
	"time"

	"github.com/linke15d/bondy-backend/internal/model"
	"github.com/linke15d/bondy-backend/internal/repository"
	"gorm.io/gorm"
)

// coupleCodeChars 邀请码字符集，去掉容易混淆的 0/O/I/1
const coupleCodeChars = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"

// CoupleService 伴侣配对业务逻辑
type CoupleService struct {
	coupleRepo *repository.CoupleRepository
	userRepo   *repository.UserRepository
}

// NewCoupleService 创建 CoupleService 实例
func NewCoupleService(coupleRepo *repository.CoupleRepository, userRepo *repository.UserRepository) *CoupleService {
	return &CoupleService{
		coupleRepo: coupleRepo,
		userRepo:   userRepo,
	}
}

// CoupleInfo 伴侣关系返回结构
type CoupleInfo struct {
	// ID 伴侣关系唯一 ID
	ID string `json:"id"`

	// Partner 对方的用户信息
	Partner model.User `json:"partner"`

	// CreatedAt 绑定时间
	CreatedAt time.Time `json:"created_at"`
}

// InviteCodeResult 生成邀请码的返回结构
type InviteCodeResult struct {
	// InviteCode 6位邀请码，有效期15分钟
	InviteCode string `json:"invite_code" example:"A3KP7X"`

	// ExpiresAt 邀请码过期时间，过期后需要重新生成
	ExpiresAt time.Time `json:"expires_at"`
}

// GenerateInviteCode 生成邀请码
// 如果用户已有伴侣关系则报错
// 如果用户已有待绑定记录则刷新邀请码
// 如果是新用户则创建待绑定记录
func (s *CoupleService) GenerateInviteCode(userID string) (*InviteCodeResult, error) {
	// 检查是否已有有效伴侣关系
	existing, err := s.coupleRepo.FindByUserID(userID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("服务器错误")
	}
	if existing != nil && existing.User2ID != "" {
		return nil, errors.New("你已经有伴侣了，请先解除当前关系")
	}

	// 生成6位随机邀请码
	code, err := generateCode(6)
	if err != nil {
		return nil, errors.New("邀请码生成失败")
	}

	expiresAt := time.Now().Add(15 * time.Minute)

	if existing != nil {
		// 已有待绑定记录，刷新邀请码
		if err := s.coupleRepo.UpdateInviteCode(existing.ID, code, expiresAt); err != nil {
			return nil, errors.New("邀请码更新失败")
		}
	} else {
		// 新建待绑定记录（User2ID 暂时为空）
		couple := &model.Couple{
			User1ID:         userID,
			InviteCode:      &code,
			InviteExpiresAt: &expiresAt,
		}
		if err := s.coupleRepo.Create(couple); err != nil {
			return nil, errors.New("邀请码创建失败")
		}
	}

	return &InviteCodeResult{
		InviteCode: code,
		ExpiresAt:  expiresAt,
	}, nil
}

// BindPartner 通过邀请码绑定伴侣
// 验证邀请码有效性后完成双向绑定
func (s *CoupleService) BindPartner(userID string, inviteCode string) (*CoupleInfo, error) {
	// 检查自己是否已有伴侣
	myCouple, err := s.coupleRepo.FindByUserID(userID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("服务器错误")
	}
	if myCouple != nil && myCouple.User2ID != "" {
		return nil, errors.New("你已经有伴侣了，请先解除当前关系")
	}

	// 查找邀请码对应的记录
	couple, err := s.coupleRepo.FindByInviteCode(inviteCode)
	if err != nil {
		return nil, errors.New("邀请码无效或已过期")
	}

	// 不能绑定自己
	if couple.User1ID == userID {
		return nil, errors.New("不能使用自己的邀请码")
	}

	// 完成绑定
	if err := s.coupleRepo.BindCouple(couple.ID, userID); err != nil {
		return nil, errors.New("绑定失败，请重试")
	}

	// 获取对方用户信息
	partner, err := s.userRepo.FindByID(couple.User1ID)
	if err != nil {
		return nil, errors.New("获取伴侣信息失败")
	}

	return &CoupleInfo{
		ID:        couple.ID,
		Partner:   *partner,
		CreatedAt: time.Now(),
	}, nil
}

// GetCoupleInfo 获取当前伴侣关系信息
func (s *CoupleService) GetCoupleInfo(userID string) (*CoupleInfo, error) {
	couple, err := s.coupleRepo.FindByUserID(userID)
	if err != nil {
		return nil, errors.New("暂无伴侣关系")
	}

	// 判断对方是 user1 还是 user2
	var partnerID string
	if couple.User1ID == userID {
		partnerID = couple.User2ID
	} else {
		partnerID = couple.User1ID
	}

	partner, err := s.userRepo.FindByID(partnerID)
	if err != nil {
		return nil, errors.New("获取伴侣信息失败")
	}

	return &CoupleInfo{
		ID:        couple.ID,
		Partner:   *partner,
		CreatedAt: couple.CreatedAt,
	}, nil
}

// Unlink 解除伴侣关系
// 软删除：记录解绑时间，保留历史数据
func (s *CoupleService) Unlink(userID string) error {
	couple, err := s.coupleRepo.FindByUserID(userID)
	if err != nil {
		return errors.New("暂无伴侣关系")
	}

	if err := s.coupleRepo.Unlink(couple.ID); err != nil {
		return errors.New("解除失败，请重试")
	}

	return nil
}

// generateCode 生成指定长度的随机邀请码
// 使用 crypto/rand 保证随机性，防止邀请码被预测
func generateCode(length int) (string, error) {
	result := make([]byte, length)
	for i := range result {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(coupleCodeChars))))
		if err != nil {
			return "", err
		}
		result[i] = coupleCodeChars[n.Int64()]
	}
	return string(result), nil
}
