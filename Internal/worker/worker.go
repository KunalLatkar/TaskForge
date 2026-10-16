package worker

import (
	"context"
	"fmt"
	"project/Internal/queue"
	"project/Internal/task"
	"time"
)

type Worker struct {
	id    int
	queue queue.Queue
}

func NewWorker(id int, q queue.Queue) *Worker {
	return &Worker{
		id:    id,
		queue: q,
	}
}

func (w *Worker) Start(ctx context.Context) {
	fmt.Printf("Worker %d started\n", w.id)

	for {
		select {
		case <-ctx.Done():
			fmt.Printf("Worker %d stoppping\n", w.id)
			return
		default:
			t, err := w.queue.Dequeue()
			// fmt.Print("this is in start", t.StartedOn)
			if err != nil {
				continue
			}

			w.processTask(t)
		}
	}
}

func (w *Worker) processTask(t *task.Task) {
	fmt.Printf("Worker %d  processing task %s\n", w.id, t.ID)

	now := time.Now()

	t.StartedOn = &now
	t.Status = task.ProcessingStatus
	w.queue.Update(t)

	result, err := w.executeTask(t)

	completed := time.Now()
	t.EndedOn = &completed

	if err != nil {
		t.Status = task.FailedStatus
		t.Error = err.Error()
	} else {
		t.Status = task.CompletedStatus
		t.Result = result
	}

	w.queue.Update(t)
	fmt.Printf("Worker %d completed task %s\n", w.id, t.ID)
	fmt.Println(result)

}

func (w *Worker) executeTask(t *task.Task) (string, error) {
	// fmt.Printf("this is in execute", t.Type)
	switch t.Type {
	case "email":
		time.Sleep(2 * time.Second)
		return "Email sent:", nil
	case "resize_image":
		time.Sleep(4 * time.Second)
		return "Successfully resized the photo", nil
	case "report":
		time.Sleep(5 * time.Second)
		return "Repoted Successfully", nil
	case "web Search":
		time.Sleep(1 * time.Second)
		return "Google pe ja....!", nil
	default:
		return "", fmt.Errorf("Unkown operation: %s", t.Type)
	}

}
