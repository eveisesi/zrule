package http

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/eveisesi/zrule"
	"github.com/go-chi/chi"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (s *server) handleGetPolicies(w http.ResponseWriter, r *http.Request) {

	var ctx = r.Context()

	user := UserFromContext(ctx)
	if user == nil {
		err := fmt.Errorf("ctx does not contain a user")
		s.logger.WithError(err).Errorln()
		s.writeError(w, http.StatusInternalServerError, nil)
		return
	}

	policies, err := s.policy.Policies(ctx, zrule.NewEqualOperator("ownerID", user.ID))
	if err != nil {
		err = fmt.Errorf("failed to fetch actions by owner id")
		s.logger.WithError(err).Errorln()
		s.writeResponse(w, http.StatusInternalServerError, nil)
		return
	}

	s.writeResponse(w, http.StatusOK, policies)
}

func (s *server) handleCreatePolicy(w http.ResponseWriter, r *http.Request) {

	var ctx = r.Context()

	var policy = new(zrule.Policy)
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(policy)
	if err != nil {
		s.writeError(w, http.StatusBadRequest, fmt.Errorf("failed to read body: %w", err))
		return
	}

	user := UserFromContext(ctx)
	if user == nil {
		err := fmt.Errorf("ctx does not contain a user")
		s.logger.WithError(err).Errorln()
		s.writeError(w, http.StatusInternalServerError, nil)
		return
	}

	policy.OwnerID = user.ID

	if len(policy.Actions) == 0 {
		msg := "Policies are required to have at least one action associated with them when they are created"
		s.logger.Error(msg)
		s.writeError(w, http.StatusBadRequest, fmt.Errorf(msg))
		return
	}

	policy, err = s.policy.CreatePolicy(ctx, policy)
	if err != nil {
		msg := "failed to insert policy document into datastore"
		s.logger.WithError(err).Error(msg)
		s.writeError(w, http.StatusBadRequest, fmt.Errorf(msg))
		return
	}

	s.writeResponse(w, http.StatusOK, policy)

}

func (s *server) handleDeletePolicy(w http.ResponseWriter, r *http.Request) {

	var ctx = r.Context()

	user := UserFromContext(ctx)
	if user == nil {
		err := fmt.Errorf("ctx does not contain a user")
		s.logger.WithError(err).Errorln()
		s.writeError(w, http.StatusInternalServerError, nil)
		return
	}

	policyID := chi.URLParam(r, "policyID")
	if policyID == "" {
		err := fmt.Errorf("policyID is a required property")
		s.logger.WithError(err).Errorln()
		s.writeError(w, http.StatusInternalServerError, nil)
		return
	}

	objectID, err := primitive.ObjectIDFromHex(policyID)
	if err != nil {
		msg := "provided action id is invalid"
		s.logger.WithError(err).Error(msg)
		s.writeError(w, http.StatusBadRequest, fmt.Errorf(msg))

	}

	err = s.policy.DeletePolicy(ctx, objectID)
	if err != nil {
		msg := "failed to delete policy"
		s.logger.WithError(err).Error(msg)
		s.writeError(w, http.StatusBadRequest, fmt.Errorf(msg))
		return
	}

	s.writeResponse(w, http.StatusNoContent, nil)
}
