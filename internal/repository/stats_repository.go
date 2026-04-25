// Package repository 数据访问层
package repository

import (
	"time"

	"github.com/linke15d/bondy-backend/internal/model"
	"gorm.io/gorm"
)

// StatsRepository 统计数据访问对象
type StatsRepository struct {
	db *gorm.DB
}

// NewStatsRepository 创建 StatsRepository 实例
func NewStatsRepository(db *gorm.DB) *StatsRepository {
	return &StatsRepository{db: db}
}

// MonthlyCount 某月记录数和平均评分
type MonthlyCount struct {
	// Year 年份
	Year int `json:"year"`
	// Month 月份
	Month int `json:"month"`
	// Count 当月记录总数
	Count int `json:"count"`
	// AvgMood 当月平均心情评分，保留一位小数
	AvgMood float64 `json:"avg_mood"`
	// AvgSatisfaction 当月平均满意度评分，保留一位小数
	AvgSatisfaction float64 `json:"avg_satisfaction"`
	// TotalDurationMins 当月累计时长（分钟）
	TotalDurationMins int `json:"total_duration_mins"`
}

// HeatmapItem 热力图单日数据
type HeatmapItem struct {
	// Date 日期，格式 2006-01-02
	Date string `json:"date"`
	// Count 当天记录次数
	Count int `json:"count"`
}

// GetMonthlyStats 获取某年每个月的统计数据
func (r *StatsRepository) GetMonthlyStats(coupleID string, year int) ([]MonthlyCount, error) {
	var results []MonthlyCount
	err := r.db.Raw(`
		SELECT
			EXTRACT(YEAR FROM happened_at)::int AS year,
			EXTRACT(MONTH FROM happened_at)::int AS month,
			COUNT(*) AS count,
			COALESCE(ROUND(AVG(mood)::numeric, 1), 0) AS avg_mood,
			COALESCE(ROUND(AVG(satisfaction)::numeric, 1), 0) AS avg_satisfaction,
			COALESCE(SUM(duration_mins), 0) AS total_duration_mins
		FROM records
		WHERE couple_id = ?
			AND is_deleted = false
			AND EXTRACT(YEAR FROM happened_at) = ?
		GROUP BY year, month
		ORDER BY month ASC
	`, coupleID, year).Scan(&results).Error
	return results, err
}

// GetYearlyStats 获取某年整体统计数据
func (r *StatsRepository) GetYearlyStats(coupleID string, year int) (*MonthlyCount, error) {
	var result MonthlyCount
	err := r.db.Raw(`
		SELECT
			? AS year,
			0 AS month,
			COUNT(*) AS count,
			COALESCE(ROUND(AVG(mood)::numeric, 1), 0) AS avg_mood,
			COALESCE(ROUND(AVG(satisfaction)::numeric, 1), 0) AS avg_satisfaction,
			COALESCE(SUM(duration_mins), 0) AS total_duration_mins
		FROM records
		WHERE couple_id = ?
			AND is_deleted = false
			AND EXTRACT(YEAR FROM happened_at) = ?
	`, year, coupleID, year).Scan(&result).Error
	return &result, err
}

// GetHeatmap 获取某年的热力图数据
// 返回每一天的记录次数，没有记录的日期不返回
func (r *StatsRepository) GetHeatmap(coupleID string, year int) ([]HeatmapItem, error) {
	var results []HeatmapItem
	err := r.db.Raw(`
		SELECT
			TO_CHAR(happened_at, 'YYYY-MM-DD') AS date,
			COUNT(*) AS count
		FROM records
		WHERE couple_id = ?
			AND is_deleted = false
			AND EXTRACT(YEAR FROM happened_at) = ?
		GROUP BY date
		ORDER BY date ASC
	`, coupleID, year).Scan(&results).Error
	return results, err
}

