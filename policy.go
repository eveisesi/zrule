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
	Display        string             `json:"display"`
	Description    string             `json:"description"`
	Category       PathCategory       `json:"category,omitempty"`
	Searchable     bool               `json:"searchable"`
	SearchEndpoint endpoint           `json:"searchEndpoint,omitempty"`
	Format         format             `json:"format"`
	Path           Path               `json:"path"`
	Comparators    []ruler.Comparator `json:"comparators"`
}

type format string
type endpoint string

const (
	endpointESI   endpoint = "esi"
	endpointAPI   endpoint = "api"
	formatString  format   = "string"
	formatNumber  format   = "number"
	formatBoolean format   = "boolean"
)

var (
	PathSolarSystemID = PathObj{
		Display:        "Solar System",
		Description:    "The Solar System that the Killmail occurred in",
		Category:       PathCategorySystems,
		Searchable:     true,
		SearchEndpoint: endpointAPI,
		Format:         formatString,
		Path:           Path("SolarSystemID"),
		Comparators:    []ruler.Comparator{ruler.EQ, ruler.NEQ},
	}
	PathConstellationID = PathObj{
		Display:        "Constellation",
		Description:    "The Constellation that the Killmail occurred in",
		Category:       PathCategoryConstellations,
		Searchable:     true,
		SearchEndpoint: endpointAPI,
		Format:         formatString,
		Path:           Path("ConstellationID"),
		Comparators:    []ruler.Comparator{ruler.EQ, ruler.NEQ},
	}
	PathRegionID = PathObj{
		Display:        "Region",
		Description:    "The Region that the Killmail occurred in",
		Category:       PathCategoryRegions,
		Searchable:     true,
		SearchEndpoint: endpointAPI,
		Format:         formatString,
		Path:           Path("RegionID"),
		Comparators:    []ruler.Comparator{ruler.EQ, ruler.NEQ},
	}
	PathZKBNPC = PathObj{
		Display:     "ZKillboard Is NPC",
		Description: "ZKillboard has labeled the killmail as an NPC Kill",
		Format:      formatBoolean,
		Path:        Path("Meta.NPC"),
		Comparators: []ruler.Comparator{ruler.EQ},
	}
	PathZKBAWOX = PathObj{
		Display:     "ZKillboard Is AWOX",
		Description: "Zkillboard has labeled the killmail as an AWOX Kill",
		Format:      formatBoolean,
		Path:        Path("Meta.AWOX"),
		Comparators: []ruler.Comparator{ruler.EQ},
	}
	PathZKBSolo = PathObj{
		Display:     "ZKillboard Is Solo",
		Description: "ZKIllboard has labeled the killmail as a Solo Kill",
		Format:      formatBoolean,
		Path:        Path("Meta.Solo"),
		Comparators: []ruler.Comparator{ruler.EQ},
	}
	PathZKBFittedValue = PathObj{
		Display:     "Zkillboard Fitted Value",
		Description: "The ISK value of all modules and ammo fitted to the ship",
		Format:      formatNumber,
		Path:        Path("Meta.FittedValue"),
		Comparators: []ruler.Comparator{ruler.EQ, ruler.GT, ruler.GTE, ruler.LT, ruler.LTE},
	}
	PathZKBTotalValue = PathObj{
		Display:     "ZKillboard Total Value",
		Description: "The ISK value of the killmail",
		Format:      formatNumber,
		Path:        Path("Meta.TotalValue"),
		Comparators: []ruler.Comparator{ruler.EQ, ruler.GT, ruler.GTE, ruler.LT, ruler.LTE},
	}
	PathWarID = PathObj{
		Display:     "War ID",
		Description: "The ID of the war that this killmail is involved in.",
		Format:      formatNumber,
		Path:        Path("WarID"),
		Comparators: []ruler.Comparator{ruler.EQ},
	}
	PathVictimAllianceID = PathObj{
		Display:        "Victim Alliance",
		Description:    "The alliance that the victim is/was apart of at the time of the kill",
		Searchable:     true,
		SearchEndpoint: endpointESI,
		Format:         formatString,
		Category:       PathCategoryAlliance,
		Path:           Path("Victim.AllianceID"),
		Comparators:    []ruler.Comparator{ruler.EQ, ruler.NEQ},
	}
	PathVictimCorporationID = PathObj{
		Display:        "Victim Corporation",
		Description:    "The Corporation that the victim is/was apart of at the time of the kill",
		Searchable:     true,
		SearchEndpoint: endpointESI,
		Format:         formatString,
		Category:       PathCategoryCorporation,
		Path:           Path("Victim.CorporationID"),
		Comparators:    []ruler.Comparator{ruler.EQ, ruler.NEQ},
	}
	PathVictimCharacterID = PathObj{
		Display:        "Victim Character",
		Description:    "The Victims Character",
		Searchable:     true,
		SearchEndpoint: endpointESI,
		Format:         formatString,
		Category:       PathCategoryCharacter,
		Path:           Path("Victim.CharacterID"),
		Comparators:    []ruler.Comparator{ruler.EQ},
	}
	PathVictimFactionID = PathObj{
		Display:        "Victim Faction",
		Description:    "The Faction that the Character belongs to",
		Searchable:     true,
		SearchEndpoint: endpointAPI,
		Format:         formatString,
		Category:       PathCategoryFaction,
		Path:           Path("Victim.FactionID"),
		Comparators:    []ruler.Comparator{ruler.EQ},
	}
	PathVictimShipTypeID = PathObj{
		Display:        "Victim Ship",
		Description:    "The ship that the victim was flying at the time of loss",
		Searchable:     true,
		SearchEndpoint: endpointAPI,
		Format:         formatString,
		Category:       PathCategoryItems,
		Path:           Path("Victim.ShipTypeID"),
		Comparators:    []ruler.Comparator{ruler.EQ, ruler.NEQ},
	}
	PathVictimShipGroupID = PathObj{
		Display:        "Victim Ship Group",
		Description:    "The group that the ship that the victim was flying belongs to",
		Searchable:     true,
		SearchEndpoint: endpointAPI,
		Format:         formatString,
		Category:       PathCategoryItemGroups,
		Path:           Path("Victim.ShipGroupID"),
		Comparators:    []ruler.Comparator{ruler.EQ, ruler.NEQ},
	}
	PathVictimDamageTaken = PathObj{
		Display:     "Victim Damange Sustained",
		Description: "The amount of damage applied to the ship",
		Format:      formatNumber,
		Category:    PathCategoryDamageTaken,
		Path:        Path("Victim.DamageTaken"),
		Comparators: []ruler.Comparator{ruler.EQ, ruler.GT, ruler.GTE, ruler.LT, ruler.LTE},
	}
	PathAttackerAllianceID = PathObj{
		Display:        "Attacker Alliance",
		Description:    "The alliance that the attacker is/was apart of at the time of the kill",
		Searchable:     true,
		SearchEndpoint: endpointESI,
		Format:         formatString,
		Category:       PathCategoryAlliance,
		Path:           Path("Attackers.AllianceID"),
		Comparators:    []ruler.Comparator{ruler.EQ, ruler.NEQ},
	}
	PathAttackerCorporationID = PathObj{
		Display:        "Attacker Corporation",
		Description:    "The corporation that the attacker is/was apart of at the time of the kill",
		Searchable:     true,
		SearchEndpoint: endpointESI,
		Format:         formatString,
		Category:       PathCategoryCorporation,
		Path:           Path("Attackers.CorporationID"),
		Comparators:    []ruler.Comparator{ruler.EQ, ruler.NEQ},
	}
	PathAttackerCharacterID = PathObj{
		Display:        "Attacker Character",
		Description:    "The Attacker Character",
		Searchable:     true,
		SearchEndpoint: endpointESI,
		Format:         formatString,
		Category:       PathCategoryCharacter,
		Path:           Path("Attackers.CharacterID"),
		Comparators:    []ruler.Comparator{ruler.EQ},
	}
	PathAttackerFactionID = PathObj{
		Display:        "Attacker Faction",
		Description:    "The faction that the attacker is/was apart of at the time of the kill",
		Searchable:     true,
		SearchEndpoint: endpointAPI,
		Format:         formatString,
		Category:       PathCategoryFaction,
		Path:           Path("Attackers.FactionID"),
		Comparators:    []ruler.Comparator{ruler.EQ},
	}
	PathAttackersDamageDone = PathObj{
		Display:     "Attacker Damage Done",
		Description: "The amount of damage this attacker dealt to the victim",
		Format:      formatNumber,
		Category:    PathCategoryDamageDone,
		Path:        Path("Attackers.DamageDone"),
		Comparators: []ruler.Comparator{ruler.GT, ruler.GTE, ruler.LT, ruler.LTE},
	}
	PathAttackerShipTypeID = PathObj{
		Display:        "Attacker Ship",
		Description:    "The Ship that the attacker was flying at the time of the kill",
		Searchable:     true,
		SearchEndpoint: endpointAPI,
		Format:         formatString,
		Category:       PathCategoryItems,
		Path:           Path("Attackers.ShipTypeID"),
		Comparators:    []ruler.Comparator{ruler.EQ, ruler.NEQ},
	}
	PathAttackerShipGroupID = PathObj{
		Display:        "Attacker Ship Group",
		Description:    "The group of the ship that the attacker was flying at the time of the kill",
		Searchable:     true,
		SearchEndpoint: endpointAPI,
		Format:         formatString,
		Category:       PathCategoryItemGroups,
		Path:           Path("Attackers.ShipGroupID"),
		Comparators:    []ruler.Comparator{ruler.EQ, ruler.NEQ},
	}
	PathAttackerWeaponTypeID = PathObj{
		Display:        "Attacker Weapon",
		Description:    "The weapon that the attacker used during the kill",
		Searchable:     true,
		SearchEndpoint: endpointAPI,
		Format:         formatString,
		Category:       PathCategoryItems,
		Path:           Path("Attackers.WeaponTypeID"),
		Comparators:    []ruler.Comparator{ruler.EQ, ruler.NEQ},
	}
	PathAttackerWeaponGroupID = PathObj{
		Display:        "Attacker Weapon Group",
		Description:    "The group of the weapon that the attacker was using at the time of the kill",
		Searchable:     true,
		SearchEndpoint: endpointAPI,
		Format:         formatString,
		Category:       PathCategoryItemGroups,
		Path:           Path("Attackers.WeaponGroupID"),
		Comparators:    []ruler.Comparator{ruler.EQ, ruler.NEQ},
	}
)

