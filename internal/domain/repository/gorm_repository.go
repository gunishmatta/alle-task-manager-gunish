package repository

import (
	"alle-task-manager-gunish/internal/common/errors"
	loggingtype "alle-task-manager-gunish/internal/common/logging"
	"alle-task-manager-gunish/internal/common/pagination"
	"alle-task-manager-gunish/internal/domain/model"
	"context"
	"gorm.io/gorm"
	"time"
)

type GormTaskRepository struct {
	db     *gorm.DB
	logger *loggingtype.Logger
}

func NewGormTaskRepository(db *gorm.DB) (*GormTaskRepository, error) {
	return &GormTaskRepository{db: db, logger: loggingtype.GetLogger()}, nil
}

func (r *GormTaskRepository) Create(_ context.Context, task *model.Task) error {
	task.CreatedAt = time.Now()
	task.UpdatedAt = task.CreatedAt

	result := r.db.Create(task)
	if result.Error != nil {
		r.logger.Error("Failed to create task", "error", result.Error)
		return result.Error
	}
	r.logger.Info("Task created successfully", "task_id", task.ID)
	return nil
}

func (r *GormTaskRepository) GetByID(_ context.Context, id string) (*model.Task, error) {
	var task model.Task
	result := r.db.First(&task, "id = ?", id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			r.logger.Warn("Task not found", "task_id", id)
			return nil, errors.ErrNotFound
		}
		r.logger.Error("Failed to get task", "task_id", id, "error", result.Error)
		return nil, result.Error
	}
	r.logger.Info("Task retrieved successfully", "task_id", id)
	return &task, nil
}

func (r *GormTaskRepository) Update(_ context.Context, task *model.Task) error {
	task.UpdatedAt = time.Now()
	result := r.db.Model(&model.Task{}).Where("id = ?", task.ID).Updates(task)
	if result.Error != nil {
		r.logger.Error("Failed to update task", "task_id", task.ID, "error", result.Error)
		return result.Error
	}
	if result.RowsAffected == 0 {
		r.logger.Warn("No task updated, task not found", "task_id", task.ID)
		return errors.ErrNotFound
	}
	r.logger.Info("Task updated successfully", "task_id", task.ID)
	return nil
}

func (r *GormTaskRepository) Delete(_ context.Context, id string) error {
	result := r.db.Delete(&model.Task{}, "id = ?", id)
	if result.Error != nil {
		r.logger.Error("Failed to delete task", "task_id", id, "error", result.Error)
		return result.Error
	}
	if result.RowsAffected == 0 {
		r.logger.Warn("No task deleted, task not found", "task_id", id)
		return errors.ErrNotFound
	}
	r.logger.Info("Task deleted successfully", "task_id", id)
	return nil
}

func (r *GormTaskRepository) List(_ context.Context, filter map[string]interface{}, page *pagination.Page) ([]*model.Task, int, error) {
	var tasks []model.Task
	var totalCount int64

	query := r.db.Model(&model.Task{})

	if filter != nil {
		if status, ok := filter["status"]; ok {
			query = query.Where("status = ?", status)
		}
	}

	if err := query.Count(&totalCount).Error; err != nil {
		r.logger.Error("Failed to get total count of tasks", "error", err)
		return nil, 0, err
	}

	if page != nil {
		query = query.Offset(page.Size * (page.Number - 1)).Limit(page.Size)
	}

	if err := query.Find(&tasks).Error; err != nil {
		r.logger.Error("Failed to list tasks", "error", err)
		return nil, 0, err
	}

	r.logger.Info("Tasks listed successfully", "count", len(tasks))
	taskPtrs := make([]*model.Task, len(tasks))
	for i := range tasks {
		taskPtrs[i] = &tasks[i]
	}
	return taskPtrs, int(totalCount), nil
}

func (r *GormTaskRepository) Close() error {
	sqlDB, err := r.db.DB()
	if err != nil {
		r.logger.Error("Failed to get underlying *sql.DB", "error", err)
		return err
	}
	r.logger.Info("Database connection closed")
	return sqlDB.Close()
}
