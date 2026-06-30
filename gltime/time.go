package gltime

import "time"

// Format 按指定 layout 格式化时间。
func Format(t time.Time, layout string) string {
	return t.Format(layout)
}

// Parse 按指定 layout 解析时间字符串。
func Parse(layout string, value string) (time.Time, error) {
	return time.Parse(layout, value)
}

// UnixMilli 返回时间对应的毫秒级 Unix 时间戳。
func UnixMilli(t time.Time) int64 {
	return t.UnixMilli()
}

// FromUnixMilli 将毫秒级 Unix 时间戳转换为时间。
func FromUnixMilli(ms int64) time.Time {
	return time.UnixMilli(ms)
}

// DayStart 返回给定时间所在日期的开始时间。
func DayStart(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

// DayEnd 返回给定时间所在日期的结束时间。
func DayEnd(t time.Time) time.Time {
	return DayStart(t).AddDate(0, 0, 1).Add(-time.Nanosecond)
}
