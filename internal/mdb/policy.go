package mdb

import (
	"context"
	"time"

	"github.com/eveisesi/zrule"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type policyRepository struct {
	policies *mongo.Collection
}

func NewPolicyRepository(d *mongo.Database) (zrule.PolicyRepository, error) {

	policies := d.Collection("policies")

	// TODO: Policies may need some indexes outside of the default unique index on _id,
	// but at the time of writing this, I couldn't think of any - ddouglas

	return &policyRepository{
		policies: policies,
	}, nil

}

func (r *policyRepository) Policy(ctx context.Context, id primitive.ObjectID) (*zrule.Policy, error) {

	policy := new(zrule.Policy)

	err := r.policies.FindOne(ctx, primitive.D{primitive.E{Key: "_id", Value: id}}).Decode(policy)

	return policy, err

}

func (r *policyRepository) Policies(ctx context.Context, operators ...*zrule.Operator) ([]*zrule.Policy, error) {

	filters := BuildFilters(operators...)
	options := BuildFindOptions(operators...)

	var policies = make([]*zrule.Policy, 0)
	result, err := r.policies.Find(ctx, filters, options)
	if err != nil {
		return policies, err
	}

	err = result.All(ctx, &policies)
	return policies, err

}

func (r *policyRepository) CreatePolicy(ctx context.Context, policy *zrule.Policy) (*zrule.Policy, error) {

	policy.CreatedAt = time.Now()
	policy.UpdatedAt = time.Now()

	result, err := r.policies.InsertOne(ctx, policy)
	if err != nil {
		return nil, err
	}

	policy.ID = result.InsertedID.(primitive.ObjectID)

	return policy, nil
}

func (r *policyRepository) UpdatePolicy(ctx context.Context, id primitive.ObjectID, policy *zrule.Policy) (*zrule.Policy, error) {

	policy.ID = id
	policy.UpdatedAt = time.Now()

	update := primitive.D{primitive.E{Key: "$set", Value: policy}}

	_, err := r.policies.UpdateOne(ctx, primitive.D{primitive.E{Key: "_id", Value: id}}, update)

	return policy, err
}

func (r *policyRepository) DeletePolicy(ctx context.Context, id primitive.ObjectID) error {

	_, err := r.policies.DeleteOne(ctx, primitive.D{primitive.E{Key: "_id", Value: id}})

	return err
}
