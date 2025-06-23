package model

import "time"

type Project struct {
	ID          int
	Name        string
	Description string
	OwnerID     int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
