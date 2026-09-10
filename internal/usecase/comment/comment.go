package comment

import (
	"context"

	"github.com/dubovayaRoshcha/posts-service/internal/entity"
	"github.com/dubovayaRoshcha/posts-service/internal/usecase"
	"github.com/dubovayaRoshcha/posts-service/internal/usecase/dto"
	"github.com/dubovayaRoshcha/posts-service/internal/usecase/utils/validator"
)

type CommentUseCase struct {
	commentRepo usecase.CommentRepo
	postRepo    usecase.PostRepo
}

func NewCommentUseCase(commentRepo usecase.CommentRepo, postRepo usecase.PostRepo) *CommentUseCase {
	return &CommentUseCase{
		commentRepo: commentRepo,
		postRepo:    postRepo,
	}
}

func (uc *CommentUseCase) GetTopListComments(ctx context.Context, params dto.CommentRequest) (*dto.CommentResponse, error) {
	err := validator.ValidateTopCommentsParams(&params)
	if err != nil {
		return nil, err
	}

	comments, err := uc.commentRepo.GetTopList(ctx, params)
	if err != nil {
		return nil, err
	}

	return &dto.CommentResponse{Comments: comments, Len: len(comments)}, nil
}

func (uc *CommentUseCase) GetCommentReplies(ctx context.Context, params dto.CommentRequest) (*dto.CommentResponse, error) {
	err := validator.ValidateRepliesParams(&params)
	if err != nil {
		return nil, err
	}

	comments, err := uc.commentRepo.GetReplies(ctx, params)
	if err != nil {
		return nil, err
	}

	return &dto.CommentResponse{Comments: comments, Len: len(comments)}, nil
}

func (uc *CommentUseCase) CreateComment(ctx context.Context, comment dto.CommentDTO) (*entity.Comment, error) {
	err := validator.ValidateComment(&comment)
	if err != nil {
		return nil, err
	}

	post, err := uc.postRepo.GetByID(ctx, comment.PostID)
	if err != nil {
		return nil, err
	}

	if !post.CommentsAllowed {
		return nil, entity.CommentsNotAllowed
	}

	commentResp, err := uc.commentRepo.Create(ctx, comment)
	if err != nil {
		return nil, err
	}

	return commentResp, nil
}
