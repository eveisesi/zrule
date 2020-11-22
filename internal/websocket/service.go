package websocket

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/eveisesi/zrule"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/sirupsen/logrus"
)

type Service interface {
	Run() error
}

type service struct {
	redis    *redis.Client
	logger   *logrus.Logger
	newrelic *newrelic.Application
}

func NewService(redis *redis.Client, logger *logrus.Logger, newrelic *newrelic.Application) Service {
	return &service{
		redis:    redis,
		logger:   logger,
		newrelic: newrelic,
	}
}

var (
	conn *websocket.Conn
	err  error
)

func (s *service) Run() error {

	for {
		for {
			txn := s.newrelic.StartTransaction("connect to zkillboard")

			ctx := newrelic.NewContext(context.Background(), txn)

			conn, err = s.connect(ctx)
			if err != nil {
				txn.NoticeError(err)
				s.logger.WithError(err).Error("failed to connect to zkillboard websocket")
				time.Sleep(time.Second * 2)
				txn.End()
				continue
			}

			txn.End()
			break
		}

		for {

			_, message, err := conn.ReadMessage()
			if err != nil {
				var werr *websocket.CloseError
				if errors.Is(err, werr) {
					if werr.Code == 1000 {
						s.logger.Info("websocket connection gracefully closed.")
						break
					}

					s.logger.WithError(err).Error("err encountered. Attempting to disconnect and reconnect")
				}
				eerr := conn.Close()
				if eerr != nil {
					s.logger.WithError(err).Fatal("failed to disconnect from websocket")
				}
				break
			}

			_, err = s.redis.ZAdd(context.Background(), zrule.QUEUES_KILLMAIL_PROCESSING, &redis.Z{Score: float64(time.Now().UnixNano()), Member: string(message)}).Result()
			if err != nil {
				s.logger.WithError(err).WithField("payload", string(message)).Error("unable to push killmail to processing queue")
				return err
			}
			s.logger.Info("payload dispatched successfully")
		}
	}

}

// func (s *service) PushPayloadToQueue(payload *zrule.Message) {
// 	txn := s.newrelic.StartTransaction("handle payload")
// 	defer txn.End()

// 	txn.AddAttribute("id", payload.ID)
// 	txn.AddAttribute("hash", payload.Hash)

// 	ctx := newrelic.NewContext(context.Background(), txn)

// 	data, err := json.Marshal(payload)
// 	if err != nil {
// 		txn.NoticeError(err)
// 		s.logger.WithContext(ctx).WithError(err).Error("unable to marshal WSSPayload")
// 		return
// 	}
// 	_, err = s.redis.ZAdd(ctx, zrule.QUEUES_KILLMAIL_PROCESSING, &redis.Z{Score: float64(payload.ID), Member: string(data)}).Result()
// 	if err != nil {
// 		txn.NoticeError(err)
// 		s.logger.WithContext(ctx).WithError(err).WithField("payload", string(data)).Error("unable to push killmail to processing queue")
// 		return
// 	}

// 	s.logger.WithContext(ctx).WithFields(logrus.Fields{
// 		"id":   payload.ID,
// 		"hash": payload.Hash,
// 	}).Info("payload dispatched successfully")
// }

func (s *service) connect(ctx context.Context) (*websocket.Conn, error) {
	address := url.URL{
		Scheme: "wss",
		Host:   "zkillboard.com",
		Path:   "/websocket/",
	}

	body := struct {
		Action  string `json:"action"`
		Channel string `json:"channel"`
	}{
		Action:  "sub",
		Channel: "killstream",
	}

	msg, err := json.Marshal(body)
	if err != nil {
		newrelic.FromContext(ctx).NoticeError(err)
		s.logger.WithContext(ctx).WithError(err).WithField("request", body).Error("Encoutered Error Attempting marshal sub message")
		return nil, err
	}

	s.logger.WithContext(ctx).WithField("address", address.String()).Info("attempting to connect to websocket")
	c, _, err := websocket.DefaultDialer.DialContext(ctx, address.String(), nil)
	if err != nil {
		newrelic.FromContext(ctx).NoticeError(err)
		return nil, err
	}

	s.logger.WithContext(ctx).Info("successfully connected to websocket. Sending Initial Msg")

	err = c.WriteMessage(websocket.TextMessage, msg)
	if err != nil {
		newrelic.FromContext(ctx).NoticeError(err)
		return nil, fmt.Errorf("failed to send initial message: %w", err)
	}

	return c, err
}
