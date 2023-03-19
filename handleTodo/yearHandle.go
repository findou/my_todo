/*
 * @Version: 1.0
 * @Date: 2023-03-17 14:06:25
 * @LastEditTime: 2023-03-19 16:45:25
 */
package handletodo

import (
	"encoding/json"
	"fmt"
	"mytodo/model"
	"mytodo/utils"
	"strconv"
	"strings"
	"time"
)

type YearHandle struct {
	cycleType     int
	cycleTypeName string
}

func NewYearHandle(cycleType int) *YearHandle {
	if cycleType == 4 {
		cycleTypeName := "YEAR"
		return &YearHandle{cycleType: cycleType, cycleTypeName: cycleTypeName}
	}
	return nil
}

func (h *YearHandle) CheckCycleArrays(checkCycleArrays string, startTime, endTime time.Time) (string, error) {
	// ["03-19,04-11"]

	// 结束时间是一年的最后一天,开始时间是年的第一天
	et := utils.GetLastDayOfYear(endTime)
	st := utils.GetFirstDayOfYear(startTime)

	if (!endTime.Equal(et)) || (!startTime.Equal(st)) {
		return "", fmt.Errorf("start time , end time error")
	}

	carr := make([]string, 0)
	newArr := make([]string, 0)
	err := json.Unmarshal([]byte(checkCycleArrays), &carr)

	if err != nil {
		return "", err
	}
	for i := range carr {
		str := carr[i]
		ss := strings.Split(str, "-")
		if len(ss) != 2 {
			continue
		} else {
			months := ss[0]
			month, _ := strconv.Atoi(months)
			if month < 0 || month > 12 {
				continue
			}
			days := ss[1]
			day, _ := strconv.Atoi(days)
			if day < 0 || day > 31 {
				continue
			}
			newArr = append(newArr, str)
		}

	}
	b, _ := json.Marshal(newArr)
	cycleArrays := (string(b))

	return cycleArrays, nil
}

// 检查待办 将待办转化为任务
func (h *YearHandle) CheckToDoToTask(todo *model.Todo, now time.Time) error {
	// ["03-19,04-11"]
	cycleArrays, err := h.getCycleArrays(todo)
	if err != nil {
		return err
	}

	_, nmonth, nday := now.Date()

	for i := range cycleArrays {
		str := cycleArrays[i]
		ss := strings.Split(str, "-")

		monthstr := ss[0]
		month, _ := strconv.Atoi(monthstr)

		daystr := ss[1]
		day, _ := strconv.Atoi(daystr)

		if month == int(nmonth) && nday == day {
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

func (h *YearHandle) getCycleArrays(todo *model.Todo) ([]string, error) {
	// ["03-19,04-11"]
	cycleArrays := make([]string, 0)
	ca := []byte(todo.CycleArrays)
	err := json.Unmarshal(ca, &cycleArrays)

	if err != nil {
		return nil, err
	}
	return cycleArrays, nil
}

// 检查 cycleArrays 字段包含的时间是否存在 在 传入的开始结束时间中
func (h *YearHandle) CheckContainTodo(todo *model.Todo, startTime, endTime time.Time) (map[string]struct{}, error) {
	taskCreatTimes := make(map[string]struct{}, 0)
	// ["03-19,04-11"]
	cycleArrays, err := h.getCycleArrays(todo)
	if err != nil {
		return nil, err
	}
	years, _, _ := startTime.Date()
	yeare, _, _ := endTime.Date()
	for yearNum := years; yearNum <= yeare; yearNum++ {
		for i := range cycleArrays {
			str := cycleArrays[i]
			ss := strings.Split(str, "-")

			monthstr := ss[0]
			month, _ := strconv.Atoi(monthstr)

			daystr := ss[1]
			day, _ := strconv.Atoi(daystr)

			now := time.Date(yearNum, time.Month(month), day, 1, 0, 0, 0, time.Now().Location())
			cTimeStr := utils.ParseTimeToStr(now)
			if !endTime.Before(now) && !now.Before(startTime) {
				// taskCreatTimes = append(taskCreatTimes, cTimeStr)
				taskCreatTimes[cTimeStr] = struct{}{}
			}
		}
	}
	return taskCreatTimes, nil
}

// 检查待办 将待办转化为历史任务，用于新生成的todo
func (h *YearHandle) CheckToDoToRandomHistoryTask(todo *model.Todo, untilTime time.Time) ([]*model.Task, error) {
	// ["03-19,04-11"]
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

	years, _, _ := startTime.Date()

	et := todo.EndTime
	endTime, err := utils.ParseStrTime(et)
	if err != nil {
		return nil, err
	}

	yeare, _, _ := endTime.Date()

	yearNum := 0
	for yearNum = years; yearNum <= yeare; yearNum++ {
		for i := range cycleArrays {
			str := cycleArrays[i]
			ss := strings.Split(str, "-")

			monthstr := ss[0]
			month, _ := strconv.Atoi(monthstr)

			daystr := ss[1]
			day, _ := strconv.Atoi(daystr)

			now := time.Date(yearNum, time.Month(month), day, 1, 0, 0, 0, time.Now().Location())

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

	return tasks, nil
}

// 获取开始时间结束时间
func (h *YearHandle) GetTimeArrange(nowTime time.Time) (startTime, endTime time.Time, err error) {
	// 结束时间是一年的最后一天,开始时间是年的第一天
	et := utils.GetLastDayOfYear(nowTime)
	st := utils.GetFirstDayOfYear(nowTime)

	return st, et, nil
}
