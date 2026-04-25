// Package service 业务逻辑层
package service

import (
	"errors"

	"github.com/linke15d/bondy-backend/internal/repository"
)

// StatsService 数据统计业务逻辑
type StatsService struct {
	statsRepo  *repository.StatsRepository
	coupleRepo *repository.CoupleRepository
}

// NewStatsService 创建 StatsService 实例
func NewStatsService(statsRepo *repository.StatsRepository, coupleRepo *repository.CoupleRepository) *StatsService {
	return &StatsService{
		statsRepo:  statsRepo,
		coupleRepo: coupleRepo,
	}
}

// YearInput 年度查询参数
type YearInput struct {
	// Year 要查询的年份，如 2024
	Year int `json:"year" binding:"required,min=2000,max=2100" example:"2024"`
}

// MonthInput 月度查询参数
type MonthInput struct {
	// Year 年份
	Year int `json:"year" binding:"required,min=2000,max=2100" example:"2024"`
	// Month 月份 1-12
	Month int `json:"month" binding:"required,min=1,max=12" example:"1"`
}

// YearlyStatsResult 年度统计返回结构
type YearlyStatsResult struct {
	// Year 查询的年份
	Year int `json:"year"`
	// Overall 年度整体数据
	Overall *repository.MonthlyCount `json:"overall"`
	// Monthly 每个月的数据，没有记录的月份不返回
	Monthly []repository.MonthlyCount `json:"monthly"`
	// Heatmap 全年热力图数据
	Heatmap []repository.HeatmapItem `json:"heatmap"`
	// TopTags 使用最多的标签 Top5
	TopTags []repository.TagStat `json:"top_tags"`
	// TopPositions 使用最多的姿势 Top5
	TopPositions []repository.PositionStat `json:"top_positions"`
	// TimeDistribution 一天中各时段分布
	TimeDistribution []repository.TimeSlot `json:"time_distribution"`
	// WeekdayDistribution 一周各星期分布
	WeekdayDistribution []repository.WeekdayStat `json:"weekday_distribution"`
}

// OverviewResult 总览统计返回结构
type OverviewResult struct {
	// Total 累计总统计
	Total *repository.TotalStats `json:"total"`
	// StreakDays 当前连续记录天数
	StreakDays int `json:"streak_days"`
	// ActiveYears 有记录的年份列表，用于年份选择器
	ActiveYears []int `json:"active_years"`
}

// GetOverview 获取总览统计
// 包含累计数据、连续天数、有记录的年份列表
func (s *StatsService) GetOverview(userID string) (*OverviewResult, error) {
	couple, err := s.coupleRepo.FindByUserID(userID)
	if err != nil {
		return nil, errors.New("暂无伴侣关系")
	}

	total, err := s.statsRepo.GetTotalStats(couple.ID)
	if err != nil {
		return nil, errors.New("获取统计数据失败")
	}

	streak, err := s.statsRepo.GetStreakDays(couple.ID)
	if err != nil {
		return nil, errors.New("获取连续天数失败")
	}

	years, err := s.statsRepo.GetActiveYears(couple.ID)
	if err != nil {
		return nil, errors.New("获取年份列表失败")
	}

	return &OverviewResult{
		Total:       total,
		StreakDays:  streak,
		ActiveYears: years,
	}, nil
}

// GetYearlyStats 获取年度统计
// 包含每月数据、热力图、标签/姿势排行、时段分布
func (s *StatsService) GetYearlyStats(userID string, input YearInput) (*YearlyStatsResult, error) {
	couple, err := s.coupleRepo.FindByUserID(userID)
	if err != nil {
		return nil, errors.New("暂无伴侣关系")
	}

	overall, err := s.statsRepo.GetYearlyStats(couple.ID, input.Year)
	if err != nil {
		return nil, errors.New("获取年度数据失败")
	}

	monthly, err := s.statsRepo.GetMonthlyStats(couple.ID, input.Year)
	if err != nil {
		return nil, errors.New("获取月度数据失败")
	}

	heatmap, err := s.statsRepo.GetHeatmap(couple.ID, input.Year)
	if err != nil {
		return nil, errors.New("获取热力图数据失败")
	}

	topTags, err := s.statsRepo.GetTopTags(couple.ID, 5)
	if err != nil {
		return nil, errors.New("获取标签排行失败")
	}

	topPositions, err := s.statsRepo.GetTopPositions(couple.ID, 5)
	if err != nil {
		return nil, errors.New("获取姿势排行失败")
	}

	timeDistribution, err := s.statsRepo.GetTimeDistribution(couple.ID)
	if err != nil {
		return nil, errors.New("获取时段分布失败")
	}

	weekdayDistribution, err := s.statsRepo.GetWeekdayDistribution(couple.ID)
	if err != nil {
		return nil, errors.New("获取星期分布失败")
	}

	return &YearlyStatsResult{
		Year:                input.Year,
		Overall:             overall,
		Monthly:             monthly,
		Heatmap:             heatmap,
		TopTags:             topTags,
		TopPositions:        topPositions,
		TimeDistribution:    timeDistribution,
		WeekdayDistribution: weekdayDistribution,
	}, nil
}

// GetMonthlyStats 获取月度统计
// 包含当月每天的记录（用于日历视图）
func (s *StatsService) GetMonthlyStats(userID string, input MonthInput) (map[string]interface{}, error) {
	couple, err := s.coupleRepo.FindByUserID(userID)
	if err != nil {
		return nil, errors.New("暂无伴侣关系")
	}

	// 当月整体统计
	monthly, err := s.statsRepo.GetMonthlyStats(couple.ID, input.Year)
	if err != nil {
		return nil, errors.New("获取月度数据失败")
	}

	// 找到对应月份的数据
	var currentMonth *repository.MonthlyCount
	for _, m := range monthly {
		if m.Month == input.Month {
			currentMonth = &m
			break
		}
	}

	// 当月每天的记录（日历视图用）
	records, err := s.statsRepo.GetRecordsByMonth(couple.ID, input.Year, input.Month)
	if err != nil {
		return nil, errors.New("获取日历数据失败")
	}

	return map[string]interface{}{
		"year":    input.Year,
		"month":   input.Month,
		"summary": currentMonth,
		"records": records,
	}, nil
}
