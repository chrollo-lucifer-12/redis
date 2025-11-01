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

func (l *listDB) RPUSH(key string, elements []string) int {
	list, ok := l.listMap.Load(key)
	if !ok {
		newList := datastructures.NewLinkedList()
		l.listMap.Store(key, newList)
		list = newList
	}

	l1 := list.(*datastructures.LinkedList)

	return l1.InsertAtTail(elements)
}

func (l *listDB) LPOP(key string, count int) []string {
	list, ok := l.listMap.Load(key)
	var res []string
	if !ok {
		return res
	}
	l1 := list.(*datastructures.LinkedList)
	res = l1.RemoveFromHead(count)
	return res
}

func (l *listDB) RPOP(key string, count int) []string {
	list, ok := l.listMap.Load(key)
	var res []string
	if !ok {
		return res
	}
	l1 := list.(*datastructures.LinkedList)
	res = l1.RemoveFromTail(count)
	return res
}
