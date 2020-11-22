package esi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/eveisesi/zrule"
)

type characterService interface {
	GetCharactersCharacterID(ctx context.Context, id uint64) (*zrule.Character, Meta)
}

func validateCharacter(character *zrule.Character) bool {
	return character.Name != ""
}

// GetCharactersCharacterID makes a HTTP GET Request to the /characters/{character_id} endpoint
// for information about the provided character
//
// Documentation: https://esi.evetech.net/ui/#/Character/get_characters_character_id
// Version: v4
// Cache: 86400 sec (24 Hour)
func (s *service) GetCharactersCharacterID(ctx context.Context, id uint64) (*zrule.Character, Meta) {

	path := fmt.Sprintf("/v4/characters/%d/", id)

	request := request{
		method:  http.MethodGet,
		path:    path,
		headers: make(map[string]string),
	}

	response, m := s.request(ctx, request)
	if m.IsErr() {
		return nil, m
	}

	character := new(zrule.Character)

	switch m.Code {
	case 200:
		err := json.Unmarshal(response, character)
		if err != nil {
			m.Msg = fmt.Errorf("unable to unmarshal response body on request %s: %w", path, err)
			return nil, m
		}

		character.ID = id

		if !validateCharacter(character) {
			m.Code = http.StatusUnprocessableEntity
			return nil, m
		}
	default:
		m.Msg = fmt.Errorf("unexpected status code received from ESI on request %s", path)

	}

	return character, m
}
