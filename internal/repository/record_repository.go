// Package repository 数据访问层
package repository

import (
	"github.com/linke15d/bondy-backend/internal/model"
	"gorm.io/gorm"
)

// RecordRepository 亲密记录数据访问对象
type RecordRepository struct {
	db *gorm.DB
}

// NewRecordRepository 创建 RecordRepository 实例
func NewRecordRepository(db *gorm.DB) *RecordRepository {
	return &RecordRepository{db: db}
}

// ListFilter 记录列表查询过滤条件
type ListFilter struct {
	// Page 页码，从 1 开始
	Page int
	// PageSize 每页数量，最大50
	PageSize int
	// Year 按年份过滤，0表示不过滤
	Year int
	// Month 按月份过滤，0表示不过滤
	Month int
}

// Create 创建亲密记录
func (r *RecordRepository) Create(record *model.Record) error {
	return r.db.Create(record).Error
}

// FindByID 通过 ID 查找记录
// 同时预加载关联的标签和姿势
func (r *RecordRepository) FindByID(id string, coupleID string) (*model.Record, error) {
	var record model.Record
	err := r.db.
		Preload("Tags").
		Preload("Positions").
		Where("id = ? AND couple_id = ? AND is_deleted = false", id, coupleID).
		First(&record).Error
	if err != nil {
		return nil, err
	}
	return &record, nil
}

// List 获取伴侣的记录列表，支持分页和按年月过滤
func (r *RecordRepository) List(coupleID string, filter ListFilter) ([]model.Record, int64, error) {
	var records []model.Record
	var total int64

	query := r.db.Model(&model.Record{}).
		Where("couple_id = ? AND is_deleted = false", coupleID)

	// 按年份过滤
	if filter.Year > 0 {
		query = query.Where("EXTRACT(YEAR FROM happened_at) = ?", filter.Year)
	}

	// 按月份过滤
	if filter.Month > 0 {
		query = query.Where("EXTRACT(MONTH FROM happened_at) = ?", filter.Month)
	}

	// 查总数
	query.Count(&total)

	// 分页查询，按发生时间倒序
	offset := (filter.Page - 1) * filter.PageSize
	err := query.
		Preload("Tags").
		Preload("Positions").
		Order("happened_at DESC").
		Offset(offset).
		Limit(filter.PageSize).
		Find(&records).Error

	return records, total, err
}

// Update 更新记录
func (r *RecordRepository) Update(record *model.Record) error {
	return r.db.Save(record).Error
}

// SoftDelete 软删除记录
func (r *RecordRepository) SoftDelete(id string, coupleID string) error {
	return r.db.Model(&model.Record{}).
		Where("id = ? AND couple_id = ?", id, coupleID).
		Update("is_deleted", true).Error
}

// FindTagsByType 获取标签列表
// 返回系统预设标签 + 当前伴侣自定义标签
func (r *RecordRepository) FindTagsByType(tagType string, coupleID string) ([]model.Tag, error) {
	var tags []model.Tag
	err := r.db.
		Where("type = ? AND (is_system = true OR couple_id = ?)", tagType, coupleID).
		Order("is_system DESC, name ASC").
		Find(&tags).Error
	return tags, err
}

// PositionsByCategory 按分类分组的姿势列表
type PositionsByCategory struct {
	// Category 分类名
	Category string `json:"category"`
	// Positions 该分类下的姿势列表
	Positions []model.Position `json:"positions"`
	// CategoryName 分类中文名，根据语言返回对应翻译
	CategoryName string `json:"category_name"`
}

// FindPositions 获取姿势列表
// FindPositionsByCategory 获取姿势列表，按分类分组，填充对应语言翻译
func (r *RecordRepository) FindPositions(coupleID string, lang string) ([]PositionsByCategory, error) {
	categories := []string{"CLASSIC", "ADVENTURE", "INTIMATE", "FUN"}
	var result []PositionsByCategory

	for _, cat := range categories {
		var positions []model.Position
		err := r.db.
			Where("(is_system = true OR couple_id = ?) AND category = ?", coupleID, cat).
			Order("is_system DESC").
			Find(&positions).Error
		if err != nil {
			return nil, err
		}

		// 填充翻译
		for i := range positions {
			// 填充 name 翻译
			var nameTrans model.Translation
			if err := r.db.Where(
				"module = 'position' AND ref_id = ? AND field = 'name' AND language_code = ?",
				positions[i].ID, lang,
			).First(&nameTrans).Error; err == nil {
				positions[i].Name = nameTrans.Content
			} else {
				// 没有翻译则用默认名称
				positions[i].Name = positions[i].DefaultName
			}
		}

		// 获取分类翻译
		categoryName := cat
		var catTrans model.Translation
		if err := r.db.Where(
			"module = 'position' AND ref_id = ? AND field = 'category_name' AND language_code = ?",
			cat, lang,
		).First(&catTrans).Error; err == nil {
			categoryName = catTrans.Content
		}

		result = append(result, PositionsByCategory{
			Category:     cat,
			CategoryName: categoryName,
			Positions:    positions,
		})
	}

	return result, nil
}

// CreateTag 创建自定义标签
func (r *RecordRepository) CreateTag(tag *model.Tag) error {
	return r.db.Create(tag).Error
}

// CreatePosition 创建自定义姿势
func (r *RecordRepository) CreatePosition(position *model.Position) error {
	return r.db.Create(position).Error
}

// FindCategories 获取所有启用的分类，填充对应语言翻译
func (r *RecordRepository) FindCategories(lang string) ([]model.PositionCategory, error) {
	var categories []model.PositionCategory
	err := r.db.
		Where("is_active = true").
		Order("sort_order ASC").
		Find(&categories).Error
	if err != nil {
		return nil, err
	}

	// 填充翻译
	for i := range categories {
		var trans model.Translation
		if err := r.db.Where(
			"module = 'position_category' AND ref_id = ? AND field = 'name' AND language_code = ?",
			categories[i].ID, lang,
		).First(&trans).Error; err == nil {
			categories[i].DefaultName = trans.Content
		}
	}

	return categories, nil
}
