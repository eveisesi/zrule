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

type regionRepository struct {
	regions *mongo.Collection
}

func NewRegionRepository(d *mongo.Database) (zrule.RegionRepository, error) {

	regions := d.Collection("regions")
	_, err := regions.Indexes().CreateOne(context.Background(), mongo.IndexModel{Keys: bsonx.Doc{{Key: "id", Value: bsonx.Int32(1)}}, Options: &options.IndexOptions{Name: newString("uniqueRegionID"), Unique: newBool(true)}})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize user repository. Error encountered configured collection indexes: %w", err)
	}

	return &regionRepository{
		regions: regions,
	}, nil

}

func (r *regionRepository) Regions(ctx context.Context, operators ...*zrule.Operator) ([]*zrule.Region, error) {

	filters := BuildFilters(operators...)
	options := BuildFindOptions(operators...)

	var regions = make([]*zrule.Region, 0)
	result, err := r.regions.Find(ctx, filters, options)
	if err != nil {
		return regions, err
	}

	err = result.All(ctx, &regions)
	return regions, err

}

func (r *regionRepository) Region(ctx context.Context, id uint) (*zrule.Region, error) {

	region := zrule.Region{}

	err := r.regions.FindOne(ctx, primitive.D{primitive.E{Key: "id", Value: id}}).Decode(&region)

	return &region, err

}

func (r *regionRepository) CreateRegion(ctx context.Context, region *zrule.Region) (*zrule.Region, error) {

	region.CreatedAt = time.Now()
	region.UpdatedAt = time.Now()

	_, err := r.regions.InsertOne(ctx, region)
	if err != nil {
		if !IsUniqueConstrainViolation(err) {
			return nil, err
		}
	}
	return region, nil

}
