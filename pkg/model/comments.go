package model

import "time"

type Comments struct {
	ID        int
	TaskID    int
	UserID    int
	Text      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
