/*
 * @Version: 1.0
 * @Date: 2023-03-19 15:33:07
 * @LastEditTime: 2023-03-19 17:28:08
 */
package data

import (
	handletodo "mytodo/handleTodo"
	"mytodo/model"
	"mytodo/utils"
)

func ImportData(data *model.DataOracles) error {
	datatodos := data.ResData.DataTodos
	for i, _ := range datatodos {
		dataTodo := datatodos[i]
		todo := &model.Todo{
			Id:          dataTodo.Id,
			Name:        dataTodo.Name,
			Priority:    dataTodo.Priority,
			CreateTime:  dataTodo.CreateTime,
			UpdateTime:  dataTodo.CreateTime,
			StartTime:   dataTodo.StartTime,
			EndTime:     dataTodo.EndTime,
			CycleType:   dataTodo.CycleType,
			CycleArrays: dataTodo.CycleArrays,
			Status:      dataTodo.Status,
		}
		err := handletodo.CheckTodo(todo)
		if err != nil {
			continue
		}
		// 保存到数据库
		err = todo.Save()
		if err != nil {
			continue
		}

		err = handleTask(dataTodo, todo)
		if err != nil {
			continue
		}
	}
	return nil
}

// 处理校验存储task
func handleTask(dataTodo *model.DataTodo, todo *model.Todo) error {
	tasks := make([]*model.Task, 0)

	// 时间
	ststr := todo.StartTime
	startTime, _ := utils.ParseStrTime(ststr)
	etstr := todo.EndTime
	endTime, _ := utils.ParseStrTime(etstr)
	createTimes, err := handletodo.CheckTodoTaskTimeRange(todo, startTime, endTime)
	if err != nil {
		return err
	}
	todotasks := dataTodo.DataTodoTasks

	for i, _ := range todotasks {
		dataTask := todotasks[i]
		dataCreate := string(dataTask.CreateTime)

		_, ex := createTimes[dataCreate]
		if ex {
			task := &model.Task{
				Id:         dataTask.Id,
				Name:       dataTask.Name,
				TodoId:     todo.Id,
				Priority:   todo.Priority,
				CreateTime: dataCreate,
				UpdateTime: dataTask.UpdateTime,
				Status:     dataTask.Status,
			}
			tasks = append(tasks, task)
		}
	}
	// 保存tasks
	ta := new(model.Task)
	err = ta.SaveAllTasks(tasks)
	if err != nil {
		return err
	}
	return nil
}
