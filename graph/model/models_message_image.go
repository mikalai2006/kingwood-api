package model

// import (
// 	"time"

// 	"go.mongodb.org/mongo-driver/bson/primitive"
// )

// type MessageImage struct {
// 	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
// 	UserID    primitive.ObjectID `json:"userId" bson:"userId"`
// 	RoomID    string             `json:"roomId" bson:"roomId"`
// 	MessageID string             `json:"messageId" bson:"messageId"`
// 	Service   string             `json:"service" bson:"service"`
// 	Path      string             `json:"path" bson:"path"`
// 	Ext       string             `json:"ext" bson:"ext"`
// 	Title     string             `json:"title" bson:"title"`
// 	Dir       string             `json:"dir" bson:"dir"`

// 	//User User `json:"user,omitempty" bson:"user,omitempty"`

// 	Description string    `json:"description" bson:"description"`
// 	CreatedAt   time.Time `json:"createdAt" bson:"createdAt"`
// 	UpdatedAt   time.Time `json:"updatedAt" bson:"updatedAt"`
// }

// type MessageImageInputMongo struct {
// 	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
// 	UserID    primitive.ObjectID `json:"userId" bson:"userId"`
// 	RoomID    primitive.ObjectID `json:"roomId" bson:"roomId"`
// 	MessageID primitive.ObjectID `json:"messageId" bson:"messageId"`
// 	Service   string             `json:"service" bson:"service"`
// 	Path      string             `json:"path" bson:"path"`
// 	Ext       string             `json:"ext" bson:"ext"`
// 	Title     string             `json:"title" bson:"title"`
// 	Dir       string             `json:"dir" bson:"dir"`

// 	Description string    `json:"description" bson:"description"`
// 	CreatedAt   time.Time `json:"createdAt" bson:"createdAt"`
// 	UpdatedAt   time.Time `json:"updatedAt" bson:"updatedAt"`
// }

// type MessageImageInput struct {
// 	UserID      string `json:"userId" bson:"userId" form:"userId" primitive:"true"`
// 	RoomID      string `json:"roomId" bson:"roomId" form:"roomId" primitive:"true"`
// 	MessageID   string `json:"messageId" bson:"messageId" form:"messageId" primitive:"true"`
// 	Service     string `json:"service" bson:"service" form:"service"`
// 	Path        string `json:"path" bson:"path"`
// 	Title       string `json:"title" bson:"title" form:"title"`
// 	Ext         string `json:"ext" bson:"ext"`
// 	Dir         string `json:"dir" bson:"dir" form:"dir"`
// 	Description string `json:"description" bson:"description" form:"description"`
// 	// Images      *multipart.FileHeader `bson:"image" form:"image"`
// }
