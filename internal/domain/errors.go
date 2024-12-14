package domain

import "errors"

var (
	ErrNodeNotFound       = errors.New("node not found")
	ErrReviewNotFound     = errors.New("review not found")
	ErrLikeExist          = errors.New("like exist")
	ErrQuestionExistValue = errors.New("question exist")
)
