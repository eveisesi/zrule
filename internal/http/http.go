package http

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/eveisesi/zrule/internal/action"
	"github.com/eveisesi/zrule/internal/policy"
	"github.com/eveisesi/zrule/internal/token"
	"github.com/eveisesi/zrule/internal/user"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-redis/redis/v8"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/sirupsen/logrus"
)

type server struct {
	port     uint
	logger   *logrus.Logger
	redis    *redis.Client
	newrelic *newrelic.Application

	token  token.Service
	user   user.Service
	action action.Service
	policy policy.Service

	server *http.Server
}

func NewServer(
	port uint,
	logger *logrus.Logger,
	redis *redis.Client,
	newrelic *newrelic.Application,
	token token.Service,
	user user.Service,
	action action.Service,
	policy policy.Service,
) *server {

	s := &server{
		port:     port,
		logger:   logger,
		redis:    redis,
		newrelic: newrelic,
		token:    token,
		user:     user,
		action:   action,
		policy:   policy,
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
		middleware.SetHeader("Content-Type", "application/json"),
		middleware.Timeout(time.Second*4),
		s.cors,
		s.monitoring,
	)

	r.Post("/auth/login", s.handlePostAuthLogin)
	r.Get("/auth/login", s.handleGetAuthLogin)
	r.Get("/auth/callback", s.handleGetAuthCallback)

	r.Group(func(r chi.Router) {

		r.Use(s.auth)
		r.Get("/paths", s.handleGetPaths)
		r.Get("/policies", s.handleGetPolicies)
		r.Post("/policies", s.handleCreatePolicy)
		r.Delete("/policies/{policyID}", s.handleDeletePolicy)
		r.Get("/actions", s.handleGetActions)
		r.Post("/actions", s.handleCreateAction)
		r.Post("/actions/{actionID}/test", s.handlePostActionTest)

		r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
			s.writeResponse(w, http.StatusOK, map[string]interface{}{
				"pong": 1,
			})
		})

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

	s.writeResponse(w, code, err)

}
