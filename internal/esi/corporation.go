package esi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/eveisesi/zrule"
)

type corporationService interface {
	GetCorporationsCorporationID(ctx context.Context, id uint) (*zrule.Corporation, Meta)
}

func validateCorporation(corporation *zrule.Corporation) bool {
	return corporation.Name != "" && corporation.Ticker != ""
}

// GetCorporationsCorporationID makes a HTTP GET Request to the /corporations/{corporation_id} endpoint
// for information about the provided corporation
//
// Documentation: https://esi.evetech.net/ui/#/Corporation/get_corporations_corporation_id
// Version: v4
// Cache: 3600 sec (1 Hour)
func (s *service) GetCorporationsCorporationID(ctx context.Context, id uint) (*zrule.Corporation, Meta) {

	path := fmt.Sprintf("/v4/corporations/%d/", id)

	request := request{
		method:  http.MethodGet,
		path:    path,
		headers: make(map[string]string),
	}

	response, m := s.request(ctx, request)
	if m.IsErr() {
		return nil, m
	}

	corporation := new(zrule.Corporation)

	switch m.Code {
	case 200:
		err := json.Unmarshal(response, corporation)
		if err != nil {
			m.Msg = fmt.Errorf("unable to unmarshal response body on request %s: %w", path, err)
			return nil, m
		}

		corporation.ID = id
		if !validateCorporation(corporation) {
			m.Code = http.StatusUnprocessableEntity
			return nil, m
		}
	default:
		m.Msg = fmt.Errorf("unexpected status code received from ESI on request %s", path)

	}

	return corporation, m

}
