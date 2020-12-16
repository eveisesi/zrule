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
	Rules     [][]*Rule            `bson:"rules" json:"rules"`
	Actions   []primitive.ObjectID `bson:"actions" json:"actions"`
	Paused    bool                 `bson:"paused" json:"paused"`
	CreatedAt time.Time            `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time            `bson:"updated_at" json:"updated_at"`
}

type Rule struct {
	Comparator string          `bson:"comparator" json:"comparator"`
	Path       Path            `bson:"path" json:"path"`
	Values     []interface{}   `bson:"values" json:"values"`
	Entities   []*SearchResult `bson:"-" json:"entities"`
}

type PathObj struct {
	Display     string             `json:"display"`
	Description string             `json:"description"`
	Category    PathCategory       `json:"category,omitempty"`
	Searchable  bool               `json:"searchable"`
	Format      string             `json:"format"`
	Path        Path               `json:"path"`
	Comparators []ruler.Comparator `json:"comparators"`
}

type Path string

var (
	PathSolarSystemID = PathObj{
		Display:     "Solar System",
		Description: "The Solar System that the Killmail occurred in",
		Category:    "system",
		Searchable:  true,
		Format:      "string",
		Path:        Path("solar_system_id"),
		Comparators: []ruler.Comparator{ruler.EQ, ruler.NEQ},
	}
	PathZKBNPC = PathObj{
		Display:     "ZKillboard Is NPC",
		Description: "ZKillboard has labeled the killmail as an NPC Kill",
		Searchable:  false,
		Format:      "boolean",
		Path:        Path("zkb.npc"),
		Comparators: []ruler.Comparator{ruler.EQ},
	}
	PathZKBAWOX = PathObj{
		Display:     "ZKillboard Is AWOX",
		Description: "Zkillboard has labeled the killmail as an AWOX Kill",
		Searchable:  false,
		Format:      "boolean",
		Path:        Path("zkb.awox"),
		Comparators: []ruler.Comparator{ruler.EQ},
	}
	PathZKBSolo = PathObj{
		Display:     "ZKillboard Is Solo",
		Description: "ZKIllboard has labeled the killmail as a Solo Kill",
		Searchable:  false,
		Format:      "boolean",
		Path:        Path("zkb.solo"),
		Comparators: []ruler.Comparator{ruler.EQ},
	}
	PathZKBFittedValue = PathObj{
		Display:     "Zkillboard Fitted Value",
		Description: "The ISK value of all modules and ammo fitted to the ship",
		Searchable:  false,
		Format:      "number",
		Path:        Path("zkb.fittedValue"),
		Comparators: []ruler.Comparator{ruler.EQ, ruler.GT, ruler.GTE, ruler.LT, ruler.LTE},
	}
	PathZKBTotalValue = PathObj{
		Display:     "ZKillboard Total Value",
		Description: "The ISK value of the killmail",
		Searchable:  false,
		Format:      "number",
		Path:        Path("zkb.totalValue"),
		Comparators: []ruler.Comparator{ruler.EQ, ruler.GT, ruler.GTE, ruler.LT, ruler.LTE},
	}
	PathVictimAllianceID = PathObj{
		Display:     "Victim Alliance",
		Description: "The alliance that the victim is apart of",
		Searchable:  true,
		Format:      "string",
		Category:    "alliance",
		Path:        Path("victim.alliance_id"),
		Comparators: []ruler.Comparator{ruler.EQ, ruler.NEQ},
	}
	PathVictimCorporationID = PathObj{
		Display:     "Victim Corporation",
		Description: "The Corporation that the victim is apart of",
		Searchable:  true,
		Format:      "string",
		Category:    "corporation",
		Path:        Path("victim.corporation_id"),
		Comparators: []ruler.Comparator{ruler.EQ, ruler.NEQ},
	}
	PathVictimCharacterID = PathObj{
		Display:     "Victim Character",
		Description: "The Victim",
		Searchable:  true,
		Format:      "string",
		Category:    "character",
		Path:        Path("victim.character_id"),
		Comparators: []ruler.Comparator{ruler.EQ},
	}
	PathVictimShipTypeID = PathObj{
		Display:     "Victim Ship",
		Description: "The ship that the victim was flying",
		Searchable:  true,
		Format:      "string",
		Category:    "item",
		Path:        Path("victim.ship_type_id"),
		Comparators: []ruler.Comparator{ruler.EQ, ruler.NEQ},
	}
	PathAttackerAllianceID = PathObj{
		Display:     "Attacker Alliance",
		Description: "The alliance that the attacker belongs to",
		Searchable:  true,
		Format:      "string",
		Category:    "alliance",
		Path:        Path("attackers.alliance_id"),
		Comparators: []ruler.Comparator{ruler.EQ, ruler.NEQ},
	}
	PathAttackerCorporationID = PathObj{
		Display:     "Attacker Corporation",
		Description: "The corporation that the attacker belongs to",
		Searchable:  true,
		Format:      "string",
		Category:    "alliance",
		Path:        Path("attackers.corporation_id"),
		Comparators: []ruler.Comparator{ruler.EQ, ruler.NEQ},
	}
	PathAttackerCharacterID = PathObj{
		Display:     "Attacker",
		Description: "The Attacker",
		Searchable:  true,
		Format:      "string",
		Category:    "character",
		Path:        Path("attackers.character_id"),
		Comparators: []ruler.Comparator{ruler.EQ},
	}
	PathAttackerShipTypeID = PathObj{
		Display:     "Attacker Ship",
		Description: "The Ship that the attacker was flying",
		Searchable:  true,
		Format:      "string",
		Category:    "item",
		Path:        Path("attackers.ship_type_id"),
		Comparators: []ruler.Comparator{ruler.EQ, ruler.NEQ},
	}
	PathAttackerWeaponTypeID = PathObj{
		Display:     "Attacker Weapon",
		Description: "The weapon that the attacker used",
		Searchable:  true,
		Format:      "string",
		Category:    "item",
		Path:        Path("attackers.weapon_type_id"),
		Comparators: []ruler.Comparator{ruler.EQ, ruler.NEQ},
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
	// PathVictimItemsTypeID,
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

type PathCategory string

const (
	PathCategoryAlliance      PathCategory = "alliance"
	PathCategoryCorporation   PathCategory = "corporation"
	PathCategoryCharacter     PathCategory = "character"
	PathCategoryRegion        PathCategory = "region"
	PathCategoryConstellation PathCategory = "constellation"
	PathCategorySystem        PathCategory = "system"
	PathCategoryItem          PathCategory = "item"
)

var AllPathCategories = []PathCategory{
	PathCategoryAlliance, PathCategoryCorporation, PathCategoryCharacter,
	PathCategoryRegion, PathCategoryConstellation, PathCategorySystem,
	PathCategoryItem,
}

func (p PathCategory) IsValid() bool {
	for _, v := range AllPathCategories {
		if v == p {
			return true
		}
	}
	return false
}

func (p PathCategory) String() string {
	return string(p)
}

type Infraction struct {
	InfrationID primitive.ObjectID `bson:"_id" json:"_id"`
	Message     string
}
