/*
 * @Version: 1.0
 * @Date: 2023-03-17 13:53:43
 * @LastEditTime: 2023-03-19 19:18:01
 */
package handletodo

import (
	"fmt"
	"mytodo/model"
	"mytodo/utils"
	"strconv"
	"time"
)

const (
	PRI_MIN = 1
	PRI_MAX = 10
)

// CycleType (DAY = 1, ...).
type CycleType int

const (
	DAY CycleType = iota + 1
	WEEK
	MONTH
	YEAR
)

type TodoHandler interface {
	CheckCycleArrays(checkCycleArrays string, startTime, endTime time.Time) (string, error)
	CheckToDoToTask(todo *model.Todo, now time.Time) error
	CheckToDoToRandomHistoryTask(todo *model.Todo, untilTime time.Time) ([]*model.Task, error)
	// 更具日期所在的周期，获取周期开始时间结束时间
	GetTimeArrange(nowTime time.Time) (startTime, endTime time.Time, err error)
	// 检查 cycleArrays 字段包含的时间是否在 传入的开始结束时间中
	CheckContainTodo(todo *model.Todo, startTime, endTime time.Time) (map[string]struct{}, error)
}

func genHandler(cycleType int) (TodoHandler, error) {
	if cycleType == 1 {
		return NewDayHandle(cycleType), nil
	} else if cycleType == 2 {
		return NewWeekHandle(cycleType), nil
	} else if cycleType == 3 {
		return NewMonthHandle(cycleType), nil
	} else if cycleType == 4 {
		return NewYearHandle(cycleType), nil
	}
	return nil, fmt.Errorf("not support this cycleType")
}

func CheckAndHandleTodo(todo *model.Todo) error {
	err := CheckTodo(todo)
	if err != nil {
		return err
	}
	return HandleTodo(todo)
}

func CheckTodo(todo *model.Todo) error {
	cycleType := todo.CycleType

	name := todo.Name
	if name == "" {
		return fmt.Errorf("name not support")
	}

	priority := todo.Priority
	if priority > PRI_MAX || priority < PRI_MIN {
		return fmt.Errorf("priority not support")
	}

	startTime := todo.StartTime
	st, err := utils.ParseStrTime(startTime)
	if err != nil {
		return err
	}

	endTime := todo.EndTime
	et, err := utils.ParseStrTime(endTime)
	if err != nil {
		return err
	}

	if et.Before(st) {
		// 结束时间需在开始时间之后
		return fmt.Errorf("time not support")
	}

	/*
		// 批量生成历史任务时，暂时解除任务现在时间限制
			now := time.Now()
			 todo.CreateTime = utils.ParseTimeToStr(now)

			 if st.Before(now) {
				// 开始时间不能在现在时间之前
				return fmt.Errorf("time not support")
			}
	*/

	handler, err := genHandler(cycleType)
	if err != nil {
		return err
	}
	checkCycleArrayStr := todo.CycleArrays

	cycleArrays, err := handler.CheckCycleArrays(checkCycleArrayStr, st, et)
	if err != nil {
		return err
	}
	todo.CycleArrays = cycleArrays

	return nil
}

func HandleTodo(todo *model.Todo) error {
	now := time.Now()
	nowstr := utils.ParseTimeToStr(now)
	newtodo := &model.Todo{
		Name:        todo.Name,
		Id:          utils.GenUuid(),
		Priority:    todo.Priority,
		CreateTime:  nowstr,
		CycleType:   todo.CycleType,
		StartTime:   todo.StartTime,
		EndTime:     todo.EndTime,
		CycleArrays: todo.CycleArrays,
		Status:      true,
	}
	// 保存到数据库
	err := newtodo.Save()
	return err
}

// 根据类型生成待办任务，并保存
func CheckToDoToTaskByType(cycleType int, todo *model.Todo, now time.Time) error {
	handler, err := genHandler(cycleType)
	if err != nil {
		return err
	}
	return handler.CheckToDoToTask(todo, now)
}

func GetTimeArrangeByCt(dtime time.Time, cycleType int) (time.Time, time.Time, error) {
	handler, err := genHandler(cycleType)
	if err != nil {
		return time.Now(), time.Now(), err
	}
	return handler.GetTimeArrange(dtime)
}

// 更具开始时间结束时间查询待办数据
func HandleDataWithRangeTime(mst, met time.Time) (*model.ResData, error) {
	// 更具开始时间结束时间查询待办
	todo := &model.Todo{}
	todos, err := todo.FindAllTodoByTimeRange(mst, met)
	if err != nil {
		return nil, err
	}
	// 处理这些todos返回新的todos
	return checkContainTodos(todos, mst, met)
}

// 根据todo类型和，起止时间 获取任务 的创建时间列表
func CheckTodoTaskTimeRange(todo *model.Todo, startTime, endTime time.Time) (map[string]struct{}, error) {
	handler, err := genHandler(todo.CycleType)
	if err != nil {
		return nil, err
	}
	return handler.CheckContainTodo(todo, startTime, endTime)
}

