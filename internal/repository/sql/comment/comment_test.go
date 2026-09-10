package comment

import (
	"context"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/dubovayaRoshcha/posts-service/internal/entity"
	"github.com/dubovayaRoshcha/posts-service/internal/usecase/dto"
	"github.com/pashagolub/pgxmock/v5"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
)

func TestGetTopList(t *testing.T) {
	createdAt := time.Now()
	expectedComments := []entity.Comment{
		{ID: 1, PostID: 10, ReplyToCommentID: nil, UserID: 20, Text: "first comment", CreatedAt: createdAt},
		{ID: 2, PostID: 10, ReplyToCommentID: nil, UserID: 30, Text: "second comment", CreatedAt: createdAt},
	}

	inputParams := dto.CommentRequest{
		PostID: 10,
		Limit:  10,
		Offset: 0,
	}

	query := regexp.QuoteMeta(`
		SELECT id, post_id, reply_to_comment_id, user_id, text, created_at
		FROM comments
		WHERE post_id = $1
			AND reply_to_comment_id IS NULL
		ORDER BY created_at DESC, id DESC
		LIMIT $2 OFFSET $3
	`)

	testErr := errors.New("database error")

	tests := []struct {
		name      string
		params    dto.CommentRequest
		setupMock func(m pgxmock.PgxPoolIface)
		want      []entity.Comment
		wantErr   bool
	}{
		{
			name:   "OK",
			params: inputParams,
			setupMock: func(m pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{
					"id", "post_id", "reply_to_comment_id", "user_id", "text", "created_at",
				}).AddRow(1, 10, nil, 20, "first comment", createdAt).
					AddRow(2, 10, nil, 30, "second comment", createdAt)

				m.ExpectQuery(query).WithArgs(10, 10, 0).WillReturnRows(rows)
			},
			want:    expectedComments,
			wantErr: false,
		},
		{
			name:   "empty_result",
			params: inputParams,
			setupMock: func(m pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{
					"id", "post_id", "reply_to_comment_id", "user_id", "text", "created_at",
				})

				m.ExpectQuery(query).WithArgs(10, 10, 0).WillReturnRows(rows)
			},
			want:    []entity.Comment{},
			wantErr: false,
		},
		{
			name:   "scan_error",
			params: inputParams,
			setupMock: func(m pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{
					"id", "post_id", "reply_to_comment_id", "user_id", "text", "created_at",
				}).AddRow("bad-id", 10, nil, 20, "first comment", createdAt)

				m.ExpectQuery(query).WithArgs(10, 10, 0).WillReturnRows(rows)
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:   "query_error",
			params: inputParams,
			setupMock: func(m pgxmock.PgxPoolIface) {
				m.ExpectQuery(query).WithArgs(10, 10, 0).WillReturnError(testErr)
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

			repo := NewCommentRepo(mock)

			got, err := repo.GetTopList(context.Background(), test.params)

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

func TestGetReplies(t *testing.T) {
	createdAt := time.Now()
	replyToCommentID := 1

	expectedComments := []entity.Comment{
		{ID: 2, PostID: 10, ReplyToCommentID: &replyToCommentID, UserID: 20, Text: "first reply", CreatedAt: createdAt},
		{ID: 3, PostID: 10, ReplyToCommentID: &replyToCommentID, UserID: 30, Text: "second reply", CreatedAt: createdAt},
	}

	inputParams := dto.CommentRequest{
		PostID:           10,
		ReplyToCommentID: &replyToCommentID,
		Limit:            10,
		Offset:           0,
	}

	query := regexp.QuoteMeta(`
		SELECT id, post_id, reply_to_comment_id, user_id, text, created_at
		FROM comments
		WHERE post_id = $1
			AND reply_to_comment_id = $2
		ORDER BY created_at DESC, id DESC
		LIMIT $3 OFFSET $4
	`)

	testErr := errors.New("database error")

	tests := []struct {
		name      string
		params    dto.CommentRequest
		setupMock func(m pgxmock.PgxPoolIface)
		want      []entity.Comment
		wantErr   bool
	}{
		{
			name:   "OK",
			params: inputParams,
			setupMock: func(m pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{
					"id", "post_id", "reply_to_comment_id", "user_id", "text", "created_at",
				}).AddRow(2, 10, lo.ToPtr(1), 20, "first reply", createdAt).
					AddRow(3, 10, lo.ToPtr(1), 30, "second reply", createdAt)

				m.ExpectQuery(query).WithArgs(10, &replyToCommentID, 10, 0).WillReturnRows(rows)
			},
			want:    expectedComments,
			wantErr: false,
		},
		{
			name:   "empty_result",
			params: inputParams,
			setupMock: func(m pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{
					"id", "post_id", "reply_to_comment_id", "user_id", "text", "created_at",
				})

				m.ExpectQuery(query).WithArgs(10, &replyToCommentID, 10, 0).WillReturnRows(rows)
			},
			want:    []entity.Comment{},
			wantErr: false,
		},
		{
			name:   "scan_error",
			params: inputParams,
			setupMock: func(m pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{
					"id", "post_id", "reply_to_comment_id", "user_id", "text", "created_at",
				}).AddRow("bad-id", 10, lo.ToPtr(1), 20, "first reply", createdAt)

				m.ExpectQuery(query).WithArgs(10, &replyToCommentID, 10, 0).WillReturnRows(rows)
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:   "query_error",
			params: inputParams,
			setupMock: func(m pgxmock.PgxPoolIface) {
				m.ExpectQuery(query).WithArgs(10, &replyToCommentID, 10, 0).WillReturnError(testErr)
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

			repo := NewCommentRepo(mock)

			got, err := repo.GetReplies(context.Background(), test.params)

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

func TestCreate(t *testing.T) {
	createdAt := time.Now()
	replyToCommentID := 1

	inputComment := dto.CommentDTO{
		PostID:           10,
		ReplyToCommentID: &replyToCommentID,
		UserID:           20,
		Text:             "new comment",
	}

	expectedComment := &entity.Comment{
		ID:               2,
		PostID:           10,
		ReplyToCommentID: &replyToCommentID,
		UserID:           20,
		Text:             "new comment",
		CreatedAt:        createdAt,
	}

	query := regexp.QuoteMeta(`
		INSERT INTO comments (
			post_id,
			reply_to_comment_id,
			user_id,
			text
		)
		VALUES ($1, $2, $3, $4)
		RETURNING id, post_id, reply_to_comment_id, user_id, text, created_at
	`)

	testErr := errors.New("database error")

	tests := []struct {
		name      string
		comment   dto.CommentDTO
		setupMock func(m pgxmock.PgxPoolIface)
		want      *entity.Comment
		wantErr   bool
	}{
		{
			name:    "OK",
			comment: inputComment,
			setupMock: func(m pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{
					"id", "post_id", "reply_to_comment_id", "user_id", "text", "created_at",
				}).AddRow(2, 10, lo.ToPtr(1), 20, "new comment", createdAt)

				m.ExpectQuery(query).WithArgs(10, &replyToCommentID, 20, "new comment").WillReturnRows(rows)
			},
			want:    expectedComment,
			wantErr: false,
		},
		{
			name:    "query_error",
			comment: inputComment,
			setupMock: func(m pgxmock.PgxPoolIface) {
				m.ExpectQuery(query).WithArgs(10, &replyToCommentID, 20, "new comment").WillReturnError(testErr)
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

			repo := NewCommentRepo(mock)

			got, err := repo.Create(context.Background(), test.comment)

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
