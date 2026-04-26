// Package admin 后台管理业务逻辑层
package admin

import (
	"github.com/linke15d/bondy-backend/internal/model"
	"github.com/linke15d/bondy-backend/internal/repository"
)

// AdminRecordService 记录管理业务逻辑
type AdminRecordService struct {
	recordRepo *repository.RecordRepository
}

// NewAdminRecordService 创建 AdminRecordService 实例
func NewAdminRecordService(recordRepo *repository.RecordRepository) *AdminRecordService {
	return &AdminRecordService{recordRepo: recordRepo}
}

// AdminRecordListInput 记录列表查询参数
type AdminRecordListInput struct {
	// CoupleID 按伴侣关系过滤
	CoupleID string `json:"couple_id"`

	// Page 页码
	Page int `json:"page" binding:"omitempty,min=1" example:"1"`

	// PageSize 每页数量
	PageSize int `json:"page_size" binding:"omitempty,min=1,max=100" example:"20"`
}

// AdminRecordListResult 记录列表返回结构
type AdminRecordListResult struct {
	// List 记录列表
	List []model.Record `json:"list"`

	// Total 总记录数
	Total int64 `json:"total"`

	// Page 当前页码
	Page int `json:"page"`

	// PageSize 每页数量
	PageSize int `json:"page_size"`
}

// ListRecords 获取记录列表
func (s *AdminRecordService) ListRecords(input AdminRecordListInput) (*AdminRecordListResult, error) {
	if input.Page == 0 {
		input.Page = 1
	}
	if input.PageSize == 0 {
		input.PageSize = 20
	}

	filter := repository.ListFilter{
		Page:     input.Page,
		PageSize: input.PageSize,
	}

	records, total, err := s.recordRepo.List(input.CoupleID, filter)
	if err != nil {
		return nil, err
	}

	return &AdminRecordListResult{
		List:     records,
		Total:    total,
		Page:     input.Page,
		PageSize: input.PageSize,
	}, nil
}
