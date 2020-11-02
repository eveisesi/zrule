package slack

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/davecgh/go-spew/spew"
	"github.com/eveisesi/zrule"
)

type service struct {
	client  *http.Client
	webhook string
}

func NewService(action *zrule.Action, client *http.Client) (zrule.Dispatcher, error) {
	return &service{
		client:  client,
		webhook: action.Endpoint,
	}, nil
}

func (s *service) Send(ctx context.Context, policy *zrule.Policy, id uint, hash string) error {

	uri := url.URL{
		Scheme: "https",
		Host:   "zkillboard.com",
		Path:   fmt.Sprintf("/kill/%d", id),
	}

	content := fmt.Sprintf("Match Found with Policy %s (%s)\n%s", policy.Name, policy.ID.Hex(), uri.String())

	data, err := json.Marshal(map[string]interface{}{
		"text":         content,
		"unfurl_links": true,
	})
	if err != nil {
		return fmt.Errorf("failed to prepare request body to post to slack: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.webhook, bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("failed to prepare request to slack: %w", err)
	}

	res, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request to slack: %w", err)
	}

	if res.StatusCode != http.StatusOK {
		defer res.Body.Close()
		data, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("failed to decode error response: %w", err)
		}

		spew.Dump(string(data))
	}

	return nil
}

func (s *service) SendTest(ctx context.Context, message string) error {
	fmt.Println("Rest Test", message)
	return nil
}