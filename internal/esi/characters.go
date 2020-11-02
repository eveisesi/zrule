package esi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/eveisesi/zrule"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
)

type Character struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
}

func (r Character) validate() bool {
	if r.Name == "" {
		return false
	}
	return true
}

// GetCharactersCharacterID makes a HTTP GET Request to the /characters/{character_id} endpoint
// for information about the provided character
//
// Documentation: https://esi.evetech.net/ui/#/Character/get_characters_character_id
// Version: v4
// Cache: 86400 sec (24 Hour)
func (s *service) GetCharactersCharacterID(ctx context.Context, id uint64, etag string) (*zrule.Character, Meta) {

	path := fmt.Sprintf("/v4/characters/%d/", id)
	headers := make(map[string]string)

	if etag != "" {
		headers["If-None-Match"] = etag
	}

	request := request{
		method:  http.MethodGet,
		path:    path,
		headers: headers,
	}

	response, m := s.request(ctx, request)
	if m.IsErr() {
		return nil, m
	}

	esiCharacter := new(Character)

	switch m.Code {
	case 200:
		err := json.Unmarshal(response, esiCharacter)
		if err != nil {
			m.Msg = errors.Wrapf(err, "unable to unmarshal response body on request %s", path)
			return nil, m
		}

		esiCharacter.ID = id

		if !esiCharacter.validate() {
			m.Code = http.StatusUnprocessableEntity
			return nil, m
		}

	}

	character := new(zrule.Character)
	err = copier.Copy(character, esiCharacter)
	if err != nil {
		m.Msg = err
		return nil, m
	}

	return character, m
}
