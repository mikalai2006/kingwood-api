package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PayTemplate struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID      primitive.ObjectID `json:"userId" bson:"userId"`
	Total       *int64             `json:"total" bson:"total"`
	Name        string             `json:"name" bson:"name"`
	Description string             `json:"description" bson:"description"`
	Enabled     *int64             `json:"enabled" bson:"enabled"`

	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}

type PayTemplateInput struct {
	ID     primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID string             `json:"userId" bson:"userId" form:"userId"`

	Total       *int64 `json:"total" bson:"total"`
	Name        string `json:"name" bson:"name"`
	Description string `json:"description" bson:"description"`
	Enabled     *int64 `json:"enabled" bson:"enabled"`

	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}
