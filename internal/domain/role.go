package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Role struct {
	ID primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	// UserID    primitive.ObjectID `json:"userId" bson:"userId"`
	Name      string    `json:"name" bson:"name"`
	Code      string    `json:"code" bson:"code"`
	Value     []string  `json:"value" bson:"value"`
	SortOrder int64     `json:"sortOrder" bson:"sortOrder"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}

type RoleInput struct {
	// UserID    string   `json:"userId" bson:"userId" form:"userId"`
	Name      string   `json:"name" bson:"name"`
	Code      string   `json:"code" bson:"code"`
	Value     []string `json:"value" bson:"value"`
	SortOrder int64    `json:"sortOrder" bson:"sortOrder"`
}

type RoleFilter struct {
	ID    []string            `json:"id,omitempty"`
	Name  []string            `json:"name,omitempty"`
	Code  []string            `json:"code,omitempty"`
	Sort  []*FilterSortParams `json:"sort,omitempty"`
	Limit *int                `json:"$limit,omitempty"`
	Skip  *int                `json:"$skip,omitempty"`
}
