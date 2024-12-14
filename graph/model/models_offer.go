package model

import (
	"time"

	"github.com/mikalai2006/kingwood-api/internal/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Offer struct {
	ID     primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID primitive.ObjectID `json:"userId" bson:"userId"`

	ProductID     primitive.ObjectID `json:"productId" bson:"productId"`
	UserProductID primitive.ObjectID `json:"userProductId" bson:"userProductId"`
	RoomId        primitive.ObjectID `json:"roomId" bson:"roomId"`
	RejectUserId  primitive.ObjectID `json:"rejectUserId" bson:"rejectUserId"`
	Message       string             `json:"message" bson:"message"`
	Status        int64              `json:"status" bson:"status"`
	Win           *int               `json:"win" bson:"win"`
	Take          *int               `json:"take" bson:"take"`
	Give          *int               `json:"give" bson:"give"`
	Cost          int64              `json:"cost" bson:"cost"`

	User domain.User `json:"user" bson:"user,omitempty"`

	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}

type OfferFilter struct {
	ID            *string               `json:"id,omitempty"`
	ProductID     []*primitive.ObjectID `json:"productId,omitempty"`
	UserProductID []*primitive.ObjectID `json:"userProductId,omitempty"`
	UserID        *primitive.ObjectID   `json:"userId,omitempty"`
	Status        *int                  `json:"status" bson:"status"`

	Sort  []*ProductFilterSortParams `json:"sort,omitempty"`
	Limit *int                       `json:"limit,omitempty"`
	Skip  *int                       `json:"skip,omitempty"`
}

type OfferInputMongo struct {
	ID     primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID primitive.ObjectID `json:"userId" bson:"userId"`

	ProductID     primitive.ObjectID `json:"productId" bson:"productId"`
	UserProductID primitive.ObjectID `json:"userProductId" bson:"userProductId"`
	RejectUserId  primitive.ObjectID `json:"rejectUserId" bson:"rejectUserId"`
	RoomId        primitive.ObjectID `json:"roomId" bson:"roomId"`
	Message       string             `json:"message" bson:"message"`
	Status        int64              `json:"status" bson:"status"`
	Cost          int64              `json:"cost" bson:"cost"`
	Win           *int               `json:"win" bson:"win"`
	Take          *int               `json:"take" bson:"take"`
	Give          *int               `json:"give" bson:"give"`

	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}

type OfferInput struct {
	ProductID     string `json:"productId" bson:"productId"`
	UserProductID string `json:"userProductId" bson:"userProductId"`
	RoomId        string `json:"roomId" bson:"roomId"`
	RejectUserId  string `json:"rejectUserId" bson:"rejectUserId"`
	Message       string `json:"message" bson:"message"`
	Status        int64  `json:"status" bson:"status"`
	Cost          int64  `json:"cost" bson:"cost"`
	Win           *int   `json:"win" bson:"win"`
	Take          *int   `json:"take" bson:"take"`
	Give          *int   `json:"give" bson:"give"`
}

type NodedataAudit struct {
	ID         primitive.ObjectID     `json:"id" bson:"_id,omitempty"`
	UserID     primitive.ObjectID     `json:"userId" bson:"userId"`
	NodedataID primitive.ObjectID     `json:"nodedataId" bson:"nodedataId"`
	Value      int64                  `json:"value" bson:"value"`
	Props      map[string]interface{} `json:"props" bson:"props"`

	User domain.User `json:"user,omitempty" bson:"user,omitempty"`

	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}

type NodedataAuditDB struct {
	ID         primitive.ObjectID     `json:"id" bson:"_id,omitempty"`
	UserID     primitive.ObjectID     `json:"userId" bson:"user_id"`
	NodedataID primitive.ObjectID     `json:"nodedataId" bson:"nodedata_id"`
	Value      int64                  `json:"value" bson:"value"`
	Props      map[string]interface{} `json:"props" bson:"props"`
	CreatedAt  time.Time              `json:"createdAt" bson:"created_at"`
	UpdatedAt  time.Time              `json:"updatedAt" bson:"updated_at"`
}

type NodedataAuditInput struct {
	ID         primitive.ObjectID     `json:"id" bson:"_id,omitempty"`
	UserID     string                 `json:"userId" bson:"user_id"`
	NodedataID string                 `json:"nodedataId" bson:"nodedata_id"`
	Value      int64                  `json:"value" bson:"value"`
	Props      map[string]interface{} `json:"props" bson:"props"`
	CreatedAt  time.Time              `json:"createdAt" bson:"created_at"`
	UpdatedAt  time.Time              `json:"updatedAt" bson:"updated_at"`
}
