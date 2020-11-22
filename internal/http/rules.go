package http

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/eveisesi/zrule/pkg/ruler"
)

func (s *server) handlePostValidateRules(w http.ResponseWriter, r *http.Request) {

	var rules ruler.Rules
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&rules)
	if err != nil {
		s.writeError(w, http.StatusBadRequest, fmt.Errorf("failed to decode request body"))
		return
	}

	err = ruler.Validate(rules)
	if err != nil {
		s.writeError(w, http.StatusBadRequest, fmt.Errorf("unable to validate rule: %w", err))
		return
	}

	s.writeResponse(w, http.StatusOK, nil)

}
