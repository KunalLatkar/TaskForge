# 🔥 TaskForge
### Distributed Task Queue System in Go

**TaskForge** is a lightweight, concurrent task queue system built in Go that processes background jobs using worker pools and channels. Perfect for learning Go concurrency patterns and building production-ready async task processing systems.

---

## 📖 What is TaskForge?

TaskForge solves the common problem of handling long-running tasks in web applications. Instead of making users wait, tasks are queued and processed in the background by concurrent workers.

### Real-World Use Cases:
- 📧 Sending emails asynchronously
- 🖼️ Image/video processing
- 📊 Report generation
- 🔍 Web Search

---

## ✨ Features

- ✅ **Concurrent Workers** - Multiple goroutines process tasks simultaneously
- ✅ **RESTful API** - Submit and track tasks via HTTP
- ✅ **In-Memory Queue** - Fast, buffered channel-based queue
- ✅ **Task Status Tracking** - Monitor task lifecycle (pending → processing → completed/failed)
- ✅ **Real-time Statistics** - View queue metrics instantly
- ✅ **Graceful Shutdown** - Workers finish current tasks before stopping
- ✅ **Thread-Safe** - All operations are goroutine-safe
- ✅ **Extensible** - Easy to add new task types

---

## 🏗️ Architecture

┌─────────────┐         ┌──────────────┐         ┌──────────────┐
│   Client    │───────▶│  HTTP API    │───────▶│    Queue     │
│             │  POST   │              │  Enqueue│  (Channel)   │
└─────────────┘  /task  └──────────────┘         └──────┬───────┘
                                                        │
                                                        │ Dequeue
                                                        ▼
                                                  ┌──────────────┐
                                                  │   Workers    │
                                                  │ (Goroutines) │
                                                  └──────┬───────┘
                                                         │
                                                         ▼
                                                  ┌──────────────┐
                                                  │   Execute    │
                                                  │     Task     │
                                                  └──────────────┘

### Components:

| Component | Purpose | Technology |
|-----------|---------|------------|
| **HTTP API** | Accept task submissions | `net/http` |
| **Queue** | Store pending tasks | Buffered channels |
| **Workers** | Process tasks concurrently | Goroutines |
| **Storage** | Track task state | In-memory map with `sync.RWMutex` |

---

## 🚀 Quick Start

### Prerequisites
- Go 1.21 or higher
- `curl` for testing (or Postman)

### Installation

