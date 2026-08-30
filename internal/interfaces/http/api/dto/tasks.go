package dto

import (
	"encoding/json"
	"time"

	"main/internal/domain/models"
)

type CreateTaskRequest struct {
	Name      string
	Type      string
	Active    bool
	Interval  string
	LastRunAt time.Time
	Config    json.RawMessage
}

type UpdateTaskRequest struct {
	Name      *string
	Type      *string
	Active    *bool
	Interval  *string
	LastRunAt *time.Time
	Config    *json.RawMessage
}

type UpdateTaskURI struct {
	ID int `binding:"required" uri:"task_id"`
}

type TaskResponse struct {
	ID        int             `json:"id"`
	Name      string          `json:"name"`
	Type      string          `json:"type"`
	Active    bool            `json:"active"`
	Interval  time.Duration   `json:"interval"`
	LastRunAt *time.Time      `json:"last_run_at"`
	Config    json.RawMessage `json:"config"`
}

func ToTaskResponse(task models.Task) TaskResponse {
	return TaskResponse{
		ID:        task.ID,
		Name:      task.Name,
		Type:      string(task.Type),
		Active:    task.Active,
		Interval:  task.Interval,
		LastRunAt: task.LastRunAt,
		Config:    task.Config,
	}
}
