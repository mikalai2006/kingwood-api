package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Pay struct {
	ID     primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID primitive.ObjectID `json:"userId" bson:"userId" primitive:"true"`

	WorkerId primitive.ObjectID `json:"workerId" bson:"workerId" primitive:"true"`
	Month    int64              `json:"month" bson:"month"`
	Year     int64              `json:"year" bson:"year"`
	Total    *int64             `json:"total" bson:"total"`
	Name     string             `json:"name" bson:"name"`
	Auto     *int               `json:"auto" bson:"auto"`

	Worker User                   `json:"worker" bson:"worker"`
	Props  map[string]interface{} `json:"props" bson:"props"`

	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}

type PayInput struct {
	ID     primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID primitive.ObjectID `json:"userId" bson:"userId"`

	WorkerId primitive.ObjectID     `json:"workerId" bson:"workerId" primitive:"true"`
	Month    *int64                 `json:"month" bson:"month"`
	Year     *int64                 `json:"year" bson:"year"`
	Total    *int64                 `json:"total" bson:"total"`
	Name     string                 `json:"name" bson:"name"`
	Auto     *int                   `json:"auto" bson:"auto"`
	Props    map[string]interface{} `json:"props" bson:"props"`

	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}

type PayFilter struct {
	ID       []string            `json:"id,omitempty"`
	WorkerId []string            `json:"workerId,omitempty"`
	Month    *int                `json:"month" bson:"month"`
	Year     *int                `json:"year" bson:"year"`
	Name     string              `json:"name,omitempty"`
	Auto     *int                `json:"auto" bson:"auto"`
	Sort     []*FilterSortParams `json:"sort,omitempty"`
	Limit    *int                `json:"$limit,omitempty"`
	Skip     *int                `json:"$skip,omitempty"`
}
