package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Object struct {
	ID     primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID primitive.ObjectID `json:"userId" bson:"userId" primitive:"true"`

	Name string `json:"name" bson:"name"`
	// Orders   []Order       `json:"orders" bson:"orders"`

	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}

type ObjectInput struct {
	ID     primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID primitive.ObjectID `json:"userId" bson:"userId"`

	Name string `json:"name" bson:"name"`

	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}

type ObjectFilter struct {
	ID    []*string           `json:"id,omitempty"`
	Name  *string             `json:"name,omitempty"`
	Sort  []*FilterSortParams `json:"sort,omitempty"`
	Limit *int                `json:"$limit,omitempty"`
	Skip  *int                `json:"$skip,omitempty"`
}

type FilterSortParams struct {
	Key   string `json:"key,omitempty"`
	Value int    `json:"value,omitempty"`
}

type ObjectInputData struct {
	ID     string `json:"id" bson:"_id" primitive:"true"`
	UserID string `json:"userId" bson:"userId" primitive:"true"`

	Name  string `json:"name" bson:"name"`
	Query string `json:"query" bson:"query"`

	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}
