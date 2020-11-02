package zrule

import "context"

type CorporationRespository interface {
	Corporation(ctx context.Context, id uint) (*Corporation, error)
	CreateCorporation(ctx context.Context, corporation *Corporation) error
	DeleteCorporation(ctx context.Context, id *uint)
}

type Corporation struct {
	ID        uint   `bson:"id" json:"id"`
	Name      string `bson:"name" json:"name"`
	Ticker    string `bson:"ticker" json:"ticker"`
	CreatedAt int64  `bson:"createdAt" json:"createdAt"`
	UpdatedAt int64  `bson:"updatedAt" json:"updatedAt"`
}
