package task

import (
	"context"
	"fmt"
	"math/rand"

	"main/assets"
	domain "main/internal/domain/requester"
)

type DoctorReminder struct {
	requester      domain.IRequester
	recipientEmail string
}

func NewDoctorReminder(
	requester domain.IRequester,
	recipientEmail string,
) *DoctorReminder {
	return &DoctorReminder{
		requester:      requester,
		recipientEmail: recipientEmail,
	}
}

func (t *DoctorReminder) Name() string {
	return "doctor_reminder"
}

func (t *DoctorReminder) Run(ctx context.Context) error {
	subject := assets.Subjects[rand.Intn(11)] // no lint
	content := assets.Messages[rand.Intn(21)] // no lint

	err := t.requester.Request(subject, content, t.recipientEmail)
	if err != nil {
		return fmt.Errorf("could not request task %v: %w", t.Name(), err)
	}

	return nil
}
