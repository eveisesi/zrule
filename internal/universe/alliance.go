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

func (s *service) Alliance(ctx context.Context, id uint) (*zrule.Alliance, error) {
	var alliance = new(zrule.Alliance)
	var key = fmt.Sprintf(zrule.CACHE_ALLIANCE, id)

	result, err := s.redis.Get(ctx, key).Bytes()
	if err != nil && err.Error() != "redis: nil" {
		return nil, err
	}

	if len(result) > 0 {

		err = json.Unmarshal(result, alliance)
		if err != nil {
			return nil, errors.Wrap(err, "unable to unmarshal alliance from redis")
		}
		return alliance, nil
	}

	alliance, err = s.AllianceRepository.Alliance(ctx, id)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return nil, errors.Wrap(err, "unable to query database for alliance")
	}

	if err == nil {
		bSlice, err := json.Marshal(alliance)
		if err != nil {
			return nil, errors.Wrap(err, "unable to marshal alliance for cache")
		}

		_, err = s.redis.Set(ctx, key, bSlice, time.Hour).Result()

		return alliance, errors.Wrap(err, "failed to cache alliance in redis")
	}

	// Alliance is not cached, the DB doesn't have this alliance, lets check ESI
	alliance, m := s.esi.GetAlliancesAllianceID(ctx, id)
	if m.IsErr() {
		return nil, m.Msg
	}

	if m.Code == http.StatusUnprocessableEntity {
		return nil, errors.New("invalid alliance received from ESI")
	}

	// ESI has the alliance. Lets insert it into the db, and cache it is redis
	_, err = s.CreateAlliance(ctx, alliance)
	if err != nil {
		return alliance, errors.Wrap(err, "unable to insert alliance into db")
	}

	byteSlice, err := json.Marshal(alliance)
	if err != nil {
		return alliance, errors.Wrap(err, "unable to marshal alliance for cache")
	}

	_, err = s.redis.Set(ctx, key, byteSlice, time.Hour).Result()

	return alliance, errors.Wrap(err, "failed to cache alliance in redis")
}

func (s *service) CreateAlliance(ctx context.Context, alliance *zrule.Alliance) (*zrule.Alliance, error) {
	return s.AllianceRepository.CreateAlliance(ctx, alliance)
}
