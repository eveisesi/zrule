package http

import (
	"fmt"
	"net/http"

	"github.com/eveisesi/zrule/internal/universe"
)

func (s *server) handleGetSearchCategories(w http.ResponseWriter, r *http.Request) {

	var list = make([]string, len(universe.ValidSearchCategories))
	i := 0
	for v := range universe.ValidSearchCategories {
		list[i] = v
		i++
	}

	s.writeResponse(w, http.StatusOK, list)

}

func (s *server) handleGetSearchName(w http.ResponseWriter, r *http.Request) {

	var ctx = r.Context()

	user := UserFromContext(ctx)
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

	category := r.URL.Query().Get("category")
	if category == "" {
		s.writeError(w, http.StatusBadRequest, fmt.Errorf("category query paramater is required to execute a search"))
		return
	}

	var valid bool
	for i, v := range universe.ValidSearchCategories {
		if category == i {
			valid = true
			category = v
		}
	}

	if !valid {
		s.writeError(w, http.StatusBadRequest, fmt.Errorf("invalid category %s submitted. Please use a valid category [/search/categories]", category))
		return
	}

	results, err := s.universe.SearchName(ctx, category, term)
	if err != nil {
		s.writeError(w, http.StatusBadRequest, err)
		return
	}

	s.writeResponse(w, http.StatusOK, results)

}
