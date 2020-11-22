package mdb

import (
	"context"
	"fmt"
	"net/url"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/eveisesi/zrule"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Connect(ctx context.Context, uri *url.URL) (*mongo.Client, error) {

	// monitor := nrmongo.NewCommandMonitor(nil)
	// .SetMonitor(monitor)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri.String()))
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to mongo db")
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to ping mongo db")
	}

	return client, err

}

// Mongo Operators
const (
	equal            string = "$eq"
	greaterthan      string = "$gt"
	greaterthanequal string = "$gte"
	in               string = "$in"
	lessthan         string = "$lt"
	lessthanequal    string = "$lte"
	notequal         string = "$ne"
	notin            string = "$nin"
	and              string = "$and"
	or               string = "$or"
	exists           string = "$exists"
)

func BuildFilters(operators ...*zrule.Operator) primitive.D {

	var ops = make(primitive.D, 0)
	for _, a := range operators {
		switch a.Operation {
		case zrule.EqualOp:
			ops = append(ops, primitive.E{Key: a.Column, Value: primitive.D{primitive.E{Key: equal, Value: a.Value}}})
		case zrule.NotEqualOp:
			ops = append(ops, primitive.E{Key: a.Column, Value: primitive.D{primitive.E{Key: notequal, Value: a.Value}}})
		case zrule.GreaterThanOp:
			ops = append(ops, primitive.E{Key: a.Column, Value: primitive.D{primitive.E{Key: greaterthan, Value: a.Value}}})
		case zrule.GreaterThanEqualToOp:
			ops = append(ops, primitive.E{Key: a.Column, Value: primitive.D{primitive.E{Key: greaterthanequal, Value: a.Value}}})
		case zrule.LessThanOp:
			ops = append(ops, primitive.E{Key: a.Column, Value: primitive.D{primitive.E{Key: lessthan, Value: a.Value}}})
		case zrule.LessThanEqualToOp:
			ops = append(ops, primitive.E{Key: a.Column, Value: primitive.D{primitive.E{Key: lessthanequal, Value: a.Value}}})
		case zrule.ExistsOp:
			ops = append(ops, primitive.E{Key: a.Column, Value: primitive.D{primitive.E{Key: exists, Value: a.Value.(bool)}}})
		case zrule.OrOp:
			switch o := a.Value.(type) {
			case []*zrule.Operator:
				arr := make(primitive.A, 0)

				for _, op := range o {
					arr = append(arr, BuildFilters(op))
				}

				ops = append(ops, primitive.E{Key: or, Value: arr})
			default:
				panic(fmt.Sprintf("valid type %#T supplied, expected one of []*zrule.Operator", o))
			}

		case zrule.AndOp:
			switch o := a.Value.(type) {
			case []*zrule.Operator:
				arr := make(primitive.A, 0)
				for _, op := range o {
					arr = append(arr, BuildFilters(op))
				}

				ops = append(ops, primitive.E{Key: and, Value: arr})
			default:
				panic(fmt.Sprintf("valid type %#T supplied, expected one of []*zrule.Operator", o))
			}

		case zrule.InOp:
			switch o := a.Value.(type) {
			case []zrule.OpValue:
				arr := make(primitive.A, 0)
				for _, value := range o {
					arr = append(arr, value)
				}

				ops = append(ops, primitive.E{Key: a.Column, Value: primitive.D{primitive.E{Key: in, Value: arr}}})
			default:
				panic(fmt.Sprintf("valid type %#T supplied, expected one of []*zruleOpValue", o))
			}
		case zrule.NotInOp:
			switch o := a.Value.(type) {
			case []zrule.OpValue:
				arr := make(primitive.A, 0)
				for _, value := range o {
					arr = append(arr, value)
				}

				ops = append(ops, primitive.E{Key: a.Column, Value: primitive.D{primitive.E{Key: notin, Value: arr}}})
			default:
				panic(fmt.Sprintf("valid type %#T supplied, expected one of []*zrule.OpValue", o))
			}
		}
	}

	return ops

}

func BuildFindOptions(ops ...*zrule.Operator) *options.FindOptions {
	var opts = options.Find()
	for _, a := range ops {
		switch a.Operation {
		case zrule.LimitOp:
			opts.SetLimit(a.Value.(int64))
		case zrule.SkipOp:
			opts.SetSkip(a.Value.(int64))
		case zrule.OrderOp:
			opts.SetSort(primitive.D{primitive.E{Key: a.Column, Value: a.Value}})
		}
	}

	return opts
}

const duplicateKeyError = 11000

func IsUniqueConstrainViolation(exception error) bool {

	var bwe mongo.BulkWriteException
	if errors.As(exception, &bwe) {

		if len(bwe.WriteErrors) == 0 {
			return false
		}
		for _, errs := range bwe.WriteErrors {
			if errs.Code == duplicateKeyError {
				return true
			}
		}
	}
	var we mongo.WriteException
	if errors.As(exception, &we) {
		if len(we.WriteErrors) == 0 {
			return false
		}
		for _, errs := range we.WriteErrors {
			if errs.Code == duplicateKeyError {
				return true
			}
		}
	}

	return false
}

func newBool(b bool) *bool {
	return &b
}
func newString(s string) *string {
	return &s
}
