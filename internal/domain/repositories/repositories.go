package repositories

import (
	"context"
	"time"

	"main/internal/domain/models"
)

type ITasks interface {
	Insert(ctx context.Context, name string) (models.Task, error)
	Update(ctx context.Context, id int, update map[string]any) (models.Task, error)
	List(ctx context.Context) ([]models.Task, error)

	FindDue(ctx context.Context, before time.Time) ([]models.Task, error)
}
