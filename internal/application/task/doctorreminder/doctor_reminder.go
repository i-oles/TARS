package doctorreminder

import (
	"context"
	"errors"
	"fmt"
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
		return nil
	}

	subjects := getSubjects()
	sKey := rand.Intn(len(subjects))

	subject, ok := getSubjects()[sKey]
	if !ok {
		return fmt.Errorf("subject with key: %d not found", sKey)
	}

	messgaes := getMessages(t.refID)
	mKey := rand.Intn(len(messgaes))

	content, ok := getMessages(t.refID)[mKey]
	if !ok {
		return fmt.Errorf("subject with key: %d not found", mKey)
	}

	err = t.requester.Request(subject, content, t.recipientEmail)
	if err != nil {
		return fmt.Errorf("could not request task %v: %w", t.Name(), err)
	}

	return nil
}
