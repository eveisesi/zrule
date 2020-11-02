package zrule

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PolicyRepository interface {
	Policy(ctx context.Context, id primitive.ObjectID) (*Policy, error)
	Policies(ctx context.Context, operators ...*Operator) ([]*Policy, error)
	CreatePolicy(ctx context.Context, policy *Policy) (*Policy, error)
	UpdatePolicy(ctx context.Context, id primitive.ObjectID, policy *Policy) (*Policy, error)
	DeletePolicy(ctx context.Context, id primitive.ObjectID) error
}

type Dispatchable struct {
	PolicyID primitive.ObjectID `json:"policyID"`
	ID       uint               `json:"id"`
	Hash     string             `json:"hash"`
}

type Policy struct {
	ID        primitive.ObjectID   `bson:"_id,omitempty" json:"_id,omitempty"`
	Name      string               `bson:"name" json:"name"`
	OwnerID   primitive.ObjectID   `bson:"ownerID" json:"ownerID"`
	Rules     []*Rule              `bson:"rules" json:"rules"`
	Actions   []primitive.ObjectID `bson:"actions" json:"actions"`
	CreatedAt time.Time            `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time            `bson:"updatedAt" json:"updatedAt"`
}

// Heavily mirror's ruler.Rule
type Rule struct {
	Comparator string      `json:"comparator"`
	Path       string      `json:"path"`
	Value      interface{} `json:"value"`
}

type Property string

const (
	PropertyMoonID                Property = "moon_id"
	PropertySolarSystemID         Property = "solar_system_od"
	PropertyWarID                 Property = "war_id"
	PropertyZKBNPC                Property = "zkb.npc"
	PropertyZKBAWOX               Property = "zkb.awox"
	PropertyZKBSolo               Property = "zkb.solo"
	PropertyZKBFittedValue        Property = "zkb.fittedValue"
	PropertyZKBTotalValue         Property = "zkb.totalValue"
	PropertyVictimAllianceID      Property = "victim.alliance_id"
	PropertyVictimCorporationID   Property = "victim.corporation_id"
	PropertyVictimFractionID      Property = "victim.fraction_id"
	PropertyVictimCharacterID     Property = "victim.character_id"
	PropertyVictimShipTypeID      Property = "victim.ship_type_id"
	PropertyVictimItemsTypeID     Property = "victim.items.item_type_id"
	PropertyAttackerAllianceID    Property = "attackers.alliance_id"
	PropertyAttackerCorporationID Property = "attackers.corporation_id"
	PropertyAttackerFractionID    Property = "attackers.fraction_id"
	PropertyAttackerCharacterID   Property = "attackers.character_id"
	PropertyAttackerShipTypeID    Property = "attackers.ship_type_id"
	PropertyAttackerWeaponTypeID  Property = "attackers.weapon_type_id"
)

var AllProperties = []Property{
	PropertyMoonID, PropertySolarSystemID, PropertyWarID,
	PropertyZKBNPC, PropertyZKBAWOX, PropertyZKBSolo, PropertyZKBFittedValue, PropertyZKBTotalValue,
	PropertyVictimAllianceID, PropertyVictimCorporationID, PropertyVictimFractionID, PropertyVictimCharacterID,
	PropertyVictimShipTypeID, PropertyVictimItemsTypeID, PropertyAttackerAllianceID, PropertyAttackerCorporationID,
	PropertyAttackerFractionID, PropertyAttackerCharacterID, PropertyAttackerShipTypeID, PropertyAttackerWeaponTypeID,
}

func (p Property) IsValid() bool {
	for _, v := range AllProperties {
		if v == p {
			return true
		}
	}
	return false
}

func (p Property) String() string {
	return string(p)
}

type Infraction struct {
	InfrationID primitive.ObjectID `bson:"_id" json:"_id"`
	Message     string
}
