package http

import (
	"context"

	"github.com/eveisesi/zrule"
)

func (s *server) restartRedisTracker(ctx context.Context) error {
	_, err := s.redis.Set(ctx, zrule.QUEUE_RESTART_TRACKER, 1, 0).Result()
	if err != nil {
		return err
	}

	return nil
}
