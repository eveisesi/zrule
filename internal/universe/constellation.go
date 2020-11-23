package universe

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/eveisesi/zrule"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
)

func (s *service) Constellation(ctx context.Context, id uint) (*zrule.Constellation, error) {
	var constellation = new(zrule.Constellation)
	var key = fmt.Sprintf(zrule.CACHE_CONSTELLATION, id)

	result, err := s.redis.Get(ctx, key).Bytes()
	if err != nil && err.Error() != "redis: nil" {
		return nil, err
	}

	if len(result) > 0 {

		err = json.Unmarshal(result, constellation)
		if err != nil {
			return nil, errors.Wrap(err, "unable to unmarshal constellation from redis")
		}
		return constellation, nil
	}

	constellation, err = s.ConstellationRepository.Constellation(ctx, id)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return nil, errors.Wrap(err, "unable to query database for constellation")
	}

	if err == nil {
		bSlice, err := json.Marshal(constellation)
		if err != nil {
			return nil, errors.Wrap(err, "unable to marshal constellation for cache")
		}

		_, err = s.redis.Set(ctx, key, bSlice, time.Hour).Result()

		return constellation, errors.Wrap(err, "failed to cache constellation in redis")
	}

	// Constellation is not cached, the DB doesn't have this constellation, lets check ESI
	constellation, m := s.esi.GetUniverseConstellationsConstellationID(ctx, id)
	if m.IsErr() {
		return nil, m.Msg
	}

	if m.Code == http.StatusUnprocessableEntity {
		return nil, errors.New("invalid constellation received from ESI")
	}

	// ESI has the constellation. Lets insert it into the db, and cache it is redis
	_, err = s.CreateConstellation(ctx, constellation)
	if err != nil {
		return constellation, errors.Wrap(err, "unable to insert constellation into db")
	}

	byteSlice, err := json.Marshal(constellation)
	if err != nil {
		return constellation, errors.Wrap(err, "unable to marshal constellation for cache")
	}

	_, err = s.redis.Set(ctx, key, byteSlice, time.Hour).Result()

	return constellation, errors.Wrap(err, "failed to cache constellation in redis")
}

func (s *service) CreateConstellation(ctx context.Context, constellation *zrule.Constellation) (*zrule.Constellation, error) {
	return s.ConstellationRepository.CreateConstellation(ctx, constellation)
}
