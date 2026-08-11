package models

import "time"

type Task struct {
	ID        int
	Name      string
	LastRunAt *time.Time
	Active    bool
}
