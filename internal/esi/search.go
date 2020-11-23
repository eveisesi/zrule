package esi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type searchService interface {
	GetSearch(
		ctx context.Context, category, term string, strict bool,
	) ([]uint64, Meta)
}

func (s *service) GetSearch(ctx context.Context, category, term string, strict bool) ([]uint64, Meta) {

	query := url.Values{}
	query.Set("categories", category)
	query.Set("search", term)

	strStrict := "false"
	if strict {
		strStrict = "true"
	}

	query.Set("strict", strStrict)

	path := "/v2/search/"

	request := request{
		method: http.MethodGet,
		path:   path,
		query:  query.Encode(),
	}

	response, m := s.request(ctx, request)
	if m.IsErr() {
		return nil, m
	}

	var results = make(map[string][]uint64)

	switch m.Code {
	case http.StatusOK:
		err := json.Unmarshal(response, &results)
		if err != nil {
			m.Msg = fmt.Errorf("unable to unmarshal response body on request %s: %w", fmt.Sprintf("%s?%s", path, query.Encode()), err)
			return nil, m
		}
	case http.StatusBadRequest:

		esiErrorResponse := struct {
			Message string `json:"error"`
		}{}

		err = json.Unmarshal(response, &esiErrorResponse)
		if err != nil {
			m.Msg = fmt.Errorf("unable to unmarshal error from request %s: %w", fmt.Sprintf("%s?%s", path, query.Encode()), err)
			return nil, m
		}

		m.Msg = fmt.Errorf(esiErrorResponse.Message)
		return nil, m
	default:
		m.Msg = fmt.Errorf("unexpected status code %d received from ESI on request %s", m.Code, path)
		return nil, m
	}

	if _, ok := results[category]; !ok {
		m.Msg = fmt.Errorf("response does not contain results for specified category: %w", m.Msg)
		return nil, m
	}

	return results[category], m

}
