package comment

import (
	"context"

	"github.com/dubovayaRoshcha/posts-service/internal/entity"
	"github.com/dubovayaRoshcha/posts-service/internal/repository/sql"
	"github.com/dubovayaRoshcha/posts-service/internal/repository/sql/utils"
	"github.com/dubovayaRoshcha/posts-service/internal/usecase/dto"
)

type CommentRepo struct {
	db sql.DB
}

func NewCommentRepo(db sql.DB) *CommentRepo {
	return &CommentRepo{
		db: db,
	}
}

func (r *CommentRepo) GetTopList(ctx context.Context, params dto.CommentRequest) ([]entity.Comment, error) {
	query := `
		SELECT id, post_id, reply_to_comment_id, user_id, text, created_at
		FROM comments
		WHERE post_id = $1
			AND reply_to_comment_id IS NULL
		ORDER BY created_at DESC, id DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(
		ctx, query, params.PostID,
		params.Limit, params.Offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	comments := make([]entity.Comment, 0, params.Limit)

	for rows.Next() {
		var comment entity.Comment

		err := rows.Scan(
			&comment.ID, &comment.PostID,
			&comment.ReplyToCommentID,
			&comment.UserID, &comment.Text,
			&comment.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		comments = append(comments, comment)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return comments, nil
}

func (r *CommentRepo) GetReplies(ctx context.Context, params dto.CommentRequest) ([]entity.Comment, error) {
	query := `
		SELECT id, post_id, reply_to_comment_id, user_id, text, created_at
		FROM comments
		WHERE post_id = $1
			AND reply_to_comment_id = $2
		ORDER BY created_at DESC, id DESC
		LIMIT $3 OFFSET $4
	`

	rows, err := r.db.Query(
		ctx, query, params.PostID,
		params.ReplyToCommentID,
		params.Limit, params.Offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	comments := make([]entity.Comment, 0, params.Limit)

	for rows.Next() {
		var comment entity.Comment

		err := rows.Scan(
			&comment.ID, &comment.PostID,
			&comment.ReplyToCommentID,
			&comment.UserID, &comment.Text,
			&comment.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		comments = append(comments, comment)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return comments, nil
}

func (r *CommentRepo) GetByID(ctx context.Context, commentID int) (*entity.Comment, error) {
	query := `
		SELECT id, post_id, reply_to_comment_id, user_id, text, created_at
		FROM comments
		WHERE id = $1
	`

	var comment entity.Comment

	err := r.db.QueryRow(ctx, query, commentID).Scan(
		&comment.ID, &comment.PostID,
		&comment.ReplyToCommentID,
		&comment.UserID, &comment.Text,
		&comment.CreatedAt,
	)
	if err != nil {
		return nil, utils.HandelPgError(err, entity.CommentNotFound)
	}

	return &comment, nil
}

func (r *CommentRepo) Create(ctx context.Context, comment dto.CommentDTO) (*entity.Comment, error) {
	query := `
		INSERT INTO comments (
			post_id,
			reply_to_comment_id,
			user_id,
			text
		)
		VALUES ($1, $2, $3, $4)
		RETURNING id, post_id, reply_to_comment_id, user_id, text, created_at
	`

	var createdComment entity.Comment

	err := r.db.QueryRow(
		ctx, query, comment.PostID,
		comment.ReplyToCommentID,
		comment.UserID, comment.Text,
	).Scan(
		&createdComment.ID, &createdComment.PostID,
		&createdComment.ReplyToCommentID,
		&createdComment.UserID, &createdComment.Text,
		&createdComment.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &createdComment, nil
}
