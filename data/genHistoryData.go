/*
 * @Version: 1.0
 * @Date: 2023-03-18 08:53:12
 * @LastEditTime: 2023-03-19 18:16:34
 */
package data

import (
	"encoding/json"
	"io/ioutil"
	"log"
	handletodo "mytodo/handleTodo"
	"mytodo/model"
	"os"
	"time"
)

// 根据json文件生成历史待办
func GenHistoryTodoData() {
	url := "../data/" + "todo.json"
	// url := "todo.json"
	// body, err := os.ReadFile(url)
	f, err := os.Open(url)
	if err != nil {
		log.Fatal(err)
	}
	body, err := ioutil.ReadAll(f)

	if err != nil {
		log.Fatal(err)
	}
	todos := make([]*model.Todo, 0)
	err = json.Unmarshal(body, &todos)
	if err != nil {
		log.Fatal(err)
	}
	for i := range todos {
		todo := todos[i]
		err := handletodo.CheckAndHandleTodo(todo)
		if err != nil {
			log.Println(err.Error())
			continue
		}
	}
}

// 生成历史待办数据，随机选择一些更新状态为已完成
func GenHistoryTask(todos []*model.Todo) error {
	tasks := make([]*model.Task, 0)
	now := time.Now()
	for i := range todos {
		todo := todos[i]
		ts, err := handletodo.CheckToDoToRandomHistoryTaskByType(todo.CycleType, todo, now)
		if err != nil {
			return err
		}
		tasks = append(tasks, ts...)
	}
	task := new(model.Task)
	err := task.SaveAllTasks(tasks)
	if err != nil {
		return err
	}
	// 随机选择一些更新状态为已完成
	err = ChooseSomeTasksToFinish(tasks)
	if err != nil {
		return err
	}
	return nil
}

// 随机选择一些任务将设置为完成
func ChooseSomeTasksToFinish(tasks []*model.Task) error {
	newTasks := make([]*model.Task, 0)
	task := new(model.Task)
	// 假如完成率为80%
	// i+1 对10 取余为 3，5 的不更新完成状态，即为未完成
	for i := range tasks {
		oneTask := tasks[i]

		num := i + 1
		if num%10 != 3 && num%10 != 5 {
			newTasks = append(newTasks, oneTask)
		}
	}
	task.UpdateAllStatus(newTasks)
	return nil
}

// 根据历史待办生成历史任务，并随机选择一些任务更改状态为已完成
func GenHistoryTaskData() {
	// task := &model.Task{}
	// 查询待办列表，生成任务
	todo := new(model.Todo)
	todos, err := todo.ListTodos()
	if err != nil {
		return
	}
	err = GenHistoryTask(todos)
	if err != nil {
		return
	}
}
