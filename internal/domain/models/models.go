package models

import (
	"time"
)

type Task struct {
	ID        int
	Name      string
	Type      TaskType
	Active    bool
	Interval  time.Duration
	LastRunAt *time.Time
	Config    []byte
}

type TaskCreation struct {
	Name     string
	Type     TaskType
	Active   bool
	Interval time.Duration
	Config   []byte
}

type TaskType string

const (
	TaskTypeCeneoCatcher   TaskType = "ceneo_catcher"
	TaskTypeDoctorReminder TaskType = "doctor_reminder"
)

type CeneoCatcherConfig struct {
	ProductID      int     `json:"product_id"`
	MaxPrice       float32 `json:"max_price"`
	RecipientEmail string  `json:"recipient_email"`
}

type DoctorReminderConfig struct {
	ReferenceID    string `json:"reference_id"`
	RecipientEmail string `json:"recipient_email"`
}
