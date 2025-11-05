package server

import (
	"sync"

	"github.com/chrollo-lucifer-12/redis/datastructures"
)

type sortedSetsDB struct {
	sortedSetsMap sync.Map
}

func newSortedSetsDb() *sortedSetsDB {
	return &sortedSetsDB{}
}

func (s *sortedSetsDB) ZADD(key string, member string, score int) {
	set, ok := s.sortedSetsMap.Load(key)
	if !ok {
		newSet := datastructures.NewSkipList()
		s.sortedSetsMap.Store(key, newSet)
		set = newSet
	}
	s1 := set.(*datastructures.SkipList)
	s1.InsertNode(member, score)
}
