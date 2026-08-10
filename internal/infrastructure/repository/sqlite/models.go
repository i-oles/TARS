package sqlite

import (
	"time"

	"main/internal/domain/models"
)

type SQLTask struct {
	ID         int    `gorm:"primaryKey"`
	Name       string `gorm:"not null"`
	Active     bool   `gorm:"not null"`
	RemindedAt time.Time
}

func (SQLTask) TableName() string {
	return "tasks"
}

func (s SQLTask) ToDomain() models.Task {
	return models.Task{
		ID:         s.ID,
		Name:       s.Name,
		RemindedAt: s.RemindedAt,
		Active:     s.Active,
	}
}

func SQLTaskFromDomain(task models.Task) SQLTask {
	return SQLTask{
		ID:         task.ID,
		Name:       task.Name,
		RemindedAt: task.RemindedAt,
		Active:     task.Active,
	}
}
