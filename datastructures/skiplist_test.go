package datastructures

import "testing"

// func TestSkipListInsert(t *testing.T) {
// 	list := newSkipList()
// 	list.InsertNode("apple", 5)
// 	list.InsertNode("banana", 3)
// 	list.InsertNode("cherry", 5)
// 	list.InsertNode("date", 7)

// 	list.PrintList()
// }

func TestSkipListZRangePrint(t *testing.T) {
	list := newSkipList()
	list.InsertNode("apple", 5)
	list.InsertNode("banana", 3)
	list.InsertNode("cherry", 5)
	list.InsertNode("date", 7)

	tests := []struct {
		start, end int
	}{
		{0, 1},
		{1, 2},
		{0, 3},
		{-2, -1},
	}

	for _, tt := range tests {
		result := list.SearchInRange(tt.start, tt.end)
		t.Logf("ZRange(%d, %d): %v", tt.start, tt.end, result)
	}
}
