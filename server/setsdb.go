package server

import (
	"sync"

	"github.com/chrollo-lucifer-12/redis/datastructures"
)

type setsDB struct {
	setsMap sync.Map
}

func newSetsDB() *setsDB {
	return &setsDB{}
}

func (s *setsDB) SADD(key string, elements []string) int {
	set, ok := s.setsMap.Load(key)
	if !ok {
		newSet := datastructures.NewHashSet()
		s.setsMap.Store(key, newSet)
		set = newSet
	}
	s1 := set.(*datastructures.HashSet)
	for _, element := range elements {
		s1.Insert(element)
	}
	return s1.Size()
}

func (s *setsDB) SREM(key string, element string) int {
	set, ok := s.setsMap.Load(key)
	if !ok {
		return 0
	}
	s1 := set.(*datastructures.HashSet)
	if s1.Erase(element) {
		return 1
	} else {
		return 0
	}
}

func (s *setsDB) SISMEMBER(key string, element string) int {
	set, ok := s.setsMap.Load(key)
	if !ok {
		return 0
	}
	s1 := set.(*datastructures.HashSet)
	if s1.Contains(element) {
		return 1
	} else {
		return 0
	}
}
