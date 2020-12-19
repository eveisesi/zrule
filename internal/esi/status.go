package esi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/eveisesi/zrule"
)

type statusService interface {
	GetStatus(ctx context.Context) (*zrule.ServerStatus, Meta)
}

func (s *service) GetStatus(ctx context.Context) (*zrule.ServerStatus, Meta) {

	path := "/v1/status"

	response, m := s.request(ctx, request{
		method: http.MethodGet,
		path:   path,
	})
	if m.IsErr() {
		return nil, m
	}

	status := new(zrule.ServerStatus)

	switch m.Code {
	case http.StatusOK:
		err := json.Unmarshal(response, status)
		if err != nil {
			m.Msg = fmt.Errorf("unable to unmarshal response body on request %s: %w", path, err)
			return nil, m
		}
	default:
		m.Msg = fmt.Errorf("non ok status code received from ESI %d, esi is unable to accomodate our initialization requests right now", m.Code)
	}

	return status, m

}
