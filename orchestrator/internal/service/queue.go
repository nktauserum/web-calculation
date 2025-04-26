package service

import (
	"sync"

	"github.com/nktauserum/web-calculation/orchestrator/pkg/task"
)

var queue task.Queue
var once sync.Once

func GetQueue() *task.Queue {
	var err error
	once.Do(func() {
		q, err_in := task.NewQueue("sqlite.db")
		if err_in != nil {
			err = err_in
		}

		queue = *q
	})

	if err != nil {
		return nil
	}

	return &queue
}
