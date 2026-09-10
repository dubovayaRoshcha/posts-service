package entity

import "errors"

var (
	InvalidInput       = errors.New("invalid data")
	CommentsNotAllowed = errors.New("comments are disabled for this post")
	PostNotFound       = errors.New("post not found")
	CommentNotFound    = errors.New("comment not found")
	MaxLengthExceeded  = errors.New("maximum length exceeded")
)
