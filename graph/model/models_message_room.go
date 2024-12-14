package model

import (
	"time"

	"github.com/mikalai2006/kingwood-api/internal/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MessageRoom struct {
	ID         primitive.ObjectID     `json:"id" bson:"_id,omitempty"`
	UserID     primitive.ObjectID     `json:"userId" bson:"userId"`
	ProductID  primitive.ObjectID     `json:"productId" bson:"productId"`
	TakeUserID primitive.ObjectID     `json:"takeUserId" bson:"takeUserId"`
	OfferID    primitive.ObjectID     `json:"offerId" bson:"offerId"`
	Status     *int                   `json:"status" bson:"status"`
	Props      map[string]interface{} `json:"props" bson:"props"`

	User domain.User `json:"user,omitempty" bson:"user,omitempty"`

	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}

type MessageRoomMongo struct {
	ID         primitive.ObjectID     `json:"id" bson:"_id,omitempty"`
	UserID     primitive.ObjectID     `json:"userId" bson:"userId"`
	ProductID  primitive.ObjectID     `json:"productId" bson:"productId"`
	OfferID    primitive.ObjectID     `json:"offerId" bson:"offerId"`
	TakeUserID primitive.ObjectID     `json:"takeUserId" bson:"takeUserId"`
	Status     *int                   `json:"status" bson:"status"`
	Props      map[string]interface{} `json:"props" bson:"props"`

	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}

type MessageRoomFilter struct {
	ID         *primitive.ObjectID        `json:"id" bson:"id"`
	UserID     *primitive.ObjectID        `json:"userId" bson:"userId"`
	ProductID  *primitive.ObjectID        `json:"productId" bson:"productId"`
	OfferID    *primitive.ObjectID        `json:"offerId" bson:"offerId"`
	TakeUserID *primitive.ObjectID        `json:"takeUserId" bson:"takeUserId"`
	Status     *int                       `json:"status" bson:"status"`
	Sort       []*ProductFilterSortParams `json:"sort" bson:"sort"`
	Limit      *int                       `json:"limit" bson:"limit"`
	Skip       *int                       `json:"skip" bson:"skip"`
}
