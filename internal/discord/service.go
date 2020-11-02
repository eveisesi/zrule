package discord

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/eveisesi/zrule"
	newrelic "github.com/newrelic/go-agent"
)

type service struct {
	dgo       *discordgo.Session
	id, token string
}

func NewService(action *zrule.Action, client *http.Client) (zrule.Dispatcher, error) {

	if action.Platform != zrule.PlatformDiscord {
		return nil, fmt.Errorf("invalid platform for discord service constructor")
	}

	session, err := discordgo.New()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize discord session: %w", err)
	}

	session.Client = client

	uri, err := url.Parse(action.Endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to parse endpount")
	}

	path := uri.Path

	parts := strings.Split(path, "/")
	if len(parts) != 5 {
		return nil, fmt.Errorf("expect path split by / to equal 4, got %d", len(parts))
	}

	var id, token string
	id = parts[3]
	token = parts[4]

	return &service{
		dgo:   session,
		id:    id,
		token: token,
	}, nil
}

func (s *service) Send(ctx context.Context, policy *zrule.Policy, id uint, hash string) error {

	seg := newrelic.StartSegment(newrelic.FromContext(ctx), "send discord message")
	defer seg.End()

	uri := url.URL{
		Scheme: "https",
		Host:   "zkillboard.com",
		Path:   fmt.Sprintf("/kill/%d", id),
	}

	content := fmt.Sprintf("Match Found with Policy %s (%s)\n%s", policy.Name, policy.ID.Hex(), uri.String())

	_, err := s.dgo.WebhookExecute(s.id, s.token, true, &discordgo.WebhookParams{
		Content: content,
	})

	return err
}

func (s *service) SendTest(ctx context.Context, message string) error {

	seg := newrelic.StartSegment(newrelic.FromContext(ctx), "send discord test message")
	defer seg.End()

	_, err := s.dgo.WebhookExecute(s.id, s.token, true, &discordgo.WebhookParams{
		Username: "zrule Test",
		Content:  message,
	})

	return err
}
