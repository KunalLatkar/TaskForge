package queue

import (
	"errors"
	"project/Internal/task"
	"sync"
)

type MemoryQueue struct {
	tasks   map[string]*task.Task
	Pending chan *task.Task // channel that maintaines the pending processes
	mu      sync.RWMutex
}

func NewMemoryQueue(bufferSize int) *MemoryQueue {
	return &MemoryQueue{
		tasks:   make(map[string]*task.Task),
		Pending: make(chan *task.Task, bufferSize),
	}
}
func (q *MemoryQueue) Enqueue(t *task.Task) error {
	q.mu.Lock()
	q.tasks[t.ID] = t
	q.mu.Unlock()

	q.Pending <- t
	return nil
}

func (q *MemoryQueue) Dequeue() (*task.Task, error) {
	t := <-q.Pending // to return the dequeued process
	return t, nil
}

func (q *MemoryQueue) Get(id string) (*task.Task, error) {
	q.mu.RLock()
	t, exits := q.tasks[id]
	q.mu.RUnlock()
	if !exits {
		return nil, errors.New("Task not found")
	}
	return t, nil
}

func (q *MemoryQueue) Update(t *task.Task) error {
	q.mu.Lock()
	q.tasks[t.ID] = t
	q.mu.Unlock()

	return nil
}

func (q *MemoryQueue) GetStats() QueueStates {
	q.mu.RLock()
	defer q.mu.RUnlock()

	stats := QueueStates{
		Pending:    0,
		Completed:  0,
		Processing: 0,
		Failed:     0,
	}

	for _, t := range q.tasks {
		switch t.Status {
		case task.PendingStatus:
			stats.Pending++
		case task.CompletedStatus:
			stats.Completed++
		case task.FailedStatus:
			stats.Failed++
		case task.ProcessingStatus:
			stats.Processing++
		}
	}

	return stats
}
