package policy

import (
	"context"
	"errors"

	"github.com/eveisesi/zrule"
	"go.mongodb.org/mongo-driver/mongo"
)

type Service interface {
	zrule.PolicyRepository
}

type service struct {
	zrule.PolicyRepository
}

func NewService(action zrule.PolicyRepository) Service {
	return &service{
		PolicyRepository: action,
	}
}

func (s *service) Policies(ctx context.Context, operators ...*zrule.Operator) ([]*zrule.Policy, error) {

	policies, err := s.PolicyRepository.Policies(ctx, operators...)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return policies, err
	}

	return policies, nil

}
