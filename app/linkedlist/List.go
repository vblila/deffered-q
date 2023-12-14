package linkedlist

type List struct {
	head   *Node
	tail   *Node
	length uint32
}

func (l *List) Length() uint32 {
	return l.length
}

func (l *List) Push(value interface{}) {
	node := &Node{prev: nil, next: nil, Value: value}
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

func (l *List) Next(curr *Node) *Node {
	if curr != nil {
		return curr.next
	} else {
		return l.head
	}
}

func (l *List) Delete(node *Node) {
	if node == l.head {
		l.head = node.next
		if l.head != nil {
			l.head.prev = nil
		}
	}
	if node == l.tail {
		l.tail = node.prev
		if l.tail != nil {
			l.tail.next = nil
		}
	}

	if node.prev != nil {
		node.prev.next = node.next
		node.next.prev = node.prev

		node.prev = nil
		node.next = nil
	}

	l.length--
}
