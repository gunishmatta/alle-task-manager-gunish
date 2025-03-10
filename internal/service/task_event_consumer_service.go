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
	default:
		loggingtype.GetLogger().Error("Unknown event type: ", "event_type", baseEvent.EventType)
		return nil
	}
}
