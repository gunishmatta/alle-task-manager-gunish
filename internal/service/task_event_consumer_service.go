package service

import (
	"alle-task-manager-gunish/internal/common/events"
	loggingtype "alle-task-manager-gunish/internal/common/logging"
	"encoding/json"
	"github.com/IBM/sarama"
)

type TaskEventConsumerService struct {
}

func NewTaskEventConsumerService() *TaskEventConsumerService {
	return &TaskEventConsumerService{}
}

func (s *TaskEventConsumerService) HandleMessage(message *sarama.ConsumerMessage) error {
	var baseEvent events.TaskEvent
	if err := json.Unmarshal(message.Value, &baseEvent); err != nil {
		return err
	}

	switch baseEvent.EventType {
	case events.EventTypeTaskCreated:
		return s.handleTaskCreated(message.Value)
	case events.EventTypeTaskUpdated:
		return s.handleTaskUpdated(message.Value)
	default:
		loggingtype.GetLogger().Error("Unknown event type: ", "event_type", baseEvent.EventType)
		return nil
	}
}

func (s *TaskEventConsumerService) handleTaskCreated(data []byte) error {
	var event events.TaskCreatedEvent
	if err := json.Unmarshal(data, &event); err != nil {
		return err
	}

	loggingtype.GetLogger().Info("Processing task created event: %s - %s\n", event.TaskID, event.Title)
	return nil
}

func (s *TaskEventConsumerService) handleTaskUpdated(data []byte) error {
	var event events.TaskUpdatedEvent
	if err := json.Unmarshal(data, &event); err != nil {
		return err
	}
	loggingtype.GetLogger().Info("Processing task updated event: %s - %s\n", event.TaskID, event.Title)
	return nil
}
