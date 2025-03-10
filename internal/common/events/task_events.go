package events

import (
	"time"
)

type TaskEvent struct {
	EventID   string    `json:"event_id"`
	TaskID    string    `json:"task_id"`
	EventType string    `json:"event_type"`
	Timestamp time.Time `json:"timestamp"`
}

type TaskCreatedEvent struct {
	TaskEvent
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

type TaskUpdatedEvent struct {
	TaskEvent
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

const (
	EventTypeTaskCreated = "TASK_CREATED"
	EventTypeTaskUpdated = "TASK_UPDATED"
)
