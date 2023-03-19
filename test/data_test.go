/*
 * @Version: 1.0
 * @Date: 2023-03-18 08:51:48
 * @LastEditTime: 2023-03-19 17:56:00
 */
package test

import (
	"mytodo/data"
	"testing"
)

func TestGenHistoryTodoData(t *testing.T) {
	// 根据json文件生成历史待办
	data.GenHistoryTodoData()
}

func TestGenHistoryTaskData(t *testing.T) {
	// 根据历史待办生成历史任务，并随机选择一些任务更改状态为已完成
	data.GenHistoryTaskData()
}
