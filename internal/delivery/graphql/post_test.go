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

func TestPosts(t *testing.T) {
	ctx := context.Background()

	existingPosts := []entity.Post{
		{ID: 1, UserID: 10, Title: "first post", Description: "first description", CommentsAllowed: true},
		{ID: 2, UserID: 20, Title: "second post", Description: "second description", CommentsAllowed: false},
	}

	postResponse := &dto.PostResponse{
		Posts: existingPosts,
		Len:   len(existingPosts),
	}

	expectedResponse := &PostResponse{
		Posts: []*Post{
			{ID: "1", UserID: 10, Title: "first post", Description: "first description", CommentsAllowed: true},
			{ID: "2", UserID: 20, Title: "second post", Description: "second description", CommentsAllowed: false},
		},
		Len: 2,
	}

	tests := []struct {
		name      string
		params    PostRequest
		setupMock func(m *mocks.MockPostUseCase)
		want      *PostResponse
		wantErr   error
	}{
		{
			name: "OK",
			params: PostRequest{
				Limit:  10,
				Offset: 0,
			},
			setupMock: func(m *mocks.MockPostUseCase) {
				m.EXPECT().GetListPosts(ctx, dto.PostRequest{
					Limit:  10,
					Offset: 0,
				}).
					Return(postResponse, nil).Times(1)
			},
			want:    expectedResponse,
			wantErr: nil,
		},
		{
			name: "empty_result",
			params: PostRequest{
				Limit:  10,
				Offset: 0,
			},
			setupMock: func(m *mocks.MockPostUseCase) {
				m.EXPECT().GetListPosts(ctx, dto.PostRequest{
					Limit:  10,
					Offset: 0,
				}).
					Return(&dto.PostResponse{
						Posts: []entity.Post{},
						Len:   0,
					}, nil).Times(1)
			},
			want: &PostResponse{
				Len:   0,
				Posts: []*Post{},
			},
			wantErr: nil,
		},
		{
			name: "usecase_invalid_input",
			params: PostRequest{
				Limit:  10,
				Offset: 0,
			},
			setupMock: func(m *mocks.MockPostUseCase) {
				m.EXPECT().GetListPosts(ctx, dto.PostRequest{
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
				test.setupMock(mockPostUC)
			}

			resolver := NewResolver(mockPostUC, mockCommentUC)
			query := &queryResolver{resolver}

			got, err := query.Posts(ctx, test.params)
			if test.wantErr != nil {
				require.ErrorIs(t, err, test.wantErr)
				return
			}

			require.NoError(t, err)
			require.Equal(t, test.want, got)
		})
	}
}

func TestPost(t *testing.T) {
	ctx := context.Background()

	existingPost := &entity.Post{
		ID:              1,
		UserID:          10,
		Title:           "first post",
		Description:     "first description",
		CommentsAllowed: true,
	}

	expectedPost := &Post{
		ID:              "1",
		UserID:          10,
		Title:           "first post",
		Description:     "first description",
		CommentsAllowed: true,
	}

	tests := []struct {
		name      string
		postID    int
		setupMock func(m *mocks.MockPostUseCase)
		want      *Post
		wantErr   error
	}{
		{
			name:   "OK",
			postID: 1,
			setupMock: func(m *mocks.MockPostUseCase) {
				m.EXPECT().GetPostByID(ctx, 1).
					Return(existingPost, nil).Times(1)
			},
			want:    expectedPost,
			wantErr: nil,
		},
		{
			name:   "usecase_post_not_found",
			postID: 1,
			setupMock: func(m *mocks.MockPostUseCase) {
				m.EXPECT().GetPostByID(ctx, 1).
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

			mockPostUC := mocks.NewMockPostUseCase(ctrl)
			mockCommentUC := mocks.NewMockCommentUseCase(ctrl)

			if test.setupMock != nil {
				test.setupMock(mockPostUC)
			}

			resolver := NewResolver(mockPostUC, mockCommentUC)
			query := &queryResolver{resolver}

			got, err := query.Post(ctx, test.postID)
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

	inputPost := CreatePostRequest{
		UserID:          10,
		Title:           "new post",
		Description:     "new description",
		CommentsAllowed: true,
	}

	existingPost := &entity.Post{
		ID:              1,
		UserID:          10,
		Title:           "new post",
		Description:     "new description",
		CommentsAllowed: true,
	}

	expectedPost := &Post{
		ID:              "1",
		UserID:          10,
		Title:           "new post",
		Description:     "new description",
		CommentsAllowed: true,
	}

	tests := []struct {
		name      string
		params    CreatePostRequest
		setupMock func(m *mocks.MockPostUseCase)
		want      *Post
		wantErr   error
	}{
		{
			name:   "OK",
			params: inputPost,
			setupMock: func(m *mocks.MockPostUseCase) {
				m.EXPECT().CreatePost(ctx, dto.PostDTO{
					UserID:          10,
					Title:           "new post",
					Description:     "new description",
					CommentsAllowed: true,
				}).
					Return(existingPost, nil).Times(1)
			},
			want:    expectedPost,
			wantErr: nil,
		},
		{
			name:   "usecase_invalid_input",
			params: inputPost,
			setupMock: func(m *mocks.MockPostUseCase) {
				m.EXPECT().CreatePost(ctx, dto.PostDTO{
					UserID:          10,
					Title:           "new post",
					Description:     "new description",
					CommentsAllowed: true,
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
				test.setupMock(mockPostUC)
			}

			resolver := NewResolver(mockPostUC, mockCommentUC)
			mutation := &mutationResolver{resolver}

			got, err := mutation.CreatePost(ctx, test.params)

			if test.wantErr != nil {
				require.ErrorIs(t, err, test.wantErr)
				return
			}

			require.NoError(t, err)
			require.Equal(t, test.want, got)
		})
	}
}
