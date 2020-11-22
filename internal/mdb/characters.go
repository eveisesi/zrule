package mdb

import (
	"context"
	"fmt"
	"time"

	"github.com/eveisesi/zrule"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

type characterRepository struct {
	characters *mongo.Collection
}

func NewCharacterRepository(d *mongo.Database) (zrule.CharacterRepository, error) {

	characters := d.Collection("characters")
	_, err := characters.Indexes().CreateOne(context.Background(), mongo.IndexModel{Keys: bsonx.Doc{{Key: "id", Value: bsonx.Int32(1)}}, Options: &options.IndexOptions{Name: newString("uniqueCharacterID"), Unique: newBool(true)}})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize user repository. Error encountered configured collection indexes: %w", err)
	}

	return &characterRepository{
		characters: characters,
	}, nil

}

func (r *characterRepository) Character(ctx context.Context, id uint64) (*zrule.Character, error) {

	character := zrule.Character{}

	err := r.characters.FindOne(ctx, primitive.D{primitive.E{Key: "id", Value: id}}).Decode(&character)

	return &character, err

}

func (r *characterRepository) CreateCharacter(ctx context.Context, character *zrule.Character) (*zrule.Character, error) {

	character.CreatedAt = time.Now()
	character.UpdatedAt = time.Now()

	_, err := r.characters.InsertOne(ctx, character)
	if err != nil {
		if !IsUniqueConstrainViolation(err) {
			return nil, err
		}
	}
	return character, nil

}
