package post

import (
	"context"

	"github.com/dubovayaRoshcha/posts-service/internal/entity"
	"github.com/dubovayaRoshcha/posts-service/internal/repository/sql"
	"github.com/dubovayaRoshcha/posts-service/internal/usecase/dto"
)

type PostRepo struct {
	db sql.DB
}

func NewPostRepo(db sql.DB) *PostRepo {
	return &PostRepo{
		db: db,
	}
}

func (r *PostRepo) GetList(ctx context.Context, params dto.PostRequest) ([]entity.Post, error) {
	query := `
		SELECT id, user_id, title, description, comments_allowed, created_at
		FROM posts
		ORDER BY created_at DESC, id DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(ctx, query, params.Limit, params.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := make([]entity.Post, 0, params.Limit)

	for rows.Next() {
		var post entity.Post

		err := rows.Scan(
			&post.ID, &post.UserID,
			&post.Title, &post.Description,
			&post.CommentsAllowed, &post.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return posts, nil
}

func (r *PostRepo) GetByID(ctx context.Context, postID int) (*entity.Post, error) {
	query := `
		SELECT id, user_id, title, description, comments_allowed, created_at
		FROM posts
		WHERE id = $1
	`

	var post entity.Post

	err := r.db.QueryRow(ctx, query, postID).Scan(
		&post.ID, &post.UserID, &post.Title,
		&post.Description, &post.CommentsAllowed,
		&post.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &post, nil
}

func (r *PostRepo) Create(ctx context.Context, post dto.PostDTO) (*entity.Post, error) {
	query := `
		INSERT INTO posts (
			user_id,
			title,
			description,
			comments_allowed
		)
		VALUES ($1, $2, $3, $4)
		RETURNING id, user_id, title, description, comments_allowed, created_at
	`

	var createdPost entity.Post

	err := r.db.QueryRow(
		ctx, query, post.UserID, post.Title,
		post.Description, post.CommentsAllowed,
	).Scan(
		&createdPost.ID, &createdPost.UserID, &createdPost.Title,
		&createdPost.Description, &createdPost.CommentsAllowed,
		&createdPost.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &createdPost, nil
}
