package http

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/eveisesi/zrule"
	"github.com/go-chi/chi"
	"github.com/newrelic/go-agent/v3/newrelic"
	"go.mongodb.org/mongo-driver/mongo"
)

// Cors middleware to allow frontend consumption
func (s *server) cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Max-Age", "600")

		// Just return for OPTIONS
		if r.Method == "OPTIONS" {
			return
		}

		next.ServeHTTP(w, r)
	})
}

// Monitoring is middleware that will start and end a newrelic transaction
func (s *server) monitoring(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		txn := s.newrelic.StartTransaction(fmt.Sprintf("%s %s", r.Method, r.URL.Path))
		txn.SetWebRequestHTTP(r)
		rw := txn.SetWebResponse(w)
		defer txn.End()

		r = newrelic.RequestWithTransactionContext(r, txn)

		next.ServeHTTP(rw, r)

		rctx := chi.RouteContext(r.Context())
		name := rctx.RoutePattern()

		// ignore invalid routes
		if name == "/*" {
			txn.Ignore()
		}

	})
}

type userCtxKey struct{}

func (s *server) auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var ctx = r.Context()

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(strings.ToLower(authHeader), "bearer") {
			s.logger.Error("rejecting request due to missing of malformed token")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		token := authHeader[len("bearer "):]
		var validationError *jwt.ValidationError
		parsed, err := s.token.ParseToken(token)
		if err != nil {
			if _, ok := err.(*jwt.ValidationError); !ok {
				s.logger.WithError(err).Error("rejecting request: failed to parse token")
				s.writeError(w, http.StatusUnauthorized, nil)
				return
			}
			validationError = err.(*jwt.ValidationError)
			// Did validation fail for any other reason except the token being expired
			if validationError.Errors != jwt.ValidationErrorExpired {
				s.logger.WithError(err).Error("rejecting request: parsed token is invalid")
				s.writeError(w, http.StatusUnauthorized, validationError)
				return
			}
		}

		userID, err := s.token.UserIDFromToken(ctx, parsed)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Println("failed to parse user id from token: ", err)
			return
		}

		user, err := s.user.User(ctx, userID)
		if err != nil {
			if !errors.Is(err, mongo.ErrNoDocuments) {
				s.logger.WithError(err).Error("encountered error search user by token")
				w.WriteHeader(http.StatusInternalServerError)
			}
			if errors.Is(err, mongo.ErrNoDocuments) {
				s.logger.WithError(err).Error("rejecting request: could not find user for token provided")
				w.WriteHeader(http.StatusUnauthorized)
			}
			return
		}

		if validationError != nil && validationError.Errors == jwt.ValidationErrorExpired {
			newToken, newRefreshToken, err := s.token.RefreshToken(ctx, user, parsed)
			if err != nil {
				s.logger.WithError(err).Error("rejecting request: failed to refresh expired token")
				s.writeError(w, http.StatusUnauthorized, fmt.Errorf("failed to refresh token"))
				return
			}

			expires, err := s.token.ExpiresFromToken(ctx, newToken)
			if err != nil {
				s.logger.WithError(err).Error("rejecting request: failed to get expirey from new token")
				s.writeError(w, http.StatusUnauthorized, fmt.Errorf("failed to get expirey from new token"))
				return
			}

			user.AccessToken = newToken.Raw
			user.Expires = expires
			if newRefreshToken != nil {
				user.RefreshToken = *newRefreshToken
			}

			_, err = s.user.UpdateUser(ctx, user.CharacterID, user)
			if err != nil {
				s.logger.WithError(err).Error("rejecting request: failed to update data store with refresh token")
				s.writeError(w, http.StatusUnauthorized, fmt.Errorf("failed to update data store with refresh token"))
				return
			}

			w.Header().Set("X-Refreshed-Token", user.AccessToken)
		}

		ctx = context.WithValue(ctx, userCtxKey{}, user)

		next.ServeHTTP(w, r.WithContext(ctx))

	})
}

func UserFromContext(ctx context.Context) *zrule.User {
	user, ok := ctx.Value(userCtxKey{}).(*zrule.User)
	if !ok {
		return nil
	}

	return user
}
