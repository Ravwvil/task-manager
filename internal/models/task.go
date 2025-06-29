package models

import (
	"time"

	"github.com/google/uuid"
)

type TaskStatus string

const (
	TaskStatusPending   TaskStatus = "pending"
	TaskStatusRunning   TaskStatus = "running"
	TaskStatusCompleted TaskStatus = "completed"
	TaskStatusFailed    TaskStatus = "failed"
	TaskStatusCancelled TaskStatus = "cancelled"
)

type Task struct {
	ID          string     `json:"id"`
	Description string     `json:"description,omitempty"`
	Status      TaskStatus `json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	StartedAt   *time.Time `json:"started_at,omitempty"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
	Duration    int64      `json:"duration"` 
	Result      string     `json:"result,omitempty"`
	Error       string     `json:"error,omitempty"`
}

func NewTask(description string) *Task {
	return &Task{
		ID:          uuid.New().String(),
		Description: description,
		Status:      TaskStatusPending,
		CreatedAt:   time.Now(),
		Duration:    0,
	}
}

func (t *Task) Start() {
	now := time.Now()
	t.Status = TaskStatusRunning
	t.StartedAt = &now
}

func (t *Task) Complete(result string) {
	now := time.Now()
	t.Status = TaskStatusCompleted
	t.CompletedAt = &now
	t.Result = result
	if t.StartedAt != nil {
		t.Duration = int64(now.Sub(*t.StartedAt).Seconds())
	}
}

func (t *Task) Fail(err string) {
	now := time.Now()
	t.Status = TaskStatusFailed
	t.CompletedAt = &now
	t.Error = err
	if t.StartedAt != nil {
		t.Duration = int64(now.Sub(*t.StartedAt).Seconds())
	}
}

func (t *Task) Cancel(reason string) {
	now := time.Now()
	t.Status = TaskStatusCancelled
	t.CompletedAt = &now
	t.Error = reason
	if t.StartedAt != nil {
		t.Duration = int64(now.Sub(*t.StartedAt).Seconds())
	}
}

type CreateTaskRequest struct {
	Description string `json:"description" validate:"required"`
}

type TaskResponse struct {
	*Task
}

type ErrorResponse struct {
	Error string `json:"error"`
}
