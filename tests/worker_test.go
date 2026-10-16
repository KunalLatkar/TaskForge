package main

import (
	"context"
	"fmt"
	"project/Internal/queue"
	"project/Internal/task"
	"project/Internal/worker"
	"testing"
)

func BenchmarkWorkerThoughput(b *testing.B) {
	q := queue.NewMemoryQueue(200)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for i := 0; i < 10; i++ {
		w := worker.NewWorker(i, q)
		go w.Start(ctx)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		task := &task.Task{
			ID:   fmt.Sprintf("task-%d", i),
			Type: "web search",
		}
		q.Enqueue(task)
	}
}

// command to run this benchmark test
// go test tests/worker_test.go -bench .
