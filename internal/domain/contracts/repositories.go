package contracts

import (
	"context"

	"main/internal/domain/models"
)

type ITasks interface {
	GetByName(ctx context.Context, name string) (models.Task, error)
	Insert(ctx context.Context, name string) (models.Task, error)
	Update(ctx context.Context, id int, update map[string]any) (models.Task, error)
	List(ctx context.Context) ([]models.Task, error)
}
