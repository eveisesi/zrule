package http

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/eveisesi/zrule/internal/search"
	"github.com/go-chi/chi"
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
	// Regions, Constellations, Systems, Items, ItemGroups, and Factions.
	// For Alliaces, Corporations, or Characters, we need to tell the consumer to use
	// the actual ESI API

	// There are all of the possible search categories per say
	allCategories := []string{"alliances", "corporations", "characters", "regions", "constellations", "systems", "items", "itemGroups", "faction"}

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
		s.writeError(w, http.StatusBadRequest, fmt.Errorf("key query paramater is required to execute a search"))
		return
	}

	if !search.Key(key).Valid() {
		s.writeError(w, http.StatusBadRequest, fmt.Errorf("key is invalid"))
		return
	}

	results, err := s.search.Search(search.Key(key), term)
	if err != nil {
		s.writeError(w, http.StatusBadRequest, err)
		return
	}

	s.writeResponse(w, http.StatusOK, results)

}

func (s *server) handleNewCategoryEntity(w http.ResponseWriter, r *http.Request) {

	var ctx = r.Context()

	category := chi.URLParam(r, "category")
	if category == "" {
		s.writeError(w, http.StatusBadRequest, fmt.Errorf("invalid or empty category received"))
		return
	}

	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		s.writeError(w, http.StatusBadRequest, fmt.Errorf("failed to coerce %v to integer: %w", chi.URLParam(r, "id"), err))
		return
	}
	if id == 0 {
		s.writeError(w, http.StatusBadRequest, fmt.Errorf("id must be greater than 0"))
		return
	}

	switch category {
	case "alliance":
		_, err := s.universe.Alliance(ctx, uint(id))
		if err != nil {
			s.writeError(w, http.StatusBadRequest, fmt.Errorf("unable to validate alliance id: %w", err))
			return
		}
	case "corporation":
		_, err := s.universe.Corporation(ctx, uint(id))
		if err != nil {
			s.writeError(w, http.StatusBadRequest, fmt.Errorf("unable to validate corporation id: %w", err))
			return
		}
	case "character":
		_, err := s.universe.Character(ctx, id)
		if err != nil {
			s.writeError(w, http.StatusBadRequest, fmt.Errorf("unable to validate character id: %w", err))
			return
		}
	default:
		s.writeError(w, http.StatusBadRequest, fmt.Errorf("%v is not a supported category", category))
		return
	}

	s.writeResponse(w, http.StatusNoContent, nil)
}
