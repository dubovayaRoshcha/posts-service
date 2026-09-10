package post

import (
	"context"
	"fmt"

	"github.com/dubovayaRoshcha/posts-service/internal/entity"
	"github.com/dubovayaRoshcha/posts-service/internal/usecase"
	"github.com/dubovayaRoshcha/posts-service/internal/usecase/dto"
	"github.com/dubovayaRoshcha/posts-service/internal/usecase/utils/validator"
)

type PostUseCase struct {
	repo usecase.PostRepo
}

func NewPostUseCase(repo usecase.PostRepo) *PostUseCase {
	return &PostUseCase{
		repo: repo,
	}
}

func (uc *PostUseCase) GetListPosts(ctx context.Context, params dto.PostRequest) (*dto.PostResponse, error) {
	err := validator.ValidatePostParams(&params)
	if err != nil {
		return nil, fmt.Errorf("validator.ValidatePostParams: %w", err)
	}

	posts, err := uc.repo.GetList(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("uc.repo.GetList: %w", err)
	}

	return &dto.PostResponse{Posts: posts, Len: len(posts)}, nil
}

func (uc *PostUseCase) GetPostByID(ctx context.Context, postID int) (*entity.Post, error) {
	err := validator.ValidateID(postID)
	if err != nil {
		return nil, fmt.Errorf("validator.ValidateID: %w", err)
	}

	post, err := uc.repo.GetByID(ctx, postID)
	if err != nil {
		return nil, fmt.Errorf("uc.repo.GetByID: %w", err)
	}

	return post, nil
}

func (uc *PostUseCase) CreatePost(ctx context.Context, post dto.PostDTO) (*entity.Post, error) {
	err := validator.ValidatePost(&post)
	if err != nil {
		return nil, fmt.Errorf("validator.ValidatePost: %w", err)
	}

	postResp, err := uc.repo.Create(ctx, post)
	if err != nil {
		return nil, fmt.Errorf("uc.repo.Create: %w", err)
	}

	return postResp, nil
}
