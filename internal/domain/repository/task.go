package repository

import (
	"alle-task-manager-gunish/internal/common/pagination"
	"alle-task-manager-gunish/internal/domain/model"
	"context"
)

type TaskRepository interface {
	Create(ctx context.Context, task *model.Task) error
	GetByID(ctx context.Context, id string) (*model.Task, error)
	Update(ctx context.Context, task *model.Task) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, filter map[string]interface{}, page *pagination.Page) ([]*model.Task, int, error)
}
