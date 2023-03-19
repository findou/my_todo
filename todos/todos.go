/*
 * @Version: 1.0
 * @Date: 2023-03-17 13:01:31
 * @LastEditTime: 2023-03-18 10:33:24
 */
package todos

import (
	"encoding/json"
	"fmt"
	"io"
	handletodo "mytodo/handleTodo"
	"mytodo/model"
	"net/http"
	"strings"
)

func TodosHandler(w http.ResponseWriter, r *http.Request) {
	m := r.Method
	if m == http.MethodPost {
		post(w, r)
		return
	}

	if m == http.MethodDelete {
		del(w, r)
		return
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
}

// 创建待办
func post(w http.ResponseWriter, r *http.Request) {
	// POST  http://127.0.0.1:8888/todo/

	todo, err := checkFormData(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = handletodo.HandleTodo(todo)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// 删除待办
func del(w http.ResponseWriter, r *http.Request) {
	// DELETE http://127.0.0.1:8888/todo/<todo_id>
	path := r.URL.EscapedPath()
	todo_id := strings.Split(path, "/")[2]
	fmt.Println(todo_id)

	todo := new(model.Todo)
	err := todo.Del(todo_id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}

func checkFormData(body io.Reader) (*model.Todo, error) {
	todo := &model.Todo{}

	err := json.NewDecoder(body).Decode(todo)
	if err != nil {
		return nil, err
	}
	err = handletodo.CheckTodo(todo)
	if err != nil {
		return nil, err
	}
	return todo, nil
}
