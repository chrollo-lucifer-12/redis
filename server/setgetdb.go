package server

import (
	"sync"
	"time"
)

type data struct {
	value  string
	expiry time.Time
}

type database struct {
	mp sync.Map
}

func newDB() *database {
	return &database{
		mp: sync.Map{},
	}
}