// GetTotalStats 获取伴侣关系建立至今的总体统计
func (r *StatsRepository) GetTotalStats(coupleID string) (*TotalStats, error) {
	var result TotalStats
	err := r.db.Raw(`
		SELECT
			COUNT(*) AS total_count,
			COALESCE(SUM(duration_mins), 0) AS total_duration_mins,
			COALESCE(ROUND(AVG(mood)::numeric, 1), 0) AS avg_mood,
			COALESCE(ROUND(AVG(satisfaction)::numeric, 1), 0) AS avg_satisfaction,
			MIN(happened_at) AS first_time,
			MAX(happened_at) AS last_time
		FROM records
		WHERE couple_id = ?
			AND is_deleted = false
	`, coupleID).Scan(&result).Error
	return &result, err
}

// GetStreakDays 获取当前连续有记录的天数（streak）
func (r *StatsRepository) GetStreakDays(coupleID string) (int, error) {
	// 获取所有有记录的日期，倒序
	var dates []string
	err := r.db.Raw(`
		SELECT DISTINCT TO_CHAR(happened_at, 'YYYY-MM-DD') AS date
		FROM records
		WHERE couple_id = ? AND is_deleted = false
		ORDER BY date DESC
	`, coupleID).Scan(&dates).Error
	if err != nil || len(dates) == 0 {
		return 0, err
	}

	// 从最近一天往前数连续天数
	streak := 0
	today := time.Now().Format("2006-01-02")
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")

	// 最近一条记录必须是今天或昨天，否则 streak 为 0
	if dates[0] != today && dates[0] != yesterday {
		return 0, nil
	}

	current, _ := time.Parse("2006-01-02", dates[0])
	for _, d := range dates {
		t, err := time.Parse("2006-01-02", d)
		if err != nil {
			break
		}
		diff := current.Sub(t)
		if diff.Hours() > 24 {
			break
		}
		streak++
		current = t
	}

	return streak, nil
}

// TotalStats 总体统计数据
type TotalStats struct {
	// TotalCount 累计记录总次数
	TotalCount int `json:"total_count"`
	// TotalDurationMins 累计总时长（分钟）
	TotalDurationMins int `json:"total_duration_mins"`
	// AvgMood 整体平均心情评分
	AvgMood float64 `json:"avg_mood"`
	// AvgSatisfaction 整体平均满意度评分
	AvgSatisfaction float64 `json:"avg_satisfaction"`
	// FirstTime 第一次记录时间
	FirstTime *time.Time `json:"first_time"`
	// LastTime 最近一次记录时间
	LastTime *time.Time `json:"last_time"`
}

// GetActiveYears 获取有记录的年份列表
// 用于前端年份选择器
func (r *StatsRepository) GetActiveYears(coupleID string) ([]int, error) {
	var years []int
	err := r.db.Raw(`
		SELECT DISTINCT EXTRACT(YEAR FROM happened_at)::int AS year
		FROM records
		WHERE couple_id = ? AND is_deleted = false
		ORDER BY year DESC
	`, coupleID).Scan(&years).Error
	return years, err
}

// GetTopTags 获取使用最多的标签 Top N
func (r *StatsRepository) GetTopTags(coupleID string, limit int) ([]TagStat, error) {
	var results []TagStat
	err := r.db.Raw(`
		SELECT
			t.id,
			t.name,
			t.type,
			COUNT(*) AS count
		FROM record_tags rt
		JOIN tags t ON t.id = rt.tag_id
		JOIN records r ON r.id = rt.record_id
		WHERE r.couple_id = ? AND r.is_deleted = false
		GROUP BY t.id, t.name, t.type
		ORDER BY count DESC
		LIMIT ?
	`, coupleID, limit).Scan(&results).Error
	return results, err
}

// TagStat 标签使用统计
type TagStat struct {
	// ID 标签 ID
	ID string `json:"id"`
	// Name 标签名称
	Name string `json:"name"`
	// Type 标签类型
	Type string `json:"type"`
	// Count 使用次数
	Count int `json:"count"`
}