// 如果当前待办的周期大于 查询周期，需要细化筛选cycleArrays字段判断当前查询周期是否包含这个待办任务
func checkContainTodos(todos []*model.Todo, mst, met time.Time) (*model.ResData, error) {
	dataTodos := make([]*model.DataTodo, 0)

	totalTaskNum := 0
	totalTodoNum := 0
	totalTaskFinishedNum := 0
	var totalTaskFinishedRatio float64
	tasksGroup := taskGroupByOrderId(mst, met)
	for i := range todos {
		todo := todos[i]
		todoCycleType := todo.CycleType
		// 如果当前待办的周期大于 查询周期，需要细化筛选cycleArrays字段判断当前查询周期是否包含这个待办任务
		//if todoCycleType > queryCt {
		handler, err := genHandler(todoCycleType)
		if err != nil {
			continue
		}
		ststr := todo.StartTime
		st, err := utils.ParseStrTime(ststr)
		if err != nil {
			continue
		}
		etstr := todo.EndTime
		et, err := utils.ParseStrTime(etstr)
		if err != nil {
			continue
		}
		startTime := mst
		if mst.Before(st) {
			startTime = st
		}
		endTime := met
		if et.Before(met) {
			endTime = et
		}

		taskCreatTimes, err := handler.CheckContainTodo(todo, startTime, endTime)
		if err != nil {
			continue
		}
		if len(taskCreatTimes) != 0 {
			tasks := tasksGroup[todo.Id]
			dataTodo := genTodoTask(todo, taskCreatTimes, tasks)
			dataTodos = append(dataTodos, dataTodo)
			totalTodoNum = totalTodoNum + 1
			totalTaskNum = totalTaskNum + dataTodo.TodoTaskNum
			totalTaskFinishedNum = totalTaskFinishedNum + dataTodo.TodoTaskFinishedNum
		}
	}
	totalTaskFinishedRatio = float64(totalTaskFinishedNum) / float64(totalTaskNum) * 100
	totalTaskFinishedRatio, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", totalTaskFinishedRatio), 64)

	resData := &model.ResData{
		TotalTaskNum:           totalTaskNum,
		TotalTodoNum:           totalTodoNum,
		TotalTaskFinishedNum:   totalTaskFinishedNum,
		TotalTaskFinishedRatio: totalTaskFinishedRatio,
		DataTodos:              dataTodos,
	}
	return resData, nil
}

// 导出数据查询待办任务时候调用
func genTodoTask(todo *model.Todo, taskCreatTimes map[string]struct{}, tasks []*model.Task) *model.DataTodo {
	ctArr := make([]string, 0)
	for k, _ := range taskCreatTimes {
		ctArr = append(ctArr, k)
	}

	dataTodoTasks := make([]*model.DataTodoTask, 0)
	todoTaskNum := 0
	todoTaskFinishedNum := 0

	for i := 0; i < len(ctArr); i++ {
		taskCreateTimeStr := ctArr[i]
		dataTodoTask := &model.DataTodoTask{
			Name:       todo.Name,
			TodoId:     todo.Id,
			Priority:   todo.Priority,
			CreateTime: taskCreateTimeStr,
			Status:     false,
		}

		for j, _ := range tasks {
			ts := tasks[j]
			if taskCreateTimeStr == (ts.CreateTime) {
				dataTodoTask.Id = ts.Id
				dataTodoTask.UpdateTime = ts.UpdateTime
				dataTodoTask.Status = ts.Status
			}
		}
		dataTodoTasks = append(dataTodoTasks, dataTodoTask)
		if dataTodoTask.Status {
			todoTaskFinishedNum = todoTaskFinishedNum + 1
		}
		todoTaskNum = todoTaskNum + 1
	}

	todoTaskFinishedRatio := float64(todoTaskFinishedNum) / float64(todoTaskNum) * 100
	todoTaskFinishedRatio, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", todoTaskFinishedRatio), 64)

	todoTask := &model.DataTodo{
		Id:                    todo.Id,
		Name:                  todo.Name,
		Priority:              todo.Priority,
		CreateTime:            todo.CreateTime,
		UpdateTime:            todo.UpdateTime,
		StartTime:             todo.StartTime,
		EndTime:               todo.EndTime,
		CycleType:             todo.CycleType,
		CycleArrays:           todo.CycleArrays,
		Status:                todo.Status,
		TodoTaskNum:           todoTaskNum,
		TodoTaskFinishedNum:   todoTaskFinishedNum,
		TodoTaskFinishedRatio: todoTaskFinishedRatio,
		DataTodoTasks:         dataTodoTasks,
	}
	return todoTask
}

func taskGroupByOrderId(startTime, endTime time.Time) map[string][]*model.Task {
	res := make(map[string][]*model.Task, 0)
	task := &model.Task{}
	newTasks, _ := task.FindAllTaskByTimeRange(startTime, endTime)
	for i, _ := range newTasks {
		task := newTasks[i]
		_, f := res[task.TodoId]
		if f {
			res[task.TodoId] = append(res[task.TodoId], task)
		} else {
			tasks := make([]*model.Task, 0)
			tasks = append(tasks, task)
			res[task.TodoId] = tasks
		}
	}
	return res
}

// 根据todo生成历史待办数据
func CheckToDoToRandomHistoryTaskByType(cycleType int, todo *model.Todo, now time.Time) ([]*model.Task, error) {
	handler, err := genHandler(cycleType)
	if err != nil {
		return nil, err
	}
	return handler.CheckToDoToRandomHistoryTask(todo, now)
}
