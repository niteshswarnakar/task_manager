package model

import (
	"time"

	"github.com/google/uuid"
)

type TaskStatus string

const (
	StatusPending   TaskStatus = "pending"
	StatusCompleted TaskStatus = "completed"
)

type Task struct {
	ID        string     `json:"id" gorm:"primary_key"`
	CreatedAt time.Time  `json:"created_at"`
	Title     string     `json:"title" gorm:"unique"`
	Status    TaskStatus `json:"status"`
}

func NewID() string {
	return uuid.New().String()
}
