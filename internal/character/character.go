package character

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
	var key = fmt.Sprintf(zrule.REDIS_CHARACTER, id)

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

	character, err = s.CharacterRespository.Character(ctx, id)
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
	character, m := s.esi.GetCharactersCharacterID(ctx, id, "")
	if m.IsErr() {
		return nil, m.Msg
	}

	if m.Code == http.StatusUnprocessableEntity {
		return nil, errors.New("invalid character received from ESI, skipping create and cache")
	}

	// ESI has the character. Lets insert it into the db, and cache it is redis
	err = s.CharacterRespository.CreateCharacter(ctx, character)
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

// func (s *service) CharactersByCharacterIDs(ctx context.Context, ids []uint64) ([]*zrule.Character, error) {

// 	var characters = make([]*zrule.Character, 0)
// 	keys := make([]string, len(ids))
// 	for i, id := range ids {
// 		keys[i] = fmt.Sprintf(zrule.REDIS_CHARACTER, id)
// 	}

// 	results, err := s.redis.MGet(ctx, keys...).Result()
// 	if err != nil && err.Error() != "redis: nil" {
// 		return nil, errors.Wrap(err, "encountered error querying redis")
// 	}

// 	for i, resultInt := range results {
// 		if resultInt == nil {
// 			continue
// 		}

// 		switch result := resultInt.(type) {
// 		case string:
// 			if len(result) > 0 {
// 				var character = new(zrule.Character)
// 				err = json.Unmarshal([]byte(result), character)
// 				if err != nil {
// 					return nil, errors.Wrap(err, "unable to unmarshal character bytes into struct")
// 				}

// 				characters = append(characters, character)
// 			}
// 		default:
// 			panic(fmt.Sprintf("unexpected type received from redis. expected string, got %#T. redis key is %s", result, keys[i]))
// 		}
// 	}

// 	if len(ids) == len(characters) {
// 		return characters, nil
// 	}

// 	var missing []zrule.OpValue
// 	for _, id := range ids {
// 		found := false
// 		for _, character := range characters {
// 			if character.ID == id {
// 				found = true
// 				break
// 			}
// 		}
// 		if !found {
// 			missing = append(missing, id)
// 		}
// 	}

// 	if len(missing) == 0 {
// 		return characters, nil
// 	}

// 	dbTypes, err := s.Characters(ctx, zrule.NewInOperator("id", missing))
// 	if err != nil {
// 		return nil, errors.Wrap(err, "failed to query db for missing type ids")
// 	}

// 	keyMap := make(map[string]interface{})

// 	for _, character := range dbTypes {
// 		characters = append(characters, character)

// 		key := fmt.Sprintf(zrule.REDIS_CHARACTER, character.ID)

// 		byteSlice, err := json.Marshal(character)
// 		if err != nil {
// 			s.logger.WithError(err).WithField("id", character.ID).Error("unable to marshal character to slice of bytes")
// 			continue
// 		}

// 		keyMap[key] = string(byteSlice)

// 	}

// 	_, err = s.redis.MSet(ctx, keyMap).Result()
// 	if err != nil {
// 		return nil, errors.Wrap(err, "unable to cache characters in redis")
// 	}

// 	for i := range keyMap {
// 		s.redis.Expire(ctx, i, time.Hour)
// 	}

// 	return characters, nil

// }

func (s *service) CreateCharacter(ctx context.Context, character *zrule.Character) error {
	return s.CharacterRespository.CreateCharacter(ctx, character)
}

// func (s *service) UpdateExpired(ctx context.Context) {

// 	for {

// 		expired, err := s.Expired(ctx)
// 		if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
// 			s.logger.WithError(err).Error("Failed to fetch expired characters")
// 			return
// 		}

// 		if len(expired) == 0 {
// 			s.logger.Info("no expired characters found")
// 			time.Sleep(time.Minute * 5)
// 			continue
// 		}

// 		s.logger.WithField("count", len(expired)).Info("updating expired characters")

// 		for _, character := range expired {

// 			// s.tracker.Watchman(ctx)
// 			entry := s.logger.WithContext(ctx).WithField("character_id", character.ID)

// 			txn := s.newrelic.StartTransaction("update-expired-characters")
// 			txn.AddAttribute("characterID", character.ID)
// 			ctx = newrelic.NewContext(ctx, txn)

// 			newCharacter, m := s.esi.GetCharactersCharacterID(ctx, character.ID, character.Etag)
// 			if m.IsErr() {
// 				txn.NoticeError(m.Msg)
// 				txn.End()
// 				entry.WithError(m.Msg).Error("failed to fetch character from esi")
// 				continue
// 			}

// 			entry = entry.WithField("status_code", m.Code)
// 			txn.AddAttribute("status_code", m.Code)
// 			switch m.Code {
// 			case http.StatusInternalServerError, http.StatusBadRequest, http.StatusNotFound, http.StatusUnprocessableEntity:
// 				err = errors.New("bad status code received from ESI")
// 				txn.NoticeError(err)
// 				entry.WithError(err).Errorln()

// 				character.CachedUntil = time.Now().Add(time.Minute * 2).Unix()
// 				character.UpdateError++

// 				err = s.UpdateCharacter(ctx, character.ID, character)
// 			case http.StatusNotModified:

// 				if character.NotModifiedCount >= 2 && character.UpdatePriority <= 3 {
// 					character.NotModifiedCount = 0
// 					character.UpdatePriority++
// 				} else {
// 					character.NotModifiedCount++
// 				}

// 				character.UpdateError = 0
// 				character.CachedUntil = time.Unix(newCharacter.CachedUntil, 0).AddDate(0, 0, int(character.UpdatePriority)).Unix()
// 				character.Etag = newCharacter.Etag

// 				err = s.UpdateCharacter(ctx, character.ID, character)
// 			case http.StatusOK:
// 				err = s.UpdateCharacter(ctx, character.ID, newCharacter)
// 			default:
// 				entry.WithField("status_code", m.Code).Error("unaccounted for status code received from esi service")
// 			}
// 			if err != nil {
// 				txn.NoticeError(err)
// 				entry.WithError(err).Error("failed to update character")
// 			}

// 			txn.End()
// 			time.Sleep(time.Millisecond * 100)
// 		}
// 		s.logger.WithContext(ctx).WithField("count", len(expired)).Info("characters successfully updated")
// 		time.Sleep(time.Second)

// 	}

// }
