package http

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/eveisesi/zrule/pkg/ruler"

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

func (s *server) handleGetPolicyByID(w http.ResponseWriter, r *http.Request) {

	var ctx = r.Context()

	policyID := chi.URLParam(r, "policyID")
	if policyID == "" {
		s.writeError(w, http.StatusBadRequest, fmt.Errorf("policyID is required"))
		return
	}

	objectID, err := primitive.ObjectIDFromHex(policyID)
	if err != nil {
		msg := "provided action id is invalid"
		s.logger.WithError(err).Error(msg)
		s.writeError(w, http.StatusBadRequest, fmt.Errorf(msg))

	}

	user := UserFromContext(ctx)
	if user == nil {
		err := fmt.Errorf("ctx does not contain a user")
		s.logger.WithError(err).Errorln()
		s.writeError(w, http.StatusInternalServerError, nil)
		return
	}

	policies, err := s.policy.Policies(ctx, zrule.NewEqualOperator("ownerID", user.ID), zrule.NewEqualOperator("_id", objectID))
	if err != nil {
		err = fmt.Errorf("failed to fetch actions by owner id")
		s.logger.WithError(err).Errorln()
		s.writeResponse(w, http.StatusInternalServerError, nil)
		return
	}

	if len(policies) == 0 {
		s.writeResponse(w, http.StatusNotFound, fmt.Errorf("failed to locate a policy with ID of %s", policyID))
		return
	}

	s.writeResponse(w, http.StatusOK, policies[0])
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

func (s *server) handleUpdatePolicy(w http.ResponseWriter, r *http.Request) {

	var ctx = r.Context()

	policyID := chi.URLParam(r, "policyID")
	if policyID == "" {
		s.writeError(w, http.StatusBadRequest, fmt.Errorf("policyID is required"))
		return
	}

	user := UserFromContext(ctx)
	if user == nil {
		err := fmt.Errorf("ctx does not contain a user")
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

	policies, err := s.policy.Policies(ctx, zrule.NewEqualOperator("ownerID", user.ID), zrule.NewEqualOperator("_id", objectID))
	if err != nil {
		err = fmt.Errorf("failed to fetch actions by owner id")
		s.logger.WithError(err).Errorln()
		s.writeResponse(w, http.StatusInternalServerError, nil)
		return
	}

	if len(policies) == 0 || len(policies) > 1 {
		s.writeResponse(w, http.StatusNotFound, fmt.Errorf("failed to locate a policy with ID of %s", policyID))
		return
	}

	policy := policies[0]

	createdAt := policy.CreatedAt

	defer r.Body.Close()
	err = json.NewDecoder(r.Body).Decode(&policy)
	if err != nil {
		s.writeError(w, http.StatusBadRequest, fmt.Errorf("failed to read body: %w", err))
		return
	}

	policy.ID = objectID
	policy.OwnerID = user.ID
	policy.CreatedAt = createdAt

	if len(policy.Actions) == 0 {
		msg := "Policies are required to have at least one action associated with them when they are created"
		s.logger.Error(msg)
		s.writeError(w, http.StatusBadRequest, fmt.Errorf(msg))
		return
	}

	if len(policy.Rules) == 0 {
		msg := "Policies are required to have at least one rule associated with them when they are created"
		s.logger.Error(msg)
		s.writeError(w, http.StatusBadRequest, fmt.Errorf(msg))
		return
	}

	err = ruler.Validate(policy.Rules)
	if err != nil {
		s.writeError(w, http.StatusBadRequest, err)
		return
	}

	policy, err = s.policy.UpdatePolicy(ctx, objectID, policy)
	if err != nil {
		s.writeError(w, http.StatusBadRequest, err)
		return
	}

	s.writeResponse(w, http.StatusOK, policy)

}

func (s *server) handleGetPolicyActions(w http.ResponseWriter, r *http.Request) {

	var ctx = r.Context()

	policyID := chi.URLParam(r, "policyID")
	if policyID == "" {
		s.writeError(w, http.StatusBadRequest, fmt.Errorf("policyID is required"))
		return
	}

	user := UserFromContext(ctx)
	if user == nil {
		err := fmt.Errorf("ctx does not contain a user")
		s.logger.WithError(err).Errorln()
		s.writeError(w, http.StatusInternalServerError, nil)
		return
	}

	objectID, err := primitive.ObjectIDFromHex(policyID)
	if err != nil {
		msg := "provided action id is invalid"
		s.logger.WithError(err).Error(msg)
		s.writeError(w, http.StatusBadRequest, fmt.Errorf(msg))
		return
	}

	policies, err := s.policy.Policies(ctx, zrule.NewEqualOperator("ownerID", user.ID), zrule.NewEqualOperator("_id", objectID))
	if err != nil {
		err = fmt.Errorf("failed to fetch policies by owner id")
		s.logger.WithError(err).Errorln()
		s.writeResponse(w, http.StatusInternalServerError, nil)
		return
	}

	policy := policies[0]

	var actionIDs = make([]zrule.OpValue, len(policy.Actions))
	for i, actionID := range policy.Actions {
		actionIDs[i] = actionID
	}
	actions, err := s.action.Actions(ctx, zrule.NewEqualOperator("ownerID", user.ID), zrule.NewInOperator("_id", actionIDs))

	if err != nil {
		err = fmt.Errorf("failed to fetch actions by ownerID and actionIDs")
		s.logger.WithError(err).Errorln()
		s.writeResponse(w, http.StatusInternalServerError, nil)
		return
	}

	s.writeResponse(w, http.StatusOK, actions)

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
