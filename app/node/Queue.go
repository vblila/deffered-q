package node

import (
	"dq/utils"
	"log"
	"sync"
	"time"
)

const taskIdLength = 8

type Queue struct {
	mu            sync.Mutex
	tasks         *LinkedTasks
	reservedTasks *PTMap
}

func (q *Queue) Init() {
	q.tasks = &LinkedTasks{}
	q.reservedTasks = &PTMap{}
}

func (q *Queue) Add(taskBody []byte, delayMs uint32) (string, error) {
	taskId, err := utils.RandomString(taskIdLength)
	if err != nil {
		log.Println("Task id generation failed", err.Error())
		return "", err
	}

	delayedTime := time.Now()
	if delayMs > 0 {
		delayedTime = delayedTime.Add(time.Millisecond * time.Duration(delayMs))
	}

	q.mu.Lock()
	if q.tasks == nil {
		q.tasks = &LinkedTasks{}
	}
	q.tasks.Append(&Task{Id: taskId, Body: taskBody, DelayedTime: delayedTime})
	defer q.mu.Unlock()

	return taskId, nil
}

func (q *Queue) Reserve() *Task {
	q.mu.Lock()
	defer q.mu.Unlock()

	var lastTaskNode *TaskNode

	for {
		// FIFO очередь тасков
		taskNode := q.tasks.Next(lastTaskNode)
		if taskNode == nil {
			break
		}

		if taskNode.Task.DelayedTime.Before(time.Now()) {
			task := taskNode.Task
			q.reservedTasks.Put([]rune(task.Id), task)
			q.tasks.Pull(taskNode)

			return task
		}
		lastTaskNode = taskNode
	}

	return nil
}

func (q *Queue) Return(taskId string, delayMs uint32) bool {
	q.mu.Lock()
	defer q.mu.Unlock()

	rawValue := q.reservedTasks.Get([]rune(taskId))
	if rawValue == nil {
		return false
	}

	task := rawValue.(*Task)

	task.DelayedTime = time.Now()
	if delayMs > 0 {
		task.DelayedTime = task.DelayedTime.Add(time.Duration(delayMs) * time.Millisecond)
	}

	q.tasks.Append(task)

	q.reservedTasks.Delete([]rune(task.Id))

	return true
}

func (q *Queue) Delete(taskId string) bool {
	q.mu.Lock()
	defer q.mu.Unlock()

	return q.reservedTasks.Delete([]rune(taskId))
}

func (q *Queue) TasksLength() int {
	q.mu.Lock()
	defer q.mu.Unlock()

	return q.tasks.Length()
}

func (q *Queue) ReservedTasksLength() int {
	q.mu.Lock()
	defer q.mu.Unlock()

	return int(q.reservedTasks.Length())
}
