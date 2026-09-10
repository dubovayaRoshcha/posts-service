package comment

import (
	"context"
	"strings"
	"testing"

	"github.com/dubovayaRoshcha/posts-service/internal/entity"
	"github.com/dubovayaRoshcha/posts-service/internal/mocks"
	"github.com/dubovayaRoshcha/posts-service/internal/usecase/dto"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

const (
	maxCommentLimit = 50
	maxTextLength   = 2000
)

func TestGetTopListComments(t *testing.T) {
	ctx := context.Background()

	existingComments := []entity.Comment{
		{ID: 1, PostID: 10, ReplyToCommentID: nil, UserID: 10, Text: "first comment"},
		{ID: 2, PostID: 10, ReplyToCommentID: nil, UserID: 20, Text: "second comment"},
	}

	expectedResponse := &dto.CommentResponse{
		Comments: existingComments,
		Len:      len(existingComments),
	}

	tests := []struct {
		name      string
		params    dto.CommentRequest
		setupMock func(m *mocks.MockCommentRepo)
		want      *dto.CommentResponse
		wantErr   error
	}{
		{
			name: "OK",
			params: dto.CommentRequest{
				PostID: 10,
				Limit:  10,
				Offset: 0,
			},
			setupMock: func(m *mocks.MockCommentRepo) {
				m.EXPECT().GetTopList(ctx, dto.CommentRequest{
					PostID: 10,
					Limit:  10,
					Offset: 0,
				}).
					Return(existingComments, nil).Times(1)
			},
			want:    expectedResponse,
			wantErr: nil,
		},
		{
			name: "invalid_limit",
			params: dto.CommentRequest{
				PostID: 10,
				Limit:  -1,
				Offset: 0,
			},
			setupMock: nil,
			want:      nil,
			wantErr:   entity.InvalidInput,
		},
		{
			name: "invalid_offset",
			params: dto.CommentRequest{
				PostID: 10,
				Limit:  10,
				Offset: -1,
			},
			setupMock: nil,
			want:      nil,
			wantErr:   entity.InvalidInput,
		},
		{
			name: "invalid_post_id",
			params: dto.CommentRequest{
				PostID: 0,
				Limit:  10,
				Offset: 0,
			},
			setupMock: nil,
			want:      nil,
			wantErr:   entity.InvalidInput,
		},
		{
			name: "limit_too_large",
			params: dto.CommentRequest{
				PostID: 10,
				Limit:  maxCommentLimit + 1,
				Offset: 0,
			},
			setupMock: func(m *mocks.MockCommentRepo) {
				m.EXPECT().GetTopList(ctx, dto.CommentRequest{
					PostID: 10,
					Limit:  maxCommentLimit,
					Offset: 0,
				}).
					Return(existingComments, nil).Times(1)
			},
			want:    expectedResponse,
			wantErr: nil,
		},
		{
			name: "error_from_repo",
			params: dto.CommentRequest{
				PostID: 10,
				Limit:  10,
				Offset: 0,
			},
			setupMock: func(m *mocks.MockCommentRepo) {
				m.EXPECT().GetTopList(ctx, dto.CommentRequest{
					PostID: 10,
					Limit:  10,
					Offset: 0,
				}).
					Return(nil, entity.CommentNotFound).Times(1)
			},
			want:    nil,
			wantErr: entity.CommentNotFound,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockCommentRepo := mocks.NewMockCommentRepo(ctrl)
			mockPostRepo := mocks.NewMockPostRepo(ctrl)

			if test.setupMock != nil {
				test.setupMock(mockCommentRepo)
			}

			uc := NewCommentUseCase(mockCommentRepo, mockPostRepo)

			got, err := uc.GetTopListComments(ctx, test.params)

			if test.wantErr != nil {
				require.ErrorIs(t, err, test.wantErr)
				return
			}

			require.NoError(t, err)
			require.Equal(t, test.want, got)
		})
	}
}

