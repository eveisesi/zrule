package zrule

import (
	"context"
	"time"

	"github.com/eveisesi/zrule/pkg/ruler"
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
	OwnerID   primitive.ObjectID   `bson:"owner_id" json:"owner_id"`
	Rules     [][]*ruler.Rule      `bson:"rules" json:"rules"`
	Actions   []primitive.ObjectID `bson:"actions" json:"actions"`
	CreatedAt time.Time            `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time            `bson:"updated_at" json:"updated_at"`
}

type PathObj struct {
	Display     string `json:"display"`
	Description string `json:"description"`
	Category    string `json:"category,omitempty"`
	Searchable  bool   `json:"searchable"`
	Format      string `json:"format"`
	Path        Path   `json:"path"`
}

type Path string

var (
	PathSolarSystemID = PathObj{
		Display:     "Solar System",
		Description: "The Solar System that the Killmail occurred in",
		Category:    "solar_system",
		Searchable:  true,
		Format:      "string",
		Path:        Path("solar_system_id"),
	}
	PathZKBNPC = PathObj{
		Display:     "ZKillboard Is NPC",
		Description: "ZKillboard has labeled the killmail as an NPC Kill",
		Searchable:  false,
		Format:      "boolean",
		Path:        Path("zkb.npc"),
	}
	PathZKBAWOX = PathObj{
		Display:     "ZKillboard Is AWOX",
		Description: "Zkillboard has labeled the killmail as an AWOX Kill",
		Searchable:  false,
		Format:      "boolean",
		Path:        Path("zkb.awox"),
	}
	PathZKBSolo = PathObj{
		Display:     "ZKillboard Is Solo",
		Description: "ZKIllboard has labeled the killmail as a Solo Kill",
		Searchable:  false,
		Format:      "boolean",
		Path:        Path("zkb.solo"),
	}
	PathZKBFittedValue = PathObj{
		Display:     "Zkillboard Fitted Value",
		Description: "The ISK value of all modules and ammo fitted to the ship",
		Searchable:  false,
		Format:      "number",
		Path:        Path("zkb.fittedValue"),
	}
	PathZKBTotalValue = PathObj{
		Display:     "ZKillboard Total Value",
		Description: "The ISK value of the killmail",
		Searchable:  false,
		Format:      "number",
		Path:        Path("zkb.totalValue"),
	}
	PathVictimAllianceID = PathObj{
		Display:     "Victim Alliance",
		Description: "The alliance that the victim is apart of",
		Searchable:  true,
		Format:      "string",
		Category:    "alliance",
		Path:        Path("victim.alliance_id"),
	}
	PathVictimCorporationID = PathObj{
		Display:     "Victim Corporation",
		Description: "The Corporation that the victim is apart of",
		Searchable:  true,
		Format:      "string",
		Category:    "corporation",
		Path:        Path("victim.corporation_id"),
	}
	PathVictimCharacterID = PathObj{
		Display:     "Victim Character",
		Description: "The Victim",
		Searchable:  true,
		Format:      "string",
		Category:    "character",
		Path:        Path("victim.character_id"),
	}
	PathVictimShipTypeID = PathObj{
		Display:     "Victim Ship",
		Description: "The ship that the victim was flying",
		Searchable:  true,
		Format:      "string",
		Category:    "item",
		Path:        Path("victim.ship_type_id"),
	}
	PathVictimItemsTypeID = PathObj{
		Display:     "Victim Items",
		Description: "The items that the victim pocessed",
		Searchable:  true,
		Format:      "string",
		Category:    "item",
		Path:        Path("victim.items.item_type_id"),
	}
	PathAttackerAllianceID = PathObj{
		Display:     "Attacker Alliance",
		Description: "The alliance that the attacker belongs to",
		Searchable:  true,
		Format:      "string",
		Category:    "alliance",
		Path:        Path("attackers.alliance_id"),
	}
	PathAttackerCorporationID = PathObj{
		Display:     "Attacker Corporation",
		Description: "The corporation that the attacker belongs to",
		Searchable:  true,
		Format:      "string",
		Category:    "alliance",
		Path:        Path("attackers.corporation_id"),
	}
	PathAttackerCharacterID = PathObj{
		Display:     "Attacker",
		Description: "The Attacker",
		Searchable:  true,
		Format:      "string",
		Category:    "character",
		Path:        Path("attackers.character_id"),
	}
	PathAttackerShipTypeID = PathObj{
		Display:     "Attacker Ship",
		Description: "The Ship that the attacker was flying",
		Searchable:  true,
		Format:      "string",
		Category:    "item",
		Path:        Path("attackers.ship_type_id"),
	}
	PathAttackerWeaponTypeID = PathObj{
		Display:     "Attacker Weapon",
		Description: "The weapon that the attacker used",
		Searchable:  true,
		Format:      "string",
		Category:    "item",
		Path:        Path("attackers.weapon_type_id"),
	}
)

var AllPaths = []PathObj{
	PathSolarSystemID,
	PathZKBNPC,
	PathZKBAWOX,
	PathZKBSolo,
	PathZKBFittedValue,
	PathZKBTotalValue,
	PathVictimAllianceID,
	PathVictimCorporationID,
	PathVictimCharacterID,
	PathVictimShipTypeID,
	PathVictimItemsTypeID,
	PathAttackerAllianceID,
	PathAttackerCorporationID,
	PathAttackerCharacterID,
	PathAttackerShipTypeID,
	PathAttackerWeaponTypeID,
}

func (p Path) IsValid() bool {
	for _, v := range AllPaths {
		if v.Path == p {
			return true
		}
	}
	return false
}

func (p Path) String() string {
	return string(p)
}

type Infraction struct {
	InfrationID primitive.ObjectID `bson:"_id" json:"_id"`
	Message     string
}
