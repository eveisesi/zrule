package esi

// func (s *service) GetStatus(ctx context.Context) (*zrule.ServerStatus, Meta) {

// 	response, m := s.request(ctx, request{
// 		method: http.MethodGet,
// 		path:   "/v1/status",
// 	})
// 	if m.IsErr() {
// 		return nil, m
// 	}

// 	status := new(zrule.ServerStatus)
// 	err := json.Unmarshal(response, status)
// 	if err != nil {
// 		m.Msg = fmt.Errorf("unable to unmarshal response body on request %s: %w", path, err)
// 		return nil, m
// 	}

// 	return status, m

// }
