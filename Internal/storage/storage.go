package storage

import "project/Internal/task"

type Storage interface {
	Save(t *task.Task) error

	Move(t *task.Task, newStatus task.Status) error

	Load(id string) (*task.Task, error)

	LoadAll() ([]*task.Task, error)

	Delete(id string) error
}
