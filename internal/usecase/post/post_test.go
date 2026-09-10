package post

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
	maxPostLimit         = 50
	maxTitleLength       = 200
	maxDescriptionLength = 10000
)

func TestGetListPosts(t *testing.T) {
	ctx := context.Background()

	existingPosts := []entity.Post{
		{ID: 1, UserID: 10, Title: "first post", Description: "first description", CommentsAllowed: true},
		{ID: 2, UserID: 20, Title: "second post", Description: "second description", CommentsAllowed: false},
	}

	expectedResponse := &dto.PostResponse{
		Posts: existingPosts,
		Len:   len(existingPosts),
	}

	tests := []struct {
		name      string
		params    dto.PostRequest
		setupMock func(m *mocks.MockPostRepo)
		want      *dto.PostResponse
		wantErr   error
	}{
		{
			name: "OK",
			params: dto.PostRequest{
				Limit:  10,
				Offset: 0,
			},
			setupMock: func(m *mocks.MockPostRepo) {
				m.EXPECT().GetList(ctx, dto.PostRequest{
					Limit:  10,
					Offset: 0,
				}).
					Return(existingPosts, nil).Times(1)
			},
			want:    expectedResponse,
			wantErr: nil,
		},
		{
			name: "invalid_limit",
			params: dto.PostRequest{
				Limit:  -1,
				Offset: 0,
			},
			setupMock: nil,
			want:      nil,
			wantErr:   entity.InvalidInput,
		},
		{
			name: "invalid_offset",
			params: dto.PostRequest{
				Limit:  10,
				Offset: -1,
			},
			setupMock: nil,
			want:      nil,
			wantErr:   entity.InvalidInput,
		},
		{
			name: "limit_too_large",
			params: dto.PostRequest{
				Limit:  maxPostLimit + 1,
				Offset: 0,
			},
			setupMock: func(m *mocks.MockPostRepo) {
				m.EXPECT().GetList(ctx, dto.PostRequest{
					Limit:  maxPostLimit,
					Offset: 0,
				}).
					Return(existingPosts, nil).Times(1)
			},
			want:    expectedResponse,
			wantErr: nil,
		},
		{
			name: "error_from_repo",
			params: dto.PostRequest{
				Limit:  10,
				Offset: 0,
			},
			setupMock: func(m *mocks.MockPostRepo) {
				m.EXPECT().GetList(ctx, dto.PostRequest{
					Limit:  10,
					Offset: 0,
				}).
					Return(nil, entity.PostNotFound).Times(1)
			},
			want:    nil,
			wantErr: entity.PostNotFound,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mocks.NewMockPostRepo(ctrl)

			if test.setupMock != nil {
				test.setupMock(mockRepo)
			}

			uc := NewPostUseCase(mockRepo)

			got, err := uc.GetListPosts(ctx, test.params)
			if test.wantErr != nil {
				require.ErrorIs(t, err, test.wantErr)
				return
			}

			require.NoError(t, err)
			require.Equal(t, test.want, got)
		})
	}
}

func TestGetPostByID(t *testing.T) {
	ctx := context.Background()

	existingPost := &entity.Post{
		ID:              1,
		UserID:          10,
		Title:           "first post",
		Description:     "first description",
		CommentsAllowed: true,
	}

	tests := []struct {
		name      string
		postID    int
		setupMock func(m *mocks.MockPostRepo)
		want      *entity.Post
		wantErr   error
	}{
		{
			name:   "OK",
			postID: 1,
			setupMock: func(m *mocks.MockPostRepo) {
				m.EXPECT().GetByID(ctx, 1).
					Return(existingPost, nil).Times(1)
			},
			want:    existingPost,
			wantErr: nil,
		},
		{
			name:      "invalid_id",
			postID:    0,
			setupMock: nil,
			want:      nil,
			wantErr:   entity.InvalidInput,
		},
		{
			name:   "error_from_repo",
			postID: 1,
			setupMock: func(m *mocks.MockPostRepo) {
				m.EXPECT().GetByID(ctx, 1).
					Return(nil, entity.PostNotFound).Times(1)
			},
			want:    nil,
			wantErr: entity.PostNotFound,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mocks.NewMockPostRepo(ctrl)

			if test.setupMock != nil {
				test.setupMock(mockRepo)
			}

			uc := NewPostUseCase(mockRepo)

			got, err := uc.GetPostByID(ctx, test.postID)
			if test.wantErr != nil {
				require.ErrorIs(t, err, test.wantErr)
				return
			}

			require.NoError(t, err)
			require.Equal(t, test.want, got)
		})
	}
}

func TestCreatePost(t *testing.T) {
	ctx := context.Background()

	inputPost := dto.PostDTO{
		UserID:          10,
		Title:           "first post",
		Description:     "first description",
		CommentsAllowed: true,
	}

	expectedPost := &entity.Post{
		ID:              1,
		UserID:          10,
		Title:           "first post",
		Description:     "first description",
		CommentsAllowed: true,
	}

	tests := []struct {
		name      string
		post      dto.PostDTO
		setupMock func(m *mocks.MockPostRepo)
		want      *entity.Post
		wantErr   error
	}{
		{
			name: "OK",
			post: inputPost,
			setupMock: func(m *mocks.MockPostRepo) {
				m.EXPECT().Create(ctx, inputPost).
					Return(expectedPost, nil).Times(1)
			},
			want:    expectedPost,
			wantErr: nil,
		},
		{
			name: "empty_title",
			post: dto.PostDTO{
				UserID:          10,
				Title:           "",
				Description:     "first description",
				CommentsAllowed: true,
			},
			setupMock: nil,
			want:      nil,
			wantErr:   entity.InvalidInput,
		},
		{
			name: "title_too_long",
			post: dto.PostDTO{
				UserID:          10,
				Title:           strings.Repeat("a", maxTitleLength+1),
				Description:     "first description",
				CommentsAllowed: true,
			},
			setupMock: nil,
			want:      nil,
			wantErr:   entity.MaxLengthExceeded,
		},
		{
			name: "description_too_long",
			post: dto.PostDTO{
				UserID:          10,
				Title:           "new post",
				Description:     strings.Repeat("a", maxDescriptionLength+1),
				CommentsAllowed: true,
			},
			setupMock: nil,
			want:      nil,
			wantErr:   entity.MaxLengthExceeded,
		},
		{
			name: "error_from_repo",
			post: inputPost,
			setupMock: func(m *mocks.MockPostRepo) {
				m.EXPECT().Create(ctx, inputPost).
					Return(nil, entity.PostNotFound).Times(1)
			},
			want:    nil,
			wantErr: entity.PostNotFound,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mocks.NewMockPostRepo(ctrl)

			if test.setupMock != nil {
				test.setupMock(mockRepo)
			}

			uc := NewPostUseCase(mockRepo)

			got, err := uc.CreatePost(ctx, test.post)
			if test.wantErr != nil {
				require.ErrorIs(t, err, test.wantErr)
				return
			}

			require.NoError(t, err)
			require.Equal(t, test.want, got)
		})
	}
}
