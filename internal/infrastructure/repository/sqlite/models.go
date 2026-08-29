package sqlite

import (
	"time"

	"main/internal/domain/models"
)

type SQLTask struct {
	ID        int           `gorm:"primaryKey"`
	Name      string        `gorm:"not null;uniqueIndex"`
	Type      string        `gorm:"not null;index"`
	Active    bool          `gorm:"not null"`
	Interval  time.Duration `gorm:"not null"`
	LastRunAt *time.Time
	Config    []byte
}

func (SQLTask) TableName() string {
	return "tasks"
}

func (s SQLTask) ToDomain() models.Task {
	return models.Task{
		ID:        s.ID,
		Name:      s.Name,
		Type:      models.TaskType(s.Type),
		Active:    s.Active,
		Interval:  s.Interval,
		LastRunAt: s.LastRunAt,
		Config:    s.Config,
	}
}

func SQLTaskFromDomain(task models.Task) SQLTask {
	return SQLTask{
		ID:        task.ID,
		Name:      task.Name,
		Type:      string(task.Type),
		Active:    task.Active,
		Interval:  task.Interval,
		LastRunAt: task.LastRunAt,
		Config:    task.Config,
	}
}
