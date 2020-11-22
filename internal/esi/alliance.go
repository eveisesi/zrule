package esi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/eveisesi/zrule"
)

type allianceService interface {
	GetAlliancesAllianceID(ctx context.Context, id uint) (*zrule.Alliance, Meta)
}

func validateAlliance(alliance *zrule.Alliance) bool {
	return alliance.Name != "" && alliance.Ticker != ""
}

// GetAlliancesAllianceID makes a HTTP GET Request to the /alliances/{alliance_id} endpoint
// for information about the provided alliance
//
// Documentation: https://esi.evetech.net/ui/#/Alliance/get_alliances_alliance_id
// Version: v3
// Cache: 3600 sec (1 Hour)
func (s *service) GetAlliancesAllianceID(ctx context.Context, id uint) (*zrule.Alliance, Meta) {

	path := fmt.Sprintf("/v3/alliances/%d/", id)

	request := request{
		method:  http.MethodGet,
		path:    path,
		headers: make(map[string]string),
	}

	response, m := s.request(ctx, request)
	if m.IsErr() {
		return nil, m
	}

	alliance := new(zrule.Alliance)

	switch m.Code {
	case 200:
		err = json.Unmarshal(response, alliance)
		if err != nil {
			m.Msg = fmt.Errorf("unable to unmarshal response body on request %s: %w", path, err)
			return nil, m
		}

		alliance.ID = id

		if !validateAlliance(alliance) {
			m.Code = http.StatusUnprocessableEntity
			return nil, m
		}
	default:
		m.Msg = fmt.Errorf("unexpected status code received from ESI on request %s", path)

	}

	return alliance, m
}
