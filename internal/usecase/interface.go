package usecase

import (
	"context"

	"github.com/dubovayaRoshcha/posts-service/internal/entity"
	"github.com/dubovayaRoshcha/posts-service/internal/usecase/dto"
)

//go:generate mockgen -source=interface.go -destination=../mocks/mocks_repo.go -package=mocks
type PostRepo interface {
	GetList(ctx context.Context, params dto.PostRequest) ([]entity.Post, error)
	GetByID(ctx context.Context, postID int) (*entity.Post, error)
	Create(ctx context.Context, post dto.PostDTO) (*entity.Post, error)
}

type CommentRepo interface {
	GetTopList(ctx context.Context, params dto.CommentRequest) ([]entity.Comment, error)
	GetReplies(ctx context.Context, params dto.CommentRequest) ([]entity.Comment, error)
	GetByID(ctx context.Context, commentID int) (*entity.Comment, error)
	Create(ctx context.Context, comment dto.CommentDTO) (*entity.Comment, error)
}
