package doctorreminder

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"math/rand"
	"time"

	"main/internal/application/email"
	"main/internal/domain/contracts"
	"main/internal/domain/models"
)

type DoctorReminderTaskRunner struct {
	emailComposer email.Composer
	mailer        email.IMailer
	tasksRepo     contracts.ITasks
}

func NewTaskRunner(
	emailComposer email.Composer,
	mailer email.IMailer,
	tasksRepo contracts.ITasks,
) *DoctorReminderTaskRunner {
	return &DoctorReminderTaskRunner{
		emailComposer: emailComposer,
		mailer:        mailer,
		tasksRepo:     tasksRepo,
	}
}

func (t *DoctorReminderTaskRunner) Run(ctx context.Context, taskID int, config []byte) error {
	var cfg models.DoctorReminderConfig

	if err := json.Unmarshal(config, &cfg); err != nil {
		return fmt.Errorf("invalid doctor reminder config: %w", err)
	}

	subject, content, err := t.getSubjectAndContent(cfg.ReferenceID)
	if err != nil {
		return fmt.Errorf("could get suject and content for task with ID %d: %w", taskID, err)
	}

	data := email.DoctorReminder{
		Subject:        subject,
		Content:        content,
		RecipientEmail: cfg.RecipientEmail,
	}

	msg, err := t.emailComposer.ComposeForDoctorReminder(data)
	if err != nil {
		return fmt.Errorf("could not compose msg for doctor reminder task: %w", err)
	}

	err = t.mailer.Send(msg)
	if err != nil {
		return fmt.Errorf("could not send msg %v: %w", msg, err)
	}

	_, err = t.tasksRepo.Update(ctx, taskID, map[string]any{"last_run_at": time.Now()})
	if err != nil {
		return fmt.Errorf("could not update task with ID %d: %w", taskID, err)
	}

	slog.Info("task finished", slog.Int("id", taskID))

	return nil
}

func (t *DoctorReminderTaskRunner) getSubjectAndContent(refID string) (string, string, error) {
	subjects := getSubjects()
	subjectKey := rand.Intn(len(subjects))

	subject, ok := getSubjects()[subjectKey]
	if !ok {
		return "", "", fmt.Errorf("subject with key: %d not found", subjectKey)
	}

	messages := getMessages(refID)
	msgKey := rand.Intn(len(messages))

	content, ok := getMessages(refID)[msgKey]
	if !ok {
		return "", "", fmt.Errorf("subject with key: %d not found", msgKey)
	}

	return subject, content, nil
}
