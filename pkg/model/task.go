package model

import "time"

type Task struct {
	ID          int
	Title       string
	Description string
	Status      string
	Priority    string
	AssignedTo  int
	ProjectID   int
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DueDate     *time.Time
}
