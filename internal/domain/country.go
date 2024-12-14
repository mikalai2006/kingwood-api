package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Country struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Publish   bool               `json:"publish" bson:"publish"`
	Flag      string             `json:"flag" bson:"flag"`
	Image     string             `json:"image" bson:"image"`
	Name      string             `json:"name" bson:"name"`
	Code      string             `json:"code" bson:"code"`
	SortOrder int64              `json:"sortOrder" bson:"sort_order"`

	Stat StatInfo `json:"stat" bson:"stat"`

	CreatedAt time.Time `json:"createdAt" bson:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updated_at"`
}

type CountryInput struct {
	Publish   bool   `json:"publish" bson:"publish" form:"publish"`
	Flag      string `json:"flag" bson:"flag" form:"flag"`
	Image     string `json:"image" bson:"image"`
	Name      string `json:"name" bson:"name" form:"name"`
	Code      string `json:"code" bson:"code" form:"code"`
	SortOrder int64  `json:"sortOrder" bson:"sort_order" form:"sortOrder"`

	Stat StatInfo `json:"stat" bson:"stat" form:"stat"`
}
