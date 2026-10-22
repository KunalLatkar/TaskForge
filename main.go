package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"project/Internal/api"
	"project/Internal/queue"
	"project/Internal/storage"
	"project/Internal/worker"
	"syscall"
	"time"
)

func main() {
	fmt.Println("TaskForge - Distributed Task Queue System")

	store, err := storage.NewFileStorage("./data")

	if err != nil {
		fmt.Printf("Failed to initialize storage: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("File storage initialized storage ./data")

	q := queue.NewMemoryQueue(100, store)
	fmt.Println("Queue is initialized (buffer: 100) ")

	if err := q.LoadFromStorage(); err != nil {
		fmt.Printf("Warning failed to load the tasks : %v\n", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	numWorkers := 5

	for i := 1; i <= numWorkers; i++ {
		w := worker.NewWorker(i, q)
		go w.Start(ctx)
	}

	server := api.NewServer(q)

	http.HandleFunc("POST /task", server.SubmitTask)
	http.HandleFunc("GET /task", server.GetTask)
	http.HandleFunc("GET /stats", server.GetStats)

	go func() {
		fmt.Println("Server starting on: 8080")
		if err := http.ListenAndServe(":8080", nil); err != nil {
			fmt.Printf("Server errror : %v\n", err)
		}
	}()

	signChan := make(chan os.Signal, 1)
	signal.Notify(signChan, os.Interrupt, syscall.SIGTERM)
	<-signChan

	fmt.Println("\nShutting down...")
	cancel() //to stop the workers
	time.Sleep(2 * time.Second)
	fmt.Println("Shutdown Complete")
}
