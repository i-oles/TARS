package doctorreminder

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"math/rand"
	"time"

	"main/internal/domain/contracts"
	"main/internal/domain/errs/api"
)

type DoctorReminder struct {
	requester      contracts.IRequester
	repoTasks      contracts.ITasks
	recipientEmail string
	refID          string
	interval       time.Duration
}

func New(
	requester contracts.IRequester,
	repoTasks contracts.ITasks,
	recipientEmail string,
	refID string,
	interval time.Duration,
) *DoctorReminder {
	return &DoctorReminder{
		requester:      requester,
		repoTasks:      repoTasks,
		recipientEmail: recipientEmail,
		refID:          refID,
		interval:       interval,
	}
}

func (t *DoctorReminder) Name() string {
	return "doctor_reminder"
}

func (t *DoctorReminder) Run(ctx context.Context) error {
	task, err := t.repoTasks.GetByName(ctx, t.Name())
	if err != nil {
		if errors.Is(err, api.ErrTaskNotFound) {
			task, err = t.repoTasks.Insert(ctx, t.Name())
			if err != nil {
				return fmt.Errorf("could insert task with name %s: %w", t.Name(), err)
			}
		} else {
			return fmt.Errorf("could not get task by %s: %w", t.Name(), err)
		}
	}

	if !task.Active {
		return nil
	}

	if task.LastRunAt != nil && time.Since(*task.LastRunAt) < t.interval {
		slog.Info("not enought time passed since lat run... ", slog.String("name", task.Name))

		return nil
	}

	subject, content, err := t.getSubjectAndContent()
	if err != nil {
		return fmt.Errorf("could get suject and content for: %s: %w", t.Name(), err)
	}

	err = t.requester.Request(subject, content, t.recipientEmail)
	if err != nil {
		return fmt.Errorf("could not request task %v: %w", t.Name(), err)
	}

	_, err = t.repoTasks.Update(ctx, task.ID, map[string]any{"last_run_at": time.Now()})
	if err != nil {
		return fmt.Errorf("could not update task %v: %w", t.Name(), err)
	}

	slog.Info("task finished", slog.String("name", task.Name))

	return nil
}

func (t *DoctorReminder) getSubjectAndContent() (string, string, error) {
	subjects := getSubjects()
	sKey := rand.Intn(len(subjects))

	subject, ok := getSubjects()[sKey]
	if !ok {
		return "", "", fmt.Errorf("subject with key: %d not found", sKey)
	}

	messgaes := getMessages(t.refID)
	mKey := rand.Intn(len(messgaes))

	content, ok := getMessages(t.refID)[mKey]
	if !ok {
		return "", "", fmt.Errorf("subject with key: %d not found", mKey)
	}

	return subject, content, nil
}
