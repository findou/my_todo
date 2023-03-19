/*
 * @Version: 1.0
 * @Date: 2023-03-17 15:56:17
 * @LastEditTime: 2023-03-19 10:15:16
 */
package test

import (
	"fmt"
	"testing"
	"time"
)

func TestTime(t *testing.T) {
	now := time.Now()
	/* w := now.Weekday()
	fmt.Printf("now.Weekday(): %v\n", w)
	d := now.Day()
	fmt.Printf("now.Day(): %v\n", d)
	m := now.Month()
	fmt.Printf("now.Month(): %v\n", m)
	y := now.Year()
	fmt.Printf("now.Year(): %v\n", y) */

	// syear, smonth, sday := now.Date()
	// fmt.Printf("syear:%d,smonth:%d,sday:%d\n", syear, smonth, sday)

	// 获取月份第一天
	// fyear, fmonth, fday := GetFirstDayOfMonth(now).Date()
	// fmt.Printf("fyear:%d,fmonth:%d,fday:%d\n", fyear, fmonth, fday)

	// 获取月份最后一天
	// lyear, lmonth, lday := GetLastDayOfMonth(now).Date()
	// fmt.Printf("lyear:%d,lmonth:%d,lday:%d\n", lyear, lmonth, lday)

	// 获取本周的第一天
	// lyear, lmonth, lday := GetFirstDayOfWeek(now).Date()
	// fmt.Printf("lyear:%d,lmonth:%d,lday:%d\n", lyear, lmonth, lday)

	// 获取本周的最后一天
	lyear, lmonth, lday := GetLastDayOfWeek(now).Date()
	fmt.Printf("lyear:%d,lmonth:%d,lday:%d\n", lyear, lmonth, lday)

	// 获取本年的第一天
	// lyear, lmonth, lday := GetFirstDayOfYear(now).Date()
	// fmt.Printf("lyear:%d,lmonth:%d,lday:%d\n", lyear, lmonth, lday)

	// 获取本年的最后一天，最后一秒
	// lyear, lmonth, lday := GetLastDayOfYear(now).Date()
	// fmt.Printf("lyear:%d,lmonth:%d,lday:%d\n", lyear, lmonth, lday)

}

func TestParseTime(t *testing.T) {
	str := "2023-03-20 00:00:00"
	t2 := ParseStrTime(str)
	year, m, day := t2.Date()
	fmt.Printf("year:%d,m:%d,day:%d\n", year, m, day)
}

func TestNumTime(t *testing.T) {
	now := time.Now()

	fmt.Printf("GetWeekNum(now): %v\n", GetWeekDayNum(now))
}

func GetWeekDayNum(d time.Time) int {
	wd := d.Weekday()
	if wd == 0 {
		return 7
	}
	return int(wd)
}

func ParseStrTime(strTime string) time.Time {
	t, _ := time.ParseInLocation("2006-01-02 15:04:05", strTime, time.Local)
	return t
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
