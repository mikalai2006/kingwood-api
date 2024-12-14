package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Language struct {
	ID      primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Publish bool               `json:"publish" bson:"publish"`
	Flag    string             `json:"flag" bson:"flag"`
	Name    string             `json:"name" bson:"name"`
	Code    string             `json:"code" bson:"code"`
	Locale  string             `json:"locale" bson:"locale"`

	Localization map[string]interface{} `json:"localization" bson:"localization"`

	SortOrder int64     `json:"sortOrder" bson:"sort_order"`
	CreatedAt time.Time `json:"createdAt" bson:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updated_at"`
}

type LanguageInput struct {
	Publish      bool                   `json:"publish" bson:"publish" form:"publish"`
	Flag         string                 `json:"flag" bson:"flag" form:"flag"`
	Name         string                 `json:"name" bson:"name" form:"name"`
	Code         string                 `json:"code" bson:"code" form:"code"`
	Locale       string                 `json:"locale" bson:"locale" form:"locale"`
	SortOrder    int64                  `json:"sortOrder" bson:"sort_order" form:"sortOrder"`
	Localization map[string]interface{} `json:"localization" bson:"localization" form:"localization"`
}

// type Category struct {
// 	ID          int64             `json:"id" bson:"_id,omitempty" form:"-"`
// 	ParentID    int64             `json:"parentId" bson:"parentId"`
// 	Title       map[string]string `json:"title" bson:"title" form:"title"`
// 	Description map[string]string `json:"description" bson:"description" form:"description"`
// 	Seo         string            `json:"seo" bson:"seo" form:"seo"`
// 	SortOrder   int64             `json:"sortOrder" bson:"sortOrder" form:"sort_order"`
// 	MPath       string            `json:"mpath" bson:"mpath" form:"mpath"`
// 	Level       string            `json:"level" bson:"level" form:"level"`
// 	Status      bool              `json:"status" bson:"status" form:"status"`
// 	CreatedAt   time.Time         `json:"createdAt" bson:"createdAt" form:"created_at"`
// 	UpdatedAt   time.Time         `json:"updatedAt" bson:"updatedAt" form:"updated_at"`
// }

// type CategoryInput struct {
// 	Title       map[string]string `json:"title" bson:"title" form:"title"`
// 	Description map[string]string `json:"description" bson:"description" form:"description"`
// 	Seo         string            `json:"seo" bson:"seo" form:"seo"`
// 	SortOrder   int64             `json:"sortOrder" bson:"sortOrder" form:"sortOrder"`
// 	MPath       string            `json:"mpath" bson:"mpath" form:"mpath"`
// 	Level       string            `json:"level" bson:"level" form:"level"`
// 	Status      bool              `json:"status" bson:"status" form:"status"`
// }
