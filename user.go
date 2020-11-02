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
	CharacterID       uint64             `bson:"characterID" json:"characterID"`
	OwnerHash         string             `bson:"ownerHash" json:"ownerHash"`
	AccessToken       string             `bson:"accessToken" json:"accessToken"`
	RefreshToken      string             `bson:"refreshToken" json:"refreshToken"`
	Expires           time.Time          `bson:"expires" json:"expires"`
	Disabled          bool               `bson:"disabled" json:"disabled"`
	DisabledReason    *string            `bson:"disabledReason" json:"disabledReason"`
	DisabledTimestamp *time.Time         `bson:"disabledTime" json:"disabledTime"`
	CreatedAt         time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt         time.Time          `bson:"updatedAt" json:"updatedAt"`
}
