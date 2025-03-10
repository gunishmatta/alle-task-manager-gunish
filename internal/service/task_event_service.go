package service

import (
	"alle-task-manager-gunish/internal/common/events"
	"alle-task-manager-gunish/internal/common/kafka"
	"alle-task-manager-gunish/internal/domain/model"
	"github.com/google/uuid"
	"time"
)

const (
	TopicTaskEvents = "task-events"
)

type TaskEventPublisher interface {
	PublishTaskCreated(task *model.Task) error
	PublishTaskUpdated(task *model.Task) error
}

type TaskEventService struct {
	producer *kafka.Producer
}

var _ TaskEventPublisher = (*TaskEventService)(nil)

func NewTaskEventService(producer *kafka.Producer) *TaskEventService {
	return &TaskEventService{
		producer: producer,
	}
}

func (s *TaskEventService) PublishTaskCreated(task *model.Task) error {
	event := &events.TaskCreatedEvent{
		TaskEvent: events.TaskEvent{
			EventID:   uuid.New().String(),
			TaskID:    task.ID,
			EventType: events.EventTypeTaskCreated,
			Timestamp: time.Now(),
		},
		Title:       task.Title,
		Description: task.Description,
		Status:      string(task.Status),
	}

	return s.producer.PublishMessage(TopicTaskEvents, task.ID, event)
}

func (s *TaskEventService) PublishTaskUpdated(task *model.Task) error {
	event := &events.TaskUpdatedEvent{
		TaskEvent: events.TaskEvent{
			EventID:   uuid.New().String(),
			TaskID:    task.ID,
			EventType: events.EventTypeTaskUpdated,
			Timestamp: time.Now(),
		},
		Title:       task.Title,
		Description: task.Description,
		Status:      string(task.Status),
	}

	return s.producer.PublishMessage(TopicTaskEvents, task.ID, event)
}
