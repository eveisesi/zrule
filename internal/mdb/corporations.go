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

type corporationRepository struct {
	corporations *mongo.Collection
}

func NewCorporationRepository(d *mongo.Database) (zrule.CorporationRepository, error) {

	corporations := d.Collection("corporations")
	_, err := corporations.Indexes().CreateOne(context.Background(), mongo.IndexModel{Keys: bsonx.Doc{{Key: "id", Value: bsonx.Int32(1)}}, Options: &options.IndexOptions{Name: newString("uniqueCorporationID"), Unique: newBool(true)}})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize user repository. Error encountered configured collection indexes: %w", err)
	}

	return &corporationRepository{
		corporations: corporations,
	}, nil

}

func (r *corporationRepository) Corporation(ctx context.Context, id uint) (*zrule.Corporation, error) {

	corporation := zrule.Corporation{}

	err := r.corporations.FindOne(ctx, primitive.D{primitive.E{Key: "id", Value: id}}).Decode(&corporation)

	return &corporation, err

}

func (r *corporationRepository) CreateCorporation(ctx context.Context, corporation *zrule.Corporation) (*zrule.Corporation, error) {

	corporation.CreatedAt = time.Now()
	corporation.UpdatedAt = time.Now()

	_, err := r.corporations.InsertOne(ctx, corporation)
	if err != nil {
		if !IsUniqueConstrainViolation(err) {
			return nil, err
		}
	}
	return corporation, nil

}
