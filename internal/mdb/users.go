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

type userRepository struct {
	users *mongo.Collection
}

func NewUserRepository(d *mongo.Database) (zrule.UserRepository, error) {

	users := d.Collection("users")
	_, err := users.Indexes().CreateOne(context.Background(), mongo.IndexModel{Keys: bsonx.Doc{{Key: "character_id", Value: bsonx.Int32(1)}}, Options: &options.IndexOptions{Name: newString("uniqueCharacter"), Unique: newBool(true)}})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize user repository. Error encountered configured collection indexes: %w", err)
	}

	return &userRepository{
		users: users,
	}, nil
}

func (r *userRepository) User(ctx context.Context, id uint64) (*zrule.User, error) {
	user := new(zrule.User)
	err := r.users.FindOne(ctx, primitive.D{primitive.E{Key: "character_id", Value: id}}).Decode(user)

	return user, err
}

func (r *userRepository) CreateUser(ctx context.Context, user *zrule.User) (*zrule.User, error) {

	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	result, err := r.users.InsertOne(ctx, user)
	if err != nil {
		if !IsUniqueConstrainViolation(err) {
			return nil, err
		}
	}

	user.ID = result.InsertedID.(primitive.ObjectID)

	return user, err

}

func (r *userRepository) UpdateUser(ctx context.Context, id uint64, user *zrule.User) (*zrule.User, error) {

	user.CharacterID = id
	user.UpdatedAt = time.Now()

	update := primitive.D{primitive.E{Key: "$set", Value: user}}

	_, err := r.users.UpdateOne(ctx, primitive.D{primitive.E{Key: "characterID", Value: id}}, update)

	return user, err
}

func (r *userRepository) DeleteUser(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.users.DeleteOne(ctx, primitive.D{primitive.E{Key: "_id", Value: id}})

	return err
}
