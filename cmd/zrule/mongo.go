package main

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/eveisesi/zrule/internal/mdb"
	"go.mongodb.org/mongo-driver/mongo"
)

func makeMongoDB(cfg config) (*mongo.Database, error) {

	q := url.Values{}
	q.Set("authMechanism", cfg.Mongo.AuthMech)
	q.Set("maxIdleTimeMS", strconv.FormatInt(int64(time.Second*10), 10))
	q.Set("connectTimeoutMS", strconv.FormatInt(int64(time.Second*4), 10))
	q.Set("serverSelectionTimeoutMS", strconv.FormatInt(int64(time.Second*4), 10))
	q.Set("socketTimeoutMS", strconv.FormatInt(int64(time.Second*4), 10))
	c := &url.URL{
		Scheme:   "mongodb",
		Host:     fmt.Sprintf("%s:%d", cfg.Mongo.Host, cfg.Mongo.Port),
		User:     url.UserPassword(cfg.Mongo.User, cfg.Mongo.Pass),
		Path:     fmt.Sprintf("/%s", cfg.Mongo.Name),
		RawQuery: q.Encode(),
	}

	mc, err := mdb.Connect(context.TODO(), c)
	if err != nil {
		return nil, err
	}

	err = mc.Ping(context.TODO(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed ping mongo db server")
	}

	mdb := mc.Database(cfg.Mongo.Name)

	return mdb, nil

}
