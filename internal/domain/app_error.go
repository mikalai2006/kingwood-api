package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AppError struct {
	ID     primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID primitive.ObjectID `json:"userId" bson:"userId" primitive:"true"`

	Error  string `json:"error" bson:"error"`
	Status *int64 `json:"status" bson:"status"`

	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}

type AppErrorInput struct {
	ID     primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID primitive.ObjectID `json:"userId" bson:"userId"`

	Error  string `json:"error" bson:"error"`
	Status *int64 `json:"status" bson:"status"`

	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}

type AppErrorFilter struct {
	ID     []string            `json:"id,omitempty"`
	UserId []string            `json:"userId,omitempty"`
	Error  string              `json:"error,omitempty"`
	Status *int                `json:"status,omitempty"`
	Sort   []*FilterSortParams `json:"sort,omitempty"`
	Limit  *int                `json:"$limit,omitempty"`
	Skip   *int                `json:"$skip,omitempty"`
}
