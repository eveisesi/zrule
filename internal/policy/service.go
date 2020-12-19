package policy

import (
	"context"
	"errors"

	"github.com/eveisesi/zrule"
	"github.com/eveisesi/zrule/internal/universe"
	"github.com/newrelic/go-agent/v3/newrelic"
	"go.mongodb.org/mongo-driver/mongo"
)

type Service interface {
	zrule.PolicyRepository
}

type service struct {
	universe universe.Service
	zrule.PolicyRepository
}

func NewService(universe universe.Service, policy zrule.PolicyRepository) Service {
	return &service{
		universe: universe,

		PolicyRepository: policy,
	}
}

func (s *service) Policies(ctx context.Context, operators ...*zrule.Operator) ([]*zrule.Policy, error) {

	policies, err := s.PolicyRepository.Policies(ctx, operators...)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return policies, err
	}

	for _, policy := range policies {
		for _, rule := range policy.Rules {
			for _, and := range rule {
				for _, pathObj := range zrule.AllPaths {
					if and.Path.String() == pathObj.Path.String() {
						if !pathObj.Searchable {
							break
						}
						and.Entities = make([]*zrule.SearchResult, len(and.Values))
						for i, v := range and.Values {
							switch t := v.(type) {
							case float64:
								switch pathObj.Category {
								case zrule.PathCategorySystems:
									system, err := s.universe.SolarSystem(ctx, uint(t))
									if err != nil {
										newrelic.FromContext(ctx).NoticeError(err)
										continue
									}

									and.Entities[i] = &zrule.SearchResult{
										ID:   uint64(system.ID),
										Name: system.Name,
									}
								case zrule.PathCategoryConstellations:
									constellation, err := s.universe.Constellation(ctx, uint(t))
									if err != nil {
										newrelic.FromContext(ctx).NoticeError(err)
										continue
									}

									and.Entities[i] = &zrule.SearchResult{
										ID:   uint64(constellation.ID),
										Name: constellation.Name,
									}
								case zrule.PathCategoryRegions:
									region, err := s.universe.Region(ctx, uint(t))
									if err != nil {
										newrelic.FromContext(ctx).NoticeError(err)
										continue
									}

									and.Entities[i] = &zrule.SearchResult{
										ID:   uint64(region.ID),
										Name: region.Name,
									}

								case zrule.PathCategoryCorporation:
									corporation, err := s.universe.Corporation(ctx, uint(t))
									if err != nil {
										newrelic.FromContext(ctx).NoticeError(err)
										continue
									}

									and.Entities[i] = &zrule.SearchResult{
										ID:   uint64(corporation.ID),
										Name: corporation.Name,
									}
								case zrule.PathCategoryAlliance:
									alliance, err := s.universe.Alliance(ctx, uint(t))
									if err != nil {
										newrelic.FromContext(ctx).NoticeError(err)
										continue
									}

									and.Entities[i] = &zrule.SearchResult{
										ID:   uint64(alliance.ID),
										Name: alliance.Name,
									}

								case zrule.PathCategoryCharacter:
									character, err := s.universe.Character(ctx, uint64(t))
									if err != nil {
										newrelic.FromContext(ctx).NoticeError(err)
										continue
									}

									and.Entities[i] = &zrule.SearchResult{
										ID:   uint64(character.ID),
										Name: character.Name,
									}
								case zrule.PathCategoryItems:
									item, err := s.universe.Item(ctx, uint(t))
									if err != nil {
										newrelic.FromContext(ctx).NoticeError(err)
										continue
									}

									and.Entities[i] = &zrule.SearchResult{
										ID:   uint64(item.ID),
										Name: item.Name,
									}
								}
							}
						}
					}
				}
			}
		}
	}

	return policies, nil

}
