package service

import (
	"sync"

	"github.com/nktauserum/web-calculation/orchestrator/pkg/task"
	"github.com/nktauserum/web-calculation/shared"
)

var queue task.Queue
var once sync.Once

func GetQueue() *task.Queue {
	once.Do(func() {
		queue = task.Queue{Tasks: make(map[int64]shared.Task), Expressions: map[int64]shared.Expression{}}
	})
	return &queue
}
