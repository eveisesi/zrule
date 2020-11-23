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

func (s *service) ItemGroup(ctx context.Context, id uint) (*zrule.ItemGroup, error) {
	var itemGroup = new(zrule.ItemGroup)
	var key = fmt.Sprintf(zrule.CACHE_ITEMGROUP, id)

	result, err := s.redis.Get(ctx, key).Bytes()
	if err != nil && err.Error() != "redis: nil" {
		return nil, err
	}

	if len(result) > 0 {

		err = json.Unmarshal(result, itemGroup)
		if err != nil {
			return nil, errors.Wrap(err, "unable to unmarshal itemGroup from redis")
		}
		return itemGroup, nil
	}

	itemGroup, err = s.ItemGroupRepository.ItemGroup(ctx, id)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return nil, errors.Wrap(err, "unable to query database for itemGroup")
	}

	if err == nil {
		bSlice, err := json.Marshal(itemGroup)
		if err != nil {
			return nil, errors.Wrap(err, "unable to marshal itemGroup for cache")
		}

		_, err = s.redis.Set(ctx, key, bSlice, time.Hour).Result()

		return itemGroup, errors.Wrap(err, "failed to cache itemGroup in redis")
	}

	// ItemGroup is not cached, the DB doesn't have this itemGroup, lets check ESI
	itemGroup, m := s.esi.GetUniverseGroupsGroupID(ctx, id)
	if m.IsErr() {
		return nil, m.Msg
	}

	if m.Code == http.StatusUnprocessableEntity {
		return nil, errors.New("invalid itemGroup received from ESI")
	}

	// ESI has the itemGroup. Lets insert it into the db, and cache it is redis
	_, err = s.CreateItemGroup(ctx, itemGroup)
	if err != nil {
		return itemGroup, errors.Wrap(err, "unable to insert itemGroup into db")
	}

	byteSlice, err := json.Marshal(itemGroup)
	if err != nil {
		return itemGroup, errors.Wrap(err, "unable to marshal itemGroup for cache")
	}

	_, err = s.redis.Set(ctx, key, byteSlice, time.Hour).Result()

	return itemGroup, errors.Wrap(err, "failed to cache itemGroup in redis")
}

func (s *service) CreateItemGroup(ctx context.Context, itemGroup *zrule.ItemGroup) (*zrule.ItemGroup, error) {
	return s.ItemGroupRepository.CreateItemGroup(ctx, itemGroup)
}
