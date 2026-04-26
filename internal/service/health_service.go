// Package service 业务逻辑层
package service

import (
	"errors"
	"time"

	"github.com/linke15d/bondy-backend/internal/model"
	"github.com/linke15d/bondy-backend/internal/repository"
)

// HealthService 健康记录业务逻辑
type HealthService struct {
	healthRepo *repository.HealthRepository
}

// NewHealthService 创建 HealthService 实例
func NewHealthService(healthRepo *repository.HealthRepository) *HealthService {
	return &HealthService{healthRepo: healthRepo}
}

// CreateHealthRecordInput 创建健康记录请求参数
type CreateHealthRecordInput struct {
	// Type 记录类型：STI_TEST（性病检测）或 VACCINE（疫苗接种）
	Type string `json:"type" binding:"required,oneof=STI_TEST VACCINE" example:"STI_TEST"`

	// ItemName 检测/接种项目名称，如 HIV检测、HPV疫苗第一针
	ItemName string `json:"item_name" binding:"required,max=100" example:"HIV检测"`

	// Result 检测结果
	// STI_TEST 可填：NEGATIVE（阴性）、POSITIVE（阳性）、UNKNOWN（待出结果）
	// VACCINE 可填：COMPLETED（已完成）
	Result string `json:"result" binding:"required,oneof=NEGATIVE POSITIVE UNKNOWN COMPLETED" example:"NEGATIVE"`

	// TestedAt 检测/接种日期
	TestedAt time.Time `json:"tested_at" binding:"required" example:"2024-01-15T00:00:00Z"`

	// NextRemindAt 下次提醒时间，可选。设置后到期会收到推送通知
	NextRemindAt *time.Time `json:"next_remind_at" example:"2024-07-15T00:00:00Z"`

	// NoteEncrypted 备注密文，由客户端加密后传入
	NoteEncrypted *string `json:"note_encrypted"`
}

// UpdateHealthRecordInput 更新健康记录请求参数
type UpdateHealthRecordInput struct {
	// ItemName 修改项目名称
	ItemName *string `json:"item_name" binding:"omitempty,max=100"`

	// Result 修改检测结果
	Result *string `json:"result" binding:"omitempty,oneof=NEGATIVE POSITIVE UNKNOWN COMPLETED"`

	// TestedAt 修改检测日期
	TestedAt *time.Time `json:"tested_at"`

	// NextRemindAt 修改提醒时间，传 null 表示取消提醒
	NextRemindAt *time.Time `json:"next_remind_at"`

	// NoteEncrypted 修改备注密文
	NoteEncrypted *string `json:"note_encrypted"`
}

// HealthListInput 获取健康记录列表请求参数
type HealthListInput struct {
	// Type 按类型过滤：STI_TEST / VACCINE，不传返回全部
	Type string `json:"type" binding:"omitempty,oneof=STI_TEST VACCINE" example:"STI_TEST"`

	// Page 页码，默认 1
	Page int `json:"page" binding:"omitempty,min=1" example:"1"`

	// PageSize 每页数量，默认 20
	PageSize int `json:"page_size" binding:"omitempty,min=1,max=50" example:"20"`
}

// HealthListResult 健康记录列表返回结构
type HealthListResult struct {
	// List 记录列表
	List []model.HealthRecord `json:"list"`

	// Total 总记录数
	Total int64 `json:"total"`

	// Page 当前页码
	Page int `json:"page"`

	// PageSize 每页数量
	PageSize int `json:"page_size"`
}

// CreateHealthRecord 创建健康记录
func (s *HealthService) CreateHealthRecord(userID string, input CreateHealthRecordInput) (*model.HealthRecord, error) {
	record := &model.HealthRecord{
		UserID:        userID,
		Type:          input.Type,
		ItemName:      input.ItemName,
		Result:        input.Result,
		TestedAt:      input.TestedAt,
		NextRemindAt:  input.NextRemindAt,
		NoteEncrypted: input.NoteEncrypted,
	}

	if err := s.healthRepo.Create(record); err != nil {
		return nil, errors.New("创建记录失败，请重试")
	}

	return record, nil
}

// GetHealthRecord 获取单条健康记录
// 只能查看自己的记录，不会返回伴侣的健康数据
func (s *HealthService) GetHealthRecord(userID string, recordID string) (*model.HealthRecord, error) {
	record, err := s.healthRepo.FindByID(recordID, userID)
	if err != nil {
		return nil, errors.New("记录不存在")
	}
	return record, nil
}

// ListHealthRecords 获取健康记录列表
func (s *HealthService) ListHealthRecords(userID string, input HealthListInput) (*HealthListResult, error) {
	if input.Page == 0 {
		input.Page = 1
	}
	if input.PageSize == 0 {
		input.PageSize = 20
	}

	filter := repository.HealthListFilter{
		Type:     input.Type,
		Page:     input.Page,
		PageSize: input.PageSize,
	}

	records, total, err := s.healthRepo.List(userID, filter)
	if err != nil {
		return nil, errors.New("获取记录失败")
	}

	return &HealthListResult{
		List:     records,
		Total:    total,
		Page:     input.Page,
		PageSize: input.PageSize,
	}, nil
}

// UpdateHealthRecord 更新健康记录
func (s *HealthService) UpdateHealthRecord(userID string, recordID string, input UpdateHealthRecordInput) (*model.HealthRecord, error) {
	record, err := s.healthRepo.FindByID(recordID, userID)
	if err != nil {
		return nil, errors.New("记录不存在")
	}

	if input.ItemName != nil {
		record.ItemName = *input.ItemName
	}
	if input.Result != nil {
		record.Result = *input.Result
	}
	if input.TestedAt != nil {
		record.TestedAt = *input.TestedAt
	}
	// NextRemindAt 允许传 nil 取消提醒
	if input.NextRemindAt != nil {
		record.NextRemindAt = input.NextRemindAt
	}
	if input.NoteEncrypted != nil {
		record.NoteEncrypted = input.NoteEncrypted
	}

	if err := s.healthRepo.Update(record); err != nil {
		return nil, errors.New("更新失败，请重试")
	}

	return record, nil
}

// DeleteHealthRecord 删除健康记录
func (s *HealthService) DeleteHealthRecord(userID string, recordID string) error {
	if err := s.healthRepo.Delete(recordID, userID); err != nil {
		return errors.New("删除失败，请重试")
	}
	return nil
}

// CancelReminder 取消提醒
func (s *HealthService) CancelReminder(userID string, recordID string) error {
	record, err := s.healthRepo.FindByID(recordID, userID)
	if err != nil {
		return errors.New("记录不存在")
	}

	record.NextRemindAt = nil
	return s.healthRepo.Update(record)
}
