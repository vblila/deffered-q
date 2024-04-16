package queue

import (
	"dq/utils"
	"sync"
	"time"
)

type Queue struct {
	mu            sync.Mutex
	tasks         linkedList
	reservedTasks map[string]*Task
}

func (q *Queue) Init() {
	q.tasks = linkedList{}
	q.reservedTasks = make(map[string]*Task, 1024)
}

func (q *Queue) Add(taskBody []byte, delayMs uint32) (string, error) {
	taskId := utils.RandomId()

	delayedTime := time.Now()
	if delayMs > 0 {
		delayedTime = delayedTime.Add(time.Millisecond * time.Duration(delayMs))
	}

	q.mu.Lock()
	q.tasks.Push(&Task{Id: taskId, Body: taskBody, DelayedTime: delayedTime})
	q.mu.Unlock()

	return taskId, nil
}

func (q *Queue) Reserve() (taskId string, taskBody []byte, stuckAttempts uint8, ok bool) {
	q.mu.Lock()
	defer q.mu.Unlock()

	var lastNode *linkedListNode

	for {
		// FIFO очередь тасков
		node := q.tasks.Next(lastNode)
		if node == nil {
			break
		}

		task := node.Value

		if task.DelayedTime.Before(time.Now()) {
			q.reservedTasks[task.Id] = task
			q.tasks.Delete(node)

			return task.Id, task.Body, task.StuckAttempts, true
		}
		lastNode = node
	}

	ok = false
	return
}

func (q *Queue) Return(taskId string, delayMs uint32, isStuckAttempt bool) bool {
	q.mu.Lock()
	defer q.mu.Unlock()

	task, ok := q.reservedTasks[taskId]
	if !ok {
		return false
	}

	task.DelayedTime = time.Now()
	if delayMs > 0 {
		task.DelayedTime = task.DelayedTime.Add(time.Duration(delayMs) * time.Millisecond)
	}

	if isStuckAttempt {
		task.StuckAttempts++
	}

	q.tasks.Push(task)

	delete(q.reservedTasks, task.Id)

	return true
}

func (q *Queue) Delete(taskId string) bool {
	q.mu.Lock()

	_, ok := q.reservedTasks[taskId]
	if ok {
		delete(q.reservedTasks, taskId)
	}

	q.mu.Unlock()
	return ok
}

func (q *Queue) TasksLength() int {
	q.mu.Lock()
	tasks := int(q.tasks.Length())
	q.mu.Unlock()

	return tasks
}

func (q *Queue) ReservedTasksLength() int {
	q.mu.Lock()
	reservedTasks := len(q.reservedTasks)
	q.mu.Unlock()

	return reservedTasks
}
