package zrule

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserRepository interface {
	User(ctx context.Context, id uint64) (*User, error)
	CreateUser(ctx context.Context, user *User) (*User, error)
	UpdateUser(ctx context.Context, id uint64, user *User) (*User, error)
	DeleteUser(ctx context.Context, id primitive.ObjectID) error
}

type User struct {
	ID                primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	CharacterID       uint64             `bson:"character_id" json:"character_id"`
	OwnerHash         string             `bson:"owner_hash" json:"owner_hash"`
	AccessToken       string             `bson:"access_token" json:"access_token"`
	RefreshToken      string             `bson:"refresh_token" json:"refresh_token"`
	Expires           time.Time          `bson:"expires" json:"expires"`
	Disabled          bool               `bson:"disabled" json:"disabled"`
	DisabledReason    *string            `bson:"disabled_reason" json:"disabled_reason"`
	DisabledTimestamp *time.Time         `bson:"disabled_time" json:"disabled_time"`
	CreatedAt         time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt         time.Time          `bson:"updated_at" json:"updated_at"`
}
