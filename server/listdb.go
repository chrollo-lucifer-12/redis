package server

import (
	"sync"

	"github.com/chrollo-lucifer-12/redis/datastructures"
)

type listDB struct {
	//listMap map[string]*datastructures.LinkedList
	listMap sync.Map
}

func newListDb() *listDB {
	return &listDB{}
}

func (l *listDB) LPUSH(key string, elements []string) int {
	list, ok := l.listMap.Load(key)
	if !ok {
		newList := datastructures.NewLinkedList()
		l.listMap.Store(key, newList)
		list = newList
	}

	l1 := list.(*datastructures.LinkedList)

	return l1.InsertAtHead(elements)
}
