package models

import "time"

type Task struct {
	ID         int
	Name       string
	RemindedAt time.Time
	Active     bool
}
