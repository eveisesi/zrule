package http

import (
	"fmt"
	"net/http"

	"github.com/eveisesi/zrule"
)

func (s *server) handleGetPaths(w http.ResponseWriter, r *http.Request) {

	var ctx = r.Context()
	user := UserFromContext(ctx)
	if user == nil {
		err := fmt.Errorf("ctx does not contain a user")
		s.logger.WithError(err).Errorln()
		s.writeError(w, http.StatusInternalServerError, nil)
		return
	}

	s.writeResponse(w, http.StatusOK, zrule.AllPaths)

}
