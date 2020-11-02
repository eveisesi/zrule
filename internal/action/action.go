package action

import (
	"context"
	"errors"

	"github.com/eveisesi/zrule"
	"go.mongodb.org/mongo-driver/mongo"
)

type Service interface {
	zrule.ActionRepository
}

type service struct {
	zrule.ActionRepository
}

func NewService(action zrule.ActionRepository) Service {
	return &service{
		ActionRepository: action,
	}
}

func (s *service) Actions(ctx context.Context, operators ...*zrule.Operator) ([]*zrule.Action, error) {

	actions, err := s.ActionRepository.Actions(ctx, operators...)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return actions, err
	}

	return actions, nil

}

func (s *service) CreateAction(ctx context.Context, action *zrule.Action) (*zrule.Action, error) {

	action.Infractions = make([]*zrule.Infraction, 0)
	return s.ActionRepository.CreateAction(ctx, action)

}
