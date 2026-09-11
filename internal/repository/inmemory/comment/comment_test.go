package comment

import (
	"context"
	"testing"
	"time"

	"github.com/dubovayaRoshcha/posts-service/internal/entity"
	"github.com/dubovayaRoshcha/posts-service/internal/usecase/dto"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
)

func TestGetTopList(t *testing.T) {
	createdAt := time.Now()

	inputParams := dto.CommentRequest{
		PostID: 10,
		Limit:  10,
		Offset: 0,
	}

	tests := []struct {
		name     string
		comments map[int]entity.Comment
		params   dto.CommentRequest
		want     []entity.Comment
	}{
		{
			name: "OK",
			comments: map[int]entity.Comment{
				1: {ID: 1, PostID: 10, ReplyToCommentID: nil, UserID: 20, Text: "first comment", CreatedAt: createdAt},
				2: {ID: 2, PostID: 10, ReplyToCommentID: nil, UserID: 30, Text: "second comment", CreatedAt: createdAt},
			},
			params: inputParams,
			want: []entity.Comment{
				{ID: 2, PostID: 10, ReplyToCommentID: nil, UserID: 30, Text: "second comment", CreatedAt: createdAt},
				{ID: 1, PostID: 10, ReplyToCommentID: nil, UserID: 20, Text: "first comment", CreatedAt: createdAt},
			},
		},
		{
			name:     "empty_result",
			comments: map[int]entity.Comment{},
			params:   inputParams,
			want:     []entity.Comment{},
		},
		{
			name: "ignores_comments_from_another_post",
			comments: map[int]entity.Comment{
				1: {ID: 1, PostID: 10, ReplyToCommentID: nil, UserID: 10, Text: "first comment", CreatedAt: createdAt},
				2: {ID: 2, PostID: 20, ReplyToCommentID: nil, UserID: 20, Text: "comment from another post", CreatedAt: createdAt},
				3: {ID: 3, PostID: 10, ReplyToCommentID: nil, UserID: 30, Text: "second comment", CreatedAt: createdAt},
			},
			params: inputParams,
			want: []entity.Comment{
				{ID: 3, PostID: 10, ReplyToCommentID: nil, UserID: 30, Text: "second comment", CreatedAt: createdAt},
				{ID: 1, PostID: 10, ReplyToCommentID: nil, UserID: 10, Text: "first comment", CreatedAt: createdAt},
			},
		},
		{
			name: "pagination",
			comments: map[int]entity.Comment{
				1: {ID: 1, PostID: 10, ReplyToCommentID: nil, UserID: 10, Text: "first comment", CreatedAt: createdAt},
				2: {ID: 2, PostID: 10, ReplyToCommentID: nil, UserID: 20, Text: "second comment", CreatedAt: createdAt},
				3: {ID: 3, PostID: 10, ReplyToCommentID: nil, UserID: 30, Text: "third comment", CreatedAt: createdAt},
				4: {ID: 4, PostID: 10, ReplyToCommentID: nil, UserID: 40, Text: "fourth comment", CreatedAt: createdAt},
			},
			params: dto.CommentRequest{
				PostID: 10,
				Limit:  2,
				Offset: 1,
			},
			want: []entity.Comment{
				{ID: 3, PostID: 10, ReplyToCommentID: nil, UserID: 30, Text: "third comment", CreatedAt: createdAt},
				{ID: 2, PostID: 10, ReplyToCommentID: nil, UserID: 20, Text: "second comment", CreatedAt: createdAt},
			},
		},
		{
			name: "offset_out_of_range",
			comments: map[int]entity.Comment{
				1: {ID: 1, PostID: 10, ReplyToCommentID: nil, UserID: 10, Text: "first comment", CreatedAt: createdAt},
				2: {ID: 2, PostID: 10, ReplyToCommentID: nil, UserID: 20, Text: "second comment", CreatedAt: createdAt},
			},
			params: dto.CommentRequest{
				PostID: 10,
				Limit:  10,
				Offset: 5,
			},
			want: []entity.Comment{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo := &CommentRepo{
				comments: test.comments,
			}

			got, err := repo.GetTopList(context.Background(), test.params)

			require.NoError(t, err)
			require.Equal(t, test.want, got)
		})
	}
}

