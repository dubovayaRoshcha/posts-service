package validator

import (
	"github.com/dubovayaRoshcha/posts-service/internal/entity"
	"github.com/dubovayaRoshcha/posts-service/internal/usecase/dto"
)

const (
	MaxPostLimit         = 50
	MaxCommentLimit      = 50
	MaxTitleLenght       = 200
	MaxDescriptionLength = 10000
	MaxTextLenght        = 2000
)

func ValidatePostParams(params *dto.PostRequest) error {
	if params.Limit <= 0 || params.Offset < 0 {
		return entity.InvalidInput
	}

	if params.Limit > MaxPostLimit {
		params.Limit = MaxPostLimit
	}

	return nil
}

func ValidateID(id int) error {
	if id <= 0 {
		return entity.InvalidInput
	}

	return nil
}

func ValidatePost(post *dto.PostDTO) error {
	runeTitle := []rune(post.Title)
	runeDescription := []rune(post.Description)
	if len(runeTitle) == 0 {
		return entity.InvalidInput
	}

	if len(runeTitle) > MaxTitleLenght {
		return entity.MaxLengthExceeded
	}

	if len(runeDescription) > MaxDescriptionLength {
		return entity.MaxLengthExceeded
	}

	return nil
}

func ValidateTopCommentsParams(params *dto.CommentRequest) error {
	if params.Limit <= 0 || params.Offset < 0 || params.PostID <= 0 {
		return entity.InvalidInput
	}

	if params.Limit > MaxCommentLimit {
		params.Limit = MaxCommentLimit
	}

	return nil
}

func ValidateRepliesParams(params *dto.CommentRequest) error {
	if params.Limit <= 0 || params.Offset < 0 || params.PostID <= 0 {
		return entity.InvalidInput
	}

	if params.ReplyToCommentID == nil || *params.ReplyToCommentID <= 0 {
		return entity.InvalidInput
	}

	if params.Limit > MaxCommentLimit {
		params.Limit = MaxCommentLimit
	}

	return nil
}

func ValidateComment(comment *dto.CommentDTO) error {
	runeText := []rune(comment.Text)
	if len(runeText) == 0 || comment.PostID <= 0 {
		return entity.InvalidInput
	}

	if len(runeText) > MaxTextLenght {
		return entity.MaxLengthExceeded
	}

	return nil
}
