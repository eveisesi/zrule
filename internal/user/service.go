package user

import (
	"context"

	"github.com/eveisesi/zrule"
	"github.com/eveisesi/zrule/internal/token"
	"github.com/eveisesi/zrule/internal/universe"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

type Service interface {
	VerifyUserRegistrationByToken(ctx context.Context, token *oauth2.Token) error

	zrule.UserRepository
}

type service struct {
	logger   *logrus.Logger
	redis    *redis.Client
	token    token.Service
	universe universe.Service

	zrule.UserRepository
}

func NewService(logger *logrus.Logger, redis *redis.Client, token token.Service, universe universe.Service, user zrule.UserRepository) Service {
	return &service{
		logger:         logger,
		redis:          redis,
		token:          token,
		universe:       universe,
		UserRepository: user,
	}
}
