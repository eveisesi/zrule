package esi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/eveisesi/zrule"
)

type regionService interface {
	GetUniverseRegions(ctx context.Context) ([]uint, Meta)
	GetUniverseRegionsRegionID(ctx context.Context, id uint) (*zrule.Region, Meta)
}

func (s *service) GetUniverseRegions(ctx context.Context) ([]uint, Meta) {

	path := "/v1/universe/regions/"

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

func (s *service) GetUniverseRegionsRegionID(ctx context.Context, id uint) (*zrule.Region, Meta) {

	path := fmt.Sprintf("/v1/universe/regions/%d/", id)

	request := request{
		method:  http.MethodGet,
		path:    path,
		headers: make(map[string]string),
	}

	response, m := s.request(ctx, request)
	if m.IsErr() {
		return nil, m
	}

	region := new(zrule.Region)

	switch m.Code {
	case http.StatusOK:
		err := json.Unmarshal(response, region)
		if err != nil {
			m.Msg = fmt.Errorf("unable to unmarshal response body on request %s: %w", path, err)
			return nil, m
		}
	default:
		m.Msg = fmt.Errorf("unexpected status code received from ESI on request %s", path)
	}

	return region, m

}
