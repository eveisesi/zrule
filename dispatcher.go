package zrule

import (
	"context"
)

type Dispatcher interface {
	Send(ctx context.Context, policy *Policy, id uint, hash string) error
	SendTest(ctx context.Context, message string) error
}
