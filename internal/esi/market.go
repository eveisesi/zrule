package esi

// type MarketPrices struct {
// 	AdjustedPrice float64 `json:"adjusted_price"`
// 	AveragePrice  float64 `json:"average_price"`
// 	TypeID        uint    `json:"type_id"`
// }

// type MarketGroup struct {
// 	MarketGroupID uint   `json:"market_group_id"`
// 	ParentGroupID uint   `json:"parent_group_id"`
// 	Name          string `json:"name"`
// 	Description   string `json:"description"`
// 	Types         []uint `json:"types"`
// }

// func (s *service) GetMarketGroups(ctx context.Context) ([]int, Meta) {

// 	response, m := s.request(ctx, request{
// 		method: http.MethodGet,
// 		path:   "/v1/markets/groups/",
// 	})
// 	if m.IsErr() {
// 		return nil, m
// 	}

// 	ids := make([]int, 0)
// 	err := json.Unmarshal(response, &ids)
// 	if err != nil {
// 		m.Msg = errors.Wrapf(err, "unable to unmarshal response body on request %s", "/v1/markets/groups/")
// 		return nil, m
// 	}

// 	return ids, m

// }

// func (s *service) GetMarketGroupsMarketGroupID(ctx context.Context, id int) (*zrule.MarketGroup, Meta) {

// 	path := fmt.Sprintf("/v1/markets/groups/%d", id)

// 	response, m := s.request(ctx, request{
// 		method: http.MethodGet,
// 		path:   path,
// 	})
// 	if m.IsErr() {
// 		return nil, m
// 	}

// 	esiGroup := new(MarketGroup)
// 	err := json.Unmarshal(response, esiGroup)
// 	if err != nil {
// 		m.Msg = errors.Wrapf(err, "unable to unmarshal response body on request %s", path)
// 		return nil, m
// 	}

// 	var group = new(zrule.MarketGroup)
// 	err = copier.Copy(group, esiGroup)
// 	if err != nil {
// 		m.Msg = errors.Wrap(err, "unable to copy esi type MarketGroup to neo MarketGroup ")
// 	}

// 	return group, m

// }

// func (s *service) GetMarketsRegionIDHistory(ctx context.Context, regionID uint, typeID uint) ([]*zrule.HistoricalRecord, Meta) {

// 	path := fmt.Sprintf("/v1/markets/%d/history/", regionID)

// 	query := url.Values{}
// 	query.Set("type_id", strconv.Itoa(int(typeID)))

// 	response, m := s.request(ctx, request{
// 		method: http.MethodGet,
// 		path:   path,
// 		query:  query.Encode(),
// 	})
// 	if m.IsErr() {
// 		return nil, m
// 	}

// 	records := make([]*zrule.HistoricalRecord, 0)

// 	err := json.Unmarshal(response, &records)
// 	if err != nil {
// 		m.Msg = errors.Wrapf(err, "unable to unmarshal response body on request %s", path)
// 		return nil, m
// 	}

// 	return records, m
// }

// func (s *service) HeadMarketsRegionIDTypes(ctx context.Context, regionID uint) Meta {

// 	_, m := s.request(ctx, request{
// 		method: http.MethodHead,
// 		path:   fmt.Sprintf("/v1/markets/%d/types/", regionID),
// 	})
// 	return m

// }

// func (s *service) GetMarketsRegionIDTypes(ctx context.Context, regionID uint, page null.String) ([]int, Meta) {

// 	path := fmt.Sprintf("/v1/markets/%d/types/", regionID)

// 	query := url.Values{}
// 	if page.Valid {
// 		query.Set("page", page.String)
// 	}

// 	response, m := s.request(ctx, request{
// 		method: http.MethodGet,
// 		path:   path,
// 		query:  query.Encode(),
// 	})
// 	if m.IsErr() {
// 		return nil, m
// 	}

// 	ids := make([]int, 0)

// 	err := json.Unmarshal(response, &ids)
// 	if err != nil {
// 		m.Msg = errors.Wrapf(err, "unable to unmarshal response body on request %s", path)
// 		return nil, m
// 	}

// 	return ids, m

// }

// func (s *service) GetMarketsPrices(ctx context.Context) ([]*zrule.MarketPrices, Meta) {

// 	path := "/v1/markets/prices/"

// 	response, m := s.request(ctx, request{
// 		method: http.MethodGet,
// 		path:   path,
// 	})
// 	if m.IsErr() {
// 		return nil, m
// 	}

// 	prices := make([]*zrule.MarketPrices, 0)
// 	err := json.Unmarshal(response, &prices)
// 	if err != nil {
// 		m.Msg = errors.Wrapf(err, "unable to unmarshal response body on request %s", path)
// 	}

// 	return prices, m

// }
