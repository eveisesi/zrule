package esi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/eveisesi/zrule"
)

type factionService interface {
	GetUniverseFactions(ctx context.Context) ([]*zrule.Faction, Meta)
}

func (s *service) GetUniverseFactions(ctx context.Context) ([]*zrule.Faction, Meta) {

	path := "/v2/universe/factions/"

	request := request{
		method:  http.MethodGet,
		path:    path,
		headers: make(map[string]string),
	}

	response, m := s.request(ctx, request)
	if m.IsErr() {
		return nil, m
	}

	factions := make([]*zrule.Faction, 0)

	switch m.Code {
	case 200:
		err := json.Unmarshal(response, &factions)
		if err != nil {
			m.Msg = fmt.Errorf("unable to unmarshal response body on request %s: %w", path, err)
			return nil, m
		}

		for _, faction := range factions {
			faction.ID = faction.ESIID
		}
	default:
		m.Msg = fmt.Errorf("unexpected status code received from ESI on request %s", path)

	}

	return factions, m

}
