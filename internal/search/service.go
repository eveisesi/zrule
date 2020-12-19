package search

import (
	"encoding/json"
	"fmt"

	"github.com/RediSearch/redisearch-go/redisearch"
	"github.com/sirupsen/logrus"
)

type Entity struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
}

type Service interface {
	InitializeAutocompleter(key Key, suggestions []*Entity) error
	Search(key Key, term string) ([]*Entity, error)
}

type service struct {
	address string
	buckets map[string]*redisearch.Autocompleter
	logger  *logrus.Logger
}

func NewService(logger *logrus.Logger, address string) Service {

	return &service{
		address,
		make(map[string]*redisearch.Autocompleter),
		logger,
	}

}

func (s *service) InitializeAutocompleter(key Key, entities []*Entity) error {

	if !key.Valid() {
		return fmt.Errorf("invalid key of %s used", key)
	}

	if existing, ok := s.buckets[key.String()]; ok {
		if existing != nil {
			err := existing.Delete()
			if err != nil {
				return fmt.Errorf("failed to flush existing autocompleter: %w", err)
			}
			delete(s.buckets, key.String())
		}
	}

	s.buckets[key.String()] = redisearch.NewAutocompleter(s.address, key.String())

	suggestions := make([]redisearch.Suggestion, len(entities))

	for i, entity := range entities {
		payload, err := json.Marshal(entity)
		if err != nil {
			return fmt.Errorf("failed to marshal entites for autocompleter: %w", err)
		}

		suggestions[i] = redisearch.Suggestion{
			Term:    entity.Name,
			Score:   1,
			Payload: string(payload),
		}

	}

	err := s.buckets[key.String()].AddTerms(suggestions...)
	if err != nil {
		return fmt.Errorf("failed to add suggestions to autocompleter %s: %w", key, err)
	}

	return nil
}

func (s *service) Search(key Key, term string) ([]*Entity, error) {

	if _, ok := s.buckets[key.String()]; !ok {
		return nil, fmt.Errorf("there is no autocompleter for that key")
	}

	if s.buckets[key.String()] == nil {
		return nil, fmt.Errorf("specified autocompleter has not been initialized")
	}

	suggestions, err := s.buckets[key.String()].SuggestOpts(term, redisearch.SuggestOptions{
		Num:          20,
		Fuzzy:        false,
		WithPayloads: true,
		WithScores:   false,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to execute search with autocompleter: %w", err)
	}

	var results = make([]*Entity, len(suggestions))
	for i, suggestion := range suggestions {
		var x = new(Entity)
		err = json.Unmarshal([]byte(suggestion.Payload), &x)
		if err != nil {
			s.logger.WithError(err).Error("failed to decode suggestion")
			continue
		}

		results[i] = x
	}

	return results, nil

}
