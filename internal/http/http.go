package http

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/eveisesi/zrule/pkg/ruler"

	"github.com/eveisesi/zrule/internal/action"
	"github.com/eveisesi/zrule/internal/policy"
	"github.com/eveisesi/zrule/internal/token"
	"github.com/eveisesi/zrule/internal/universe"
	"github.com/eveisesi/zrule/internal/user"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-redis/redis/v8"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

type server struct {
	port     uint
	db       *mongo.Database
	logger   *logrus.Logger
	redis    *redis.Client
	newrelic *newrelic.Application

	token    token.Service
	user     user.Service
	action   action.Service
	policy   policy.Service
	universe universe.Service

	server *http.Server
}

func NewServer(
	port uint,
	db *mongo.Database,
	logger *logrus.Logger,
	redis *redis.Client,
	newrelic *newrelic.Application,
	token token.Service,
	user user.Service,
	action action.Service,
	policy policy.Service,
	universe universe.Service,
) *server {

	s := &server{
		port:     port,
		db:       db,
		logger:   logger,
		redis:    redis,
		newrelic: newrelic,
		token:    token,
		user:     user,
		action:   action,
		policy:   policy,
		universe: universe,
	}

	s.server = &http.Server{
		Addr:              fmt.Sprintf(":%d", port),
		ReadTimeout:       time.Second * 5,
		WriteTimeout:      time.Second * 5,
		ReadHeaderTimeout: time.Second * 5,
		IdleTimeout:       time.Second * 5,
		ErrorLog:          log.New(logger.Writer(), "", 0),
		Handler:           s.buildRouter(),
	}

	return s

}

func (s *server) Run() error {
	s.logger.WithField("port", s.port).Info("starting http server")
	return s.server.ListenAndServe()
}

