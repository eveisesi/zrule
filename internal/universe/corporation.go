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

func (s *service) Corporation(ctx context.Context, id uint) (*zrule.Corporation, error) {
	var corporation = new(zrule.Corporation)
	var key = fmt.Sprintf(zrule.CACHE_CORPORATION, id)

	result, err := s.redis.Get(ctx, key).Bytes()
	if err != nil && err.Error() != "redis: nil" {
		return nil, err
	}

	if len(result) > 0 {

		err = json.Unmarshal(result, corporation)
		if err != nil {
			return nil, errors.Wrap(err, "unable to unmarshal corporation from redis")
		}
		return corporation, nil
	}

	corporation, err = s.CorporationRepository.Corporation(ctx, id)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return nil, errors.Wrap(err, "unable to query database for corporation")
	}

	if err == nil {
		bSlice, err := json.Marshal(corporation)
		if err != nil {
			return nil, errors.Wrap(err, "unable to marshal corporation for cache")
		}

		_, err = s.redis.Set(ctx, key, bSlice, time.Hour).Result()

		return corporation, errors.Wrap(err, "failed to cache corporation in redis")
	}

	// Corporation is not cached, the DB doesn't have this corporation, lets check ESI
	corporation, m := s.esi.GetCorporationsCorporationID(ctx, id)
	if m.IsErr() {
		return nil, m.Msg
	}

	if m.Code == http.StatusUnprocessableEntity {
		return nil, errors.New("invalid corporation received from ESI")
	}

	// ESI has the corporation. Lets insert it into the db, and cache it is redis
	_, err = s.CreateCorporation(ctx, corporation)
	if err != nil {
		return corporation, errors.Wrap(err, "unable to insert corporation into db")
	}

	byteSlice, err := json.Marshal(corporation)
	if err != nil {
		return corporation, errors.Wrap(err, "unable to marshal corporation for cache")
	}

	_, err = s.redis.Set(ctx, key, byteSlice, time.Hour).Result()

	return corporation, errors.Wrap(err, "failed to cache corporation in redis")
}

func (s *service) CreateCorporation(ctx context.Context, corporation *zrule.Corporation) (*zrule.Corporation, error) {
	return s.CorporationRepository.CreateCorporation(ctx, corporation)
}
