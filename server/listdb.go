package server

import "github.com/chrollo-lucifer-12/redis/datastructures"

type listDB struct {
	listMap map[string]*datastructures.LinkedList
}

func newListDb() *listDB {
	return &listDB{
		listMap: make(map[string]*datastructures.LinkedList),
	}
}

func (l *listDB) LPUSH(key string, elements []string) int {
	_, ok := l.listMap[key]
	if !ok {
		l.listMap[key] = datastructures.NewLinkedList()
	}
	return l.listMap[key].InsertAtHead(elements)
}
