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

type actionRepository struct {
	actions *mongo.Collection
}

func NewActionRepository(d *mongo.Database) (zrule.ActionRepository, error) {

	actions := d.Collection("actions")
	_, err := actions.Indexes().CreateOne(context.Background(), mongo.IndexModel{Keys: bsonx.Doc{{Key: "ownerID", Value: bsonx.Int32(1)}}, Options: &options.IndexOptions{Name: newString("ownerIDIdx")}})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize action repository. Error encountered configuring ownerIDIdx on collection: %w", err)
	}

	_, err = actions.Indexes().CreateOne(context.Background(), mongo.IndexModel{Keys: bsonx.Doc{{Key: "endpoint", Value: bsonx.Int32(1)}}, Options: &options.IndexOptions{Name: newString("endpointUniqueIdx"), Unique: newBool(true)}})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize action repository. Error encountered configuring endpointUniqueIdx on collection: %w", err)
	}

	return &actionRepository{
		actions: actions,
	}, nil
}

func (r *actionRepository) Action(ctx context.Context, id primitive.ObjectID) (*zrule.Action, error) {

	action := new(zrule.Action)
	err := r.actions.FindOne(ctx, primitive.D{{Key: "_id", Value: id}}).Decode(action)
	return action, err

}

func (r *actionRepository) Actions(ctx context.Context, operators ...*zrule.Operator) ([]*zrule.Action, error) {

	filters := BuildFilters(operators...)
	options := BuildFindOptions(operators...)

	var actions = make([]*zrule.Action, 0)
	result, err := r.actions.Find(ctx, filters, options)
	if err != nil {
		return actions, err
	}

	err = result.All(ctx, &actions)
	return actions, err

}

func (r *actionRepository) CreateAction(ctx context.Context, action *zrule.Action) (*zrule.Action, error) {

	action.CreatedAt = time.Now()
	action.UpdatedAt = time.Now()

	result, err := r.actions.InsertOne(ctx, action)
	if err != nil {
		return nil, err
	}

	action.ID = result.InsertedID.(primitive.ObjectID)

	return action, err

}

func (r *actionRepository) UpdateAction(ctx context.Context, id primitive.ObjectID, action *zrule.Action) (*zrule.Action, error) {
	action.ID = id
	action.UpdatedAt = time.Now()

	update := primitive.D{primitive.E{Key: "$set", Value: action}}

	_, err := r.actions.UpdateOne(ctx, primitive.D{primitive.E{Key: "_id", Value: id}}, update)

	return action, err
}

func (r *actionRepository) DeleteAction(ctx context.Context, id primitive.ObjectID) error {

	_, err := r.actions.DeleteOne(ctx, primitive.D{primitive.E{Key: "_id", Value: id}})

	return err

}