```bash
# Clone the repository
git clone https://github.com/kunallatkar/taskforge.git
cd taskforge

# Download dependencies
go mod download

# Run the server
go run .

You should see:
🎯 TaskForge - Distributed Task Queue System
================================================
✅ Queue initialized (buffer: 100)
🔧 Starting 3 workers...
🚀 Worker 1 started
🚀 Worker 2 started
🚀 Worker 3 started

🌐 HTTP Server starting on http://localhost:8080

---

📚 API Documentation

1. Submit a Task

Endpoint: POST /task

Request Body:
{
  "type": "email",
  "data": {
    "email": "user@example.com"
  },
  "priority": 1
}

Response:
{
  "task_id": "f47ac10b-58cc-4372-a567-0e02b2c3d479",
  "status": "pending"
}

Example:
curl -X POST http://localhost:8080/task \
  -H "Content-Type: application/json" \
  -d '{"type":"email", "data":{"email":"test@example.com"}, "priority":1}'

---

2. Get Task Status

Endpoint: GET /task/{task_id}

Response:
{
  "id": "f47ac10b-58cc-4372-a567-0e02b2c3d479",
  "type": "email",
  "data": {
    "email": "test@example.com"
  },
  "priority": 1,
  "status": "completed",
  "result": "Email sent to test@example.com",
  "created_at": "2024-10-16T10:00:00Z",
  "started_at": "2024-10-16T10:00:01Z",
  "completed_at": "2024-10-16T10:00:03Z",
  "retry_count": 0
}

Example:
curl http://localhost:8080/task/f47ac10b-58cc-4372-a567-0e02b2c3d479

---

3. Get Queue Statistics

Endpoint: GET /stats

Response:
{
  "pending": 5,
  "processing": 2,
  "completed": 143,
  "failed": 1
}

Example:
curl http://localhost:8080/stats

---

🎯 Supported Task Types

┌──────────────┬─────────────────────────┬──────────┬─────────────────────────────────┐
│  Task Type   │       Description       │ Duration │        Example
├──────────────┼─────────────────────────┼──────────┼─────────────────────────────────┤
│ email        │ Send email notification │ 2s       │ Welcome emails
├──────────────┼─────────────────────────┼──────────┼─────────────────────────────────┤
│ resize_image │ Image processing        │ 4s       │ Profile picture
├──────────────┼─────────────────────────┼──────────┼─────────────────────────────────┤
│ report       │ Generate reports        │ 5s       │ Monthly analyt
├──────────────┼─────────────────────────┼──────────┼─────────────────────────────────┤
│ web Search   │ this is simple task     │ 1s       │ Simple Sreach, clean
└──────────────┴─────────────────────────┴──────────┴─────────────────────────────────┘

Adding Custom Task Types

Edit internal/worker/worker.go:

func (w *Worker) executeTask(t *task.Task) (string, error) {
    switch t.Type {
    case "your_custom_type":
        // Your logic here
        return "Success message", nil

    // ... existing cases
    }
}

---

🧪 Testing

Manual Testing

1. Start the server:
go run .
2. Submit multiple tasks:
for i in {1..10}; do
  curl -X POST http://localhost:8080/task \
    -H "Content-Type: application/json" \
    -d "{\"type\":\"simple\", \"data\":{\"num\":$i}}" &
done
3. Check statistics:
curl http://localhost:8080/stats

Unit Tests

# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run benchmarks
go test -bench=. ./...

---

📁 Project Structure

taskforge/
├── main.go                    # Application entry point
├── go.mod                     # Go module definition
├── go.sum                     # Dependency checksums
├── README.md                  # This file
│
├── internal/
│   ├── task/
│   │   └── task.go           # Task struct and types
│   │
│   ├── queue/
│   │   ├── queue.go          # Queue interface
│   │   └── memory.go         # In-memory queue implementation
│   │
│   ├── worker/
│   │   └── worker.go         # Worker pool and task execution
│   │
│   └── api/
│       └── server.go         # HTTP handlers and routes
│
└── tests/                     # Test files (future)
    ├── queue_test.go
    └──worker_test.go

---

🎓 Learning Outcomes

By building TaskForge, you'll master:

Go Fundamentals

- ✅ Package structure and imports
- ✅ Interfaces and polymorphism
- ✅ Structs and methods
- ✅ Error handling
- ✅ JSON marshaling/unmarshaling

Concurrency

- ✅ Goroutines (lightweight threads)
- ✅ Channels (buffered vs unbuffered)
- ✅ Context for cancellation
- ✅ Mutexes (sync.RWMutex)
- ✅ WaitGroups
- ✅ Deadlock prevention

Web Development

- ✅ HTTP server (net/http)
- ✅ RESTful API design
- ✅ Request handling
- ✅ JSON responses

System Design

- ✅ Queue architectures
- ✅ Worker pool pattern
- ✅ Task lifecycle management
- ✅ Graceful shutdown

---

🔧 Configuration

Environment Variables (Future)

export TASKFORGE_PORT=8080
export TASKFORGE_WORKERS=5
export TASKFORGE_QUEUE_SIZE=100

Worker Pool Size

Edit main.go:
numWorkers := 3  // Change this value

More workers = higher throughput (but more CPU usage)


🎯 Use Cases & Examples

1. Email Service Integration

curl -X POST http://localhost:8080/task \
  -H "Content-Type: application/json" \
  -d '{
    "type": "email",
    "data": {
      "email": "user@example.com",
      "subject": "Welcome!",
      "template": "welcome"
    },
    "priority": 1
  }'

2. Image Processing Pipeline

curl -X POST http://localhost:8080/task \
  -H "Content-Type: application/json" \
  -d '{
    "type": "resize_image",
    "data": {
      "filename": "profile.jpg",
      "width": 200,
      "height": 200
    },
    "priority": 2
  }'

3. Report Generation

curl -X POST http://localhost:8080/task \
  -H "Content-Type: application/json" \
  -d '{
    "type": "report",
    "data": {
      "type": "monthly_sales",
      "month": "October",
      "format": "pdf"
    },
    "priority": 3
  }'

---

🐛 Troubleshooting

Port Already in Use

# Find process using port 8080
lsof -i :8080

# Kill the process
kill -9 <PID>

Workers Not Processing

- Check if numWorkers > 0
- Ensure queue buffer size > 0
- Verify task type is supported

Deadlock Detected

- Make sure channels are properly closed
- Check for goroutine leaks
- Review context.Done() usage

---

🤝 Contributing

Contributions are welcome! Here's how:

1. Fork the repository
2. Create a feature branch (git checkout -b feature/amazing-feature)
3. Commit your changes (git commit -m 'Add amazing feature')
4. Push to the branch (git push origin feature/amazing-feature)
5. Open a Pull Request

---

📝 License

This project is licensed under the MIT License - see the LICENSE file for details.

---