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

type allianceRepository struct {
	c *mongo.Collection
}

func NewRepository(db *mongo.Database) (zrule.AllianceRepository, error) {

	alliances := db.Collection("alliances")

	// Create a unique index on the ID property so looks up on faster
	_, err := alliances.Indexes().CreateOne(
		context.TODO(),
		mongo.IndexModel{
			Keys: bsonx.Doc{
				{
					Key:   "id",
					Value: bsonx.Int32(1),
				},
			},
			Options: &options.IndexOptions{
				Name:   newString("uniqueAllianceID"),
				Unique: newBool(true),
			},
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize alliance repository. Error encountered configured collection indexes: %w", err)
	}

	return &allianceRepository{
		c: alliances,
	}, nil

}

func (r *allianceRepository) Alliance(ctx context.Context, id uint) (*zrule.Alliance, error) {

	var alliance = new(zrule.Alliance)

	err := r.c.FindOne(ctx, primitive.D{primitive.E{Key: "id", Value: id}}).Decode(alliance)
	return alliance, err

}
func (r *allianceRepository) CreateAlliance(ctx context.Context, alliance *zrule.Alliance) (*zrule.Alliance, error) {

	alliance.CreatedAt = time.Now()
	alliance.UpdatedAt = time.Now()

	_, err := r.c.InsertOne(ctx, alliance)
	if err != nil {
		if !IsUniqueConstrainViolation(err) {
			return nil, err
		}
	}

	return alliance, nil

}
