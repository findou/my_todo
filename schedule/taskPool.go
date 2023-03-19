/*
 * @Version: 1.0
 * @Date: 2023-03-18 13:19:26
 * @LastEditTime: 2023-03-18 14:25:24
 */
package schedule

import "sync"

type Task interface {
	Execute(num int)
}
type TaskPool struct {
	Capacity int
	TaskChan chan Task
	wg       *sync.WaitGroup
}

func NewTaskPool(capacity int) *TaskPool {
	return &TaskPool{
		Capacity: capacity,
		TaskChan: make(chan Task, 10),
		wg:       &sync.WaitGroup{},
	}
}

func (p *TaskPool) Put(task Task) {
	p.wg.Add(1)
	p.TaskChan <- task
}

func (p *TaskPool) Wait() {
	for i := 0; i < p.Capacity; i++ {
		// 开Capacity个协程并发执行任务
		go func(num int) {
			for task := range p.TaskChan {
				task.Execute(num)
				p.wg.Done()
			}
		}(i)
	}
	p.wg.Wait()
}
