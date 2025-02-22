package model

import "time"

type TaskStatus string

const (
	StatusPending   TaskStatus = "pending"
	StatusCompleted TaskStatus = "completed"
)

type Task struct {
	ID          string     `json:"id" gorm:"primary_key"`
	CreatedAt   time.Time  `json:"created_at"`
	Title       string     `json:"title" gorm:"unique"`
	Description string     `json:"description"`
	Status      TaskStatus `json:"status"`
}
