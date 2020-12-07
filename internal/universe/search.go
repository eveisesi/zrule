package universe

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/eveisesi/zrule"
	"github.com/newrelic/go-agent/v3/newrelic"
)

var ValidSearchCategories = map[string]string{
	"character":     "character",
	"corporation":   "corporation",
	"alliance":      "alliance",
	"region":        "region",
	"constellation": "constellation",
	"system":        "solar_system",
	"item":          "inventory_type",
}

func (s *service) SearchName(ctx context.Context, category, term string, strict bool) ([]*zrule.SearchResult, error) {

	ids, m := s.esi.GetSearch(ctx, category, term, strict)
	if m.IsErr() {
		return nil, m.Msg
	}

	wg := &sync.WaitGroup{}
	resultChan := make(chan *zrule.SearchResult, 1)
	results := make([]*zrule.SearchResult, 0)

	for _, v := range ids {
		wg.Add(1)
		go func(ctx context.Context, wg *sync.WaitGroup, id uint64, resultChan chan *zrule.SearchResult) {

			txn := newrelic.FromContext(ctx)
			txn = txn.NewGoroutine()

			ctx = newrelic.NewContext(ctx, txn)

			defer wg.Done()
			switch category {
			case "character":
				character, err := s.Character(ctx, id)
				if err != nil {
					txn.NoticeError(err)
					return
				}
				resultChan <- &zrule.SearchResult{
					ID:   character.ID,
					Name: character.Name,
				}
			case "corporation":
				corporation, err := s.Corporation(ctx, uint(id))
				if err != nil {
					txn.NoticeError(err)
					return
				}
				resultChan <- &zrule.SearchResult{
					ID:   uint64(corporation.ID),
					Name: corporation.Name,
				}
			case "alliance":
				alliance, err := s.Alliance(ctx, uint(id))
				if err != nil {
					txn.NoticeError(err)
					return
				}
				resultChan <- &zrule.SearchResult{
					ID:   uint64(alliance.ID),
					Name: alliance.Name,
				}
			case "region":
				region, err := s.Region(ctx, uint(id))
				if err != nil {
					txn.NoticeError(err)
					return
				}
				resultChan <- &zrule.SearchResult{
					ID:   uint64(region.ID),
					Name: region.Name,
				}
			case "constellation":
				constellation, err := s.Constellation(ctx, uint(id))
				if err != nil {
					txn.NoticeError(err)
					return
				}
				resultChan <- &zrule.SearchResult{
					ID:   uint64(constellation.ID),
					Name: constellation.Name,
				}
			case "solar_system":
				system, err := s.SolarSystem(ctx, uint(id))
				if err != nil {
					txn.NoticeError(err)
					return
				}
				resultChan <- &zrule.SearchResult{
					ID:   uint64(system.ID),
					Name: system.Name,
				}
			case "inventory_type":
				item, err := s.Item(ctx, uint(id))
				if err != nil {
					txn.NoticeError(err)
					return
				}
				resultChan <- &zrule.SearchResult{
					ID:   uint64(item.ID),
					Name: item.Name,
				}
			}

		}(ctx, wg, v, resultChan)
	}

	ticker := time.NewTicker(time.Millisecond * 100)
	ticks := 0
	for {
		if ticks > 4 {
			break
		}
		select {
		case <-ticker.C:
			fmt.Println("tick")
			ticks++
		case result := <-resultChan:
			ticks = 0
			results = append(results, result)
		}

	}

	return results, nil

}
