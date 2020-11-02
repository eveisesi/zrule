package zrule

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ResourceRepository interface{}

type Resource struct {
	ID        primitive.ObjectID `bson:"_id" json:"_id"`
	OwnerID   uint64
	Type      string
	Endpoint  string
	CreatedAt time.Time
	UpdatedAt time.Time
}
