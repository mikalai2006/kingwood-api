package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Order struct {
	ID     primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID primitive.ObjectID `json:"userId" bson:"userId" primitive:"true"`

	Number          int                `json:"number" bson:"number"`
	Name            string             `json:"name" bson:"name" form:"name"`
	Description     string             `json:"description" bson:"description" form:"description"`
	ConstructorId   primitive.ObjectID `json:"constructorId" bson:"constructorId" primitive:"true"`
	ObjectId        primitive.ObjectID `json:"objectId" bson:"objectId" form:"objectId"`
	Term            time.Time          `json:"term" bson:"term" form:"term"`
	DateStart       time.Time          `json:"dateStart" bson:"dateStart" form:"dateStart"`
	TermMontaj      time.Time          `json:"termMontaj" bson:"termMontaj" form:"termMontaj"`
	Priority        *int64             `json:"priority" bson:"priority" form:"priority"`
	StolyarComplete *int64             `json:"stolyarComplete" bson:"stolyarComplete" form:"stolyarComplete"`
	MalyarComplete  *int64             `json:"malyarComplete" bson:"malyarComplete" form:"malyarComplete"`
	ShlifComplete   *int64             `json:"shlifComplete" bson:"shlifComplete" form:"shlifComplete"`
	GoComplete      *int64             `json:"goComplete" bson:"goComplete" form:"goComplete"`
	DateOtgruzka    time.Time          `json:"dateOtgruzka" bson:"dateOtgruzka" form:"dateOtgruzka"`
	MontajComplete  *int64             `json:"montajComplete" bson:"montajComplete" form:"montajComplete"`
	// NeedMontaj    *int64             `json:"needMontaj" bson:"needMontaj" form:"needMontaj"`
	Status *int64   `json:"status" bson:"status" form:"status"`
	Group  []string `json:"group" bson:"group" form:"group"`

	// User User `json:"user" bson:"user"`
	Object Object `json:"object" bson:"object"`
	Tasks  []Task `json:"tasks" bson:"tasks" form:"tasks"`

	Year      *int      `json:"year" bson:"year"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}

type OrderInput struct {
	ID     primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID primitive.ObjectID `json:"userId" bson:"userId"`

	Number          int64              `json:"number" bson:"number"`
	Name            string             `json:"name" bson:"name" form:"name"`
	Description     string             `json:"description" bson:"description" form:"description"`
	ObjectId        primitive.ObjectID `json:"objectId" bson:"objectId" form:"objectId"`
	ConstructorId   primitive.ObjectID `json:"constructorId" bson:"constructorId" primitive:"true"`
	Term            time.Time          `json:"term" bson:"term" form:"term"`
	DateStart       time.Time          `json:"dateStart" bson:"dateStart" form:"dateStart"`
	TermMontaj      time.Time          `json:"termMontaj" bson:"termMontaj" form:"termMontaj"`
	Priority        *int64             `json:"priority" bson:"priority" form:"priority"`
	StolyarComplete *int64             `json:"stolyarComplete" bson:"stolyarComplete" form:"stolyarComplete"`
	MalyarComplete  *int64             `json:"malyarComplete" bson:"malyarComplete" form:"malyarComplete"`
	ShlifComplete   *int64             `json:"shlifComplete" bson:"shlifComplete" form:"shlifComplete"`
	GoComplete      *int64             `json:"goComplete" bson:"goComplete" form:"goComplete"`
	DateOtgruzka    time.Time          `json:"dateOtgruzka" bson:"dateOtgruzka" form:"dateOtgruzka"`
	MontajComplete  *int64             `json:"montajComplete" bson:"montajComplete" form:"montajComplete"`
	// NeedMontaj    *int64             `json:"needMontaj" bson:"needMontaj" form:"needMontaj"`
	Status *int64   `json:"status" bson:"status" form:"status"`
	Group  []string `json:"group" bson:"group" form:"group"`

	Year      int       `json:"year" bson:"year"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}

// type OrderInputData struct {
// 	ID     string   `json:"id" bson:"_id" primitive:"true"`
// 	Name   string   `json:"name" bson:"name"  form:"name"`
// 	Status *int64   `json:"status" bson:"status" form:"status"`
// 	Group  []string `json:"group" bson:"group" form:"group"`

// 	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
// 	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
// }

type OrderFilter struct {
	ID              []string   `json:"id,omitempty"`
	Name            string     `json:"name,omitempty"`
	Group           []string   `json:"group,omitempty"`
	Status          *int64     `json:"status"`
	Number          *int       `json:"number"`
	ObjectId        []string   `json:"objectId" bson:"objectId" form:"objectId"`
	StolyarComplete *int64     `json:"stolyarComplete"`
	ShlifComplete   *int64     `json:"shlifComplete"`
	MalyarComplete  *int64     `json:"malyarComplete"`
	GoComplete      *int64     `json:"goComplete"`
	MontajComplete  *int64     `json:"montajComplete"`
	Year            *int       `json:"year"`
	From            *time.Time `json:"from,omitempty"`
	To              *time.Time `json:"to,omitempty"`
	Date            *time.Time `json:"date,omitempty"`
	// NeedMontaj *int64              `json:"needMontaj" bson:"needMontaj" form:"needMontaj"`
	Sort  []*FilterSortParams `json:"$sort,omitempty"`
	Limit *int                `json:"$limit,omitempty"`
	Skip  *int                `json:"$skip,omitempty"`
}
