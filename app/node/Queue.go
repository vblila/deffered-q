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
	tasks         *linkedlist.List
	reservedTasks *trie.Trie
}

func (q *Queue) Init() {
	q.tasks = &linkedlist.List{}
	q.reservedTasks = &trie.Trie{}
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
		q.tasks = &linkedlist.List{}
	}
	q.tasks.Push(&Task{Id: taskId, Body: taskBody, DelayedTime: delayedTime})
	defer q.mu.Unlock()

	return taskId, nil
}

func (q *Queue) Reserve() *Task {
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

			return task
		}
		lastNode = node
	}

	return nil
}

func (q *Queue) Return(taskId string, delayMs uint32) bool {
	q.mu.Lock()
	defer q.mu.Unlock()

	trieNode := q.reservedTasks.Get([]rune(taskId))
	if trieNode == nil {
		return false
	}

	task := trieNode.Value.(*Task)

	task.DelayedTime = time.Now()
	if delayMs > 0 {
		task.DelayedTime = task.DelayedTime.Add(time.Duration(delayMs) * time.Millisecond)
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