func (s *server) buildRouter() *chi.Mux {

	r := chi.NewRouter()

	r.Use(
		// middleware.Timeout(time.Second*4),
		s.cors,
		// s.monitoring,
	)

	r.Get(newrelic.WrapHandleFunc(s.newrelic, "/auth/callback", s.handleGetAuthCallback))

	r.Group(func(r chi.Router) {
		r.Use(
			middleware.SetHeader("Content-Type", "application/json"),
		)
		r.Post(newrelic.WrapHandleFunc(s.newrelic, "/auth/login", s.handlePostAuthLogin))
		r.Get(newrelic.WrapHandleFunc(s.newrelic, "/auth/login/{state}", s.handleGetAuthLogin))
		r.Get(newrelic.WrapHandleFunc(s.newrelic, "/auth/url/{state}", s.handleGetAuthURL))
		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {

			err := s.db.Client().Ping(r.Context(), nil)
			if err != nil {
				s.writeError(w, http.StatusInternalServerError, fmt.Errorf("failed to ping mongo db server"))
				return
			}

		})
		r.Group(func(r chi.Router) {

			r.Use(s.auth)
			r.Get(newrelic.WrapHandleFunc(s.newrelic, "/user", s.handleGetUser))
			r.Get(newrelic.WrapHandleFunc(s.newrelic, "/paths", s.handleGetPaths))
			r.Get(newrelic.WrapHandleFunc(s.newrelic, "/policies", s.handleGetPolicies))
			r.Get(newrelic.WrapHandleFunc(s.newrelic, "/policies/{policyID}", s.handleGetPolicyByID))
			r.Get(newrelic.WrapHandleFunc(s.newrelic, "/policies/{policyID}/actions", s.handleGetPolicyActions))
			r.Post(newrelic.WrapHandleFunc(s.newrelic, "/policies", s.handleCreatePolicy))
			r.Patch(newrelic.WrapHandleFunc(s.newrelic, "/policies/{policyID}", s.handleUpdatePolicy))
			r.Delete(newrelic.WrapHandleFunc(s.newrelic, "/policies/{policyID}", s.handleDeletePolicy))
			r.Get(newrelic.WrapHandleFunc(s.newrelic, "/actions", s.handleGetActions))
			r.Post(newrelic.WrapHandleFunc(s.newrelic, "/actions", s.handleCreateAction))
			r.Post(newrelic.WrapHandleFunc(s.newrelic, "/actions/{actionID}/test", s.handlePostActionTest))
			r.Post(newrelic.WrapHandleFunc(s.newrelic, "/rules/validate", s.handlePostValidateRules))

			r.Get(newrelic.WrapHandleFunc(s.newrelic, "/search", s.handleGetSearchName))
			r.Get(newrelic.WrapHandleFunc(s.newrelic, "/search/categories", s.handleGetSearchCategories))

			r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
				s.writeResponse(w, http.StatusNoContent, nil)
			})

			r.Get(newrelic.WrapHandleFunc(s.newrelic, "/ruler/comparators", func(w http.ResponseWriter, r *http.Request) {
				s.writeResponse(w, http.StatusOK, ruler.AllComparators)
			}))

			r.Get(newrelic.WrapHandleFunc(s.newrelic, "/example", func(w http.ResponseWriter, r *http.Request) {

				_, _ = w.Write([]byte(`
				{"attackers":[{"alliance_id":99005381,"character_id":1950120006,"corporation_id":994783381,"damage_done":5132,"final_blow":true,"security_status":-9.9,"ship_type_id":29340,"weapon_type_id":2205}],"killmail_id":88486806,"killmail_time":"2020-11-09T00:39:43Z","solar_system_id":30002758,"victim":{"alliance_id":99008228,"character_id":1466390582,"corporation_id":795045209,"damage_taken":5132,"items":[{"flag":14,"item_type_id":8263,"quantity_dropped":1,"singleton":0},{"flag":13,"item_type_id":8263,"quantity_dropped":1,"singleton":0},{"flag":21,"item_type_id":4435,"quantity_dropped":1,"singleton":0},{"flag":92,"item_type_id":34268,"quantity_destroyed":1,"singleton":0},{"flag":20,"item_type_id":4435,"quantity_dropped":1,"singleton":0},{"flag":12,"item_type_id":8263,"quantity_dropped":1,"singleton":0},{"flag":11,"item_type_id":8263,"quantity_dropped":1,"singleton":0},{"flag":16,"item_type_id":8263,"quantity_dropped":1,"singleton":0},{"flag":15,"item_type_id":8263,"quantity_dropped":1,"singleton":0},{"flag":19,"item_type_id":35657,"quantity_dropped":1,"singleton":0},{"flag":22,"item_type_id":4435,"quantity_dropped":1,"singleton":0}],"position":{"x":108720300043.89568,"y":-311116434845.82336,"z":-12767733739.166466},"ship_type_id":19744},"zkb":{"locationID":40175118,"hash":"3c9ed419c00f123ff8d46475b05b70616d1381d8","fittedValue":1511799.92,"totalValue":2807175.66,"points":6,"npc":false,"solo":true,"awox":false,"esi":"https://esi.evetech.net/latest/killmails/88486806/3c9ed419c00f123ff8d46475b05b70616d1381d8/","url":"https://zkillboard.com/kill/88486806/"}}

				`))

			}))

		})
	})

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	return r

}

// GracefullyShutdown gracefully shuts down the HTTP API.
func (s *server) GracefullyShutdown(ctx context.Context) error {
	s.logger.Info("attempting to shutdown server gracefully")
	return s.server.Shutdown(ctx)
}

func (s *server) writeResponse(w http.ResponseWriter, code int, data interface{}) {

	if code != http.StatusOK {
		// if code >= http.StatusBadRequest {

		// }

		w.WriteHeader(code)

	}

	if data != nil {
		_ = json.NewEncoder(w).Encode(data)
	}
}

func (s *server) writeError(w http.ResponseWriter, code int, err error) {

	// If err is not nil, actually pass in a map so that the output to the wire is {"error": "text...."} else just let it fall through
	if err != nil {
		s.writeResponse(w, code, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	s.writeResponse(w, code, nil)

}
