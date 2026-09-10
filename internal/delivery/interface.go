package delivery

import (
	"context"

	"github.com/dubovayaRoshcha/posts-service/internal/entity"
	"github.com/dubovayaRoshcha/posts-service/internal/usecase/dto"
)

//go:generate mockgen -source=interface.go -destination=../mocks/mocks_usecase.go -package=mocks
type PostUseCase interface {
	GetListPosts(ctx context.Context, params dto.PostRequest) (*dto.PostResponse, error)
	GetPostByID(ctx context.Context, postID int) (*entity.Post, error)
	CreatePost(ctx context.Context, post dto.PostDTO) (*entity.Post, error)
}

type CommentUseCase interface {
	GetTopListComments(ctx context.Context, params dto.CommentRequest) (*dto.CommentResponse, error)
	GetCommentReplies(ctx context.Context, params dto.CommentRequest) (*dto.CommentResponse, error)
	CreateComment(ctx context.Context, comment dto.CommentDTO) (*entity.Comment, error)
}