func TestGetReplies(t *testing.T) {
	createdAt := time.Now()
	replyToCommentID := 1

	inputParams := dto.CommentRequest{
		PostID:           10,
		ReplyToCommentID: &replyToCommentID,
		Limit:            10,
		Offset:           0,
	}

	tests := []struct {
		name     string
		comments map[int]entity.Comment
		params   dto.CommentRequest
		want     []entity.Comment
	}{
		{
			name: "OK",
			comments: map[int]entity.Comment{
				1: {ID: 1, PostID: 10, ReplyToCommentID: nil, UserID: 10, Text: "parent comment", CreatedAt: createdAt},
				2: {ID: 2, PostID: 10, ReplyToCommentID: &replyToCommentID, UserID: 20, Text: "first reply", CreatedAt: createdAt},
				3: {ID: 3, PostID: 10, ReplyToCommentID: &replyToCommentID, UserID: 30, Text: "second reply", CreatedAt: createdAt},
			},
			params: inputParams,
			want: []entity.Comment{
				{ID: 3, PostID: 10, ReplyToCommentID: &replyToCommentID, UserID: 30, Text: "second reply", CreatedAt: createdAt},
				{ID: 2, PostID: 10, ReplyToCommentID: &replyToCommentID, UserID: 20, Text: "first reply", CreatedAt: createdAt},
			},
		},
		{
			name:     "empty_result",
			comments: map[int]entity.Comment{},
			params:   inputParams,
			want:     []entity.Comment{},
		},
		{
			name: "returns_only_replies_to_comment",
			comments: map[int]entity.Comment{
				1: {ID: 1, PostID: 10, ReplyToCommentID: nil, UserID: 10, Text: "first: top comment", CreatedAt: createdAt},
				2: {ID: 2, PostID: 10, ReplyToCommentID: &replyToCommentID, UserID: 20, Text: "reply to first comment", CreatedAt: createdAt},
				3: {ID: 3, PostID: 10, ReplyToCommentID: lo.ToPtr(4), UserID: 30, Text: "reply to fourth comment", CreatedAt: createdAt},
				4: {ID: 4, PostID: 10, ReplyToCommentID: nil, UserID: 40, Text: "fourth: top comment", CreatedAt: createdAt},
				5: {ID: 5, PostID: 10, ReplyToCommentID: lo.ToPtr(2), UserID: 50, Text: "reply to second comment", CreatedAt: createdAt},
			},
			params: inputParams,
			want: []entity.Comment{
				{ID: 2, PostID: 10, ReplyToCommentID: &replyToCommentID, UserID: 20, Text: "reply to first comment", CreatedAt: createdAt},
			},
		},
		{
			name: "ignores_replies_from_another_post",
			comments: map[int]entity.Comment{
				1: {ID: 1, PostID: 10, ReplyToCommentID: nil, UserID: 10, Text: "parent comment", CreatedAt: createdAt},
				2: {ID: 2, PostID: 10, ReplyToCommentID: &replyToCommentID, UserID: 20, Text: "valid reply", CreatedAt: createdAt},
				3: {ID: 3, PostID: 20, ReplyToCommentID: &replyToCommentID, UserID: 30, Text: "reply from another post", CreatedAt: createdAt},
			},
			params: inputParams,
			want: []entity.Comment{
				{ID: 2, PostID: 10, ReplyToCommentID: &replyToCommentID, UserID: 20, Text: "valid reply", CreatedAt: createdAt},
			},
		},
		{
			name: "pagination",
			comments: map[int]entity.Comment{
				1: {ID: 1, PostID: 10, ReplyToCommentID: nil, UserID: 10, Text: "parent comment", CreatedAt: createdAt},
				2: {ID: 2, PostID: 10, ReplyToCommentID: &replyToCommentID, UserID: 20, Text: "first reply", CreatedAt: createdAt},
				3: {ID: 3, PostID: 10, ReplyToCommentID: &replyToCommentID, UserID: 30, Text: "second reply", CreatedAt: createdAt},
				4: {ID: 4, PostID: 10, ReplyToCommentID: &replyToCommentID, UserID: 40, Text: "third reply", CreatedAt: createdAt},
				5: {ID: 5, PostID: 10, ReplyToCommentID: &replyToCommentID, UserID: 50, Text: "fourth reply", CreatedAt: createdAt},
			},
			params: dto.CommentRequest{
				PostID:           10,
				ReplyToCommentID: &replyToCommentID,
				Limit:            2,
				Offset:           1,
			},
			want: []entity.Comment{
				{ID: 4, PostID: 10, ReplyToCommentID: &replyToCommentID, UserID: 40, Text: "third reply", CreatedAt: createdAt},
				{ID: 3, PostID: 10, ReplyToCommentID: &replyToCommentID, UserID: 30, Text: "second reply", CreatedAt: createdAt},
			},
		},
		{
			name: "offset_out_of_range",
			comments: map[int]entity.Comment{
				1: {ID: 1, PostID: 10, ReplyToCommentID: nil, UserID: 10, Text: "parent comment", CreatedAt: createdAt},
				2: {ID: 2, PostID: 10, ReplyToCommentID: &replyToCommentID, UserID: 20, Text: "first reply", CreatedAt: createdAt},
				3: {ID: 3, PostID: 10, ReplyToCommentID: &replyToCommentID, UserID: 30, Text: "second reply", CreatedAt: createdAt},
			},
			params: dto.CommentRequest{
				PostID:           10,
				ReplyToCommentID: &replyToCommentID,
				Limit:            10,
				Offset:           5,
			},
			want: []entity.Comment{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo := &CommentRepo{
				comments: test.comments,
			}

			got, err := repo.GetReplies(context.Background(), test.params)

			require.NoError(t, err)
			require.Equal(t, test.want, got)
		})
	}
}

