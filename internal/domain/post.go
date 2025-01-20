package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Post struct {
	ID          primitive.ObjectID     `json:"id" bson:"_id,omitempty"`
	UserID      primitive.ObjectID     `json:"userId" bson:"userId"`
	Name        string                 `json:"name" bson:"name"`
	Description string                 `json:"description" bson:"description"`
	Props       map[string]interface{} `json:"props" bson:"props"`
	Color       string                 `json:"color" bson:"color"`
	SortOrder   int64                  `json:"sortOrder" bson:"sortOrder"`
	Hidden      int                    `json:"hidden" bson:"hidden"`
	CreatedAt   time.Time              `json:"createdAt" bson:"createdAt"`
	UpdatedAt   time.Time              `json:"updatedAt" bson:"updatedAt"`
}

type PostInput struct {
	ID          primitive.ObjectID     `json:"id" bson:"_id,omitempty"`
	UserID      string                 `json:"userId" bson:"userId" form:"userId"`
	Name        string                 `json:"name" bson:"name" form:"name"`
	Description string                 `json:"description" bson:"description" form:"description"`
	Props       map[string]interface{} `json:"props" bson:"props"`
	Color       string                 `json:"color" bson:"color" form:"color"`
	SortOrder   int64                  `json:"sortOrder" bson:"sortOrder" form:"sortOrder"`
	Hidden      *int                   `json:"hidden" bson:"hidden" form:"hidden"`
	CreatedAt   time.Time              `json:"createdAt" bson:"createdAt"`
	UpdatedAt   time.Time              `json:"updatedAt" bson:"updatedAt"`
}
