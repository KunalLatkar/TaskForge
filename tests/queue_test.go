package main

import (
	"project/Internal/queue"
	"project/Internal/task"
	"testing"
)

func TestQeueEnqueue(t *testing.T) {
	q := queue.NewMemoryQueue(4, nil)

	task := &task.Task{
		ID:   "343-testing",
		Type: "email",
	}

	err := q.Enqueue(task)
	if err != nil {
		t.Fatalf("Enqueue function failed: %v", err)
	}

	retrieve, err := q.Dequeue()
	if err != nil {
		t.Fatalf("Dequeue function failed: %v", err)
	}

	if retrieve.ID != task.ID {
		t.Errorf("Expected: %v, Got: %v\n", task.ID, retrieve.ID)
	}
}

// command to run this test
// go test tests/queue_test.go
