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
	ID          string     `json:"id" gorm:"primaryKey"`
	Title       string     `json:"title" gorm:"not null"`
	Description string     `json:"description"`
	Status      TaskStatus `json:"status" gorm:"not null"`
	DueDate     *time.Time `json:"due_date,omitempty"`
	CreatedAt   time.Time  `json:"created_at" gorm:"not null"`
	UpdatedAt   time.Time  `json:"updated_at" gorm:"not null"`
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

func (Task) TableName() string {
	return "tasks"
}
