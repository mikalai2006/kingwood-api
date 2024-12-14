package model

import (
	"time"

	"github.com/mikalai2006/kingwood-api/internal/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Ticket struct {
	ID     primitive.ObjectID `json:"_id" bson:"_id,omitempty" primitive:"true"`
	UserID primitive.ObjectID `json:"userId" bson:"user_id" primitive:"true"`
	Title  string             `json:"title" bson:"title" form:"title"`
	// Props     map[string]string  `json:"props" bson:"props"`
	Progress  int             `json:"progress" bson:"progress" form:"progress"`
	Status    int             `json:"status" bson:"status" form:"status"`
	Messages  []TicketMessage `json:"messages" bson:"messages"`
	User      domain.User     `json:"user" bson:"user,omitempty"`
	CreatedAt time.Time       `json:"createdAt" bson:"created_at"`
	UpdatedAt time.Time       `json:"updatedAt" bson:"updated_at"`
}

type TicketMongo struct {
	ID     primitive.ObjectID `json:"_id" bson:"_id,omitempty" primitive:"true"`
	UserID primitive.ObjectID `json:"userId" bson:"user_id" primitive:"true"`
	Title  string             `json:"title" bson:"title" form:"title"`
	// Props     map[string]string  `json:"props" bson:"props"`
	Progress  int       `json:"progress" bson:"progress" form:"progress"`
	Status    int       `json:"status" bson:"status" form:"status"`
	CreatedAt time.Time `json:"createdAt" bson:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updated_at"`
}

type TicketInput struct {
	ID     primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	UserID string             `json:"userId" bson:"user_id"`
	Title  string             `json:"title" bson:"title"`
	// Props     map[string]string  `json:"props" bson:"props"`
	Progress  int       `json:"progress" bson:"progress"`
	Status    int       `json:"status" bson:"status"`
	CreatedAt time.Time `json:"createdAt" bson:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updated_at"`
}

type TicketMessage struct {
	ID        primitive.ObjectID `json:"_id" bson:"_id,omitempty" primitive:"true"`
	UserID    primitive.ObjectID `json:"userId" bson:"user_id" primitive:"true"`
	TicketID  primitive.ObjectID `json:"ticketId" bson:"ticket_id" primitive:"true"`
	Text      string             `json:"text" bson:"text" form:"text"`
	Images    []domain.Image     `json:"images" bson:"images"`
	Status    int                `json:"status" bson:"status" form:"status"`
	User      domain.User        `json:"user" bson:"user,omitempty"`
	CreatedAt time.Time          `json:"createdAt" bson:"created_at"`
	UpdatedAt time.Time          `json:"updatedAt" bson:"updated_at"`
}

type TicketMessageMongo struct {
	ID        primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	UserID    primitive.ObjectID `json:"userId" bson:"user_id"`
	TicketID  primitive.ObjectID `json:"ticketId" bson:"ticket_id"`
	Text      string             `json:"text" bson:"text"`
	Status    int                `json:"status" bson:"status"`
	CreatedAt time.Time          `json:"createdAt" bson:"created_at"`
	UpdatedAt time.Time          `json:"updatedAt" bson:"updated_at"`
}
