/*
 * @Version: 1.0
 * @Date: 2023-03-18 01:03:02
 * @LastEditTime: 2023-03-19 13:11:13
 */
package test

import (
	"fmt"
	"mytodo/model"
	"mytodo/utils"
	"testing"
)

func TestListTodyTask(t *testing.T) {
	ti, _ := utils.ParseStrTime("2023-03-19 01:00:00")
	task := new(model.Task)
	tasks, _ := task.ListTodyTask(ti)

	fmt.Println(tasks)
}

func TestUpdateStatus(t *testing.T) {
	id := "1c4f8fd0-0168-4bc1-9d48-bb5ab16c26bf"
	task := new(model.Task)
	err := task.UpdateStatus(id)

	fmt.Println(err)
}
