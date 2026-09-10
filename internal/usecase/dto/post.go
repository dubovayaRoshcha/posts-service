package dto

import "github.com/dubovayaRoshcha/posts-service/internal/entity"

type PostRequest struct {
	Limit  int
	Offset int
}

type PostResponse struct {
	Posts []entity.Post
	Len   int
}

type PostDTO struct {
	UserID          int
	Title           string
	Description     string
	CommentsAllowed bool
}
