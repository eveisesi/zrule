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

type solarSystemRepository struct {
	solarSystems *mongo.Collection
}

func NewSolarSystemRepository(d *mongo.Database) (zrule.SolarSystemRepository, error) {

	solarSystems := d.Collection("solarSystems")
	_, err := solarSystems.Indexes().CreateOne(context.Background(), mongo.IndexModel{Keys: bsonx.Doc{{Key: "id", Value: bsonx.Int32(1)}}, Options: &options.IndexOptions{Name: newString("uniqueSolarSystemID"), Unique: newBool(true)}})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize user repository. Error encountered configured collection indexes: %w", err)
	}

	return &solarSystemRepository{
		solarSystems: solarSystems,
	}, nil

}

func (r *solarSystemRepository) SolarSystems(ctx context.Context, operators ...*zrule.Operator) ([]*zrule.SolarSystem, error) {

	filters := BuildFilters(operators...)
	options := BuildFindOptions(operators...)

	var solarSystems = make([]*zrule.SolarSystem, 0)
	result, err := r.solarSystems.Find(ctx, filters, options)
	if err != nil {
		return solarSystems, err
	}

	err = result.All(ctx, &solarSystems)
	return solarSystems, err

}

func (r *solarSystemRepository) SolarSystem(ctx context.Context, id uint) (*zrule.SolarSystem, error) {

	solarSystem := zrule.SolarSystem{}

	err := r.solarSystems.FindOne(ctx, primitive.D{primitive.E{Key: "id", Value: id}}).Decode(&solarSystem)

	return &solarSystem, err

}

func (r *solarSystemRepository) CreateSolarSystem(ctx context.Context, solarSystem *zrule.SolarSystem) (*zrule.SolarSystem, error) {

	solarSystem.CreatedAt = time.Now()
	solarSystem.UpdatedAt = time.Now()

	_, err := r.solarSystems.InsertOne(ctx, solarSystem)
	if err != nil {
		if !IsUniqueConstrainViolation(err) {
			return nil, err
		}
	}
	return solarSystem, nil

}
