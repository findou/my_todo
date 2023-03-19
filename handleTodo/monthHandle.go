/*
 * @Version: 1.0
 * @Date: 2023-03-17 14:06:25
 * @LastEditTime: 2023-03-19 18:07:10
 */
package handletodo

import (
	"encoding/json"
	"fmt"
	"mytodo/model"
	"mytodo/utils"
	"time"
)

type MonthHandle struct {
	cycleType     int
	cycleTypeName string
}

func NewMonthHandle(cycleType int) *MonthHandle {
	if cycleType == 3 {
		cycleTypeName := "MONTH"
		return &MonthHandle{cycleType: cycleType, cycleTypeName: cycleTypeName}
	}
	return nil
}

func (h *MonthHandle) CheckCycleArrays(checkCycleArrays string, startTime, endTime time.Time) (string, error) {
	// [1,9,21]

	// 结束时间是一个月的最后一天,开始时间是一个月的第一天
	et := utils.GetLastDayOfMonth(endTime)
	st := utils.GetFirstDayOfMonth(startTime)

	if (!endTime.Equal(et)) || (!startTime.Equal(st)) {
		return "", fmt.Errorf("start time , end time error")
	}

	_, _, enlastday := et.Date()
	_, _, stlastday := st.Date()
	max := enlastday
	if stlastday > enlastday {
		max = stlastday
	}

	carr := make([]int, 0)
	newArr := make([]int, 0)
	err := json.Unmarshal([]byte(checkCycleArrays), &carr)

	if err != nil {
		return "", err
	}
	for i := range carr {
		n := carr[i]
		if n < 1 || n > max {
			continue
		} else {
			newArr = append(newArr, n)
		}
	}
	b, _ := json.Marshal(newArr)
	cycleArrays := (string(b))

	return cycleArrays, nil
}

// 检查待办 将待办转化为任务
func (h *MonthHandle) CheckToDoToTask(todo *model.Todo, now time.Time) error {
	// // [1,9,21] //数字 范围 1-31
	cycleArrays, err := h.getCycleArrays(todo)
	if err != nil {
		return err
	}

	_, _, day := now.Date()
	_, _, lastday := utils.GetLastDayOfMonth(now).Date()

	for i := range cycleArrays {
		n := cycleArrays[i]
		if n < lastday && day == n {
			task := &model.Task{
				Name:       todo.Name,
				Id:         utils.GenUuid(),
				Priority:   todo.Priority,
				TodoId:     todo.Id,
				CreateTime: utils.ParseTimeToStr(now),
				Status:     false,
			}
			err := task.Save()
			return err
		}
	}

	return nil
}

func (h *MonthHandle) getCycleArrays(todo *model.Todo) ([]int, error) {
	cycleArrays := make([]int, 0)
	ca := []byte(todo.CycleArrays)
	err := json.Unmarshal(ca, &cycleArrays)

	if err != nil {
		return nil, err
	}
	return cycleArrays, nil
}

// 检查 cycleArrays 字段包含的时间是否存在 在 传入的开始结束时间中
func (h *MonthHandle) CheckContainTodo(todo *model.Todo, startTime, endTime time.Time) (map[string]struct{}, error) {
	taskCreatTimes := make(map[string]struct{}, 0)
	// [1,9,21] //数字 范围 1-31
	cycleArrays, err := h.getCycleArrays(todo)
	if err != nil {
		return nil, err
	}

	years, months, _ := startTime.Date()
	yeare, monthe, _ := endTime.Date()
	for yearNum := years; yearNum <= yeare; yearNum++ {
		for monthNum := months; monthNum <= monthe; monthNum++ {
			for i := range cycleArrays {
				day := cycleArrays[i]
				now := time.Date(yearNum, time.Month(monthNum), day, 1, 0, 0, 0, time.Now().Location())
				cTimeStr := utils.ParseTimeToStr(now)
				if !endTime.Before(now) && !now.Before(startTime) {
					// taskCreatTimes = append(taskCreatTimes, cTimeStr)
					taskCreatTimes[cTimeStr] = struct{}{}
				}
			}
		}
	}
	return taskCreatTimes, nil
}

// 检查待办 将待办转化为历史任务，用于新生成的todo
func (h *MonthHandle) CheckToDoToRandomHistoryTask(todo *model.Todo, untilTime time.Time) ([]*model.Task, error) {
	// [1,9,21] //数字 范围 1-31
	tasks := make([]*model.Task, 0)

	cycleArrays, err := h.getCycleArrays(todo)
	if err != nil {
		return nil, err
	}

	st := todo.StartTime
	startTime, err := utils.ParseStrTime(st)
	if err != nil {
		return nil, err
	}

	years, months, _ := startTime.Date()

	et := todo.EndTime
	endTime, err := utils.ParseStrTime(et)
	if err != nil {
		return nil, err
	}

	yeare, monthe, _ := endTime.Date()

	for yearNum := years; yearNum <= yeare; yearNum++ {
		for monthNum := months; monthNum <= monthe; monthNum++ {
			for i := range cycleArrays {
				day := cycleArrays[i]
				now := time.Date(yearNum, time.Month(monthNum), day, 1, 0, 0, 0, time.Now().Location())
				if !untilTime.Before(now) {
					task := &model.Task{
						Name:       todo.Name,
						Id:         utils.GenUuid(),
						Priority:   todo.Priority,
						TodoId:     todo.Id,
						CreateTime: utils.ParseTimeToStr(now),
						Status:     false,
					}
					tasks = append(tasks, task)
				}
			}

		}
	}

	return tasks, nil
}

// 获取开始时间结束时间
func (h *MonthHandle) GetTimeArrange(nowTime time.Time) (startTime, endTime time.Time, err error) {
	// 结束时间是一个月的最后一天,开始时间是一个月的第一天
	et := utils.GetLastDayOfMonth(nowTime)
	st := utils.GetFirstDayOfMonth(nowTime)

	return st, et, nil
}
