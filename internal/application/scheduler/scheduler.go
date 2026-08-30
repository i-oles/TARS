package scheduler

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"main/internal/application/task/ceneocatcher"
	"main/internal/application/task/doctorreminder"
	"main/internal/domain/contracts"
	"main/internal/domain/models"
)

type Scheduler struct {
	tasksRepo                contracts.ITasks
	ceneoCatcherTaskRunner   *ceneocatcher.CeneoCatcherTaskRunner
	doctorReminderTaskRunner *doctorreminder.DoctorReminderTaskRunner
	interval                 time.Duration
}

func New(
	tasksRepo contracts.ITasks,
	ceneoCatcherTaskRunner *ceneocatcher.CeneoCatcherTaskRunner,
	doctorReminderTaskRunner *doctorreminder.DoctorReminderTaskRunner,
	interval time.Duration,
) *Scheduler {
	return &Scheduler{
		tasksRepo:                tasksRepo,
		ceneoCatcherTaskRunner:   ceneoCatcherTaskRunner,
		doctorReminderTaskRunner: doctorReminderTaskRunner,
		interval:                 interval,
	}
}

func (s *Scheduler) Run(ctx context.Context) error {
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	for {
		slog.Info("--------------------------------------")
		slog.Info("scheduler started...")

		if err := s.runTasks(ctx); err != nil {
			return fmt.Errorf("failed to run tasks: %w", err)
		}

		slog.Info("scheduler finished, sleep...")

		select {
		case <-ctx.Done():
			return fmt.Errorf("context failed with error: %w", ctx.Err())
		case <-ticker.C:
		}
	}
}

func (s *Scheduler) runTasks(ctx context.Context) error {
	tasks, err := s.tasksRepo.List(ctx)
	if err != nil {
		return fmt.Errorf("failed to list tasks: %w", err)
	}

	for _, task := range tasks {
		if !shouldRunTask(task) {
			continue
		}

		switch task.Type {
		case models.TaskTypeCeneoCatcher:
			err := s.ceneoCatcherTaskRunner.Run(ctx, task.ID, task.Config)
			if err != nil {
				return fmt.Errorf("failed to run ceneo catcher task: %w", err)
			}
		case models.TaskTypeDoctorReminder:
			if isWeekend(time.Now()) {
				continue
			}

			err := s.doctorReminderTaskRunner.Run(ctx, task.ID, task.Config)
			if err != nil {
				return fmt.Errorf("failed to run doctor reminder task: %w", err)
			}
		}
	}

	return nil
}

func isWeekend(t time.Time) bool {
	return t.Weekday() == time.Saturday || t.Weekday() == time.Sunday
}

func shouldRunTask(task models.Task) bool {
	if !task.Active {
		return false
	}

	if task.LastRunAt != nil && time.Since(*task.LastRunAt) < task.Interval {
		slog.Info("not enought time passed since lat run... ", slog.String("name", task.Name))

		return false
	}

	return true
}
