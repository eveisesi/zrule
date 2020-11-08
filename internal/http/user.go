package http

import "net/http"

func (s *server) handleGetUser(w http.ResponseWriter, r *http.Request) {

	var ctx = r.Context()

	user := UserFromContext(ctx)
	if user == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	character, err := s.character.Character(ctx, user.CharacterID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	s.writeResponse(w, http.StatusOK, character)

}
