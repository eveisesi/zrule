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

type itemRepository struct {
	items *mongo.Collection
}

func NewItemRepository(d *mongo.Database) (zrule.ItemRepository, error) {

	items := d.Collection("items")
	_, err := items.Indexes().CreateOne(context.Background(), mongo.IndexModel{Keys: bsonx.Doc{{Key: "id", Value: bsonx.Int32(1)}}, Options: &options.IndexOptions{Name: newString("uniqueItemID"), Unique: newBool(true)}})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize user repository. Error encountered configured collection indexes: %w", err)
	}

	return &itemRepository{
		items: items,
	}, nil

}

func (r *itemRepository) Item(ctx context.Context, id uint) (*zrule.Item, error) {

	item := zrule.Item{}

	err := r.items.FindOne(ctx, primitive.D{primitive.E{Key: "id", Value: id}}).Decode(&item)

	return &item, err

}

func (r *itemRepository) CreateItem(ctx context.Context, item *zrule.Item) (*zrule.Item, error) {

	item.CreatedAt = time.Now()
	item.UpdatedAt = time.Now()

	_, err := r.items.InsertOne(ctx, item)
	if err != nil {
		if !IsUniqueConstrainViolation(err) {
			return nil, err
		}
	}
	return item, nil

}
