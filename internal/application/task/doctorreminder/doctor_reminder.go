package doctorreminder

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"math/rand"
	"time"

	"main/internal/application/email"
	"main/internal/domain/contracts"
	"main/internal/domain/errs/api"
	"main/internal/domain/models"
)

type DoctorReminder struct {
	emailComposer  email.Composer
	mailer         email.IMailer
	tasksRepo      contracts.ITasks
	recipientEmail string
	refID          string
	taskInterval   time.Duration
}

func New(
	emailComposer email.Composer,
	mailer email.IMailer,
	tasksRepo contracts.ITasks,
	recipientEmail string,
	refID string,
	taskInterval time.Duration,
) *DoctorReminder {
	return &DoctorReminder{
		emailComposer:  emailComposer,
		mailer:         mailer,
		tasksRepo:      tasksRepo,
		recipientEmail: recipientEmail,
		refID:          refID,
		taskInterval:   taskInterval,
	}
}

func (t *DoctorReminder) Name() string {
	return "doctor_reminder"
}

func (t *DoctorReminder) Run(ctx context.Context) error {
	if isWeekend(time.Now()) {
		return nil
	}

	task, err := t.getOrInsertTask(ctx)
	if err != nil {
		return fmt.Errorf("could get or insert task: %s: %w", t.Name(), err)
	}

	if !task.Active {
		return nil
	}

	if task.LastRunAt != nil && time.Since(*task.LastRunAt) < t.taskInterval {
		slog.Info("not enought time passed since lat run... ", slog.String("name", task.Name))

		return nil
	}

	subject, content, err := t.getSubjectAndContent()
	if err != nil {
		return fmt.Errorf("could get suject and content for: %s: %w", t.Name(), err)
	}

	data := email.DoctorReminder{
		Subject:        subject,
		Content:        content,
		RecipientEmail: t.recipientEmail,
	}

	msg, err := t.emailComposer.ComposeForDoctorReminder(data)
	if err != nil {
		return fmt.Errorf("could not request task %v: %w", t.Name(), err)
	}

	err = t.mailer.Send(msg)
	if err != nil {
		return fmt.Errorf("could not send msg %v: %w", t.Name(), err)
	}

	_, err = t.tasksRepo.Update(ctx, task.ID, map[string]any{"last_run_at": time.Now()})
	if err != nil {
		return fmt.Errorf("could not update task %v: %w", t.Name(), err)
	}

	slog.Info("task finished", slog.String("name", task.Name))

	return nil
}

func isWeekend(t time.Time) bool {
	return t.Weekday() == time.Saturday || t.Weekday() == time.Sunday
}

func (t *DoctorReminder) getOrInsertTask(ctx context.Context) (models.Task, error) {
	task, err := t.tasksRepo.GetByName(ctx, t.Name())
	if err != nil {
		if errors.Is(err, api.ErrTaskNotFound) {
			task, err = t.tasksRepo.Insert(ctx, t.Name())
			if err != nil {
				return models.Task{}, fmt.Errorf("could insert task with name %s: %w", t.Name(), err)
			}
		} else {
			return models.Task{}, fmt.Errorf("could not get task by %s: %w", t.Name(), err)
		}
	}

	return task, nil
}

func (t *DoctorReminder) getSubjectAndContent() (string, string, error) {
	subjects := getSubjects()
	subjectKey := rand.Intn(len(subjects))

	subject, ok := getSubjects()[subjectKey]
	if !ok {
		return "", "", fmt.Errorf("subject with key: %d not found", subjectKey)
	}

	messages := getMessages(t.refID)
	msgKey := rand.Intn(len(messages))

	content, ok := getMessages(t.refID)[msgKey]
	if !ok {
		return "", "", fmt.Errorf("subject with key: %d not found", msgKey)
	}

	return subject, content, nil
}
