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

type factionRepository struct {
	factions *mongo.Collection
}

func NewFactionRepository(d *mongo.Database) (zrule.FactionRepository, error) {

	factions := d.Collection("factions")
	_, err := factions.Indexes().CreateOne(context.Background(), mongo.IndexModel{Keys: bsonx.Doc{{Key: "id", Value: bsonx.Int32(1)}}, Options: &options.IndexOptions{Name: newString("uniqueFactionID"), Unique: newBool(true)}})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize user repository. Error encountered configured collection indexes: %w", err)
	}

	return &factionRepository{
		factions: factions,
	}, nil

}

func (r *factionRepository) Factions(ctx context.Context, operators ...*zrule.Operator) ([]*zrule.Faction, error) {

	filters := BuildFilters(operators...)
	options := BuildFindOptions(operators...)

	var factions = make([]*zrule.Faction, 0)
	result, err := r.factions.Find(ctx, filters, options)
	if err != nil {
		return factions, err
	}

	err = result.All(ctx, &factions)
	return factions, err

}

func (r *factionRepository) Faction(ctx context.Context, id uint) (*zrule.Faction, error) {

	faction := zrule.Faction{}

	err := r.factions.FindOne(ctx, primitive.D{primitive.E{Key: "id", Value: id}}).Decode(&faction)

	return &faction, err

}

func (r *factionRepository) CreateFaction(ctx context.Context, faction *zrule.Faction) (*zrule.Faction, error) {

	faction.CreatedAt = time.Now()
	faction.UpdatedAt = time.Now()

	_, err := r.factions.InsertOne(ctx, faction)
	if err != nil {
		if !IsUniqueConstrainViolation(err) {
			return nil, err
		}
	}
	return faction, nil

}
