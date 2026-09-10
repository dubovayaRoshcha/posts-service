package post

import (
	"context"
	"sync"
	"time"

	"github.com/dubovayaRoshcha/posts-service/internal/entity"
	"github.com/dubovayaRoshcha/posts-service/internal/usecase/dto"
)

type PostRepo struct {
	mu    sync.RWMutex
	posts map[int]entity.Post
}

func NewPostRepo() *PostRepo {
	return &PostRepo{
		posts: make(map[int]entity.Post),
	}
}

func (r *PostRepo) GetList(ctx context.Context, params dto.PostRequest) ([]entity.Post, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	posts := make([]entity.Post, 0, params.Limit)

	start := params.Offset
	end := params.Offset + params.Limit

	skip := 0

	for id := len(r.posts); id >= 1; id-- {
		post, ok := r.posts[id]
		if !ok {
			continue
		}

		if skip < start {
			skip++
			continue
		}

		if skip >= end {
			break
		}

		posts = append(posts, post)
		skip++
	}

	return posts, nil
}

func (r *PostRepo) GetByID(ctx context.Context, postID int) (*entity.Post, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	post, ok := r.posts[postID]
	if !ok {
		return nil, entity.PostNotFound
	}

	return &post, nil
}

func (r *PostRepo) Create(ctx context.Context, post dto.PostDTO) (*entity.Post, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	id := len(r.posts) + 1

	createdPost := entity.Post{
		ID:              id,
		UserID:          post.UserID,
		Title:           post.Title,
		Description:     post.Description,
		CommentsAllowed: post.CommentsAllowed,
		CreatedAt:       time.Now(),
	}

	r.posts[id] = createdPost

	return &createdPost, nil
}
