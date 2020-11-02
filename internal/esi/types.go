package esi

// // Type is an object representing the database table.
// type Type struct {
// 	ID            uint             `json:"type_id"`
// 	GroupID       uint             `json:"group_id"`
// 	Name          string           `json:"name"`
// 	Description   string           `json:"description"`
// 	Published     bool             `json:"published"`
// 	MarketGroupID null.Uint        `json:"marketGroupID"`
// 	Attributes    []*TypeAttribute `json:"dogma_attributes"`
// }

// type TypeAttribute struct {
// 	AttributeID uint    `json:"attribute_id"`
// 	Value       float64 `json:"value"`
// }

// func (s *service) GetUniverseTypesTypeID(ctx context.Context, id uint) (*zrule.Type, []*zrule.TypeAttribute, Meta) {

// 	var esitype = new(Type)

// 	path := fmt.Sprintf("/v3/universe/types/%d/", id)

// 	request := request{
// 		method: http.MethodGet,
// 		path:   path,
// 	}

// 	response, m := s.request(ctx, request)
// 	if m.IsErr() {
// 		return nil, nil, m
// 	}

// 	err = json.Unmarshal(response, esitype)
// 	if err != nil {
// 		m.Msg = errors.Wrapf(err, "unable to unmarshal response body on request %s", path)
// 		return nil, nil, m
// 	}

// 	var attributes = make([]*zrule.TypeAttribute, 0)
// 	for _, v := range esitype.Attributes {
// 		attributes = append(attributes, &zrule.TypeAttribute{
// 			TypeID:      id,
// 			AttributeID: v.AttributeID,
// 			Value:       int64(v.Value),
// 		})
// 	}

// 	return &zrule.Type{
// 		ID:          esitype.ID,
// 		GroupID:     esitype.GroupID,
// 		Name:        esitype.Name,
// 		Description: esitype.Description,
// 		Published:   esitype.Published,
// 		// MarketGroupID: esitype.MarketGroupID,
// 	}, attributes, m

// }
