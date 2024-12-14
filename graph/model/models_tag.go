package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Tag struct {
	ID          primitive.ObjectID           `json:"id" bson:"_id,omitempty"`
	UserID      primitive.ObjectID           `json:"userId" bson:"user_id"`
	Key         string                       `json:"key" bson:"key"`
	Type        string                       `json:"type" bson:"type"`
	Title       string                       `json:"title" bson:"title"`
	Description string                       `json:"description" bson:"description"`
	Props       map[string]interface{}       `json:"props" bson:"props"`
	Locale      map[string]map[string]string `json:"locale" bson:"locale"`
	MultiOpt    int64                        `json:"multiOpt" bson:"multi_opt"`
	IsFilter    bool                         `json:"isFilter" bson:"is_filter"`
	// Filter        int                          `json:"filter" bson:"filter"`
	Options       []Question `json:"options,omitempty" bson:"options,omitempty"`
	Multilanguage bool       `json:"multilanguage" bson:"multilanguage"`
	CountItem     int        `json:"countItem" bson:"countItem"`
	// TagoptID      []string  `json:"tagoptId" bson:"tagopt_id"`
	SortOrder int       `json:"sortOrder" bson:"sort_order"`
	CreatedAt time.Time `json:"createdAt" bson:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updated_at"`
}

type TagInput struct {
	ID          primitive.ObjectID           `json:"id" bson:"_id,omitempty"`
	UserID      string                       `json:"userId" bson:"user_id" form:"userId"`
	Key         string                       `json:"key" bson:"key" form:"key"`
	Type        string                       `json:"type" bson:"type" form:"type"`
	Title       string                       `json:"title" bson:"title" form:"title"`
	Description string                       `json:"description" bson:"description" form:"description"`
	Props       map[string]interface{}       `json:"props" bson:"props"`
	Locale      map[string]map[string]string `json:"locale" bson:"locale"`
	MultiOpt    int64                        `json:"multiOpt" bson:"multi_opt"`
	IsFilter    bool                         `json:"isFilter" bson:"is_filter"`
	// Filter        int                          `json:"filter" bson:"filter"`
	Multilanguage bool `json:"multilanguage" bson:"multilanguage"`
	SortOrder     int  `json:"sortOrder" bson:"sort_order"`
	// TagoptID      []string  `json:"tagoptId" bson:"tagopt_id"`
	CreatedAt time.Time `json:"createdAt" bson:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updated_at"`
}
