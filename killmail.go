package zrule

import (
	"time"
)

// type KillmailRepository interface {
// 	Killmail(ctx context.Context, id uint) (*Killmail, error)
// 	Killmails(ctx context.Context, operators ...*Operator) ([]*Killmail, error)

// 	CreateKillmail(ctx context.Context, killmail *Killmail) error
// }

type KillHash struct {
	ID   uint      `bson:"id" json:"id"`
	Hash string    `bson:"hash" json:"hash"`
	Date time.Time `bson:"date" json:"date"`
}

type Killmail struct {
	ID              uint      `json:"killmail_id"`       // bson:"killmail_id"
	Hash            string    `json:"killmail_hash"`     // bson:"killmail_hash"
	MoonID          *uint     `json:"moon_id,omitempty"` // bson:"moon_id,o
	SolarSystemID   uint      `json:"solar_system_id"`   // bson:"solar_system_id"
	ConstellationID uint      `json:"constellation_id"`
	RegionID        uint      `json:"region_id"`
	WarID           *uint     `json:"war_id,omitempty"` // bson:"war_id,o
	KillmailTime    time.Time `json:"killmail_time"`    // bson:"killmail_time"

	Attackers []*KillmailAttacker `json:"attackers"` // bson:"attackers"
	Victim    *KillmailVictim     `json:"victim"`    // bson:"victim"
	Meta      *Meta               `json:"zkb"`
}

type Meta struct {
	LocationID  uint    `json:"locationID"`
	Hash        string  `json:"hash"`
	FittedValue float64 `json:"fittedValue"`
	TotalValue  float64 `json:"totalValue"`
	Points      uint    `json:"points"`
	NPC         bool    `json:"npc"`
	Solo        bool    `json:"bool"`
	Awox        bool    `json:"awox"`
	ESI         string  `json:"esi"`
	URL         string  `json:"url"`
}

type KillmailAttacker struct {
	AllianceID     *uint   `json:"alliance_id"`     // bson:"alliance_id"
	CharacterID    *uint64 `json:"character_id"`    // bson:"character_id"
	CorporationID  *uint   `json:"corporation_id"`  // bson:"corporation_id"
	FactionID      *uint   `json:"faction_id"`      // bson:"faction_id"
	DamageDone     uint    `json:"damage_done"`     // bson:"damage_done"
	FinalBlow      bool    `json:"final_blow"`      // bson:"final_blow"
	SecurityStatus float64 `json:"security_status"` // bson:"security_status"
	ShipTypeID     *uint   `json:"ship_type_id"`    // bson:"ship_type_id"
	ShipGroupID    *uint   `json:"shipGroupID"`     // bson:"shipGroupID"
	WeaponTypeID   *uint   `json:"weapon_type_id"`  // bson:"weapon_type_id"
	WeaponGroupID  *uint   `json:"weaponGroupID"`   // bson:"weaponGroupID"
}

type KillmailVictim struct {
	AllianceID    *uint   `json:"alliance_id"`    // bson:"alliance_id"
	CharacterID   *uint64 `json:"character_id"`   // bson:"character_id"
	CorporationID *uint   `json:"corporation_id"` // bson:"corporation_id"
	FactionID     *uint   `json:"faction_id"`     // bson:"faction_id"
	DamageTaken   uint    `json:"damage_taken"`   // bson:"damage_taken"
	ShipTypeID    uint    `json:"ship_type_id"`   // bson:"ship_type_id"
	ShipGroupID   uint    `json:"ship_group_id"`  // bson:"ship_group_id"
}

type Position struct {
	X float64 `bson:"x" json:"x"`
	Y float64 `bson:"y" json:"y"`
	Z float64 `bson:"z" json:"z"`
}
