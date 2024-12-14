package model

import (
	"time"

	"github.com/mikalai2006/kingwood-api/internal/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Product struct {
	ID         primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID     primitive.ObjectID `json:"userId" bson:"user_id"`
	CategoryID primitive.ObjectID `json:"categoryId" bson:"category_id" form:"categoryId" primitive:"true"`
	// RoomID string                 `json:"roomId" bson:"roomId" form:"roomId"`
	Title       string                 `json:"title" bson:"title" form:"title"`
	Description string                 `json:"description" bson:"description" form:"description"`
	Status      int64                  `json:"status" bson:"status" form:"status"`
	Props       map[string]interface{} `json:"props" bson:"props" form:"props"`
	Cost        *int64                 `json:"cost" bson:"cost" form:"cost"`
	Actions     []int                  `json:"actions" bson:"actions" form:"actions"`
	Lon         float64                `json:"lon" bson:"lon"`
	Lat         float64                `json:"lat" bson:"lat"`
	AddressId   primitive.ObjectID     `json:"addressId" bson:"addressId"`

	// Data   []Nodedata `json:"data" bson:"data,omitempty"`
	Images  []domain.Image `json:"images" bson:"images,omitempty"`
	Address Address        `json:"address" bson:"address,omitempty"`
	User    domain.User    `json:"user" bson:"user,omitempty"`
	Offers  []Offer        `json:"offers" bson:"offers,omitempty"`

	CreatedAt time.Time `json:"createdAt" bson:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updated_at"`
}

type ProductInputData struct {
	UserID      string                 `json:"userId" bson:"user_id" form:"userId" primitive:"true"`
	CategoryID  string                 `json:"categoryId" bson:"category_id" form:"categoryId" primitive:"true"`
	Title       string                 `json:"title" bson:"title" form:"title"`
	Description string                 `json:"description" bson:"description" form:"description"`
	Props       map[string]interface{} `json:"props" bson:"props" form:"props"`
	Status      int64                  `json:"status" bson:"status" form:"status"`
	Cost        *int64                 `json:"cost" bson:"cost" form:"cost"`
	Actions     []int                  `json:"actions" bson:"actions" form:"actions"`
	Lon         float64                `json:"lon" bson:"lon" form:"lon"`
	Lat         float64                `json:"lat" bson:"lat" form:"lat"`
	AddressId   string                 `json:"addressId" bson:"addressId" form:"addressId"`

	CreatedAt time.Time `json:"createdAt" bson:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updated_at"`
}

type ProductFilter struct {
	ID         []*string                  `json:"id,omitempty"`
	ProductID  []*primitive.ObjectID      `json:"productId,omitempty"`
	Query      *string                    `json:"query,omitempty"`
	CategoryID []*string                  `json:"categoryId,omitempty"`
	UserID     *string                    `json:"userId,omitempty"`
	Cost       *int                       `json:"cost,omitempty"`
	Actions    []*int                     `json:"actions,omitempty"`
	Lon        float64                    `json:"lon" bson:"lon"`
	Lat        float64                    `json:"lat" bson:"lat"`
	AddressId  *string                    `json:"addressId" bson:"addressId"`
	Sort       []*ProductFilterSortParams `json:"sort,omitempty"`
	Limit      *int                       `json:"limit,omitempty"`
	Skip       *int                       `json:"skip,omitempty"`
}
type ProductFilterSortParams struct {
	Key   *string `json:"key,omitempty"`
	Value *int    `json:"value,omitempty"`
}

type ProductInput struct {
	ID          primitive.ObjectID     `json:"id" bson:"_id,omitempty"`
	UserID      primitive.ObjectID     `json:"userId" bson:"user_id"`
	CategoryID  primitive.ObjectID     `json:"categoryId" bson:"category_id"`
	Title       string                 `json:"title" bson:"title"`
	Description string                 `json:"description" bson:"description"`
	Status      int64                  `json:"status" bson:"status"` // 1 - exchange, 2 - gift
	Props       map[string]interface{} `json:"props" bson:"props"`
	Cost        *int64                 `json:"cost" bson:"cost"`
	Actions     []int                  `json:"actions" bson:"actions"`
	Lon         float64                `json:"lon" bson:"lon"`
	Lat         float64                `json:"lat" bson:"lat"`
	AddressId   primitive.ObjectID     `json:"addressId" bson:"addressId"`

	CreatedAt time.Time `json:"createdAt" bson:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updated_at"`
}