func TestGetCommentReplies(t *testing.T) {
	ctx := context.Background()

	replyToCommentID := 1

	existingComments := []entity.Comment{
		{ID: 2, PostID: 10, ReplyToCommentID: &replyToCommentID, UserID: 20, Text: "first reply"},
		{ID: 3, PostID: 10, ReplyToCommentID: &replyToCommentID, UserID: 30, Text: "second reply"},
	}

	expectedResponse := &dto.CommentResponse{
		Comments: existingComments,
		Len:      len(existingComments),
	}

	tests := []struct {
		name      string
		params    dto.CommentRequest
		setupMock func(m *mocks.MockCommentRepo)
		want      *dto.CommentResponse
		wantErr   error
	}{
		{
			name: "OK",
			params: dto.CommentRequest{
				PostID:           10,
				ReplyToCommentID: &replyToCommentID,
				Limit:            10,
				Offset:           0,
			},
			setupMock: func(m *mocks.MockCommentRepo) {
				m.EXPECT().GetReplies(ctx, dto.CommentRequest{
					PostID:           10,
					ReplyToCommentID: &replyToCommentID,
					Limit:            10,
					Offset:           0,
				}).
					Return(existingComments, nil).Times(1)
			},
			want:    expectedResponse,
			wantErr: nil,
		},
		{
			name: "invalid_limit",
			params: dto.CommentRequest{
				PostID:           10,
				ReplyToCommentID: &replyToCommentID,
				Limit:            -1,
				Offset:           0,
			},
			setupMock: nil,
			want:      nil,
			wantErr:   entity.InvalidInput,
		},
		{
			name: "invalid_offset",
			params: dto.CommentRequest{
				PostID:           10,
				ReplyToCommentID: &replyToCommentID,
				Limit:            10,
				Offset:           -1,
			},
			setupMock: nil,
			want:      nil,
			wantErr:   entity.InvalidInput,
		},
		{
			name: "invalid_post_id",
			params: dto.CommentRequest{
				PostID:           0,
				ReplyToCommentID: &replyToCommentID,
				Limit:            10,
				Offset:           0,
			},
			setupMock: nil,
			want:      nil,
			wantErr:   entity.InvalidInput,
		},
		{
			name: "invalid_reply_to_comment_id",
			params: dto.CommentRequest{
				PostID:           10,
				ReplyToCommentID: nil,
				Limit:            10,
				Offset:           0,
			},
			setupMock: nil,
			want:      nil,
			wantErr:   entity.InvalidInput,
		},
		{
			name: "limit_too_large",
			params: dto.CommentRequest{
				PostID:           10,
				ReplyToCommentID: &replyToCommentID,
				Limit:            maxCommentLimit + 1,
				Offset:           0,
			},
			setupMock: func(m *mocks.MockCommentRepo) {
				m.EXPECT().GetReplies(ctx, dto.CommentRequest{
					PostID:           10,
					ReplyToCommentID: &replyToCommentID,
					Limit:            maxCommentLimit,
					Offset:           0,
				}).
					Return(existingComments, nil).Times(1)
			},
			want:    expectedResponse,
			wantErr: nil,
		},
		{
			name: "error_from_repo",
			params: dto.CommentRequest{
				PostID:           10,
				ReplyToCommentID: &replyToCommentID,
				Limit:            10,
				Offset:           0,
			},
			setupMock: func(m *mocks.MockCommentRepo) {
				m.EXPECT().GetReplies(ctx, dto.CommentRequest{
					PostID:           10,
					ReplyToCommentID: &replyToCommentID,
					Limit:            10,
					Offset:           0,
				}).
					Return(nil, entity.CommentNotFound).Times(1)
			},
			want:    nil,
			wantErr: entity.CommentNotFound,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockCommentRepo := mocks.NewMockCommentRepo(ctrl)
			mockPostRepo := mocks.NewMockPostRepo(ctrl)

			if test.setupMock != nil {
				test.setupMock(mockCommentRepo)
			}

			uc := NewCommentUseCase(mockCommentRepo, mockPostRepo)

			got, err := uc.GetCommentReplies(ctx, test.params)

			if test.wantErr != nil {
				require.ErrorIs(t, err, test.wantErr)
				return
			}

			require.NoError(t, err)
			require.Equal(t, test.want, got)
		})
	}
}

