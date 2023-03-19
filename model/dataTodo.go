/*
 * @Version: 1.0
 * @Date: 2023-03-19 13:43:59
 * @LastEditTime: 2023-03-19 19:16:00
 */
package model

type ResData struct {
	TotalTaskNum           int         `json:"totalTaskNum"`           // 待办查询周期范围内的任务总量
	TotalTodoNum           int         `json:"totalTodoNum"`           // 待办查询周期范围内的待办总量
	TotalTaskFinishedNum   int         `json:"totalTaskFinishedNum"`   // 待办查询周期范围内的任务已经完成的总量
	TotalTaskFinishedRatio float64     `json:"TotalTaskFinishedRatio"` // 待办查询周期范围内的任务的完成率,百分率,保留小数点后两位
	DataTodos              []*DataTodo `json:"dataTodos"`              //查询周期内包含的待办列表
}

type DataTodo struct {
	Name                  string          `json:"name"`                  // 待办名称
	Id                    string          `json:"id"`                    // 待办id
	Priority              int             `json:"priority"`              // 待办优先级
	CreateTime            string          `json:"createTime"`            // 待办创建时间
	UpdateTime            string          `json:"updateTime"`            // 待办更新时间
	StartTime             string          `json:"startTime"`             // 待办开始时间
	EndTime               string          `json:"endTime"`               // 待办结束时间
	CycleType             int             `json:"cycleType"`             // 周期类型，1日，2周，3月，4年
	CycleArrays           string          `json:"cycleArrays"`           // 周期内的值，更具类型不同设置
	Status                bool            `json:"status"`                // 待办状态 true 启用，false 未启用
	TodoTaskNum           int             `json:"todoTaskNum"`           // 待办查询周期范围内的当前待办的任务总量
	TodoTaskFinishedNum   int             `json:"todoTaskFinishedNum"`   // 待办查询周期范围内的当前待办的任务已经完成的总量
	TodoTaskFinishedRatio float64         `json:"todoTaskFinishedRatio"` // 待办查询周期范围内的当前待办的任务的完成率,百分率，保留小数点后两位
	DataTodoTasks         []*DataTodoTask // 待办查询周期范围内的任务，包括已完成的/未完成的，未完成的包括数据库定时任务已经生成的和未生成的
}

type DataTodoTask struct {
	Name       string `json:"name"`       // 任务名称同待办名称
	Id         string `json:"id"`         // 任务id
	TodoId     string `json:"todoId"`     // 关联的待办id
	Priority   int    `json:"priority"`   // 优先级，继承关联的待办的优先级
	CreateTime string `json:"createTime"` // 创建时间
	UpdateTime string `json:"updateTime"` // 更新时间
	Status     bool   `json:"status"`     // 状态
}
