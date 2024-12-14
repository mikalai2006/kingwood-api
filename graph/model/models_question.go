package model

import (
	"time"

	"github.com/mikalai2006/kingwood-api/internal/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Question struct {
	ID            primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID        primitive.ObjectID `json:"userId" bson:"userId"`
	ProductID     primitive.ObjectID `json:"productId" bson:"productId" binding:"required" primitive:"true"`
	UserProductID primitive.ObjectID `json:"userProductId" bson:"userProductId"`
	Question      string             `json:"question" bson:"question"`
	Answer        string             `json:"answer" bson:"answer"`
	Status        *int               `json:"status" bson:"status"`

	User domain.User `json:"user" bson:"user,omitempty"`

	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}

type QuestionInput struct {
	ID            primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID        primitive.ObjectID `json:"userId" bson:"userId"`
	ProductID     primitive.ObjectID `json:"productId" bson:"productId"`
	UserProductID primitive.ObjectID `json:"userProductId" bson:"userProductId"`
	Question      string             `json:"question" bson:"question"`
	Answer        string             `json:"answer" bson:"answer"`
	Status        *int               `json:"status" bson:"status"`
	CreatedAt     time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt     time.Time          `json:"updatedAt" bson:"updatedAt"`
}

type QuestionFilter struct {
	ID            *string               `json:"id,omitempty"`
	ProductID     []*primitive.ObjectID `json:"productId,omitempty"`
	UserProductID []*primitive.ObjectID `json:"userProductId,omitempty"`
	UserID        []*primitive.ObjectID `json:"userId,omitempty"`
	Status        *int                  `json:"status" bson:"status"`

	Sort  []*ProductFilterSortParams `json:"sort,omitempty"`
	Limit *int                       `json:"limit,omitempty"`
	Skip  *int                       `json:"skip,omitempty"`
}
