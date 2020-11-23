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

func (s *service) SolarSystem(ctx context.Context, id uint) (*zrule.SolarSystem, error) {
	var solarSystem = new(zrule.SolarSystem)
	var key = fmt.Sprintf(zrule.CACHE_SOLARSYSTEM, id)

	result, err := s.redis.Get(ctx, key).Bytes()
	if err != nil && err.Error() != "redis: nil" {
		return nil, err
	}

	if len(result) > 0 {

		err = json.Unmarshal(result, solarSystem)
		if err != nil {
			return nil, errors.Wrap(err, "unable to unmarshal solarSystem from redis")
		}
		return solarSystem, nil
	}

	solarSystem, err = s.SolarSystemRepository.SolarSystem(ctx, id)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return nil, errors.Wrap(err, "unable to query database for solarSystem")
	}

	if err == nil {
		bSlice, err := json.Marshal(solarSystem)
		if err != nil {
			return nil, errors.Wrap(err, "unable to marshal solarSystem for cache")
		}

		_, err = s.redis.Set(ctx, key, bSlice, time.Hour).Result()

		return solarSystem, errors.Wrap(err, "failed to cache solarSystem in redis")
	}

	// SolarSystem is not cached, the DB doesn't have this solarSystem, lets check ESI
	solarSystem, m := s.esi.GetUniverseSolarSystemsSolarSystemID(ctx, id)
	if m.IsErr() {
		return nil, m.Msg
	}

	if m.Code == http.StatusUnprocessableEntity {
		return nil, errors.New("invalid solarSystem received from ESI")
	}

	// ESI has the solarSystem. Lets insert it into the db, and cache it is redis
	_, err = s.CreateSolarSystem(ctx, solarSystem)
	if err != nil {
		return solarSystem, errors.Wrap(err, "unable to insert solarSystem into db")
	}

	byteSlice, err := json.Marshal(solarSystem)
	if err != nil {
		return solarSystem, errors.Wrap(err, "unable to marshal solarSystem for cache")
	}

	_, err = s.redis.Set(ctx, key, byteSlice, time.Hour).Result()

	return solarSystem, errors.Wrap(err, "failed to cache solarSystem in redis")
}

func (s *service) CreateSolarSystem(ctx context.Context, solarSystem *zrule.SolarSystem) (*zrule.SolarSystem, error) {
	return s.SolarSystemRepository.CreateSolarSystem(ctx, solarSystem)
}
