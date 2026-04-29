// Package service 业务逻辑层
package service

import (
	"errors"
	"time"

	"github.com/linke15d/bondy-backend/internal/model"
	"github.com/linke15d/bondy-backend/internal/repository"
)

// RecordService 亲密记录业务逻辑
type RecordService struct {
	recordRepo *repository.RecordRepository
	coupleRepo *repository.CoupleRepository
}

// NewRecordService 创建 RecordService 实例
func NewRecordService(recordRepo *repository.RecordRepository, coupleRepo *repository.CoupleRepository) *RecordService {
	return &RecordService{
		recordRepo: recordRepo,
		coupleRepo: coupleRepo,
	}
}

// CreateRecordInput 创建记录请求参数
type CreateRecordInput struct {
	// HappenedAt 实际发生时间，支持回填过去的记录，格式 RFC3339
	HappenedAt time.Time `json:"happened_at" binding:"required" example:"2024-01-15T20:00:00Z"`

	// DurationMins 持续时长（分钟），可选
	DurationMins *int `json:"duration_mins" example:"30"`

	// Mood 心情评分 1-5，1最差5最好，可选
	Mood *int `json:"mood" binding:"omitempty,min=1,max=5" example:"4"`

	// Satisfaction 满意度评分 1-5，1最差5最好，可选
	Satisfaction *int `json:"satisfaction" binding:"omitempty,min=1,max=5" example:"5"`

	// NoteEncrypted 备注密文，由客户端加密后传入，后端原样存储不解密
	NoteEncrypted *string `json:"note_encrypted" example:"U2FsdGVkX1..."`

	// TagIDs 标签 ID 列表，可同时关联多个标签
	TagIDs []string `json:"tag_ids" example:"['uuid1','uuid2']"`

	// PositionIDs 姿势 ID 列表，可同时关联多个姿势
	PositionIDs []string `json:"position_ids" example:"['uuid1']"`
}

// UpdateRecordInput 更新记录请求参数
// 所有字段均为可选，只传需要修改的字段
type UpdateRecordInput struct {
	// HappenedAt 修改发生时间
	HappenedAt *time.Time `json:"happened_at" example:"2024-01-15T21:00:00Z"`

	// DurationMins 修改时长
	DurationMins *int `json:"duration_mins" example:"45"`

	// Mood 修改心情评分
	Mood *int `json:"mood" binding:"omitempty,min=1,max=5" example:"3"`

	// Satisfaction 修改满意度评分
	Satisfaction *int `json:"satisfaction" binding:"omitempty,min=1,max=5" example:"4"`

	// NoteEncrypted 修改备注密文
	NoteEncrypted *string `json:"note_encrypted"`

	// TagIDs 重新设置标签列表（全量替换）
	TagIDs []string `json:"tag_ids"`

	// PositionIDs 重新设置姿势列表（全量替换）
	PositionIDs []string `json:"position_ids"`
}

// RecordListInput 获取记录列表请求参数
type RecordListInput struct {
	// Page 页码，从 1 开始，默认 1
	Page int `json:"page" binding:"omitempty,min=1" example:"1"`

	// PageSize 每页数量，默认 20，最大 50
	PageSize int `json:"page_size" binding:"omitempty,min=1,max=50" example:"20"`

	// Year 按年份过滤，不传则返回全部
	Year int `json:"year" example:"2024"`

	// Month 按月份过滤，不传则返回全部，需要和 year 一起使用
	Month int `json:"month" binding:"omitempty,min=1,max=12" example:"1"`
}

// RecordListResult 记录列表返回结构
type RecordListResult struct {
	// List 记录列表
	List []model.Record `json:"list"`

	// Total 总记录数
	Total int64 `json:"total"`

	// Page 当前页码
	Page int `json:"page"`

	// PageSize 每页数量
	PageSize int `json:"page_size"`
}

// CreateRecord 创建亲密记录
// 只有已绑定伴侣的用户才能创建记录
func (s *RecordService) CreateRecord(userID string, input CreateRecordInput) (*model.Record, error) {
	// 获取伴侣关系，未绑定则无法创建
	couple, err := s.coupleRepo.FindByUserID(userID)
	if err != nil {
		return nil, errors.New("请先绑定伴侣后再创建记录")
	}

	record := &model.Record{
		CoupleID:      couple.ID,
		CreatedByID:   userID,
		HappenedAt:    input.HappenedAt,
		DurationMins:  input.DurationMins,
		Mood:          input.Mood,
		Satisfaction:  input.Satisfaction,
		NoteEncrypted: input.NoteEncrypted,
	}

	// 处理标签关联
	if len(input.TagIDs) > 0 {
		var tags []model.Tag
		for _, id := range input.TagIDs {
			tags = append(tags, model.Tag{ID: id})
		}
		record.Tags = tags
	}

	// 处理姿势关联
	if len(input.PositionIDs) > 0 {
		var positions []model.Position
		for _, id := range input.PositionIDs {
			positions = append(positions, model.Position{ID: id})
		}
		record.Positions = positions
	}

	if err := s.recordRepo.Create(record); err != nil {
		return nil, errors.New("创建记录失败，请重试")
	}

	return record, nil
}

