package zrule

import "context"

type AllianceRespository interface {
	Alliance(ctx context.Context, id uint) (*Alliance, error)
	DeleteAlliance(ctx context.Context, id uint) error
	CreateAlliance(ctx context.Context, alliance *Alliance) error
}

// Alliance is an object representing the database table.
type Alliance struct {
	ID        uint   `bson:"id" json:"id"`
	Name      string `bson:"name" json:"name"`
	Ticker    string `bson:"ticker" json:"ticker"`
	CreatedAt int64  `bson:"createdAt" json:"createdAt"`
	UpdatedAt int64  `bson:"updatedAt" json:"updatedAt"`
}
