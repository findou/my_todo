/*
 * @Version: 1.0
 * @Date: 2023-03-17 22:47:15
 * @LastEditTime: 2023-03-19 16:22:10
 */
package schedule

import (
	"log"
	handletodo "mytodo/handleTodo"
	"mytodo/model"
	"mytodo/utils"
	"time"
)

type ScheduleTask struct {
	todo    *model.Todo
	nowTime time.Time
}

func (t *ScheduleTask) Execute(num int) {
	log.Printf("当前协程: %q 开始执行 生成待办任务 ,待办id: %q ,待办名称: %s ", num, t.todo.Id, t.todo.Name)
	ct := t.todo.CycleType
	err := handletodo.CheckToDoToTaskByType(ct, t.todo, t.nowTime)
	if err != nil {
		return
	}
}

func HandleScheduleTask() {
	// 每天 1点0分触发
	dur := SetTime(1, 0, 0)
	timer := time.NewTimer(dur)
	defer timer.Stop()
	for t := range timer.C {
		ti := utils.ParseTimeToStr(t)
		log.Printf("开始执行每天生成当日待办的定时任务，当前时间：%s", ti)
		timer.Reset(time.Hour * 24)
		// 执行定时任务
		ExcuteTodoTask(t)
	}
}

func SetTime(hour, min, second int) time.Duration {
	now := time.Now()
	setTime := time.Date(now.Year(), now.Month(), now.Day(), hour, min, second, 0, now.Location())
	dur := setTime.Sub(now)
	if dur > 0 {
		return dur // 大于当前时间返回
	}
	return dur + time.Hour*24 // 否则为第二天的时间
}

func ExcuteTodoTask(now time.Time) {
	// 查询待办列表，生成任务
	todo := new(model.Todo)
	todos, err := todo.ListTodos()
	if err != nil {
		return
	}

	GenTodayTask(todos, now)
}

// 生成任务,放到任务池里执行
func GenTodayTask(todos []*model.Todo, now time.Time) error {
	taskPool := NewTaskPool(5)
	for i := range todos {
		t := todos[i]
		st, err := utils.ParseStrTime(t.StartTime)
		if err != nil {
			continue
		}
		et, err := utils.ParseStrTime(t.EndTime)
		if err != nil {
			continue
		}
		if !now.Before(st) && !et.Before(now) {
			//放入n个task（可以阻塞） 如
			taskPool.Put(&ScheduleTask{todo: t, nowTime: now})
		}
	}
	//等待任务执行完毕
	taskPool.Wait()
	return nil
}
