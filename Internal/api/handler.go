package api

import (
	"encoding/json"
	"net/http"
	"project/Internal/queue"
	"project/Internal/task"
	"time"

	"github.com/google/uuid"
)

// would need the uuid package use following command to get
// go get github.com/google/uuid

type Server struct {
	queue queue.Queue
}

func NewServer(q queue.Queue) *Server {
	return &Server{queue: q}
}

// POST Request
func (s *Server) SubmitTask(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Type     string                 `json: "type`
		Data     map[string]interface{} `json: "data"`
		Priority int                    `json: "priority"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid json", http.StatusBadRequest)
		return
	}

	t := &task.Task{
		ID:       uuid.New().String(),
		Type:     req.Type,
		Data:     req.Data,
		Priority: req.Priority,
		Status:   task.PendingStatus,
		Created:  time.Now(),
	}

	// fmt.Println("GOT the json")

	if err := s.queue.Enqueue(t); err != nil {
		http.Error(w, "Failed to queue the task", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"task_id": t.ID,
		"status":  string(t.Status),
	})
}

// GET Request
func (s *Server) GetTask(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/task/"):]

	t, err := s.queue.Get(id)
	if err != nil {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(t)
}

// GET Statistics
func (s *Server) GetStats(w http.ResponseWriter, r *http.Request) {
	stats := s.queue.GetStats()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}
