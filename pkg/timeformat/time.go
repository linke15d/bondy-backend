// Package timeformat 自定义时间格式
package timeformat

import (
	"database/sql/driver"
	"fmt"
	"time"
)

// LocalTime 自定义时间类型，JSON 序列化格式为 yyyy-MM-dd HH:mm:ss
type LocalTime time.Time

const timeLayout = "2006-01-02 15:04:05"

// MarshalJSON 序列化为 JSON 时格式化为 yyyy-MM-dd HH:mm:ss
func (t LocalTime) MarshalJSON() ([]byte, error) {
	formatted := fmt.Sprintf("\"%s\"", time.Time(t).Format(timeLayout))
	return []byte(formatted), nil
}

// UnmarshalJSON 从 JSON 反序列化
func (t *LocalTime) UnmarshalJSON(data []byte) error {
	str := string(data)
	// 去掉引号
	str = str[1 : len(str)-1]
	parsed, err := time.ParseInLocation(timeLayout, str, time.Local)
	if err != nil {
		// 兼容 RFC3339 格式
		parsed, err = time.Parse(time.RFC3339, str)
		if err != nil {
			return err
		}
	}
	*t = LocalTime(parsed)
	return nil
}

// Value 写入数据库时转换
func (t LocalTime) Value() (driver.Value, error) {
	return time.Time(t), nil
}

// Scan 从数据库读取时转换
func (t *LocalTime) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	v, ok := value.(time.Time)
	if !ok {
		return fmt.Errorf("无法将 %T 转换为 LocalTime", value)
	}
	*t = LocalTime(v)
	return nil
}

// String 返回格式化字符串
func (t LocalTime) String() string {
	return time.Time(t).Format(timeLayout)
}
