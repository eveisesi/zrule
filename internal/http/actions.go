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

	actions, err := s.action.Actions(ctx, zrule.NewEqualOperator("owner_id", user.ID))
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

	entry := s.logger.WithField("actionID", actionID)

	user := UserFromContext(ctx)
	if user == nil {
		err := fmt.Errorf("ctx does not contain a user")
		entry.WithError(err).Errorln()
		s.writeError(w, http.StatusInternalServerError, nil)
		return
	}

	objectID, err := primitive.ObjectIDFromHex(actionID)
	if err != nil {
		msg := "provided action id is invalid"
		entry.WithError(err).Error(msg)
		s.writeError(w, http.StatusBadRequest, fmt.Errorf(msg))
		return
	}

	entry = entry.WithField("ownerID", user.ID)

	actions, err := s.action.Actions(ctx, zrule.NewEqualOperator("owner_id", user.ID), zrule.NewEqualOperator("_id", objectID))
	if err != nil {
		entry.WithError(err).Error("failed to find action for provided actionID")
		s.writeError(w, http.StatusBadRequest, fmt.Errorf("failed to find action for provided actionID"))
		return
	}

	if len(actions) == 0 || len(actions) > 1 {
		s.logger.WithError(err).WithField("actionID", actionID).Error("matched multiple actions for query")
		s.writeError(w, http.StatusInternalServerError, fmt.Errorf("matched multiple actions for query"))
		return
	}

	action := actions[0]

	err = s.dispatcher.SendTestMessage(ctx, action, "Is this thing on?!?!")
	if err != nil {
		entry.WithError(err).Error("failed to send test message")
		s.writeError(w, http.StatusBadRequest, err)
		return
	}

	if !action.Tested {
		action.Tested = true
		_, err = s.action.UpdateAction(ctx, action.ID, action)
		if err != nil {
			entry.WithError(err).Error("failed to update action")
			s.writeError(w, http.StatusBadRequest, err)
			return
		}
	}

	s.writeResponse(w, http.StatusAccepted, nil)

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

func (s *server) handleUpdateAction(w http.ResponseWriter, r *http.Request) {

	var ctx = r.Context()

	actionID := chi.URLParam(r, "actionID")
	if actionID == "" {
		s.writeError(w, http.StatusBadRequest, fmt.Errorf("actionID is required to perform an update"))
		return
	}

	entry := s.logger.WithField("actionID", actionID)

	user := UserFromContext(ctx)
	if user == nil {
		err := fmt.Errorf("ctx does not contain a user")
		entry.WithError(err).Errorln()
		s.writeError(w, http.StatusInternalServerError, nil)
		return
	}

	objectID, err := primitive.ObjectIDFromHex(actionID)
	if err != nil {
		msg := "provided action id is invalid"
		entry.WithError(err).Error(msg)
		s.writeError(w, http.StatusBadRequest, fmt.Errorf(msg))
		return
	}

	entry = entry.WithField("ownerID", user.ID)

	actions, err := s.action.Actions(ctx, zrule.NewEqualOperator("owner_id", user.ID), zrule.NewEqualOperator("_id", objectID))
	if err != nil {
		entry.WithError(err).Error("failed to find action for provided actionID")
		s.writeError(w, http.StatusBadRequest, fmt.Errorf("failed to find action for provided actionID"))
		return
	}

	if len(actions) == 0 || len(actions) > 1 {
		entry.WithError(err).Error("matched multiple actions for query")
		s.writeError(w, http.StatusInternalServerError, fmt.Errorf("matched multiple actions for query"))
		return
	}

	action := actions[0]

	err = json.NewDecoder(r.Body).Decode(&action)
	if err != nil {
		entry.WithError(err).Error("failed to decode response")
		s.writeError(w, http.StatusBadRequest, fmt.Errorf("failed to decode response: %w", err))
		return
	}

	action, err = s.action.UpdateAction(ctx, objectID, action)
	if err != nil {
		entry.WithError(err).Error("failed to update action")
		s.writeError(w, http.StatusBadRequest, fmt.Errorf("failed to update action: %w", err))
		return
	}

	s.writeResponse(w, http.StatusOK, action)

}

func (s *server) handleDeleteAction(w http.ResponseWriter, r *http.Request) {

	var ctx = r.Context()

	actionID := chi.URLParam(r, "actionID")
	if actionID == "" {
		s.writeError(w, http.StatusBadRequest, fmt.Errorf("actionID is required to perform test"))
		return
	}

	entry := s.logger.WithField("actionID", actionID)

	user := UserFromContext(ctx)
	if user == nil {
		err := fmt.Errorf("ctx does not contain a user")
		entry.WithError(err).Errorln()
		s.writeError(w, http.StatusInternalServerError, nil)
		return
	}

	objectID, err := primitive.ObjectIDFromHex(actionID)
	if err != nil {
		msg := "provided action id is invalid"
		entry.WithError(err).Error(msg)
		s.writeError(w, http.StatusBadRequest, fmt.Errorf(msg))
		return
	}

	entry = entry.WithField("ownerID", user.ID)

	err = s.action.DeleteAction(ctx, objectID)
	if err != nil {
		entry.WithError(err).Error("failed to delete action")
		s.writeError(w, http.StatusBadRequest, fmt.Errorf("failed to delete action: %w", err))
		return
	}

}
