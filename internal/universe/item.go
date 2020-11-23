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

func (s *service) Item(ctx context.Context, id uint) (*zrule.Item, error) {
	var item = new(zrule.Item)
	var key = fmt.Sprintf(zrule.CACHE_ITEM, id)

	result, err := s.redis.Get(ctx, key).Bytes()
	if err != nil && err.Error() != "redis: nil" {
		return nil, err
	}

	if len(result) > 0 {

		err = json.Unmarshal(result, item)
		if err != nil {
			return nil, errors.Wrap(err, "unable to unmarshal item from redis")
		}
		return item, nil
	}

	item, err = s.ItemRepository.Item(ctx, id)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return nil, errors.Wrap(err, "unable to query database for item")
	}

	if err == nil {
		bSlice, err := json.Marshal(item)
		if err != nil {
			return nil, errors.Wrap(err, "unable to marshal item for cache")
		}

		_, err = s.redis.Set(ctx, key, bSlice, time.Hour).Result()

		return item, errors.Wrap(err, "failed to cache item in redis")
	}

	// Item is not cached, the DB doesn't have this item, lets check ESI
	item, m := s.esi.GetUniverseTypesTypeID(ctx, id)
	if m.IsErr() {
		return nil, m.Msg
	}

	if m.Code == http.StatusUnprocessableEntity {
		return nil, errors.New("invalid item received from ESI")
	}

	// ESI has the item. Lets insert it into the db, and cache it is redis
	_, err = s.CreateItem(ctx, item)
	if err != nil {
		return item, errors.Wrap(err, "unable to insert item into db")
	}

	byteSlice, err := json.Marshal(item)
	if err != nil {
		return item, errors.Wrap(err, "unable to marshal item for cache")
	}

	_, err = s.redis.Set(ctx, key, byteSlice, time.Hour).Result()

	return item, errors.Wrap(err, "failed to cache item in redis")
}

func (s *service) CreateItem(ctx context.Context, item *zrule.Item) (*zrule.Item, error) {
	return s.ItemRepository.CreateItem(ctx, item)
}
