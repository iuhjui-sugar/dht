package dht

import (
	"sync"
)

type TaskQueue struct {
	size int
	mu   sync.Mutex
}

func NewTaskQueue(size int) *TaskQueue {
	tq := new(TaskQueue)
	tq.size = size
	return tq
}

func (tq *TaskQueue) ExecGo(fn func(index int)) {
	var wg sync.WaitGroup
	for i := 0; i < tq.size; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			fn(index)
		}(i)
	}
	wg.Wait()
}
