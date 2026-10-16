package task

import "time"

type Status string

const (
	PendingStatus    Status = "Pending"
	ProcessingStatus Status = "Processing"
	CompletedStatus  Status = "Completed"
	FailedStatus     Status = "Failed"
)

type Task struct {
	ID         string                 `json: "id"`
	Type       string                 `json: "type"`
	Created    time.Time              `json: "created on"`
	RetryCount int                    `json: "RetryCount"`
	Data       map[string]interface{} `json: "data"`
	StartedOn  *time.Time             `json: "Started On,omitempty"`
	EndedOn    *time.Time             `json: "Ended On,omitempty"`
	Status     Status                 `json: "Status"`
	Error      string                 `json: "Error, omitempty"`
	Result     string                 `json: "result,omitempty"`
	Priority   int                    `json: "Priority"`
}
