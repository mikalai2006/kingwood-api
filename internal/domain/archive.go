package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ArchiveMeta struct {
	Author primitive.ObjectID `json:"author" bson:"author"`

	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
}