// GetTopPositions 获取使用最多的姿势 Top N
func (r *StatsRepository) GetTopPositions(coupleID string, limit int) ([]PositionStat, error) {
	var results []PositionStat
	err := r.db.Raw(`
		SELECT
			p.id,
			p.name,
			COUNT(*) AS count
		FROM record_positions rp
		JOIN positions p ON p.id = rp.position_id
		JOIN records r ON r.id = rp.record_id
		WHERE r.couple_id = ? AND r.is_deleted = false
		GROUP BY p.id, p.name
		ORDER BY count DESC
		LIMIT ?
	`, coupleID, limit).Scan(&results).Error
	return results, err
}

// PositionStat 姿势使用统计
type PositionStat struct {
	// ID 姿势 ID
	ID string `json:"id"`
	// Name 姿势名称
	Name string `json:"name"`
	// Count 使用次数
	Count int `json:"count"`
}

// GetTimeDistribution 获取一天中各时段的记录分布
// 将24小时分为6个时段统计
func (r *StatsRepository) GetTimeDistribution(coupleID string) ([]TimeSlot, error) {
	var results []TimeSlot
	err := r.db.Raw(`
		SELECT
			CASE
				WHEN EXTRACT(HOUR FROM happened_at) BETWEEN 6 AND 11 THEN '早晨 06-12'
				WHEN EXTRACT(HOUR FROM happened_at) BETWEEN 12 AND 13 THEN '中午 12-14'
				WHEN EXTRACT(HOUR FROM happened_at) BETWEEN 14 AND 17 THEN '下午 14-18'
				WHEN EXTRACT(HOUR FROM happened_at) BETWEEN 18 AND 21 THEN '傍晚 18-22'
				WHEN EXTRACT(HOUR FROM happened_at) BETWEEN 22 AND 23 THEN '深夜 22-24'
				ELSE '凌晨 00-06'
			END AS slot,
			COUNT(*) AS count
		FROM records
		WHERE couple_id = ? AND is_deleted = false
		GROUP BY slot
		ORDER BY count DESC
	`, coupleID).Scan(&results).Error
	return results, err
}

// TimeSlot 时段统计
type TimeSlot struct {
	// Slot 时段描述
	Slot string `json:"slot"`
	// Count 该时段记录次数
	Count int `json:"count"`
}

// GetWeekdayDistribution 获取一周各星期的记录分布
func (r *StatsRepository) GetWeekdayDistribution(coupleID string) ([]WeekdayStat, error) {
	var results []WeekdayStat
	err := r.db.Raw(`
		SELECT
			TO_CHAR(happened_at, 'Day') AS weekday,
			EXTRACT(DOW FROM happened_at)::int AS dow,
			COUNT(*) AS count
		FROM records
		WHERE couple_id = ? AND is_deleted = false
		GROUP BY weekday, dow
		ORDER BY dow ASC
	`, coupleID).Scan(&results).Error
	return results, err
}

// WeekdayStat 星期统计
type WeekdayStat struct {
	// Weekday 星期名称
	Weekday string `json:"weekday"`
	// Dow 星期数字 0=周日 1=周一 ... 6=周六
	Dow int `json:"dow"`
	// Count 该星期记录次数
	Count int `json:"count"`
}

// GetRecordsByMonth 获取某月的所有记录（用于日历视图）
func (r *StatsRepository) GetRecordsByMonth(coupleID string, year, month int) ([]model.Record, error) {
	var records []model.Record
	err := r.db.
		Where(`couple_id = ?
			AND is_deleted = false
			AND EXTRACT(YEAR FROM happened_at) = ?
			AND EXTRACT(MONTH FROM happened_at) = ?`,
			coupleID, year, month).
		Order("happened_at ASC").
		Find(&records).Error
	return records, err
}
