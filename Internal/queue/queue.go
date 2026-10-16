package queue

import "project/Internal/task"

type Queue interface {
	Enqueue(t *task.Task) error
	Dequeue() (*task.Task, error)
	Get(id string) (*task.Task, error)
	Update(t *task.Task) error
	GetStats() QueueStates
}

type QueueStates struct {
	Pending    int
	Processing int
	Completed  int
	Failed     int
}
