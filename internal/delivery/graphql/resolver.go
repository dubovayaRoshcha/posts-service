package graphql

import "github.com/dubovayaRoshcha/posts-service/internal/delivery"

type Resolver struct {
	postUC    delivery.PostUseCase
	commentUC delivery.CommentUseCase
}

func NewResolver(postUC delivery.PostUseCase, commentUC delivery.CommentUseCase) *Resolver {
	return &Resolver{
		postUC:    postUC,
		commentUC: commentUC,
	}
}
