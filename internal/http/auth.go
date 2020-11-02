package http

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/eveisesi/zrule"
)

const loginURI = "https://login.eveonline.com/v2/oauth/authorize/?response_type=code&redirect_uri=http://zrule.local:42000/auth/callback&client_id=26ee94c69dc4459dbf25c7c0cd03d03b&state=%s"

func (s *server) handleGetAuthLogin(w http.ResponseWriter, r *http.Request) {
	var ctx = r.Context()

	cookie, _ := r.Cookie(zrule.COOKIE_zrule_AUTH_ATTEMPT)
	if cookie == nil {
		s.writeResponse(w, http.StatusBadRequest, zrule.AuthStatus{
			Status: zrule.StatusInvalid,
		})
		return
	}

	tokenKey := fmt.Sprintf(zrule.REDIS_zrule_AUTH_TOKEN, cookie.Value)
	attemptKey := fmt.Sprintf(zrule.REDIS_zrule_AUTH_ATTEMPT, cookie.Value)
	token, err := s.redis.Get(ctx, tokenKey).Result()
	if err != nil && err.Error() != "redis: nil" {
		s.logger.WithError(err).WithField("key", tokenKey).Error("unexpected error encountered from redis")
		s.writeResponse(w, http.StatusInternalServerError, nil)
	}

	attempt, err := s.redis.Get(ctx, attemptKey).Int64()
	if err != nil && err.Error() != "redis: nil" {
		s.logger.WithError(err).WithField("key", attemptKey).Error("unexpected error encountered from redis")
		s.writeResponse(w, http.StatusInternalServerError, nil)
	}

	if token == "" && attempt == 0 {
		s.writeResponse(w, http.StatusOK, zrule.AuthStatus{
			Status: zrule.StatusInvalid,
		})
		return
	} else if token == "" && attempt == 1 {
		s.writeResponse(w, http.StatusOK, zrule.AuthStatus{
			Status: zrule.StatusPending,
		})
		return
	}

	s.redis.Del(ctx, tokenKey, attemptKey)

	s.writeResponse(w, http.StatusOK, zrule.AuthStatus{
		Status: zrule.StatusCompleted,
		Token:  token,
	})

}

func (s *server) handlePostAuthLogin(w http.ResponseWriter, r *http.Request) {
	var ctx = r.Context()
	h := hmac.New(sha256.New, nil)
	_, _ = h.Write([]byte(time.Now().String()))
	hash := fmt.Sprintf("%x", h.Sum(nil))
	key := fmt.Sprintf(zrule.REDIS_zrule_AUTH_ATTEMPT, hash)
	duration := time.Minute * 5
	s.redis.Set(ctx, key, true, duration)

	http.SetCookie(w, &http.Cookie{
		Name:    zrule.COOKIE_zrule_AUTH_ATTEMPT,
		Value:   hash,
		Expires: time.Now().Add(duration),
		Path:    "/",
	})
	fmt.Printf(fmt.Sprintf("%s\n", loginURI), hash)

	s.writeResponse(w, http.StatusOK, zrule.AuthStatus{
		Status: zrule.StatusCreated,
		State:  hash,
	})

}

func (s *server) handleGetAuthCallback(w http.ResponseWriter, r *http.Request) {

	var ctx = r.Context()

	code, state, err := s.parseCodeAndStateFromURL(r.URL)
	if err != nil {
		s.writeResponse(w, http.StatusBadRequest, map[string]interface{}{"error": err})
		return
	}

	key := fmt.Sprintf(zrule.REDIS_zrule_AUTH_ATTEMPT, state)
	exists, err := s.redis.Exists(ctx, key).Result()
	if err != nil && err.Error() != "redis: nil" {
		s.logger.WithError(err).WithField("key", key).Error("unexpected error encountered from redis")
		s.writeResponse(w, http.StatusInternalServerError, nil)
	}
	if exists == 0 {
		s.writeResponse(w, http.StatusBadRequest, zrule.AuthStatus{
			Status: zrule.StatusInvalid,
		})
		return
	}

	bearer, err := s.token.BearerForCode(ctx, code)
	if err != nil {
		msg := "failed to exchange state and code for token"
		s.logger.WithError(err).Error(msg)

		return
	}

	err = s.user.VerifyUserRegistrationByToken(ctx, bearer)
	if err != nil {
		s.logger.WithError(err).Error("failed to verify user")
		s.writeResponse(w, http.StatusBadRequest, zrule.AuthStatus{
			Status: zrule.StatusInvalid,
		})
		return
	}

	s.redis.Set(
		ctx,
		fmt.Sprintf(zrule.REDIS_zrule_AUTH_TOKEN, state),
		bearer.AccessToken,
		time.Minute*5,
	)

	s.writeResponse(w, http.StatusOK, map[string]interface{}{})

}

func (s *server) parseCodeAndStateFromURL(uri *url.URL) (code, state string, err error) {

	code = uri.Query().Get("code")
	state = uri.Query().Get("state")
	if code == "" || state == "" {
		return "", "", fmt.Errorf("required paramter missing from request")
	}

	return code, state, nil

}
