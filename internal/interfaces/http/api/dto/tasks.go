package dto

import (
	"time"

	"main/internal/domain/models"
)

type UpdateTaskRequest struct {
	Active    *bool
	LastRunAt *time.Time
}

type UpdateTaskURI struct {
	ID int `binding:"required" uri:"task_id"`
}

type TaskResponse struct {
	ID        int        `json:"id"`
	Name      string     `json:"name"`
	Active    bool       `json:"active"`
	LastRunAt *time.Time `json:"last_run_at"`
}

func ToTaskResponse(task models.Task) TaskResponse {
	return TaskResponse{
		ID:        task.ID,
		Name:      task.Name,
		Active:    task.Active,
		LastRunAt: task.LastRunAt,
	}
}
