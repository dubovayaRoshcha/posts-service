package entity

import "time"

type Post struct {
	ID              int
	UserID          int
	Title           string
	Description     string
	CommentsAllowed bool
	CreatedAt       time.Time
}
