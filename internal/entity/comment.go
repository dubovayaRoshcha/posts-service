package entity

import "time"

type Comment struct {
	ID               int
	PostID           int
	ReplyToCommentID *int
	Text             string
	UserID           int
	CreatedAt        time.Time
}
