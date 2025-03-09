package service

import (
	"alle-task-manager-gunish/internal/common/errors"
	"alle-task-manager-gunish/internal/common/pagination"
	"alle-task-manager-gunish/internal/domain/model"
	"alle-task-manager-gunish/internal/domain/repository"
	"context"
	"time"
)

type TaskService struct {
	repo repository.TaskRepository
}

func NewTaskService(repo repository.TaskRepository) *TaskService {
	return &TaskService{
		repo: repo,
	}
}

type CreateTaskInput struct {
	Title       string     `json:"title" validate:"required"`
	Description string     `json:"description"`
	DueDate     *time.Time `json:"due_date,omitempty"`
}

func (s *TaskService) CreateTask(ctx context.Context, input CreateTaskInput) (*model.Task, error) {
	task := model.NewTask(input.Title, input.Description)
	if input.DueDate != nil {
		task.DueDate = input.DueDate
	}

	if err := s.repo.Create(ctx, task); err != nil {
		return nil, err
	}

	return task, nil
}

func (s *TaskService) GetTask(ctx context.Context, id string) (*model.Task, error) {
	return s.repo.GetByID(ctx, id)
}

type UpdateTaskInput struct {
	Title       *string    `json:"title,omitempty"`
	Description *string    `json:"description,omitempty"`
	Status      *string    `json:"status,omitempty"`
	DueDate     *time.Time `json:"due_date,omitempty"`
}

func (s *TaskService) UpdateTask(ctx context.Context, id string, input UpdateTaskInput) (*model.Task, error) {
	task, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if input.Title != nil {
		task.Title = *input.Title
	}

	if input.Description != nil {
		task.Description = *input.Description
	}

	if input.Status != nil {
		status := model.TaskStatus(*input.Status)
		if status != model.Pending && status != model.InProgress && status != model.Completed {
			return nil, errors.ErrInvalidStatus
		}
		task.Status = status
	}

	if input.DueDate != nil {
		task.DueDate = input.DueDate
	}

	task.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, task); err != nil {
		return nil, err
	}

	return task, nil
}

func (s *TaskService) DeleteTask(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *TaskService) ListTasks(ctx context.Context, status string, page *pagination.Page) ([]*model.Task, *pagination.PageInfo, error) {
	filter := make(map[string]interface{})
	if status != "" {
		filter["status"] = status
	}

	tasks, total, err := s.repo.List(ctx, filter, page)
	if err != nil {
		return nil, nil, err
	}

	pageInfo := &pagination.PageInfo{
		Page:       page.Number,
		PageSize:   page.Size,
		TotalItems: total,
		TotalPages: (total + page.Size - 1) / page.Size,
	}

	return tasks, pageInfo, nil
}
