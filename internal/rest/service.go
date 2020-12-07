package rest

import (
	"bytes"
	"context"
	"fmt"
	"net/http"

	"github.com/eveisesi/zrule"
)

type service struct {
	client *http.Client
	action *zrule.Action
}

func NewService(action *zrule.Action, client *http.Client) (zrule.Dispatcher, error) {
	return &service{
		client: client,
		action: action,
	}, nil
}

func (s *service) Send(ctx context.Context, policy *zrule.Policy, id uint, hash string) error {

	fmt.Println("Rest", id)

	return nil
}

func (s *service) SendTest(ctx context.Context, message string) error {

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.action.Endpoint, bytes.NewBuffer([]byte(fmt.Sprintf(`{"message": %s}`, message))))
	if err != nil {
		return fmt.Errorf("Failed to build request for test message: %w", err)
	}

	res, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request for test message: %w", err)
	}

	if res.StatusCode > http.StatusBadRequest {
		return fmt.Errorf("Endpoint returned an invalid status code %d", res.StatusCode)
	}

	if res.StatusCode != http.StatusNoContent {
		return fmt.Errorf("Endpoint did not return the required status code of 204 No Content")
	}

	return nil
}
