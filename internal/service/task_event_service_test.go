package service

import (
	"alle-task-manager-gunish/internal/common/kafka"
	"alle-task-manager-gunish/internal/domain/model"
	"github.com/IBM/sarama"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

type MockSyncProducer struct {
	mock.Mock
}

func (m *MockSyncProducer) SendMessage(msg *sarama.ProducerMessage) (partition int32, offset int64, err error) {
	args := m.Called(msg)
	return args.Get(0).(int32), args.Get(1).(int64), args.Error(2)
}

func (m *MockSyncProducer) SendMessages(msgs []*sarama.ProducerMessage) error {
	args := m.Called(msgs)
	return args.Error(0)
}

func (m *MockSyncProducer) Close() error {
	args := m.Called()
	return args.Error(0)
}

func TestTaskEventService(t *testing.T) {
	mockSyncProducer := new(MockSyncProducer)
	producer := &kafka.Producer{Producer: mockSyncProducer}
	service := NewTaskEventService(producer)

	t.Run("PublishTaskCreated", func(t *testing.T) {
		task := &model.Task{
			ID:          "test-id",
			Title:       "Test Task",
			Description: "Test Description",
			Status:      model.Pending,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		mockSyncProducer.On("SendMessage", mock.Anything).
			Return(int32(0), int64(1), nil).Once()

		err := service.PublishTaskCreated(task)
		require.NoError(t, err)

		mockSyncProducer.AssertExpectations(t)
	})

	t.Run("PublishTaskUpdated", func(t *testing.T) {
		task := &model.Task{
			ID:          "test-id",
			Title:       "Updated Task",
			Description: "Updated Description",
			Status:      model.InProgress,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		mockSyncProducer.On("SendMessage", mock.Anything).
			Return(int32(0), int64(0), nil)

		err := service.PublishTaskUpdated(task)
		require.NoError(t, err)

		mockSyncProducer.AssertExpectations(t)
	})

}
