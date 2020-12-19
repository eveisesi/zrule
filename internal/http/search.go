package http

import (
	"fmt"
	"net/http"

	"github.com/eveisesi/zrule/internal/search"
)

type searchable struct {
	Category string `json:"category"`
	UseESI   bool   `json:"useESI"`
	UseAPI   bool   `json:"useAPI"`
}

func (s *server) handleGetSearchCategories(w http.ResponseWriter, r *http.Request) {

	// Search is provided for a subset of searchable entities within Eve Online
	// We cannot facilitate search for Alliaces, Corporations, or Characters
	// There are just to many in Eve Online, that is why the search service only provides
	// Regions, Constellations, Systems, Items, and ItemGroups.
	// For Alliaces, Corporations, or Characters, we need to tell the consumer to use
	// the actual ESI API

	// There are all of the possible search categories per say
	allCategories := []string{"alliances", "corporations", "characters", "regions", "constellations", "systems", "items", "itemGroups"}

	// Get the ones that we provide
	allKeys := search.AllKeys

	results := make([]*searchable, len(allCategories))
	for i, category := range allCategories {
		result := new(searchable)
		result.Category = category
		for _, key := range allKeys {
			if key.String() == category {
				result.UseAPI = true
			}
		}
		if !result.UseAPI {
			result.UseESI = true
		}
		results[i] = result
	}

	s.writeResponse(w, http.StatusOK, results)
}

func (s *server) handleGetSearchName(w http.ResponseWriter, r *http.Request) {

	user := UserFromContext(r.Context())
	if user == nil {
		err := fmt.Errorf("ctx does not contain a user")
		s.logger.WithError(err).Errorln()
		s.writeError(w, http.StatusInternalServerError, nil)
		return
	}

	term := r.URL.Query().Get("term")
	if term == "" {
		s.writeError(w, http.StatusBadRequest, fmt.Errorf("term query paramater is required to execute a search"))
		return
	}

	key := r.URL.Query().Get("key")
	if key == "" {
		s.writeError(w, http.StatusBadRequest, fmt.Errorf("category query paramater is required to execute a search"))
		return
	}

	results, err := s.search.Search(search.Key(key), term)
	if err != nil {
		s.writeError(w, http.StatusBadRequest, err)
		return
	}

	s.writeResponse(w, http.StatusOK, results)

}
