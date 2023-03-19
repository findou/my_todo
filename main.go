/*
 * @Version: 1.0
 * @Date: 2023-03-17 12:19:49
 * @LastEditTime: 2023-03-19 19:58:59
 */
package main

import (
	"log"
	dataoracles "mytodo/dataOracles"
	"mytodo/schedule"
	"mytodo/statisViews"
	"mytodo/tasks"
	"mytodo/todos"
	"net/http"
)

const (
	LISTEN_ADDRESS = "127.0.0.1:8888"
)

func main() {
	// 定时任务,每天1点根据待办项生成今日任务
	go schedule.HandleScheduleTask()

	http.HandleFunc("/todo/", todos.TodosHandler)
	http.HandleFunc("/statisView/", statisViews.CycleviewsHandler)
	http.HandleFunc("/task/", tasks.TasksHandler)
	http.HandleFunc("/dataOracle", dataoracles.DataOraclesHandler)
	log.Printf("【待办任务服务】启动成功，监听地址：%q\n", LISTEN_ADDRESS)
	err := http.ListenAndServe(LISTEN_ADDRESS, nil)

	// err := http.ListenAndServe("0.0.0.0:9999", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
