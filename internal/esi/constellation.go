package esi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/eveisesi/zrule"
)

type constellationService interface {
	GetUniverseConstellations(ctx context.Context) ([]uint, Meta)
	GetUniverseConstellationsConstellationID(ctx context.Context, id uint) (*zrule.Constellation, Meta)
}

func (s *service) GetUniverseConstellations(ctx context.Context) ([]uint, Meta) {

	path := "/v1/universe/constellations/"

	request := request{
		method:  http.MethodGet,
		path:    path,
		headers: make(map[string]string),
	}

	response, m := s.request(ctx, request)
	if m.IsErr() {
		return nil, m
	}

	ids := make([]uint, 0)

	switch m.Code {
	case http.StatusOK:
		err := json.Unmarshal(response, &ids)
		if err != nil {
			m.Msg = fmt.Errorf("unable to unmarshal response body on request %s: %w", path, err)
			return nil, m
		}
	default:
		m.Msg = fmt.Errorf("unexpected status code received from ESI on request %s", path)
	}

	return ids, m

}

func (s *service) GetUniverseConstellationsConstellationID(ctx context.Context, id uint) (*zrule.Constellation, Meta) {

	path := fmt.Sprintf("/v1/universe/constellations/%d/", id)

	request := request{
		method:  http.MethodGet,
		path:    path,
		headers: make(map[string]string),
	}

	response, m := s.request(ctx, request)
	if m.IsErr() {
		return nil, m
	}

	constellation := new(zrule.Constellation)

	switch m.Code {
	case http.StatusOK:
		err := json.Unmarshal(response, constellation)
		if err != nil {
			m.Msg = fmt.Errorf("unable to unmarshal response body on request %s: %w", path, err)
			return nil, m
		}

		constellation.ID = id

	default:
		m.Msg = fmt.Errorf("unexpected status code received from ESI on request %s", path)
	}

	return constellation, m

}
