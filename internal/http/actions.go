package http

import (
	"encoding/json"
	"fmt"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/eveisesi/zrule"
	"github.com/go-chi/chi"
)

func (s *server) handleGetActions(w http.ResponseWriter, r *http.Request) {

	var ctx = r.Context()

	user := UserFromContext(ctx)
	if user == nil {
		err := fmt.Errorf("ctx does not contain a user")
		s.logger.WithError(err).Errorln()
		s.writeResponse(w, http.StatusInternalServerError, nil)
		return
	}

	actions, err := s.action.Actions(ctx, zrule.NewEqualOperator("ownerID", user.ID))
	if err != nil {
		err = fmt.Errorf("failed to fetch actions by owner id")
		s.logger.WithError(err).Errorln()
		s.writeResponse(w, http.StatusInternalServerError, nil)
		return
	}

	s.writeResponse(w, http.StatusOK, actions)

}

func (s *server) handlePostActionTest(w http.ResponseWriter, r *http.Request) {

	var ctx = r.Context()

	actionID := chi.URLParam(r, "actionID")
	if actionID == "" {
		s.writeError(w, http.StatusBadRequest, fmt.Errorf("actionID is required to perform test"))
		return
	}

	user := UserFromContext(ctx)
	if user == nil {
		err := fmt.Errorf("ctx does not contain a user")
		s.logger.WithError(err).Errorln()
		s.writeError(w, http.StatusInternalServerError, nil)
		return
	}

	objectID, err := primitive.ObjectIDFromHex(actionID)
	if err != nil {
		msg := "provided action id is invalid"
		s.logger.WithError(err).Error(msg)
		s.writeError(w, http.StatusBadRequest, fmt.Errorf(msg))

	}
	actions, err := s.action.Actions(ctx, zrule.NewEqualOperator("ownerID", user.CharacterID), zrule.NewEqualOperator("_id", objectID))
	if err != nil {
		s.logger.WithError(err).WithField("actionID", actionID).Error("failed to find action for provided actionID")
		s.writeError(w, http.StatusBadRequest, fmt.Errorf("failed to find action for provided actionID"))
		return
	}

	if len(actions) == 0 || len(actions) > 1 {
		s.logger.WithError(err).WithField("actionID", actionID).Error("matched multiple actions for query")
		s.writeError(w, http.StatusInternalServerError, fmt.Errorf("matched multiple actions for query"))
		return
	}

	// action := actions[0]

	// messenger := message.NewService(action)
	// if err != nil {
	// 	s.logger.WithError(err).WithField("actionID", actionID).Error("failed to initialize messenger service")
	// 	s.writeError(w, http.StatusInternalServerError, fmt.Errorf("failed to initialize messenger service"))
	// 	return
	// }

	// messenger.SendTest(ctx, "Hello From the Test endpoint of zrule")

}

func (s *server) handleCreateAction(w http.ResponseWriter, r *http.Request) {

	var ctx = r.Context()

	var action = &zrule.Action{}
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(action)
	if err != nil {
		s.writeError(w, http.StatusBadRequest, fmt.Errorf("failed to read body: %w", err))
		return
	}

	if action.Tested {
		action.Tested = false
	}
	if action.IsDisabled {
		action.IsDisabled = false
	}

	if err := action.IsValid(); err != nil {
		s.writeError(w, http.StatusBadRequest, err)
		return
	}

	user := UserFromContext(ctx)
	if user == nil {
		err = fmt.Errorf("ctx does not contain a user")
		s.logger.WithError(err).Errorln()
		s.writeResponse(w, http.StatusInternalServerError, nil)
		return
	}

	action.OwnerID = user.ID

	_, err = s.action.CreateAction(ctx, action)
	if err != nil {
		s.logger.WithError(err).Error("failed to save action")
		s.writeError(w, http.StatusBadRequest, fmt.Errorf("failed to save action"))
		return
	}

	s.writeResponse(w, http.StatusCreated, nil)

}
