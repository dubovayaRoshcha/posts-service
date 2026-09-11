package graphql

import (
	"context"
	"testing"

	"github.com/dubovayaRoshcha/posts-service/internal/entity"
	"github.com/dubovayaRoshcha/posts-service/internal/mocks"
	"github.com/dubovayaRoshcha/posts-service/internal/usecase/dto"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestComments(t *testing.T) {
	ctx := context.Background()

	existingComments := []entity.Comment{
		{ID: 1, PostID: 10, ReplyToCommentID: nil, UserID: 10, Text: "first comment"},
		{ID: 2, PostID: 10, ReplyToCommentID: nil, UserID: 20, Text: "second comment"},
	}

	commentResponse := &dto.CommentResponse{
		Comments: existingComments,
		Len:      len(existingComments),
	}

	expectedResponse := &CommentResponse{
		Comments: []*Comment{
			{ID: "1", PostID: 10, ReplyToCommentID: nil, UserID: 10, Text: "first comment"},
			{ID: "2", PostID: 10, ReplyToCommentID: nil, UserID: 20, Text: "second comment"},
		},
		Len: 2,
	}

	tests := []struct {
		name      string
		params    CommentRequest
		setupMock func(m *mocks.MockCommentUseCase)
		want      *CommentResponse
		wantErr   error
	}{
		{
			name: "OK",
			params: CommentRequest{
				PostID: 10,
				Limit:  10,
				Offset: 0,
			},
			setupMock: func(m *mocks.MockCommentUseCase) {
				m.EXPECT().GetTopListComments(ctx, dto.CommentRequest{
					PostID: 10,
					Limit:  10,
					Offset: 0,
				}).
					Return(commentResponse, nil).Times(1)
			},
			want:    expectedResponse,
			wantErr: nil,
		},
		{
			name: "empty_result",
			params: CommentRequest{
				PostID: 10,
				Limit:  10,
				Offset: 0,
			},
			setupMock: func(m *mocks.MockCommentUseCase) {
				m.EXPECT().GetTopListComments(ctx, dto.CommentRequest{
					PostID: 10,
					Limit:  10,
					Offset: 0,
				}).
					Return(&dto.CommentResponse{
						Comments: []entity.Comment{},
						Len:      0,
					}, nil).Times(1)
			},
			want: &CommentResponse{
				Len:      0,
				Comments: []*Comment{},
			},
			wantErr: nil,
		},
		{
			name: "error_from_usecase",
			params: CommentRequest{
				PostID: 10,
				Limit:  10,
				Offset: 0,
			},
			setupMock: func(m *mocks.MockCommentUseCase) {
				m.EXPECT().GetTopListComments(ctx, dto.CommentRequest{
					PostID: 10,
					Limit:  10,
					Offset: 0,
				}).
					Return(nil, entity.InvalidInput).Times(1)
			},
			want:    nil,
			wantErr: entity.InvalidInput,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockPostUC := mocks.NewMockPostUseCase(ctrl)
			mockCommentUC := mocks.NewMockCommentUseCase(ctrl)

			if test.setupMock != nil {
				test.setupMock(mockCommentUC)
			}

			resolver := NewResolver(mockPostUC, mockCommentUC)
			query := &queryResolver{resolver}

			got, err := query.Comments(ctx, test.params)

			if test.wantErr != nil {
				require.ErrorIs(t, err, test.wantErr)
				return
			}

			require.NoError(t, err)
			require.Equal(t, test.want, got)
		})
	}
}

