package node

import (
	"dq/linkedlist"
	"dq/trie"
	"dq/utils"
	"log"
	"sync"
	"time"
)

const taskIdLength = 8

type Queue struct {
	mu            sync.Mutex
	tasks         linkedlist.List
	reservedTasks trie.Trie
}

func (q *Queue) Init() {
	q.tasks = linkedlist.List{}
	q.reservedTasks = trie.Trie{}
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
	q.tasks.Push(&Task{Id: taskId, Body: taskBody, DelayedTime: delayedTime})
	defer q.mu.Unlock()

	return taskId, nil
}

func (q *Queue) Reserve() (taskId string, taskBody []byte, stuckAttempts uint8, ok bool) {
	q.mu.Lock()
	defer q.mu.Unlock()

	var lastNode *linkedlist.Node

	for {
		// FIFO очередь тасков
		node := q.tasks.Next(lastNode)
		if node == nil {
			break
		}

		task := node.Value.(*Task)

		if task.DelayedTime.Before(time.Now()) {
			q.reservedTasks.Put([]rune(task.Id), task)
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

	rawValue := q.reservedTasks.Get([]rune(taskId))
	if rawValue == nil {
		return false
	}

	task := rawValue.(*Task)

	task.DelayedTime = time.Now()
	if delayMs > 0 {
		task.DelayedTime = task.DelayedTime.Add(time.Duration(delayMs) * time.Millisecond)
	}

	if isStuckAttempt {
		task.StuckAttempts++
	}

	q.tasks.Push(task)

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

	return int(q.tasks.Length())
}

func (q *Queue) ReservedTasksLength() int {
	q.mu.Lock()
	defer q.mu.Unlock()

	return int(q.reservedTasks.Length())
}
