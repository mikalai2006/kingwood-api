package domain

import "errors"

var (
	ErrNodeNotFound       = errors.New("node not found")
	ErrReviewNotFound     = errors.New("review not found")
	ErrLikeExist          = errors.New("like exist")
	ErrQuestionExistValue = errors.New("question exist")

	ErrNotItemMongo     = errors.New("Не найдена запись для вашего запроса!")
	ErrBlocked          = errors.New("Ваша учетная запись заблокирована!")
	ErrNotRole          = errors.New("Нет прав для данной операции!")
	ErrExistNumberOrder = errors.New("Заказ с указанным номером существует в базе!")
)
