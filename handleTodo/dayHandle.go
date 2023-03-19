/*
 * @Version: 1.0
 * @Date: 2023-03-17 13:54:50
 * @LastEditTime: 2023-03-19 17:28:34
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

type DayHandle struct {
	cycleType     int
	cycleTypeName string
}

func NewDayHandle(cycleType int) *DayHandle {
	if cycleType == 1 {
		cycleTypeName := "DAY"
		return &DayHandle{cycleType: cycleType, cycleTypeName: cycleTypeName}
	}
	return nil
}

// 检查 cycleArrays 字段是否符合规范，将不符合规范的移除，并检查开始，结束时间是否符合规范
func (h *DayHandle) CheckCycleArrays(checkCycleArrays string, startTime, endTime time.Time) (string, error) {

	// 开始时间必须是一天的开始 ，结束时间必须是一天结束
	et := utils.GetLastTime(endTime)
	st := utils.GetZeroTime(startTime)

	if (!endTime.Equal(et)) || (!startTime.Equal(st)) {
		return "", fmt.Errorf("")
	}

	// ["8:30","6:50"]
	carr := make([]string, 0)
	newArr := make([]string, 0)
	err := json.Unmarshal([]byte(checkCycleArrays), &carr)
	if err != nil {
		return "", err
	}
	for i := range carr {
		str := carr[i]
		strarr := strings.Split(str, ":")
		if len(strarr) != 2 {
			continue
		} else {
			hour := strarr[0]
			h, _ := strconv.ParseInt(hour, 0, 0)
			if h < 0 || h > 24 {
				continue
			}
			min := strarr[1]
			m, _ := strconv.ParseInt(min, 0, 0)
			if m < 0 || m > 60 {
				continue
			}
			newArr = append(newArr, str)
		}
	}
	b, _ := json.Marshal(newArr)

	cycleArrays := (string(b))
	return cycleArrays, nil
}

// 检查待办 将待办转化为任务，并保存
func (h *DayHandle) CheckToDoToTask(todo *model.Todo, now time.Time) error {
	// ["8:30","6:50"]
	cycleArrays, err := h.getCycleArrays(todo)
	if err != nil {
		return err
	}

	for i := range cycleArrays {
		st, err := utils.ParseStrTime(todo.StartTime)
		if err != nil {
			continue
		}
		et, err := utils.ParseStrTime(todo.EndTime)
		if err != nil {
			continue
		}

		cstr := cycleArrays[i]
		strarr := strings.Split(cstr, ":")

		hourstr := strarr[0]
		hour, _ := strconv.Atoi(hourstr)

		minstr := strarr[1]
		min, _ := strconv.Atoi(minstr)
		ti := time.Date(now.Year(), now.Month(), now.Day(), hour, min, 0, 0, now.Location())
		if st.Before(ti) && ti.Before(et) {
			task := &model.Task{
				Name:       todo.Name,
				Id:         utils.GenUuid(),
				Priority:   todo.Priority,
				TodoId:     todo.Id,
				CreateTime: utils.ParseTimeToStr(ti),
				Status:     false,
			}
			err := task.Save()
			return err
		}
	}

	return nil

}

func (h *DayHandle) getCycleArrays(todo *model.Todo) ([]string, error) {
	cycleArrays := make([]string, 0)
	ca := []byte(todo.CycleArrays)
	err := json.Unmarshal(ca, &cycleArrays)

	if err != nil {
		return nil, err
	}
	return cycleArrays, nil
}

// 检查 cycleArrays 字段包含的时间是否存在 在 传入的开始结束时间中
func (h *DayHandle) CheckContainTodo(todo *model.Todo, startTime, endTime time.Time) (map[string]struct{}, error) {
	taskCreatTimes := make(map[string]struct{}, 0)
	// ["8:30","6:50"]
	cycleArrays, err := h.getCycleArrays(todo)
	if err != nil {
		return nil, err
	}

	years, months, days := startTime.Date()
	yeare, monthe, daye := endTime.Date()

	for yearNum := years; yearNum <= yeare; yearNum++ {
		for monthNum := months; monthNum <= monthe; monthNum++ {
			for dayNum := days; dayNum <= daye; dayNum++ {
				for i := range cycleArrays {
					cyclestr := cycleArrays[i]
					strarr := strings.Split(cyclestr, ":")

					hourstr := strarr[0]
					hour, _ := strconv.Atoi(hourstr)

					minstr := strarr[1]
					min, _ := strconv.Atoi(minstr)
					cTime := time.Date(yearNum, time.Month(monthNum), dayNum, hour, min, 0, 0, time.Now().Location())

					cTimeStr := utils.ParseTimeToStr(cTime)
					if !endTime.Before(cTime) && !cTime.Before(startTime) {
						//taskCreatTimes = append(taskCreatTimes, cTimeStr)
						taskCreatTimes[cTimeStr] = struct{}{}
					}
				}
			}
		}
	}
	return taskCreatTimes, nil
}

// 检查待办 将待办转化为历史任务，用于新生成的todo
func (h *DayHandle) CheckToDoToRandomHistoryTask(todo *model.Todo, untilTime time.Time) ([]*model.Task, error) {
	// ["8:30","6:50"]
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

	years, months, days := startTime.Date()

	et := todo.EndTime
	endTime, err := utils.ParseStrTime(et)
	if err != nil {
		return nil, err
	}

	yeare, monthe, daye := endTime.Date()

	for yearNum := years; yearNum <= yeare; yearNum++ {
		for monthNum := months; monthNum <= monthe; monthNum++ {
			for dayNum := days; dayNum <= daye; dayNum++ {
				for i := range cycleArrays {
					//now := time.Date(yearNum, time.Month(monthNum), dayNum, 1, 0, 0, 0, time.Now().Location())

					cyclestr := cycleArrays[i]
					strarr := strings.Split(cyclestr, ":")

					hourstr := strarr[0]
					hour, _ := strconv.Atoi(hourstr)

					minstr := strarr[1]
					min, _ := strconv.Atoi(minstr)
					cTime := time.Date(yearNum, time.Month(monthNum), dayNum, 1, 0, 0, 0, time.Now().Location())
					rTime := time.Date(yearNum, time.Month(monthNum), dayNum, hour, min, 0, 0, time.Now().Location())
					if !untilTime.Before(cTime) {
						task := &model.Task{
							Name:       todo.Name,
							Id:         utils.GenUuid(),
							Priority:   todo.Priority,
							TodoId:     todo.Id,
							CreateTime: utils.ParseTimeToStr(rTime),
							Status:     false,
						}
						tasks = append(tasks, task)
						//continue
					}
					//fmt.Println(cTime)
				}
			}
		}
	}

	return tasks, nil
}

// 获取开始时间结束时间
func (h *DayHandle) GetTimeArrange(nowTime time.Time) (startTime, endTime time.Time, err error) {
	// 开始时间是一天的开始 ，结束时间是一天结束
	et := utils.GetLastTime(nowTime)
	st := utils.GetZeroTime(nowTime)

	return st, et, nil
}
