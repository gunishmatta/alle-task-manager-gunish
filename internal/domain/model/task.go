package model

import (
	"github.com/google/uuid"
	"time"
)

type TaskStatus string

const (
	Pending    TaskStatus = "pending"
	InProgress TaskStatus = "in_progress"
	Completed  TaskStatus = "completed"
)

type Task struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      TaskStatus `json:"status"`
	DueDate     *time.Time `json:"due_date,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func NewTask(title, description string) *Task {
	now := time.Now()
	return &Task{
		ID:          uuid.New().String(),
		Title:       title,
		Description: description,
		Status:      Pending,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}
