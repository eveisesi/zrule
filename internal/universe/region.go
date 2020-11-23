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

func (s *service) Region(ctx context.Context, id uint) (*zrule.Region, error) {
	var region = new(zrule.Region)
	var key = fmt.Sprintf(zrule.CACHE_REGION, id)

	result, err := s.redis.Get(ctx, key).Bytes()
	if err != nil && err.Error() != "redis: nil" {
		return nil, err
	}

	if len(result) > 0 {

		err = json.Unmarshal(result, region)
		if err != nil {
			return nil, errors.Wrap(err, "unable to unmarshal region from redis")
		}
		return region, nil
	}

	region, err = s.RegionRepository.Region(ctx, id)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return nil, errors.Wrap(err, "unable to query database for region")
	}

	if err == nil {
		bSlice, err := json.Marshal(region)
		if err != nil {
			return nil, errors.Wrap(err, "unable to marshal region for cache")
		}

		_, err = s.redis.Set(ctx, key, bSlice, time.Hour).Result()

		return region, errors.Wrap(err, "failed to cache region in redis")
	}

	// Region is not cached, the DB doesn't have this region, lets check ESI
	region, m := s.esi.GetUniverseRegionsRegionID(ctx, id)
	if m.IsErr() {
		return nil, m.Msg
	}

	if m.Code == http.StatusUnprocessableEntity {
		return nil, errors.New("invalid region received from ESI")
	}

	// ESI has the region. Lets insert it into the db, and cache it is redis
	_, err = s.CreateRegion(ctx, region)
	if err != nil {
		return region, errors.Wrap(err, "unable to insert region into db")
	}

	byteSlice, err := json.Marshal(region)
	if err != nil {
		return region, errors.Wrap(err, "unable to marshal region for cache")
	}

	_, err = s.redis.Set(ctx, key, byteSlice, time.Hour).Result()

	return region, errors.Wrap(err, "failed to cache region in redis")
}

func (s *service) CreateRegion(ctx context.Context, region *zrule.Region) (*zrule.Region, error) {
	return s.RegionRepository.CreateRegion(ctx, region)
}
