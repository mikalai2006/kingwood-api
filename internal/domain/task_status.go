package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TaskStatus struct {
	ID          primitive.ObjectID     `json:"id" bson:"_id,omitempty"`
	UserID      primitive.ObjectID     `json:"userId" bson:"userId"`
	Name        string                 `json:"name" bson:"name"`
	Description string                 `json:"description" bson:"description"`
	Props       map[string]interface{} `json:"props" bson:"props"`
	Color       string                 `json:"color" bson:"color"`
	Enabled     *int64                 `json:"enabled" bson:"enabled"`
	Icon        string                 `json:"icon" bson:"icon"`
	Animate     string                 `json:"animate" bson:"animate"`
	Start       *int64                 `json:"start" bson:"start"`
	Finish      *int64                 `json:"finish" bson:"finish"`
	Process     *int64                 `json:"process" bson:"process"`
	Status      string                 `json:"status" bson:"status"`

	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}

type TaskStatusInput struct {
	ID          primitive.ObjectID     `json:"id" bson:"_id,omitempty"`
	UserID      string                 `json:"userId" bson:"userId" form:"userId"`
	Name        string                 `json:"name" bson:"name" form:"name"`
	Description string                 `json:"description" bson:"description" form:"description"`
	Props       map[string]interface{} `json:"props" bson:"props"`
	Color       string                 `json:"color" bson:"color" form:"color"`
	Enabled     *int64                 `json:"enabled" bson:"enabled" form:"enabled"`
	Icon        string                 `json:"icon" bson:"icon" form:"icon"`
	Animate     *string                `json:"animate" bson:"animate" form:"animate"`
	Start       *int64                 `json:"start" bson:"start"`
	Finish      *int64                 `json:"finish" bson:"finish"`
	Process     *int64                 `json:"process" bson:"process"`
	Status      string                 `json:"status" bson:"status"`
	CreatedAt   time.Time              `json:"createdAt" bson:"createdAt"`
	UpdatedAt   time.Time              `json:"updatedAt" bson:"updatedAt"`
}