func TestCreateComment(t *testing.T) {
	ctx := context.Background()

	replyToCommentID := 1

	inputComment := dto.CommentDTO{
		PostID:           10,
		ReplyToCommentID: &replyToCommentID,
		UserID:           20,
		Text:             "new comment",
	}

	existingPost := &entity.Post{
		ID:              10,
		UserID:          10,
		Title:           "first post",
		Description:     "first description",
		CommentsAllowed: true,
	}

	expectedComment := &entity.Comment{
		ID:               1,
		PostID:           10,
		ReplyToCommentID: &replyToCommentID,
		UserID:           20,
		Text:             "new comment",
	}

	tests := []struct {
		name      string
		comment   dto.CommentDTO
		setupMock func(commentMock *mocks.MockCommentRepo, postMock *mocks.MockPostRepo)
		want      *entity.Comment
		wantErr   error
	}{
		{
			name:    "OK",
			comment: inputComment,
			setupMock: func(commentMock *mocks.MockCommentRepo, postMock *mocks.MockPostRepo) {
				postMock.EXPECT().GetByID(ctx, 10).
					Return(existingPost, nil).Times(1)

				commentMock.EXPECT().Create(ctx, inputComment).
					Return(expectedComment, nil).Times(1)
			},
			want:    expectedComment,
			wantErr: nil,
		},
		{
			name: "empty_text",
			comment: dto.CommentDTO{
				PostID:           10,
				ReplyToCommentID: &replyToCommentID,
				UserID:           20,
				Text:             "",
			},
			setupMock: nil,
			want:      nil,
			wantErr:   entity.InvalidInput,
		},
		{
			name: "invalid_post_id",
			comment: dto.CommentDTO{
				PostID:           0,
				ReplyToCommentID: &replyToCommentID,
				UserID:           20,
				Text:             "new comment",
			},
			setupMock: nil,
			want:      nil,
			wantErr:   entity.InvalidInput,
		},
		{
			name: "text_too_long",
			comment: dto.CommentDTO{
				PostID:           10,
				ReplyToCommentID: &replyToCommentID,
				UserID:           20,
				Text:             strings.Repeat("a", maxTextLength+1),
			},
			setupMock: nil,
			want:      nil,
			wantErr:   entity.MaxLengthExceeded,
		},
		{
			name:    "error_from_post_repo",
			comment: inputComment,
			setupMock: func(commentMock *mocks.MockCommentRepo, postMock *mocks.MockPostRepo) {
				postMock.EXPECT().GetByID(ctx, 10).
					Return(nil, entity.PostNotFound).Times(1)
			},
			want:    nil,
			wantErr: entity.PostNotFound,
		},
		{
			name:    "comments_not_allowed",
			comment: inputComment,
			setupMock: func(commentMock *mocks.MockCommentRepo, postMock *mocks.MockPostRepo) {
				postMock.EXPECT().GetByID(ctx, 10).
					Return(&entity.Post{
						ID:              10,
						UserID:          10,
						Title:           "first post",
						Description:     "first description",
						CommentsAllowed: false,
					}, nil).Times(1)
			},
			want:    nil,
			wantErr: entity.CommentsNotAllowed,
		},
		{
			name:    "error_from_comment_repo",
			comment: inputComment,
			setupMock: func(commentMock *mocks.MockCommentRepo, postMock *mocks.MockPostRepo) {
				postMock.EXPECT().GetByID(ctx, 10).
					Return(existingPost, nil).Times(1)

				commentMock.EXPECT().Create(ctx, inputComment).
					Return(nil, entity.CommentNotFound).Times(1)
			},
			want:    nil,
			wantErr: entity.CommentNotFound,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockCommentRepo := mocks.NewMockCommentRepo(ctrl)
			mockPostRepo := mocks.NewMockPostRepo(ctrl)

			if test.setupMock != nil {
				test.setupMock(mockCommentRepo, mockPostRepo)
			}

			uc := NewCommentUseCase(mockCommentRepo, mockPostRepo)

			got, err := uc.CreateComment(ctx, test.comment)

			if test.wantErr != nil {
				require.ErrorIs(t, err, test.wantErr)
				return
			}

			require.NoError(t, err)
			require.Equal(t, test.want, got)
		})
	}
}
