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
	ProductID int `json:"product_id"`
	MaxPrice  int `json:"max_price"`
}

type DoctorReminderConfig struct {
	ReferenceID    string
	RecipientEmail string
}
