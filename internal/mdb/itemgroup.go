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

type itemGroupRepository struct {
	itemGroups *mongo.Collection
}

func NewItemGroupRepository(d *mongo.Database) (zrule.ItemGroupRepository, error) {

	itemGroups := d.Collection("itemGroups")
	_, err := itemGroups.Indexes().CreateOne(context.Background(), mongo.IndexModel{Keys: bsonx.Doc{{Key: "id", Value: bsonx.Int32(1)}}, Options: &options.IndexOptions{Name: newString("uniqueItemGroupID"), Unique: newBool(true)}})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize user repository. Error encountered configured collection indexes: %w", err)
	}

	return &itemGroupRepository{
		itemGroups: itemGroups,
	}, nil

}

func (r *itemGroupRepository) ItemGroup(ctx context.Context, id uint) (*zrule.ItemGroup, error) {

	itemGroup := zrule.ItemGroup{}

	err := r.itemGroups.FindOne(ctx, primitive.D{primitive.E{Key: "id", Value: id}}).Decode(&itemGroup)

	return &itemGroup, err

}

func (r *itemGroupRepository) CreateItemGroup(ctx context.Context, itemGroup *zrule.ItemGroup) (*zrule.ItemGroup, error) {

	itemGroup.CreatedAt = time.Now()
	itemGroup.UpdatedAt = time.Now()

	_, err := r.itemGroups.InsertOne(ctx, itemGroup)
	if err != nil {
		if !IsUniqueConstrainViolation(err) {
			return nil, err
		}
	}
	return itemGroup, nil

}