func TestGetByID(t *testing.T) {
	createdAt := time.Now()

	tests := []struct {
		name      string
		comments  map[int]entity.Comment
		commentID int
		want      *entity.Comment
		wantErr   error
	}{
		{
			name: "OK",
			comments: map[int]entity.Comment{
				1: {
					ID:               1,
					PostID:           10,
					ReplyToCommentID: nil,
					UserID:           20,
					Text:             "first comment",
					CreatedAt:        createdAt,
				},
			},
			commentID: 1,
			want: &entity.Comment{
				ID:               1,
				PostID:           10,
				ReplyToCommentID: nil,
				UserID:           20,
				Text:             "first comment",
				CreatedAt:        createdAt,
			},
			wantErr: nil,
		},
		{
			name:      "comment_not_found",
			comments:  map[int]entity.Comment{},
			commentID: 1,
			want:      nil,
			wantErr:   entity.CommentNotFound,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo := &CommentRepo{
				comments: test.comments,
			}

			got, err := repo.GetByID(context.Background(), test.commentID)
			if test.wantErr != nil {
				require.ErrorIs(t, err, test.wantErr)
				return
			}

			require.NoError(t, err)
			require.Equal(t, test.want, got)
		})
	}
}

func TestCreate(t *testing.T) {
	replyToCommentID := 1

	inputComment := dto.CommentDTO{
		PostID:           10,
		ReplyToCommentID: &replyToCommentID,
		UserID:           20,
		Text:             "first comment",
	}

	tests := []struct {
		name    string
		comment dto.CommentDTO
		want    *entity.Comment
	}{
		{
			name:    "OK",
			comment: inputComment,
			want: &entity.Comment{
				ID:               1,
				PostID:           10,
				ReplyToCommentID: &replyToCommentID,
				UserID:           20,
				Text:             "first comment",
				CreatedAt:        time.Time{},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo := NewCommentRepo()

			got, err := repo.Create(context.Background(), test.comment)

			require.NoError(t, err)
			got.CreatedAt = time.Time{}
			require.Equal(t, test.want, got)
		})
	}
}
