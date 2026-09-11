package post

import (
	"context"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/dubovayaRoshcha/posts-service/internal/entity"
	"github.com/dubovayaRoshcha/posts-service/internal/usecase/dto"
	"github.com/pashagolub/pgxmock/v5"
	"github.com/stretchr/testify/require"
)

func TestGetList(t *testing.T) {
	createdAt := time.Now()
	expectedPosts := []entity.Post{
		{ID: 1, UserID: 10, Title: "first post", Description: "first description", CommentsAllowed: true, CreatedAt: createdAt},
		{ID: 2, UserID: 20, Title: "second post", Description: "second description", CommentsAllowed: false, CreatedAt: createdAt},
	}

	inputParams := dto.PostRequest{
		Limit:  10,
		Offset: 0,
	}

	query := regexp.QuoteMeta(`
		SELECT id, user_id, title, description, comments_allowed, created_at
		FROM posts
		ORDER BY created_at DESC, id DESC
		LIMIT $1 OFFSET $2
	`)

	testErr := errors.New("database error")

	tests := []struct {
		name      string
		params    dto.PostRequest
		setupMock func(m pgxmock.PgxPoolIface)
		want      []entity.Post
		wantErr   bool
	}{
		{
			name:   "OK",
			params: inputParams,
			setupMock: func(m pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{
					"id", "user_id", "title", "description", "comments_allowed", "created_at",
				}).AddRow(1, 10, "first post", "first description", true, createdAt).
					AddRow(2, 20, "second post", "second description", false, createdAt)

				m.ExpectQuery(query).WithArgs(10, 0).WillReturnRows(rows)
			},
			want:    expectedPosts,
			wantErr: false,
		},
		{
			name:   "empty_result",
			params: inputParams,
			setupMock: func(m pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{
					"id", "user_id", "title", "description", "comments_allowed", "created_at",
				})

				m.ExpectQuery(query).WithArgs(10, 0).WillReturnRows(rows)
			},
			want:    []entity.Post{},
			wantErr: false,
		},
		{
			name:   "scan_error",
			params: inputParams,
			setupMock: func(m pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{
					"id", "user_id", "title", "description", "comments_allowed", "created_at",
				}).AddRow("bad-id", 10, "first post", "first description", true, createdAt)

				m.ExpectQuery(query).WithArgs(10, 0).WillReturnRows(rows)
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:   "query_error",
			params: inputParams,
			setupMock: func(m pgxmock.PgxPoolIface) {
				m.ExpectQuery(query).WithArgs(10, 0).WillReturnError(testErr)
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mock, err := pgxmock.NewPool()
			require.NoError(t, err)
			defer mock.Close()

			test.setupMock(mock)

			repo := NewPostRepo(mock)

			got, err := repo.GetList(context.Background(), test.params)

			if test.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, test.want, got)
			require.NoError(t, mock.ExpectationsWereMet())
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

	query := regexp.QuoteMeta(`
		SELECT id, user_id, title, description, comments_allowed, created_at
		FROM posts
		WHERE id = $1
	`)

	testErr := errors.New("database error")

	tests := []struct {
		name      string
		postID    int
		setupMock func(m pgxmock.PgxPoolIface)
		want      *entity.Post
		wantErr   error
	}{
		{
			name:   "OK",
			postID: 1,
			setupMock: func(m pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{
					"id", "user_id", "title", "description", "comments_allowed", "created_at",
				}).AddRow(1, 10, "first post", "first description", true, createdAt)

				m.ExpectQuery(query).WithArgs(1).WillReturnRows(rows)
			},
			want:    expectedPost,
			wantErr: nil,
		},
		{
			name:   "not_found",
			postID: 1,
			setupMock: func(m pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{
					"id", "user_id", "title", "description", "comments_allowed", "created_at",
				})

				m.ExpectQuery(query).WithArgs(1).WillReturnRows(rows)
			},
			want:    nil,
			wantErr: entity.PostNotFound,
		},
		{
			name:   "query_error",
			postID: 1,
			setupMock: func(m pgxmock.PgxPoolIface) {
				m.ExpectQuery(query).WithArgs(1).WillReturnError(testErr)
			},
			want:    nil,
			wantErr: testErr,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mock, err := pgxmock.NewPool()
			require.NoError(t, err)
			defer mock.Close()

			test.setupMock(mock)

			repo := NewPostRepo(mock)

			got, err := repo.GetByID(context.Background(), test.postID)

			if test.wantErr != nil {
				require.ErrorIs(t, err, test.wantErr)
				return
			}

			require.NoError(t, err)
			require.Equal(t, test.want, got)
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestCreate(t *testing.T) {
	createdAt := time.Now()
	inputPost := dto.PostDTO{
		UserID:          10,
		Title:           "new post",
		Description:     "new description",
		CommentsAllowed: true,
	}

	expectedPost := &entity.Post{
		ID:              1,
		UserID:          10,
		Title:           "new post",
		Description:     "new description",
		CommentsAllowed: true,
		CreatedAt:       createdAt,
	}

	query := regexp.QuoteMeta(`
		INSERT INTO posts (
			user_id,
			title,
			description,
			comments_allowed
		)
		VALUES ($1, $2, $3, $4)
		RETURNING id, user_id, title, description, comments_allowed, created_at
	`)

	testErr := errors.New("database error")

	tests := []struct {
		name      string
		post      dto.PostDTO
		setupMock func(m pgxmock.PgxPoolIface)
		want      *entity.Post
		wantErr   bool
	}{
		{
			name: "OK",
			post: inputPost,
			setupMock: func(m pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{
					"id", "user_id", "title", "description", "comments_allowed", "created_at",
				}).AddRow(1, 10, "new post", "new description", true, createdAt)

				m.ExpectQuery(query).WithArgs(10, "new post", "new description", true).WillReturnRows(rows)
			},
			want:    expectedPost,
			wantErr: false,
		},
		{
			name: "query_error",
			post: inputPost,
			setupMock: func(m pgxmock.PgxPoolIface) {
				m.ExpectQuery(query).WithArgs(10, "new post", "new description", true).WillReturnError(testErr)
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mock, err := pgxmock.NewPool()
			require.NoError(t, err)
			defer mock.Close()

			test.setupMock(mock)

			repo := NewPostRepo(mock)

			got, err := repo.Create(context.Background(), test.post)

			if test.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, test.want, got)
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
