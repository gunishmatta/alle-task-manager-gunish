package service

import (
	"alle-task-manager-gunish/internal/common/errors"
	"alle-task-manager-gunish/internal/common/pagination"
	"alle-task-manager-gunish/internal/domain/model"
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type MockTaskRepository struct {
	mock.Mock
}

func (m *MockTaskRepository) Create(ctx context.Context, task *model.Task) error {
	args := m.Called(ctx, task)
	return args.Error(0)
}

func (m *MockTaskRepository) GetByID(ctx context.Context, id string) (*model.Task, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Task), args.Error(1)
}

func (m *MockTaskRepository) Update(ctx context.Context, task *model.Task) error {
	args := m.Called(ctx, task)
	return args.Error(0)
}

func (m *MockTaskRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockTaskRepository) List(ctx context.Context, filter map[string]interface{}, page *pagination.Page) ([]*model.Task, int, error) {
	args := m.Called(ctx, filter, page)
	return args.Get(0).([]*model.Task), args.Int(1), args.Error(2)
}

type MockTaskEventService struct {
	mock.Mock
}

func (m *MockTaskEventService) PublishTaskCreated(task *model.Task) error {
	args := m.Called(task)
	return args.Error(0)
}

func (m *MockTaskEventService) PublishTaskUpdated(task *model.Task) error {
	args := m.Called(task)
	return args.Error(0)
}

func TestTaskService_CreateTask(t *testing.T) {
	mockRepo := new(MockTaskRepository)
	mockEventSvc := new(MockTaskEventService)
	service := NewTaskService(mockRepo, mockEventSvc)
	ctx := context.Background()

	t.Run("successful creation", func(t *testing.T) {
		input := CreateTaskInput{
			Title:       "Test Task",
			Description: "Test Description",
		}

		mockRepo.On("Create", ctx, mock.AnythingOfType("*model.Task")).Return(nil).Once()
		mockEventSvc.On("PublishTaskCreated", mock.AnythingOfType("*model.Task")).Return(nil).Once()

		task, err := service.CreateTask(ctx, input)

		assert.NoError(t, err)
		assert.NotNil(t, task)
		assert.Equal(t, input.Title, task.Title)
		assert.Equal(t, input.Description, task.Description)
		assert.Equal(t, model.Pending, task.Status)

		mockRepo.AssertExpectations(t)
		mockEventSvc.AssertExpectations(t)
	})

	t.Run("repository error", func(t *testing.T) {
		input := CreateTaskInput{
			Title:       "Test Task",
			Description: "Test Description",
		}

		mockRepo.On("Create", ctx, mock.AnythingOfType("*model.Task")).Return(errors.ErrDuplicateEntity).Once()

		task, err := service.CreateTask(ctx, input)

		assert.Error(t, err)
		assert.Nil(t, task)
		assert.Equal(t, errors.ErrDuplicateEntity, err)

		mockRepo.AssertExpectations(t)
	})
}

func TestTaskService_UpdateTask(t *testing.T) {
	mockRepo := new(MockTaskRepository)
	mockEventSvc := new(MockTaskEventService)
	service := NewTaskService(mockRepo, mockEventSvc)
	ctx := context.Background()

	existingTask := &model.Task{
		ID:          "test-id",
		Title:       "Original Title",
		Description: "Original Description",
		Status:      model.Pending,
	}

	t.Run("successful update", func(t *testing.T) {
		newTitle := "Updated Title"
		newStatus := string(model.InProgress)
		input := UpdateTaskInput{
			Title:  &newTitle,
			Status: &newStatus,
		}

		mockRepo.On("GetByID", ctx, existingTask.ID).Return(existingTask, nil).Once()
		mockRepo.On("Update", ctx, mock.AnythingOfType("*model.Task")).Return(nil).Once()
		mockEventSvc.On("PublishTaskUpdated", mock.AnythingOfType("*model.Task")).Return(nil).Once()

		updatedTask, err := service.UpdateTask(ctx, existingTask.ID, input)

		assert.NoError(t, err)
		assert.NotNil(t, updatedTask)
		assert.Equal(t, newTitle, updatedTask.Title)
		assert.Equal(t, model.TaskStatus(newStatus), updatedTask.Status)

		mockRepo.AssertExpectations(t)
		mockEventSvc.AssertExpectations(t)
	})

	t.Run("invalid status", func(t *testing.T) {
		invalidStatus := "invalid_status"
		input := UpdateTaskInput{
			Status: &invalidStatus,
		}

		mockRepo.On("GetByID", ctx, existingTask.ID).Return(existingTask, nil).Once()

		updatedTask, err := service.UpdateTask(ctx, existingTask.ID, input)

		assert.Error(t, err)
		assert.Nil(t, updatedTask)
		assert.Equal(t, errors.ErrInvalidStatus, err)

		mockRepo.AssertExpectations(t)
	})
}

func TestTaskService_ListTasks(t *testing.T) {
	mockRepo := new(MockTaskRepository)
	mockEventSvc := new(MockTaskEventService)
	service := NewTaskService(mockRepo, mockEventSvc)
	ctx := context.Background()

	t.Run("successful listing", func(t *testing.T) {
		tasks := []*model.Task{
			{ID: "1", Title: "Task 1", Status: model.Pending},
			{ID: "2", Title: "Task 2", Status: model.InProgress},
		}
		page := &pagination.Page{Number: 1, Size: 10}
		totalItems := 2

		expectedFilter := map[string]interface{}{"status": string(model.Pending)}
		mockRepo.On("List", ctx, expectedFilter, page).Return(tasks, totalItems, nil).Once()

		resultTasks, pageInfo, err := service.ListTasks(ctx, string(model.Pending), page)

		assert.NoError(t, err)
		assert.EqualValues(t, tasks, resultTasks)
		assert.NotNil(t, pageInfo)
		assert.Equal(t, 1, pageInfo.Page)
		assert.Equal(t, 10, pageInfo.PageSize)
		assert.Equal(t, totalItems, pageInfo.TotalItems)
		assert.Equal(t, 1, pageInfo.TotalPages)

		mockRepo.AssertExpectations(t)
	})

}

func TestTaskService_DeleteTask(t *testing.T) {
	mockRepo := new(MockTaskRepository)
	mockEventSvc := new(MockTaskEventService)
	service := NewTaskService(mockRepo, mockEventSvc)
	ctx := context.Background()

	t.Run("successful deletion", func(t *testing.T) {
		taskID := "test-id"
		mockRepo.On("Delete", ctx, taskID).Return(nil).Once()

		err := service.DeleteTask(ctx, taskID)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("not found error", func(t *testing.T) {
		taskID := "non-existent-id"
		mockRepo.On("Delete", ctx, taskID).Return(errors.ErrNotFound).Once()

		err := service.DeleteTask(ctx, taskID)

		assert.Error(t, err)
		assert.Equal(t, errors.ErrNotFound, err)
		mockRepo.AssertExpectations(t)
	})
}
