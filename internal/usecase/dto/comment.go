package dto

import "github.com/dubovayaRoshcha/posts-service/internal/entity"

type CommentRequest struct {
	PostID           int
	ReplyToCommentID *int
	Limit            int
	Offset           int
}

type CommentResponse struct {
	Comments []entity.Comment
	Len      int
}

type CommentDTO struct {
	PostID           int
	ReplyToCommentID *int
	UserID           int
	Text             string
}
