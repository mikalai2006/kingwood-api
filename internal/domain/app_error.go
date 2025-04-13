package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AppError struct {
	ID     primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID primitive.ObjectID `json:"userId" bson:"userId" primitive:"true"`

	Error  string `json:"error" bson:"error"`
	Code   string `json:"code" bson:"code"`
	Stack  string `json:"stack" bson:"stack"`
	Status int64  `json:"status" bson:"status"`

	User User `json:"user,omitempty" bson:"user,omitempty"`

	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}

type AppErrorInput struct {
	ID     primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID primitive.ObjectID `json:"userId" bson:"userId"`

	Error  string `json:"error" bson:"error"`
	Code   string `json:"code" bson:"code"`
	Status *int64 `json:"status" bson:"status"`
	Stack  string `json:"stack" bson:"stack"`

	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}

type AppErrorFilter struct {
	ID     []string `json:"id,omitempty"`
	UserId []string `json:"userId,omitempty"`
	Error  string   `json:"error,omitempty"`
	Status *int     `json:"status,omitempty"`
	Stack  string   `json:"stack,omitempty"`

	Code  string              `json:"code,omitempty"`
	From  *time.Time          `json:"from,omitempty"`
	To    *time.Time          `json:"to,omitempty"`
	Sort  []*FilterSortParams `json:"sort,omitempty"`
	Limit *int                `json:"$limit,omitempty"`
	Skip  *int                `json:"$skip,omitempty"`
}
