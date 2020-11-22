package zrule

import (
	"context"
	"time"
)

type AllianceRepository interface {
	Alliance(ctx context.Context, id uint) (*Alliance, error)
	CreateAlliance(ctx context.Context, alliance *Alliance) (*Alliance, error)
}

// Alliance is an object representing the database table.
type Alliance struct {
	ID        uint      `bson:"id" json:"id"`
	Name      string    `bson:"name" json:"name"`
	Ticker    string    `bson:"ticker" json:"ticker"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

type CharacterRepository interface {
	Character(ctx context.Context, id uint64) (*Character, error)
	CreateCharacter(ctx context.Context, character *Character) (*Character, error)
}

type Character struct {
	ID        uint64    `bson:"id" json:"id"`
	Name      string    `bson:"name" json:"name"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

type ConstellationRepository interface {
	Constellation(ctx context.Context, id uint) (*Constellation, error)
	CreateConstellation(ctx context.Context, id uint) (*Constellation, error)
}

// Constellation is an object representing the database table.
type Constellation struct {
	ID        uint      `bson:"id" json:"id"`
	Name      string    `bson:"name" json:"name"`
	RegionID  uint      `bson:"region_id" json:"region_id"`
	FactionID *int      `bson:"faction_id" json:"faction_id"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`

	Region *Region `bson:"-" json:"-"`
}

type CorporationRepository interface {
	Corporation(ctx context.Context, id uint) (*Corporation, error)
	CreateCorporation(ctx context.Context, corporation *Corporation) (*Corporation, error)
}

type Corporation struct {
	ID        uint      `bson:"id" json:"id"`
	Name      string    `bson:"name" json:"name"`
	Ticker    string    `bson:"ticker" json:"ticker"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

type SolarSystemRepository interface {
	SolarSystem(ctx context.Context, id uint) (*SolarSystem, error)
	CreateSolarSystem(ctx context.Context, system *SolarSystem) (*SolarSystem, error)
}

type SolarSystem struct {
	ID              uint      `bson:"id" json:"id"`
	Name            string    `bson:"name" json:"name"`
	SecurityStatus  float64   `bson:"security_status" json:"security_status"`
	ConstellationID uint      `bson:"constellation_id" json:"constellation_id"`
	CreatedAt       time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt       time.Time `bson:"updated_at" json:"updated_at"`

	Constellation *Constellation `bson:"-" json:"-"`
}

type RegionRepository interface {
	Region(ctx context.Context, id uint) (*Region, error)
	CreateRegion(ctx context.Context, region *Region) (*Region, error)
}

// Region is an object representing the database table.
type Region struct {
	ID        uint      `bson:"id" json:"id"`
	Name      string    `bson:"name" json:"name"`
	FactionID *uint     `bson:"faction_id" json:"faction_id"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

type ItemRepository interface {
	Item(ctx context.Context, id uint) (*Item, error)
	CreateItem(ctx context.Context, item *Item) (*Item, error)
}

// Item is an object representing the database table.
type Item struct {
	ID            uint      `bson:"id" json:"id"`
	GroupID       uint      `bson:"group_id" json:"group_id"`
	Name          string    `bson:"name" json:"name"`
	Description   string    `bson:"description" json:"description"`
	Published     bool      `bson:"published" json:"published"`
	MarketGroupID *uint     `bson:"marketGroupID" json:"marketGroupID"`
	CreatedAt     time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt     time.Time `bson:"updated_at" json:"updated_at"`

	Group *ItemGroup `bson:"-" json:"-"`
}

type ItemGroupRepository interface {
	ItemGroup(ctx context.Context, id uint) (*ItemGroup, error)
	CreateItemGroup(ctx context.Context, group *ItemGroup) (*ItemGroup, error)
}

// ItemGroup is an object representing the database table.
type ItemGroup struct {
	ID         uint      `bson:"id" json:"id"`
	CategoryID uint      `bson:"category_id" json:"category_id"`
	Name       string    `bson:"name" json:"name"`
	Published  bool      `bson:"published" json:"published"`
	CreatedAt  time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time `bson:"updated_at" json:"updated_at"`
}
