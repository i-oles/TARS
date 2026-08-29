package sqlite

import (
	"context"
	"errors"
	"fmt"

	"main/internal/domain/errs/api"
	"main/internal/domain/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type tasksRepo struct {
	db *gorm.DB
}

func NewTasksRepo(db *gorm.DB) *tasksRepo {
	return &tasksRepo{
		db: db,
	}
}

func (r *tasksRepo) GetByName(ctx context.Context, name string) (models.Task, error) {
	var sqlTask SQLTask

	if err := r.db.WithContext(ctx).
		Where("name = ?", name).
		First(&sqlTask).
		Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.Task{}, api.ErrTaskNotFound
		}

		return models.Task{}, fmt.Errorf("could get task: %w", err)
	}

	return sqlTask.ToDomain(), nil
}

func (r *tasksRepo) Insert(
	ctx context.Context,
	taskCreation models.TaskCreation,
) (models.Task, error) {
	sqlTask := SQLTask{
		Name:     taskCreation.Name,
		Type:     string(taskCreation.Type),
		Active:   taskCreation.Active,
		Interval: taskCreation.Interval,
		Config:   taskCreation.Config,
	}

	if err := r.db.WithContext(ctx).
		Create(&sqlTask).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return models.Task{}, api.ErrTaskAlreadyExist
		}

		return models.Task{}, fmt.Errorf("could not insert task: %w", err)
	}

	return sqlTask.ToDomain(), nil
}

func (r *tasksRepo) Update(
	ctx context.Context, id int, update map[string]any,
) (models.Task, error) {
	var task SQLTask

	if err := r.db.WithContext(ctx).
		Model(&task).
		Clauses(clause.Returning{}).
		Where("id = ?", id).
		Updates(update).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.Task{}, api.ErrTaskNotFound
		}

		return models.Task{},
			fmt.Errorf("could not update task with id: %v with data: %v, %w", id, update, err)
	}

	return task.ToDomain(), nil
}

func (r *tasksRepo) List(ctx context.Context) ([]models.Task, error) {
	var sqlTasks []SQLTask

	if err := r.db.WithContext(ctx).Find(&sqlTasks).Error; err != nil {
		return nil, fmt.Errorf("could not get tasks: %w", err)
	}

	result := make([]models.Task, len(sqlTasks))

	for i, task := range sqlTasks {
		result[i] = task.ToDomain()
	}

	return result, nil
}
