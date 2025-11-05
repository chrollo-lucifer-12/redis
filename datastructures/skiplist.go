package datastructures

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

const MAX_LEVEL = 6

type skipListNode struct {
	score   int
	key     string
	forward [MAX_LEVEL]*skipListNode
	span    [MAX_LEVEL]int
}

type SkipList struct {
	header *skipListNode
	level  int
}

func newSkipListNode(key string, score int, level int) *skipListNode {
	node := skipListNode{key: key, score: score}
	for i := 0; i < level; i++ {
		node.forward[i] = nil
	}
	return &node
}

func NewSkipList() *SkipList {
	list := SkipList{}
	list.header = newSkipListNode("", math.MinInt, MAX_LEVEL)
	list.level = 0
	return &list
}

func randomLevel() int {
	rand.Seed(time.Now().UnixNano())
	level := 0
	for rand.Intn(2) == 0 && level < MAX_LEVEL-1 {
		level++
	}
	return level
}

func (s *SkipList) InsertNode(key string, score int) {
	current := s.header
	update := make([]*skipListNode, MAX_LEVEL)
	rank := make([]int, MAX_LEVEL)
	for i := 0; i < MAX_LEVEL; i++ {
		update[i] = s.header
		rank[i] = 0
	}

	for i := s.level; i >= 0; i-- {
		if i < s.level {
			rank[i] = rank[i+1]
		}
		for current.forward[i] != nil && (current.forward[i].score < score || (current.forward[i].score == score && current.forward[i].key < key)) {
			current = current.forward[i]
			rank[i] += current.span[i]
		}
		update[i] = current
	}

	rlevel := randomLevel()
	if rlevel > s.level {
		for i := s.level + 1; i <= rlevel; i++ {
			update[i] = s.header
			rank[i] = 0
		}
		s.level = rlevel
	}

	newNode := newSkipListNode(key, score, rlevel+1)
	for i := 0; i <= rlevel; i++ {
		newNode.forward[i] = update[i].forward[i]
		update[i].forward[i] = newNode
		newNode.span[i] = update[i].span[i] - (rank[0] - rank[i])
		update[i].span[i] = (rank[0] - rank[i]) + 1
	}

	for i := rlevel + 1; i <= s.level; i++ {
		update[i].span[i]++
	}
	s.PrintList()
}

func (s *SkipList) SearchInRange(start, stop int) []string {
	result := []string{}
	current := s.header
	rank := -1

	for i := s.level; i >= 0; i-- {
		for current.forward[i] != nil && rank+current.span[i] < start {
			rank += current.span[i]
			current = current.forward[i]
		}
	}

	current = current.forward[0]
	rank++

	for current != nil && rank <= stop {
		result = append(result, current.key)
		current = current.forward[0]
		rank++
	}

	return result
}

func (s *SkipList) PrintList() {
	fmt.Println("Skip List:")
	for i := s.level; i >= 0; i-- {
		current := s.header.forward[i]
		fmt.Printf("Level %d: ", i)
		for current != nil {
			fmt.Printf("(%s,%d) ", current.key, current.score)
			current = current.forward[i]
		}
		fmt.Println()
	}
}
