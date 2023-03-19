/*
 * @Version: 1.0
 * @Date: 2023-03-17 14:06:25
 * @LastEditTime: 2023-03-19 16:44:47
 */
package handletodo

import (
	"encoding/json"
	"fmt"
	"mytodo/model"
	"mytodo/utils"
	"time"
)

type WeekHandle struct {
	cycleType     int
	cycleTypeName string
}

func NewWeekHandle(cycleType int) *WeekHandle {
	if cycleType == 2 {
		cycleTypeName := "WEEK"
		return &WeekHandle{cycleType: cycleType, cycleTypeName: cycleTypeName}
	}
	return nil
}

func (h *WeekHandle) CheckCycleArrays(checkCycleArrays string, startTime, endTime time.Time) (string, error) {
	// [1,3,5] //数字 范围 1-7

	// 开始时间是一周的第一天，结束时间是一周的最后一天
	et := utils.GetLastDayOfWeek(endTime)
	st := utils.GetFirstDayOfWeek(startTime)

	if (!endTime.Equal(et)) || (!startTime.Equal(st)) {
		return "", fmt.Errorf("start time , end time error")
	}

	carr := make([]int, 0)
	newArr := make([]int, 0)
	err := json.Unmarshal([]byte(checkCycleArrays), &carr)

	if err != nil {
		return "", err
	}
	for i := range carr {
		n := carr[i]
		if n < 1 || n > 7 {
			continue
		} else {
			newArr = append(newArr, n)
		}
	}
	b, _ := json.Marshal(newArr)
	cycleArrays := string(b)

	return cycleArrays, nil
}

// 检查待办 将待办转化为任务
func (h *WeekHandle) CheckToDoToTask(todo *model.Todo, now time.Time) error {
	// [1,3,5] //数字 范围 1-7
	cycleArrays, err := h.getCycleArrays(todo)
	if err != nil {
		return err
	}

	num := utils.GetWeekDayNum(now)

	for i := range cycleArrays {
		n := cycleArrays[i]
		if n == num {
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

func (h *WeekHandle) getCycleArrays(todo *model.Todo) ([]int, error) {
	cycleArrays := make([]int, 0)
	ca := []byte(todo.CycleArrays)
	err := json.Unmarshal(ca, &cycleArrays)

	if err != nil {
		return nil, err
	}
	return cycleArrays, nil
}

// 检查 cycleArrays 字段包含的时间是否存在 在 传入的开始结束时间中
func (h *WeekHandle) CheckContainTodo(todo *model.Todo, startTime, endTime time.Time) (map[string]struct{}, error) {
	taskCreatTimes := make(map[string]struct{}, 0)
	// [1,3,5] //数字 范围 1-7
	cycleArrays, err := h.getCycleArrays(todo)
	if err != nil {
		return nil, err
	}

	st := utils.GetFirstDayOfWeek(startTime)
	et := utils.GetLastDayOfWeek(endTime)

	//w := startTime.Weekday()
	sTime := time.Date(st.Year(), st.Month(), st.Day(), 1, 0, 0, 0, time.Now().Location())
	now := sTime
	for !et.Before(now) {
		for i := range cycleArrays {
			n := cycleArrays[i]
			nowTime := now.AddDate(0, 0, n-1)
			cTimeStr := utils.ParseTimeToStr(nowTime)
			if !endTime.Before(nowTime) && !nowTime.Before(startTime) {
				// taskCreatTimes = append(taskCreatTimes, cTimeStr)
				taskCreatTimes[cTimeStr] = struct{}{}
			}
		}
		t := now.AddDate(0, 0, 7)
		now = t
	}
	return taskCreatTimes, nil
}

// 检查待办 将待办转化为历史任务，用于新生成的todo
func (h *WeekHandle) CheckToDoToRandomHistoryTask(todo *model.Todo, untilTime time.Time) ([]*model.Task, error) {
	// [1,3,5] //数字 范围 1-7
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

	et := todo.EndTime
	endTime, err := utils.ParseStrTime(et)
	if err != nil {
		return nil, err
	}

	sTime := time.Date(startTime.Year(), startTime.Month(), startTime.Day(), 1, 0, 0, 0, time.Now().Location())
	now := sTime
	for !endTime.Before(now) && !untilTime.Before(now) {
		for i := range cycleArrays {
			n := cycleArrays[i]
			nowTime := now.AddDate(0, 0, n-1)

			if !untilTime.Before(nowTime) {
				task := &model.Task{
					Name:       todo.Name,
					Id:         utils.GenUuid(),
					Priority:   todo.Priority,
					TodoId:     todo.Id,
					CreateTime: utils.ParseTimeToStr(nowTime),
					Status:     false,
				}
				tasks = append(tasks, task)
			}
		}
		t := now.AddDate(0, 0, 7)
		now = t
	}

	return tasks, nil
}

// 获取开始时间结束时间
func (h *WeekHandle) GetTimeArrange(nowTime time.Time) (startTime, endTime time.Time, err error) {
	// 开始时间是一周的第一天，结束时间是一周的最后一天
	et := utils.GetLastDayOfWeek(nowTime)
	st := utils.GetFirstDayOfWeek(nowTime)

	return st, et, nil
}