var AllPaths = []PathObj{
	PathSolarSystemID,
	PathConstellationID,
	PathRegionID,
	PathZKBNPC,
	PathZKBAWOX,
	PathZKBSolo,
	PathZKBFittedValue,
	PathZKBTotalValue,
	PathWarID,
	PathVictimAllianceID,
	PathVictimCorporationID,
	PathVictimCharacterID,
	PathVictimFactionID,
	PathVictimShipTypeID,
	PathAttackerAllianceID,
	PathAttackerCorporationID,
	PathAttackerCharacterID,
	PathAttackerFactionID,
	PathAttackerShipTypeID,
	PathAttackerShipGroupID,
	PathAttackerWeaponTypeID,
	PathAttackerWeaponGroupID,
}

type Path string

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
	PathCategoryAlliance       PathCategory = "alliance"
	PathCategoryCorporation    PathCategory = "corporation"
	PathCategoryCharacter      PathCategory = "character"
	PathCategoryFaction        PathCategory = "factions"
	PathCategoryRegions        PathCategory = "regions"
	PathCategoryConstellations PathCategory = "constellations"
	PathCategorySystems        PathCategory = "systems"
	PathCategoryItems          PathCategory = "items"
	PathCategoryItemGroups     PathCategory = "itemGroups"
	PathCategoryDamageTaken    PathCategory = "damageTaken"
	PathCategoryDamageDone     PathCategory = "damageDone"
)

var AllPathCategories = []PathCategory{
	PathCategoryAlliance, PathCategoryCorporation, PathCategoryCharacter,
	PathCategoryRegions, PathCategoryConstellations, PathCategorySystems,
	PathCategoryItems, PathCategoryItemGroups,
	PathCategoryDamageDone, PathCategoryDamageTaken,
	PathCategoryFaction,
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
