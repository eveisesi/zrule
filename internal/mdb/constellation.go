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

type constellationRepository struct {
	constellations *mongo.Collection
}

func NewConstellationRepository(d *mongo.Database) (zrule.ConstellationRepository, error) {

	constellations := d.Collection("constellations")
	_, err := constellations.Indexes().CreateOne(context.Background(), mongo.IndexModel{Keys: bsonx.Doc{{Key: "id", Value: bsonx.Int32(1)}}, Options: &options.IndexOptions{Name: newString("uniqueConstellationID"), Unique: newBool(true)}})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize user repository. Error encountered configured collection indexes: %w", err)
	}

	return &constellationRepository{
		constellations: constellations,
	}, nil

}

func (r *constellationRepository) Constellations(ctx context.Context, operators ...*zrule.Operator) ([]*zrule.Constellation, error) {

	filters := BuildFilters(operators...)
	options := BuildFindOptions(operators...)

	var constellations = make([]*zrule.Constellation, 0)
	result, err := r.constellations.Find(ctx, filters, options)
	if err != nil {
		return constellations, err
	}

	err = result.All(ctx, &constellations)
	return constellations, err

}

func (r *constellationRepository) Constellation(ctx context.Context, id uint) (*zrule.Constellation, error) {

	constellation := zrule.Constellation{}

	err := r.constellations.FindOne(ctx, primitive.D{primitive.E{Key: "id", Value: id}}).Decode(&constellation)

	return &constellation, err

}

func (r *constellationRepository) CreateConstellation(ctx context.Context, constellation *zrule.Constellation) (*zrule.Constellation, error) {

	constellation.CreatedAt = time.Now()
	constellation.UpdatedAt = time.Now()

	_, err := r.constellations.InsertOne(ctx, constellation)
	if err != nil {
		if !IsUniqueConstrainViolation(err) {
			return nil, err
		}
	}
	return constellation, nil

}
