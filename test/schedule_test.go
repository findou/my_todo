/*
 * @Version: 1.0
 * @Date: 2023-03-17 23:58:52
 * @LastEditTime: 2023-03-19 13:11:19
 */
package test

import (
	"mytodo/schedule"
	"mytodo/utils"
	"testing"
)

func TestExcuteTodoTask(t *testing.T) {
	// 2023-03-20 01:00:00
	// 模拟这个时间触发执行任务
	str := "2023-03-21 01:00:00"
	ti, _ := utils.ParseStrTime(str)
	schedule.ExcuteTodoTask(ti)
}

func TestHandleScheduleTask(t *testing.T) {
	schedule.HandleScheduleTask()
}