func TestReplies(t *testing.T) {
	ctx := context.Background()

	replyToCommentID := 1

	existingComments := []entity.Comment{
		{ID: 2, PostID: 10, ReplyToCommentID: &replyToCommentID, UserID: 20, Text: "first reply"},
		{ID: 3, PostID: 10, ReplyToCommentID: &replyToCommentID, UserID: 30, Text: "second reply"},
	}

	commentResponse := &dto.CommentResponse{
		Comments: existingComments,
		Len:      len(existingComments),
	}

	expectedResponse := &CommentResponse{
		Comments: []*Comment{
			{ID: "2", PostID: 10, ReplyToCommentID: &replyToCommentID, UserID: 20, Text: "first reply"},
			{ID: "3", PostID: 10, ReplyToCommentID: &replyToCommentID, UserID: 30, Text: "second reply"},
		},
		Len: 2,
	}

	tests := []struct {
		name      string
		params    CommentRepliesRequest
		setupMock func(m *mocks.MockCommentUseCase)
		want      *CommentResponse
		wantErr   error
	}{
		{
			name: "OK",
			params: CommentRepliesRequest{
				PostID:           10,
				ReplyToCommentID: replyToCommentID,
				Limit:            10,
				Offset:           0,
			},
			setupMock: func(m *mocks.MockCommentUseCase) {
				m.EXPECT().GetCommentReplies(ctx, dto.CommentRequest{
					PostID:           10,
					ReplyToCommentID: &replyToCommentID,
					Limit:            10,
					Offset:           0,
				}).
					Return(commentResponse, nil).Times(1)
			},
			want:    expectedResponse,
			wantErr: nil,
		},
		{
			name: "empty_result",
			params: CommentRepliesRequest{
				PostID:           10,
				ReplyToCommentID: replyToCommentID,
				Limit:            10,
				Offset:           0,
			},
			setupMock: func(m *mocks.MockCommentUseCase) {
				m.EXPECT().GetCommentReplies(ctx, dto.CommentRequest{
					PostID:           10,
					ReplyToCommentID: &replyToCommentID,
					Limit:            10,
					Offset:           0,
				}).
					Return(&dto.CommentResponse{
						Comments: []entity.Comment{},
						Len:      0,
					}, nil).Times(1)
			},
			want: &CommentResponse{
				Len:      0,
				Comments: []*Comment{},
			},
			wantErr: nil,
		},
		{
			name: "usecase_invalid_input",
			params: CommentRepliesRequest{
				PostID:           10,
				ReplyToCommentID: replyToCommentID,
				Limit:            10,
				Offset:           0,
			},
			setupMock: func(m *mocks.MockCommentUseCase) {
				m.EXPECT().GetCommentReplies(ctx, dto.CommentRequest{
					PostID:           10,
					ReplyToCommentID: &replyToCommentID,
					Limit:            10,
					Offset:           0,
				}).
					Return(nil, entity.InvalidInput).Times(1)
			},
			want:    nil,
			wantErr: entity.InvalidInput,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockPostUC := mocks.NewMockPostUseCase(ctrl)
			mockCommentUC := mocks.NewMockCommentUseCase(ctrl)

			if test.setupMock != nil {
				test.setupMock(mockCommentUC)
			}

			resolver := NewResolver(mockPostUC, mockCommentUC)
			query := &queryResolver{resolver}

			got, err := query.Replies(ctx, test.params)

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

	inputComment := CreateCommentRequest{
		PostID:           10,
		ReplyToCommentID: &replyToCommentID,
		UserID:           20,
		Text:             "new comment",
	}

	existingComment := &entity.Comment{
		ID:               1,
		PostID:           10,
		ReplyToCommentID: &replyToCommentID,
		UserID:           20,
		Text:             "new comment",
	}

	expectedComment := &Comment{
		ID:               "1",
		PostID:           10,
		ReplyToCommentID: &replyToCommentID,
		UserID:           20,
		Text:             "new comment",
	}

	tests := []struct {
		name      string
		params    CreateCommentRequest
		setupMock func(m *mocks.MockCommentUseCase)
		want      *Comment
		wantErr   error
	}{
		{
			name:   "OK",
			params: inputComment,
			setupMock: func(m *mocks.MockCommentUseCase) {
				m.EXPECT().CreateComment(ctx, dto.CommentDTO{
					PostID:           10,
					ReplyToCommentID: &replyToCommentID,
					UserID:           20,
					Text:             "new comment",
				}).
					Return(existingComment, nil).Times(1)
			},
			want:    expectedComment,
			wantErr: nil,
		},
		{
			name:   "usecase_comments_not_allowed",
			params: inputComment,
			setupMock: func(m *mocks.MockCommentUseCase) {
				m.EXPECT().CreateComment(ctx, dto.CommentDTO{
					PostID:           10,
					ReplyToCommentID: &replyToCommentID,
					UserID:           20,
					Text:             "new comment",
				}).
					Return(nil, entity.CommentsNotAllowed).Times(1)
			},
			want:    nil,
			wantErr: entity.CommentsNotAllowed,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockPostUC := mocks.NewMockPostUseCase(ctrl)
			mockCommentUC := mocks.NewMockCommentUseCase(ctrl)

			if test.setupMock != nil {
				test.setupMock(mockCommentUC)
			}

			resolver := NewResolver(mockPostUC, mockCommentUC)
			mutation := &mutationResolver{resolver}

			got, err := mutation.CreateComment(ctx, test.params)

			if test.wantErr != nil {
				require.ErrorIs(t, err, test.wantErr)
				return
			}

			require.NoError(t, err)
			require.Equal(t, test.want, got)
		})
	}
}

func TestCommentAdded(t *testing.T) {
	tests := []struct {
		name   string
		postID int
	}{
		{
			name:   "OK",
			postID: 10,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockPostUC := mocks.NewMockPostUseCase(ctrl)
			mockCommentUC := mocks.NewMockCommentUseCase(ctrl)

			resolver := NewResolver(mockPostUC, mockCommentUC)
			subscription := &subscriptionResolver{resolver}

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			got, err := subscription.CommentAdded(ctx, test.postID)

			require.NoError(t, err)
			require.NotNil(t, got)
			require.Len(t, resolver.channelList[test.postID], 1)
		})
	}
}
