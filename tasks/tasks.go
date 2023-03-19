/*
 * @Version: 1.0
 * @Date: 2023-03-17 13:06:21
 * @LastEditTime: 2023-03-18 01:42:09
 */
package tasks

import (
	"encoding/json"
	"fmt"
	"mytodo/model"
	"net/http"
	"strings"
	"time"
)

func TasksHandler(w http.ResponseWriter, r *http.Request) {
	m := r.Method
	if m == http.MethodGet {
		get(w, r)
		return
	}

	if m == http.MethodPut {
		put(w, r)
		return
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
}

// 查看今日任务
func get(w http.ResponseWriter, r *http.Request) {
	// GET http://127.0.0.1:8888/task/
	task := new(model.Task)
	ti := time.Now()
	//ti := utils.ParseStrTime("2023-03-20 01:00:00")
	tasks, err := task.ListTodyTask(ti)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	body, err := json.Marshal(tasks)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Write(body)
	w.WriteHeader(http.StatusOK)
}

// 修改任务完成状态为已完成
func put(w http.ResponseWriter, r *http.Request) {
	// PUT http://127.0.0.1:8888/task/<task_id>
	path := r.URL.EscapedPath()
	task_id := strings.Split(path, "/")[2]

	fmt.Println(task_id)

	task := new(model.Task)
	err := task.UpdateStatus(task_id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}
