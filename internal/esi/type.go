package esi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/eveisesi/zrule"
)

type itemService interface {
	GetUniverseGroups(ctx context.Context, page *uint) ([]uint, Meta)
	GetUniverseGroupsGroupID(ctx context.Context, id uint) (*zrule.ItemGroup, Meta)
	GetUniverseTypes(ctx context.Context, page *uint) ([]uint, Meta)
	GetUniverseTypesTypeID(ctx context.Context, id uint) (*zrule.Item, Meta)
}

func (s *service) GetUniverseGroups(ctx context.Context, page *uint) ([]uint, Meta) {
	path := "/v1/universe/groups/"

	request := request{
		method: http.MethodGet,
		path:   path,
	}

	if page != nil {
		values := url.Values{}
		values.Add("page", strconv.Itoa(int(*page)))
		request.query = values.Encode()
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

func (s *service) GetUniverseGroupsGroupID(ctx context.Context, id uint) (*zrule.ItemGroup, Meta) {

	path := fmt.Sprintf("/v1/universe/groups/%d/", id)

	request := request{
		method: http.MethodGet,
		path:   path,
	}

	response, m := s.request(ctx, request)
	if m.IsErr() {
		return nil, m
	}

	group := new(zrule.ItemGroup)

	switch m.Code {
	case http.StatusOK:
		err := json.Unmarshal(response, group)
		if err != nil {
			m.Msg = fmt.Errorf("unable to unmarshal response body on request %s: %w", path, err)
			return nil, m
		}

		group.ID = id
	default:
		m.Msg = fmt.Errorf("unexpected status code received from ESI on request %s", path)
	}

	return group, m

}

func (s *service) GetUniverseTypes(ctx context.Context, page *uint) ([]uint, Meta) {
	path := "/v1/universe/types/"

	request := request{
		method: http.MethodGet,
		path:   path,
	}

	if page != nil {
		values := url.Values{}
		values.Add("page", strconv.Itoa(int(*page)))
		request.query = values.Encode()
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

func (s *service) GetUniverseTypesTypeID(ctx context.Context, id uint) (*zrule.Item, Meta) {

	path := fmt.Sprintf("/v3/universe/types/%d/", id)

	request := request{
		method: http.MethodGet,
		path:   path,
	}

	response, m := s.request(ctx, request)
	if m.IsErr() {
		return nil, m
	}

	item := new(zrule.Item)

	switch m.Code {
	case http.StatusOK:
		err := json.Unmarshal(response, item)
		if err != nil {
			m.Msg = fmt.Errorf("unable to unmarshal response body on request %s: %w", path, err)
			return nil, m
		}

		item.ID = id
	default:
		m.Msg = fmt.Errorf("unexpected status code received from ESI on request %s", path)
	}

	return item, m

}
