package zrule

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ActionRepository interface {
	Action(ctx context.Context, id primitive.ObjectID) (*Action, error)
	Actions(ctx context.Context, ops ...*Operator) ([]*Action, error)
	CreateAction(ctx context.Context, action *Action) (*Action, error)
	UpdateAction(ctx context.Context, id primitive.ObjectID, action *Action) (*Action, error)
	DeleteAction(ctx context.Context, id primitive.ObjectID) error
}

type Action struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	OwnerID        primitive.ObjectID `bson:"owner_id" json:"owner_id"`
	Label          string             `bson:"label" json:"label"`
	Platform       Platform           `bson:"platform" json:"platform"`
	Endpoint       string             `bson:"endpoint" json:"endpoint"`
	Tested         bool               `bson:"tested" json:"tested"`
	IsDisabled     bool               `bson:"is_disabled" json:"is_disabled"`
	DisabledReason *string            `bson:"disabled_reason" json:"disabled_reason"`
	Infractions    []*Infraction      `bson:"infractions" json:"infractions"`
	CreatedAt      time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt      time.Time          `bson:"updated_at" json:"updated_at"`
}

func (a *Action) IsValid() error {
	uri, err := url.Parse(a.Endpoint)
	if err != nil {
		return fmt.Errorf("failed to validate structure of endpoint. ")
	}

	if uri.Scheme != "http" && uri.Scheme != "https" {
		return fmt.Errorf("invalid url scheme detected. Please use http or https")
	}

	switch uri.Host {
	case HostSlack.String():
		a.Platform = PlatformSlack
	case HostDiscord.String(), HostDiscordApp.String():
		a.Platform = PlatformDiscord
	default:
		a.Platform = PlatformRest
	}

	if a.Platform == PlatformDiscord && !strings.HasPrefix(uri.Path, "/api/webhooks") {
		return fmt.Errorf("malformed discord webhook")
	} else if a.Platform == PlatformSlack && !strings.HasPrefix(uri.Path, "/services") {
		return fmt.Errorf("malformed slack webhook")
	}

	return nil
}

type Host string

const HostSlack Host = "hooks.slack.com"
const HostDiscordApp Host = "discordapp.com"
const HostDiscord Host = "discord.com"

func (h Host) String() string {
	return string(h)
}

type Platform string

const (
	PlatformSlack   Platform = "slack"
	PlatformDiscord Platform = "discord"
	PlatformRest    Platform = "rest"
)

var AllPlatforms = []Platform{PlatformSlack, PlatformDiscord, PlatformRest}

func (p Platform) IsValid() bool {
	for _, v := range AllPlatforms {
		if v == p {
			return true
		}
	}

	return false
}

func (p Platform) String() string { return string(p) }
