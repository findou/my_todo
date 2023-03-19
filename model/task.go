/*
 * @Version: 1.0
 * @Date: 2023-03-17 22:05:11
 * @LastEditTime: 2023-03-19 03:05:39
 */
package model

import (
	"log"
	"mytodo/utils"
	"time"
)

type Task struct {
	Name       string `json:"name" gorm:"Column:name"`
	Id         string `json:"id" gorm:"Column:id"`
	TodoId     string `json:"todoId" gorm:"Column:todo_id"`
	Priority   int    `json:"priority" gorm:"Column:priority"`
	CreateTime string `json:"createTime" gorm:"Column:create_time"`
	UpdateTime string `json:"updateTime" gorm:"Column:update_time"`
	Status     bool   `json:"status" gorm:"Column:status"`
}

func (t *Task) TableName() string {
	return "task"
}

func (t *Task) Save() error {
	DB.AutoMigrate(&Task{})
	t.Status = false // 未完成
	if err := DB.Create(t).Error; err != nil {
		log.Fatalf("save task error: %v", err)
		return err
	}
	log.Printf("今日待办任务生成成功 ,待办id: %q, 任务id: %q, 待办任务名称: %s", t.TodoId, t.Id, t.Name)
	return nil
}

func (t *Task) SaveAllTasks(tasks []*Task) error {
	DB.AutoMigrate(&Task{})
	if err := DB.Create(&tasks).Error; err != nil {
		log.Fatalf("save tasks error: %v", err)
		return err
	}
	log.Println("所有待办任务生成成功 ")
	return nil
}

func (t *Task) ListTodyTask(tody time.Time) ([]*Task, error) {
	startTime := utils.GetZeroTime(tody)
	endTime := utils.GetLastTime(tody)

	tasks := make([]*Task, 0)

	if err := DB.Where("create_time BETWEEN ? AND ?", startTime, endTime).Find(&tasks).Order("create_time desc").Error; err != nil {
		log.Fatalf("update task error: %v", err)
		return nil, err

	}

	return tasks, nil
}

func (t *Task) UpdateAllStatus(tasks []*Task) error {

	updateTime := utils.ParseTimeToStr(time.Now())
	if err := DB.Model(&tasks).Updates(
		Task{UpdateTime: updateTime, Status: true}).Error; err != nil {
		log.Fatalf("update tasks error: %v", err)
		return err

	}
	return nil

}
func (t *Task) UpdateStatus(id string) error {

	updateTime := utils.ParseTimeToStr(time.Now())

	if err := DB.Model(&Task{Id: id, Status: false}).Updates(
		Task{UpdateTime: updateTime, Status: true}).Error; err != nil {
		log.Fatalf("update task error: %v", err)
		return err

	}
	return nil
}

func (t *Task) FindAllTaskByTimeRange(mystartTime, myendTime time.Time) ([]*Task, error) {
	tasks := make([]*Task, 0)
	// create_time 在开始结束时间之间
	if err := DB.Where("create_time BETWEEN ? AND ?", mystartTime, myendTime).
		Order("create_time").Find(&tasks).Error; err != nil {
		log.Fatalf("FindAllTaskByTimeRange error: %v", err)
		return nil, err
	}
	return tasks, nil
}

func (t *Task) FindTasks(todoId string, startTime, endTime time.Time) ([]*Task, error) {
	tasks := make([]*Task, 0)
	// create_time 在开始结束时间之间
	if err := DB.Where("create_time <= ?", startTime).Where("create_time >=  ?", endTime).
		Where("todo_id = ?", todoId).
		Order("create_time ").Find(&tasks).Error; err != nil {
		log.Fatalf("FindTasks tasks error: %v", err)
		return nil, err
	}
	return tasks, nil
}
