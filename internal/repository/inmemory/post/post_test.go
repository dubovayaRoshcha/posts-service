package post

import (
	"context"
	"testing"
	"time"

	"github.com/dubovayaRoshcha/posts-service/internal/entity"
	"github.com/dubovayaRoshcha/posts-service/internal/usecase/dto"
	"github.com/stretchr/testify/require"
)

func TestGetList(t *testing.T) {
	createdAt := time.Now()

	inputParams := dto.PostRequest{
		Limit:  10,
		Offset: 0,
	}

	tests := []struct {
		name   string
		posts  map[int]entity.Post
		params dto.PostRequest
		want   []entity.Post
	}{
		{
			name: "OK",
			posts: map[int]entity.Post{
				1: {ID: 1, UserID: 10, Title: "first post", Description: "first description", CommentsAllowed: true, CreatedAt: createdAt},
				2: {ID: 2, UserID: 20, Title: "second post", Description: "second description", CommentsAllowed: false, CreatedAt: createdAt},
			},
			params: inputParams,
			want: []entity.Post{
				{ID: 2, UserID: 20, Title: "second post", Description: "second description", CommentsAllowed: false, CreatedAt: createdAt},
				{ID: 1, UserID: 10, Title: "first post", Description: "first description", CommentsAllowed: true, CreatedAt: createdAt},
			},
		},
		{
			name:   "empty_result",
			posts:  map[int]entity.Post{},
			params: inputParams,
			want:   []entity.Post{},
		},
		{
			name: "pagination",
			posts: map[int]entity.Post{
				1: {ID: 1, UserID: 10, Title: "first post", Description: "first description", CommentsAllowed: true, CreatedAt: createdAt},
				2: {ID: 2, UserID: 20, Title: "second post", Description: "second description", CommentsAllowed: true, CreatedAt: createdAt},
				3: {ID: 3, UserID: 30, Title: "third post", Description: "third description", CommentsAllowed: true, CreatedAt: createdAt},
				4: {ID: 4, UserID: 40, Title: "fourth post", Description: "fourth description", CommentsAllowed: true, CreatedAt: createdAt},
			},
			params: dto.PostRequest{
				Limit:  2,
				Offset: 1,
			},
			want: []entity.Post{
				{ID: 3, UserID: 30, Title: "third post", Description: "third description", CommentsAllowed: true, CreatedAt: createdAt},
				{ID: 2, UserID: 20, Title: "second post", Description: "second description", CommentsAllowed: true, CreatedAt: createdAt},
			},
		},
		{
			name: "offset_out_of_range",
			posts: map[int]entity.Post{
				1: {ID: 1, UserID: 10, Title: "first post", Description: "first description", CommentsAllowed: true, CreatedAt: createdAt},
				2: {ID: 2, UserID: 20, Title: "second post", Description: "second description", CommentsAllowed: true, CreatedAt: createdAt},
			},
			params: dto.PostRequest{
				Limit:  10,
				Offset: 5,
			},
			want: []entity.Post{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo := &PostRepo{
				posts: test.posts,
			}

			got, err := repo.GetList(context.Background(), test.params)

			require.NoError(t, err)
			require.Equal(t, test.want, got)
		})
	}
}

func TestGetByID(t *testing.T) {
	createdAt := time.Now()

	expectedPost := &entity.Post{
		ID:              1,
		UserID:          10,
		Title:           "first post",
		Description:     "first description",
		CommentsAllowed: true,
		CreatedAt:       createdAt,
	}

	tests := []struct {
		name    string
		posts   map[int]entity.Post
		postID  int
		want    *entity.Post
		wantErr bool
	}{
		{
			name: "OK",
			posts: map[int]entity.Post{
				1: {
					ID:              1,
					UserID:          10,
					Title:           "first post",
					Description:     "first description",
					CommentsAllowed: true,
					CreatedAt:       createdAt,
				},
			},
			postID:  1,
			want:    expectedPost,
			wantErr: false,
		},
		{
			name:    "not_found",
			posts:   map[int]entity.Post{},
			postID:  1,
			want:    nil,
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo := &PostRepo{
				posts: test.posts,
			}

			got, err := repo.GetByID(context.Background(), test.postID)
			if test.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, test.want, got)
		})
	}
}

func TestCreate(t *testing.T) {
	inputPost := dto.PostDTO{
		UserID:          10,
		Title:           "first post",
		Description:     "first description",
		CommentsAllowed: true,
	}

	tests := []struct {
		name string
		post dto.PostDTO
		want *entity.Post
	}{
		{
			name: "OK",
			post: inputPost,
			want: &entity.Post{
				ID:              1,
				UserID:          10,
				Title:           "first post",
				Description:     "first description",
				CommentsAllowed: true,
				CreatedAt:       time.Time{},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo := NewPostRepo()

			got, err := repo.Create(context.Background(), test.post)

			require.NoError(t, err)
			got.CreatedAt = time.Time{}
			require.Equal(t, test.want, got)
		})
	}
}
