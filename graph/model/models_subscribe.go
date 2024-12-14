package model

import (
	"time"

	"github.com/mikalai2006/kingwood-api/internal/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Subscribe struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID    primitive.ObjectID `json:"userId" bson:"userId"`
	SubUserID primitive.ObjectID `json:"subUserId" bson:"subUserId"`
	Status    int                `json:"status" bson:"status"`

	User    domain.User `json:"user" bson:"user,omitempty"`
	SubUser domain.User `json:"subUser" bson:"subUser,omitempty"`

	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}

type SubscribeInput struct {
	UserID    primitive.ObjectID `json:"userId" bson:"userId"`
	SubUserID primitive.ObjectID `json:"subUserId" bson:"subUserId"`
	Status    int                `json:"status" bson:"status" form:"status"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time          `json:"updatedAt" bson:"updatedAt"`
}

type SubscribeFilter struct {
	ID        *string               `json:"id,omitempty"`
	SubUserID []*primitive.ObjectID `json:"subUserId,omitempty"`
	UserID    *primitive.ObjectID   `json:"userId,omitempty"`
	Status    *int                  `json:"status" bson:"status"`

	Sort  []*ProductFilterSortParams `json:"sort,omitempty"`
	Limit *int                       `json:"limit,omitempty"`
	Skip  *int                       `json:"skip,omitempty"`
}
