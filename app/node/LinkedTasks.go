package node

type TaskNode struct {
	prev *TaskNode
	next *TaskNode
	Task *Task
}

type LinkedTasks struct {
	head   *TaskNode
	tail   *TaskNode
	length int
}

func (l *LinkedTasks) Append(task *Task) {
	node := &TaskNode{prev: nil, next: nil, Task: task}
	if l.tail == nil {
		l.tail = node
		l.head = l.tail
	} else {
		l.tail.next = node
		node.prev = l.tail
		l.tail = node
	}

	l.length++
}

func (l *LinkedTasks) Next(curr *TaskNode) *TaskNode {
	if curr != nil {
		return curr.next
	} else {
		return l.head
	}
}

func (l *LinkedTasks) Pull(taskNode *TaskNode) {
	if taskNode.Task.Id == l.head.Task.Id {
		l.head = taskNode.next
	}
	if taskNode.Task.Id == l.tail.Task.Id {
		l.tail = taskNode.prev
	}

	if taskNode.prev != nil {
		taskNode.prev.next = taskNode.next
	}

	l.length--
}

func (l *LinkedTasks) Length() int {
	return l.length
}
