package esi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/eveisesi/zrule"
)

type solarSystemService interface {
	GetUniverseSolarSystems(ctx context.Context) ([]uint, Meta)
	GetUniverseSolarSystemsSolarSystemID(ctx context.Context, id uint) (*zrule.SolarSystem, Meta)
}

func (s *service) GetUniverseSolarSystems(ctx context.Context) ([]uint, Meta) {

	var path = "/v1/universe/systems/"

	request := request{
		method:  http.MethodGet,
		path:    path,
		headers: make(map[string]string),
	}

	response, m := s.request(ctx, request)
	if m.IsErr() {
		return nil, m
	}

	var ids = make([]uint, 0)

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

func (s *service) GetUniverseSolarSystemsSolarSystemID(ctx context.Context, id uint) (*zrule.SolarSystem, Meta) {

	var path = fmt.Sprintf("/v4/universe/systems/%d/", id)

	request := request{
		method: http.MethodGet,
		path:   path,
	}

	response, m := s.request(ctx, request)
	if m.IsErr() {
		return nil, m
	}

	var solarSystem = new(zrule.SolarSystem)

	switch m.Code {
	case http.StatusOK:
		err := json.Unmarshal(response, solarSystem)
		if err != nil {
			m.Msg = fmt.Errorf("unable to unmarshal response body on request %s: %w", path, err)
			return nil, m
		}

		solarSystem.ID = id
	default:
		m.Msg = fmt.Errorf("unexpected status code received from ESI on request %s", path)
	}

	return solarSystem, m
}
