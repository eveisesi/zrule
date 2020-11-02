package zrule

import (
	"context"
)

type CharacterRespository interface {
	Character(ctx context.Context, id uint64) (*Character, error)
	CreateCharacter(ctx context.Context, character *Character) error
	DeleteCharacter(ctx context.Context, id uint64) error
}

type Character struct {
	ID        uint64 `bson:"id" json:"id"`
	Name      string `bson:"name" json:"name"`
	CreatedAt int64  `bson:"createdAt" json:"createdAt"`
	UpdatedAt int64  `bson:"updatedAt" json:"updatedAt"`
}
