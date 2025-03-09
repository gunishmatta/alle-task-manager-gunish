package repository

import (
	"alle-task-manager-gunish/internal/common/errors"
	"alle-task-manager-gunish/internal/common/pagination"
	"alle-task-manager-gunish/internal/domain/model"
	"context"
	"sync"
	"time"
)

type MemoryTaskRepository struct {
	tasks map[string]*model.Task
	mutex sync.RWMutex
}

func NewMemoryTaskRepository() *MemoryTaskRepository {
	return &MemoryTaskRepository{
		tasks: make(map[string]*model.Task),
	}
}

func (r *MemoryTaskRepository) Create(ctx context.Context, task *model.Task) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.tasks[task.ID]; exists {
		return errors.ErrDuplicateEntity
	}

	r.tasks[task.ID] = task
	return nil
}

func (r *MemoryTaskRepository) GetByID(ctx context.Context, id string) (*model.Task, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	task, exists := r.tasks[id]
	if !exists {
		return nil, errors.ErrNotFound
	}

	return cloneTask(task), nil
}

func (r *MemoryTaskRepository) Update(ctx context.Context, task *model.Task) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.tasks[task.ID]; !exists {
		return errors.ErrNotFound
	}

	task.UpdatedAt = time.Now()
	r.tasks[task.ID] = task
	return nil
}

func (r *MemoryTaskRepository) Delete(ctx context.Context, id string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.tasks[id]; !exists {
		return errors.ErrNotFound
	}

	delete(r.tasks, id)
	return nil
}

func (r *MemoryTaskRepository) List(ctx context.Context, filter map[string]interface{}, page *pagination.Page) ([]*model.Task, int, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var filteredTasks []*model.Task

	for _, task := range r.tasks {
		if matchesFilter(task, filter) {
			filteredTasks = append(filteredTasks, cloneTask(task))
		}
	}

	totalCount := len(filteredTasks)

	if page != nil {
		start, end := page.GetLimits(totalCount)
		if start < totalCount {
			if end > totalCount {
				end = totalCount
			}
			filteredTasks = filteredTasks[start:end]
		} else {
			filteredTasks = []*model.Task{}
		}
	}

	return filteredTasks, totalCount, nil
}

func matchesFilter(task *model.Task, filter map[string]interface{}) bool {
	if filter == nil {
		return true
	}

	for key, value := range filter {
		switch key {
		case "status":
			if statusFilter, ok := value.(string); ok && string(task.Status) != statusFilter {
				return false
			}
		}
	}

	return true
}

func cloneTask(t *model.Task) *model.Task {
	clone := *t
	return &clone
}
