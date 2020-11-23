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

func (s *service) Character(ctx context.Context, id uint64) (*zrule.Character, error) {
	var character = new(zrule.Character)
	var key = fmt.Sprintf(zrule.CACHE_CHARACTER, id)

	result, err := s.redis.Get(ctx, key).Bytes()
	if err != nil && err.Error() != "redis: nil" {
		return nil, err
	}

	if len(result) > 0 {

		err = json.Unmarshal(result, character)
		if err != nil {
			return nil, errors.Wrap(err, "unable to unmarshal character from redis")
		}
		return character, nil
	}

	character, err = s.CharacterRepository.Character(ctx, id)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return nil, errors.Wrap(err, "unable to query database for character")
	}

	if err == nil {
		bSlice, err := json.Marshal(character)
		if err != nil {
			return nil, errors.Wrap(err, "unable to marshal character for cache")
		}

		_, err = s.redis.Set(ctx, key, bSlice, time.Hour).Result()

		return character, errors.Wrap(err, "failed to cache character in redis")
	}

	// Character is not cached, the DB doesn't have this character, lets check ESI
	character, m := s.esi.GetCharactersCharacterID(ctx, id)
	if m.IsErr() {
		return nil, m.Msg
	}

	if m.Code == http.StatusUnprocessableEntity {
		return nil, errors.New("invalid character received from ESI")
	}

	// ESI has the character. Lets insert it into the db, and cache it is redis
	_, err = s.CharacterRepository.CreateCharacter(ctx, character)
	if err != nil {
		return character, errors.Wrap(err, "unable to insert character into db")
	}

	byteSlice, err := json.Marshal(character)
	if err != nil {
		return character, errors.Wrap(err, "unable to marshal character for cache")
	}

	_, err = s.redis.Set(ctx, key, byteSlice, time.Hour).Result()

	return character, errors.Wrap(err, "failed to cache character in redis")
}

func (s *service) CreateCharacter(ctx context.Context, character *zrule.Character) (*zrule.Character, error) {
	return s.CharacterRepository.CreateCharacter(ctx, character)
}