// GetRecord 获取单条记录详情
// 只能查看自己伴侣关系下的记录
func (s *RecordService) GetRecord(userID string, recordID string) (*model.Record, error) {
	couple, err := s.coupleRepo.FindByUserID(userID)
	if err != nil {
		return nil, errors.New("暂无伴侣关系")
	}

	record, err := s.recordRepo.FindByID(recordID, couple.ID)
	if err != nil {
		return nil, errors.New("记录不存在")
	}

	return record, nil
}

// ListRecords 获取记录列表
func (s *RecordService) ListRecords(userID string, input RecordListInput) (*RecordListResult, error) {
	couple, err := s.coupleRepo.FindByUserID(userID)
	if err != nil {
		return nil, errors.New("暂无伴侣关系")
	}

	// 设置默认值
	if input.Page == 0 {
		input.Page = 1
	}
	if input.PageSize == 0 {
		input.PageSize = 20
	}

	filter := repository.ListFilter{
		Page:     input.Page,
		PageSize: input.PageSize,
		Year:     input.Year,
		Month:    input.Month,
	}

	records, total, err := s.recordRepo.List(couple.ID, filter)
	if err != nil {
		return nil, errors.New("获取记录失败")
	}

	return &RecordListResult{
		List:     records,
		Total:    total,
		Page:     input.Page,
		PageSize: input.PageSize,
	}, nil
}

// UpdateRecord 更新记录
// 只能更新自己伴侣关系下的记录
func (s *RecordService) UpdateRecord(userID string, recordID string, input UpdateRecordInput) (*model.Record, error) {
	couple, err := s.coupleRepo.FindByUserID(userID)
	if err != nil {
		return nil, errors.New("暂无伴侣关系")
	}

	record, err := s.recordRepo.FindByID(recordID, couple.ID)
	if err != nil {
		return nil, errors.New("记录不存在")
	}

	// 只更新传入的字段
	if input.HappenedAt != nil {
		record.HappenedAt = *input.HappenedAt
	}
	if input.DurationMins != nil {
		record.DurationMins = input.DurationMins
	}
	if input.Mood != nil {
		record.Mood = input.Mood
	}
	if input.Satisfaction != nil {
		record.Satisfaction = input.Satisfaction
	}
	if input.NoteEncrypted != nil {
		record.NoteEncrypted = input.NoteEncrypted
	}

	// 全量替换标签
	if input.TagIDs != nil {
		var tags []model.Tag
		for _, id := range input.TagIDs {
			tags = append(tags, model.Tag{ID: id})
		}
		record.Tags = tags
	}

	// 全量替换姿势
	if input.PositionIDs != nil {
		var positions []model.Position
		for _, id := range input.PositionIDs {
			positions = append(positions, model.Position{ID: id})
		}
		record.Positions = positions
	}

	if err := s.recordRepo.Update(record); err != nil {
		return nil, errors.New("更新失败，请重试")
	}

	return record, nil
}

// DeleteRecord 删除记录（软删除）
func (s *RecordService) DeleteRecord(userID string, recordID string) error {
	couple, err := s.coupleRepo.FindByUserID(userID)
	if err != nil {
		return errors.New("暂无伴侣关系")
	}

	if err := s.recordRepo.SoftDelete(recordID, couple.ID); err != nil {
		return errors.New("删除失败，请重试")
	}

	return nil
}

// GetTags 获取标签列表
func (s *RecordService) GetTags(userID string, tagType string) ([]model.Tag, error) {
	couple, err := s.coupleRepo.FindByUserID(userID)
	if err != nil {
		return nil, errors.New("暂无伴侣关系")
	}
	return s.recordRepo.FindTagsByType(tagType, couple.ID)
}

// GetPositions 获取姿势列表
func (s *RecordService) GetPositions(userID string, lang string) ([]repository.PositionsByCategory, error) {
	couple, err := s.coupleRepo.FindByUserID(userID)
	if err != nil {
		return nil, errors.New("暂无伴侣关系")
	}
	return s.recordRepo.FindPositionsByCategory(couple.ID, lang)
}

// GetPositionCategories 获取姿势分类列表
func (s *RecordService) GetPositionCategories(lang string) ([]model.PositionCategory, error) {
	return s.recordRepo.FindCategories(lang)
}
