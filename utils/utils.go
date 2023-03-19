/*
 * @Version: 1.0
 * @Date: 2023-03-17 14:50:01
 * @LastEditTime: 2023-03-19 13:04:13
 */
package utils

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

func GenUuid() string {
	uuid := uuid.NewV4()
	return uuid.String()
}

// 字符串时间转化为时间
func ParseStrTime(strTime string) (time.Time, error) {
	t, err := time.ParseInLocation("2006-01-02 15:04:05", strTime, time.Local)
	return t, err
}

// 时间转化为字符串时间
func ParseTimeToStr(time time.Time) string {
	s := time.Format("2006-01-02 15:04:05")
	return s
}

// 时间转化为字符串时间,没有时间
func ParseStrShortTime(strTime string) (time.Time, error) {
	t, err := time.ParseInLocation("2006-01-02", strTime, time.Local)
	return t, err
}

// 时间转化为字符串时间,没有时间
func ParseTimeToShortStr(time time.Time) string {
	s := time.Format("2006-01-02")
	return s
}

// 获取本年的第一天
func GetFirstDayOfYear(d time.Time) time.Time {
	t := time.Date(d.Year(), time.January, 1, 0, 0, 0, 0, d.Location())
	return t
}

// 获取本年的最后一天，最后一秒
func GetLastDayOfYear(d time.Time) time.Time {
	t := time.Date(d.Year()+1, time.January, 1, 23, 59, 59, 0, d.Location())

	rt := t.AddDate(0, 0, -1)
	return rt
}

// 获取本周的第一天
func GetFirstDayOfWeek(d time.Time) time.Time {
	wd := d.Weekday()
	offSet := int(time.Sunday - wd + 1)
	if offSet > 0 {
		offSet = -6
	}

	t := GetZeroTime(d)
	rt := t.AddDate(0, 0, offSet)
	return rt
}

// 获取本周的最后一天，最后一秒
func GetLastDayOfWeek(d time.Time) time.Time {
	wd := d.Weekday()
	offSet := int(time.Sunday - wd)
	if offSet < 0 {
		offSet = offSet + 7
	}
	t := d.AddDate(0, 0, offSet)
	rt := GetLastTime(t)
	return rt
}

// 获取月份第一天
func GetFirstDayOfMonth(d time.Time) time.Time {
	day := d.Day()
	num := -day + 1
	t := d.AddDate(0, 0, num)
	return GetZeroTime(t)
}

// 获取月份最后一天,最后一秒
func GetLastDayOfMonth(d time.Time) time.Time {
	t1 := GetFirstDayOfMonth(d)
	t2 := t1.AddDate(0, 1, -1)
	rt := GetLastTime(t2)
	return rt
}

// 获取某天的0点时间
func GetZeroTime(d time.Time) time.Time {
	return time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, d.Location())
}

// 获取某天的最后一秒
func GetLastTime(d time.Time) time.Time {
	return time.Date(d.Year(), d.Month(), d.Day(), 23, 59, 59, 0, d.Location())
}

// 获取当前日期是周几
func GetWeekDayNum(d time.Time) int {
	wd := d.Weekday()
	if wd == 0 {
		return 7
	}
	return int(wd)
}
