package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"project/Internal/api"
	"project/Internal/queue"
	"project/Internal/worker"
	"syscall"
	"time"
)

func main() {
	q := queue.NewMemoryQueue(100)

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
