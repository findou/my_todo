/*
 * @Version: 1.0
 * @Date: 2023-03-18 19:08:06
 * @LastEditTime: 2023-03-18 19:11:44
 */
package test

import (
	"mytodo/model"
	"mytodo/utils"
	"testing"

	"time"
)

func TestFindAllTodoByTimeRange(t *testing.T) {
	now := time.Now()
	// 开始时间是一周的第一天，结束时间是一周的最后一天
	et := utils.GetLastDayOfWeek(now)
	st := utils.GetFirstDayOfWeek(now)

	// 更具开始时间结束时间查询待办
	todo := &model.Todo{}
	todo.FindAllTodoByTimeRange(st, et)
}
