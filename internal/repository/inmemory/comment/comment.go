package comment

import (
	"context"
	"sync"
	"time"

	"github.com/dubovayaRoshcha/posts-service/internal/entity"
	"github.com/dubovayaRoshcha/posts-service/internal/usecase/dto"
)

type CommentRepo struct {
	mu       sync.RWMutex
	comments map[int]entity.Comment
}

func NewCommentRepo() *CommentRepo {
	return &CommentRepo{
		comments: make(map[int]entity.Comment),
	}
}

func (r *CommentRepo) GetTopList(ctx context.Context, params dto.CommentRequest) ([]entity.Comment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	comments := make([]entity.Comment, 0, params.Limit)

	start := params.Offset
	end := params.Offset + params.Limit

	skip := 0

	for id := len(r.comments); id >= 1; id-- {
		comment, ok := r.comments[id]
		if !ok {
			continue
		}

		if comment.PostID != params.PostID || comment.ReplyToCommentID != nil {
			continue
		}

		if skip < start {
			skip++
			continue
		}

		if skip >= end {
			break
		}

		comments = append(comments, comment)
		skip++
	}

	return comments, nil
}

func (r *CommentRepo) GetReplies(ctx context.Context, params dto.CommentRequest) ([]entity.Comment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	comments := make([]entity.Comment, 0, params.Limit)

	start := params.Offset
	end := params.Offset + params.Limit

	skip := 0

	for id := len(r.comments); id >= 1; id-- {
		comment, ok := r.comments[id]
		if !ok {
			continue
		}

		if comment.PostID != params.PostID || comment.ReplyToCommentID == nil {
			continue
		}

		if *comment.ReplyToCommentID != *params.ReplyToCommentID {
			continue
		}

		if skip < start {
			skip++
			continue
		}

		if skip >= end {
			break
		}

		comments = append(comments, comment)
		skip++
	}

	return comments, nil
}

func (r *CommentRepo) GetByID(ctx context.Context, commentID int) (*entity.Comment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	comment, ok := r.comments[commentID]
	if !ok {
		return nil, entity.CommentNotFound
	}

	return &comment, nil
}

func (r *CommentRepo) Create(ctx context.Context, comment dto.CommentDTO) (*entity.Comment, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	id := len(r.comments) + 1

	createdComment := entity.Comment{
		ID:               id,
		PostID:           comment.PostID,
		ReplyToCommentID: comment.ReplyToCommentID,
		UserID:           comment.UserID,
		Text:             comment.Text,
		CreatedAt:        time.Now(),
	}

	r.comments[id] = createdComment

	return &createdComment, nil
}
