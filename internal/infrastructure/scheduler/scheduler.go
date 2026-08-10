package scheduler

import (
	"context"
	"fmt"
	"time"

	"main/internal/application"
)

type Scheduler struct {
	tasks    []application.ITask
	interval time.Duration
}

func New(
	interval time.Duration,
	tasks ...application.ITask,
) *Scheduler {
	return &Scheduler{
		tasks:    tasks,
		interval: interval,
	}
}

func (s *Scheduler) Run(ctx context.Context) error {
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	for {
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
	}

	return nil
}
