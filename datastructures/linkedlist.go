package datastructures

type node struct {
	value string
	next  *node
	prev  *node
}

func newNode(value string, next *node, prev *node) *node {
	return &node{
		value: value,
		next:  next,
		prev:  prev,
	}
}

type LinkedList struct {
	len  int
	head *node
	tail *node
}

func NewLinkedList() *LinkedList {
	return &LinkedList{
		len:  0,
		head: nil,
		tail: nil,
	}
}

func (l *LinkedList) RemoveFromHead(count int) []string {
	var res []string
	if l.len == 0 {
		return res
	}
	for count > 0 && l.len > 0 {
		oldHead := l.head
		res = append(res, oldHead.value)
		l.head = oldHead.next
		oldHead.next = nil
		l.len--
		count--
	}
	return res
}

func (l *LinkedList) RemoveFromTail(count int) []string {
	var res []string
	if l.len == 0 {
		return res
	}
	for count > 0 && l.len > 0 {
		oldTail := l.tail
		res = append(res, oldTail.value)
		l.tail = oldTail.prev
		oldTail.prev = nil
		l.len--
		count--
	}
	return res
}

func (l *LinkedList) InsertAtHead(elements []string) int {
	for _, element := range elements {
		nodeToInsert := newNode(element, nil, nil)
		if l.head == nil {
			l.head = nodeToInsert
			l.tail = nodeToInsert
		} else {
			oldHead := l.head
			nodeToInsert.next = oldHead
			oldHead.prev = nodeToInsert
			l.head = nodeToInsert
		}
		l.len++
	}
	return l.Length()
}

func (l *LinkedList) InsertAtTail(elements []string) int {
	for _, element := range elements {
		nodeToInsert := newNode(element, nil, nil)
		if l.head == nil {
			l.head = nodeToInsert
			l.tail = nodeToInsert
		} else {
			oldTail := l.tail
			oldTail.next = nodeToInsert
			nodeToInsert.prev = oldTail
			l.tail = nodeToInsert
		}
		l.len++
	}
	return l.Length()
}

func (l *LinkedList) GetElementsInRange(start, stop int) []string {
	length := l.len
	var res []string

	if start < 0 {
		start = length + start
	}
	if stop < 0 {
		stop = length + stop
	}

	if start < 0 {
		start = 0
	}
	if stop >= length {
		stop = length - 1
	}

	if start >= length || start > stop {
		return res
	}

	idx := 0
	temp := l.head
	for idx < start {
		temp = temp.next
		idx++
	}
	for idx <= stop && temp != nil {
		res = append(res, temp.value)
		temp = temp.next
		idx++
	}

	return res
}

func (l *LinkedList) Length() int {
	return l.len
}
