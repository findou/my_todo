/*
 * @Version: 1.0
 * @Date: 2023-03-17 13:19:53
 * @LastEditTime: 2023-03-18 21:19:31
 */
package model

import (
	"log"
	"mytodo/utils"
	"time"
)

type Todo struct {
	Name        string `json:"name" gorm:"Column:name"`
	Id          string `json:"id" gorm:"Column:id"`
	Priority    int    `json:"priority" gorm:"Column:priority"`
	CreateTime  string `json:"createTime" gorm:"Column:create_time"`
	UpdateTime  string `json:"updateTime" gorm:"Column:update_time"`
	StartTime   string `json:"startTime" gorm:"Column:start_time"`
	EndTime     string `json:"endTime" gorm:"Column:end_time"`
	CycleType   int    `json:"cycleType" gorm:"Column:cycle_type"`
	CycleArrays string `json:"cycleArrays" gorm:"Column:cycle_arrays"`
	Status      bool   `json:"status" gorm:"Column:status"`
}

func (t *Todo) TableName() string {
	return "todo"
}

func (t *Todo) Save() error {
	DB.AutoMigrate(&Todo{})
	if err := DB.Create(t).Error; err != nil {
		log.Fatalf("save todo error: %v", err)
		return err
	}
	return nil
}

func (t *Todo) Del(id string) error {
	t.Id = id
	updateTime := utils.ParseTimeToStr(time.Now())

	if err := DB.Model(&Todo{Id: id, Status: true}).Updates(
		Todo{UpdateTime: updateTime, Status: false}).Error; err != nil {
		log.Fatalf("update todo error: %v", err)
		return err

	}
	return nil
}

// 查询所有待办列表
func (t *Todo) ListTodos() ([]*Todo, error) {
	todos := make([]*Todo, 0)
	todo := &Todo{Status: true}
	if err := DB.Model(todo).Find(&todos).Error; err != nil {
		log.Fatalf("list todo error: %v", err)
		return nil, err
	}

	return todos, nil
}

// 更具时间周期查询待办
func (t *Todo) FindAllTodoByTimeRange(mystartTime, myendTime time.Time) ([]*Todo, error) {
	todos := make([]*Todo, 0)
	// 开始时间小于等于区间的结束时间 并且 结束时间大于等于区间的开始时间
	if err := DB.Where("start_time <= ?", myendTime).Where("end_time >=  ?", mystartTime).
		Where("status = 1").Order("priority desc").Find(&todos).Error; err != nil {
		log.Fatalf("FindAllTodoByTimeRange error: %v", err)
		return nil, err
	}
	return todos, nil
}
