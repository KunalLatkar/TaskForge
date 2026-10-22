package queue

import (
	"errors"
	"fmt"
	"project/Internal/storage"
	"project/Internal/task"
	"sync"
)

type MemoryQueue struct {
	tasks   map[string]*task.Task
	Pending chan *task.Task // channel that maintaines the pending processes
	mu      sync.RWMutex
	storage storage.Storage
}

func NewMemoryQueue(bufferSize int, store storage.Storage) *MemoryQueue {
	return &MemoryQueue{
		tasks:   make(map[string]*task.Task),
		Pending: make(chan *task.Task, bufferSize),
		storage: store,
	}
}
func (q *MemoryQueue) Enqueue(t *task.Task) error {
	q.mu.Lock()
	q.tasks[t.ID] = t
	q.mu.Unlock()

	if q.storage != nil {
		if err := q.storage.Save(t); err != nil {
			return err
		}
	}

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
	oldTask, _ := q.tasks[t.ID]
	q.mu.Unlock()

	fmt.Println("this is old status: ", oldTask.Status)
	if q.storage != nil {
		if err := q.storage.Move(t, t.Status); err != nil {
			return err
		} else {
			if err := q.storage.Save(t); err != nil {
				return err
			}
		}
	}
	fmt.Println("this is new status: ", t.Status)
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

func (q *MemoryQueue) LoadFromStorage() error {
	if q.storage == nil {
		return nil
	}

	tasks, err := q.storage.LoadAll()
	if err != nil {
		return err
	}

	if len(tasks) == 0 {
		return nil
	}
	q.mu.Lock()
	defer q.mu.Unlock()

	for _, t := range tasks {
		// q.tasks[t.ID] = t

		// re enqueue the incomplete process
		if t.Status == task.PendingStatus {
			q.Pending <- t
		}
		if t.Status == task.ProcessingStatus {
			t.Status = task.PendingStatus
			q.Pending <- t
		}

		// println("inside the func", q)

	}

	return nil
}
