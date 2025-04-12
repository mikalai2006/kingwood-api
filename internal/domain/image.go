package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Image struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID    primitive.ObjectID `json:"userId" bson:"userId"`
	ServiceID string             `json:"serviceId" bson:"serviceId"`
	Service   string             `json:"service" bson:"service"`
	Path      string             `json:"path" bson:"path"`
	Ext       string             `json:"ext" bson:"ext"`
	Title     string             `json:"title" bson:"title"`
	Dir       string             `json:"dir" bson:"dir"`

	//User User `json:"user,omitempty" bson:"user,omitempty"`

	Description string    `json:"description" bson:"description"`
	CreatedAt   time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt" bson:"updatedAt"`
}

type ImageInputMongo struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID    primitive.ObjectID `json:"userId" bson:"userId"`
	ServiceID string             `json:"serviceId" bson:"serviceId"`
	Service   string             `json:"service" bson:"service"`
	Path      string             `json:"path" bson:"path"`
	Ext       string             `json:"ext" bson:"ext"`
	Title     string             `json:"title" bson:"title"`
	Dir       string             `json:"dir" bson:"dir"`

	Description string    `json:"description" bson:"description"`
	CreatedAt   time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt" bson:"updatedAt"`
}

// type ImageSize struct {
// 	Url30   string `json:"url30" bson:"url30"`
// 	Url320  string `json:"url320" bson:"url320"`
// 	Url768  string `json:"url768" bson:"url768"`
// 	Url1024 string `json:"url1024" bson:"url1024"`
// 	Url1280 string `json:"url1280" bson:"url1280"`
// }

type ImageInput struct {
	UserID      string `json:"userId" bson:"userId" form:"userId" primitive:"true"`
	ServiceID   string `json:"serviceId" bson:"serviceId" form:"serviceId" primitive:"true"`
	Service     string `json:"service" bson:"service" form:"service"`
	Path        string `json:"path" bson:"path"`
	Description string `json:"description" bson:"description" form:"description"`
	Title       string `json:"title" bson:"title" form:"title"`
	Dir         string `json:"dir" bson:"dir" form:"dir"`
	Ext         string `json:"ext" bson:"ext"`
	// Images      *multipart.FileHeader `bson:"image" form:"image"`
}

type IImagePaths struct {
	Ext  string
	Path string
}

type ImageFilter struct {
	ID        []string            `json:"id" bson:"id"`
	UserId    []string            `json:"userId" bson:"userId"`
	ServiceId []string            `json:"serviceId" bson:"serviceId"`
	Service   []string            `json:"service" bson:"service"`
	Sort      []*FilterSortParams `json:"sort,omitempty"`
	Limit     *int                `json:"$limit,omitempty"`
	Skip      *int                `json:"$skip,omitempty"`
}
