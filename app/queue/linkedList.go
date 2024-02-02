package queue

type linkedListNode struct {
	prev  *linkedListNode
	next  *linkedListNode
	Value *Task
}

type linkedList struct {
	head   *linkedListNode
	tail   *linkedListNode
	length uint32
}

func (l *linkedList) Length() uint32 {
	return l.length
}

func (l *linkedList) Push(value *Task) {
	node := &linkedListNode{prev: nil, next: nil, Value: value}
	if l.tail == nil {
		l.tail = node
		l.head = node
	} else {
		l.tail.next = node
		node.prev = l.tail
		l.tail = node
	}

	l.length++
}

func (l *linkedList) Next(curr *linkedListNode) *linkedListNode {
	if curr != nil {
		return curr.next
	} else {
		return l.head
	}
}

func (l *linkedList) Delete(node *linkedListNode) {
	if node == l.head {
		l.head = node.next
	}

	if node == l.tail {
		l.tail = node.prev
	}

	if node.prev != nil {
		node.prev.next = node.next
	}

	if node.next != nil {
		node.next.prev = node.prev
	}

	node.prev = nil
	node.next = nil

	l.length--
}
