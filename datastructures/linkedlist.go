package datastructures

type node struct {
	value string
	next  *node
}

func newNode(value string, next *node) *node {
	return &node{
		value: value,
		next:  next,
	}
}

type LinkedList struct {
	len  int
	head *node
}

func NewLinkedList() *LinkedList {
	return &LinkedList{
		len:  0,
		head: nil,
	}
}

func (l *LinkedList) InsertAtHead(elements []string) int {
	for _, element := range elements {
		nodeToInsert := newNode(element, nil)
		if l.head == nil {
			l.head = nodeToInsert
		} else {
			nodeToInsert.next = l.head
			l.head = nodeToInsert
		}
		l.len++
	}
	return l.Length()
}

func (l *LinkedList) Length() int {
	return l.len
}
