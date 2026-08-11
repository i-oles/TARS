package scheduler

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"main/internal/application"
)

type Scheduler struct {
	interval time.Duration
	tasks    []application.ITask
}

func New(
	interval time.Duration,
	tasks ...application.ITask,
) *Scheduler {
	return &Scheduler{
		interval: interval,
		tasks:    tasks,
	}
}

func (s *Scheduler) Run(ctx context.Context) error {
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	for {
		slog.Info("scheduler started...")

		if err := s.runTasks(ctx); err != nil {
			return fmt.Errorf("failed to run tasks: %w", err)
		}

		select {
		case <-ctx.Done():
			return fmt.Errorf("context failed with error: %w", ctx.Err())
		case <-ticker.C:
		}
	}
}

func (s *Scheduler) runTasks(ctx context.Context) error {
	for _, task := range s.tasks {
		if err := task.Run(ctx); err != nil {
			return fmt.Errorf("failed to run task: %s: %w", task.Name(), err)
		}

		slog.Info("task run", slog.String("name", task.Name()))
	}

	return nil
}
