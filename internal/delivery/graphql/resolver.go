package graphql

import (
	"sync"

	"github.com/dubovayaRoshcha/posts-service/internal/delivery"
)

type Resolver struct {
	mu          sync.RWMutex
	channelList map[int][]chan *Comment
	postUC      delivery.PostUseCase
	commentUC   delivery.CommentUseCase
}

func NewResolver(postUC delivery.PostUseCase, commentUC delivery.CommentUseCase) *Resolver {
	return &Resolver{
		postUC:      postUC,
		commentUC:   commentUC,
		channelList: make(map[int][]chan *Comment),
	}
}

func (r *Resolver) removeChannel(postID int, ch chan *Comment) {
	r.mu.Lock()
	defer r.mu.Unlock()

	channels := r.channelList[postID]

	for i, channel := range channels {
		if channel == ch {
			r.channelList[postID] = append(channels[:i], channels[i+1:]...)
			break
		}
	}

	close(ch)
}
