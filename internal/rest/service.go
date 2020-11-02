package rest

import (
	"context"
	"fmt"

	"github.com/eveisesi/zrule"
)

type service struct{}

func NewService(action *zrule.Action) (zrule.Dispatcher, error) {
	return &service{}, nil
}

func (s *service) Send(ctx context.Context, id uint, hash string) error {

	fmt.Println("Rest", id)

	return nil
}

func (s *service) SendTest(ctx context.Context, message string) error {
	fmt.Println("Rest Test", message)
	return nil
}
