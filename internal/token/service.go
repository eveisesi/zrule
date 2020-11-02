package token

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/eveisesi/zrule"
	"github.com/go-redis/redis/v8"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

type Service interface {
	BearerForCode(ctx context.Context, code string) (*oauth2.Token, error)
	ParseToken(tokenString string) (*jwt.Token, error)
	UserIDFromToken(ctx context.Context, token *jwt.Token) (uint64, error)
	OwnerHashFromToken(ctx context.Context, token *jwt.Token) (string, error)
	ExpiresFromToken(ctx context.Context, token *jwt.Token) (time.Time, error)
	RefreshToken(ctx context.Context, user *zrule.User, token *jwt.Token) (*jwt.Token, *string, error)
}

type service struct {
	client  *http.Client
	oauth   *oauth2.Config
	logger  *logrus.Logger
	redis   *redis.Client
	jwksURL string
}

func NewService(
	client *http.Client,
	oauth *oauth2.Config,
	logger *logrus.Logger,
	redis *redis.Client,
	jwksURL string,
) Service {
	return &service{
		client:  client,
		oauth:   oauth,
		logger:  logger,
		redis:   redis,
		jwksURL: jwksURL,
	}
}

func (s *service) RefreshToken(ctx context.Context, user *zrule.User, token *jwt.Token) (*jwt.Token, *string, error) {

	if user.RefreshToken == "" {
		return nil, nil, fmt.Errorf("invalid refresh token provided")
	}

	// Locally copy a configuration of our oauth client
	config := s.oauth

	expires, err := s.ExpiresFromToken(ctx, token)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse expires fromt token")
	}

	// Rebuild a Token Model for the oauth lib
	restoreToken := &oauth2.Token{
		AccessToken:  token.Raw,
		RefreshToken: user.RefreshToken,
		Expiry:       expires,
		TokenType:    "Bearer",
	}

	// Fetch a TokenSource using rebuilt model
	tokenSource := config.TokenSource(context.Background(), restoreToken)

	// Refresh the token
	refreshedBearer, err := tokenSource.Token()
	if err != nil {
		return nil, nil, err
	}

	// Parse the new token that we got back
	parsed, err := s.ParseToken(refreshedBearer.AccessToken)
	if err != nil {
		s.logger.WithError(err).Error("failed to parse refreshed access token")
		return nil, nil, fmt.Errorf("failed to parse refreshed access token")
	}

	// Return the parsed token and potentionally a new refresh token,
	return parsed, &refreshedBearer.RefreshToken, nil

}

func (s *service) BearerForCode(ctx context.Context, code string) (*oauth2.Token, error) {

	return s.oauth.Exchange(ctx, code)

}

func (s *service) ParseToken(tokenString string) (*jwt.Token, error) {
	parser := new(jwt.Parser)
	parser.UseJSONNumber = true

	return parser.Parse(tokenString, s.getSignatureKey)

}

func (s *service) UserIDFromToken(ctx context.Context, token *jwt.Token) (uint64, error) {

	if _, ok := token.Claims.(jwt.MapClaims); !ok {
		return 0, fmt.Errorf("invalid structure to claims, expect jwt.MapClaimsm, got %T", token.Claims)
	}

	claims := token.Claims.(jwt.MapClaims)

	if _, ok := claims["sub"]; !ok {
		return 0, fmt.Errorf("claims missing sub")
	}

	if _, ok := claims["sub"].(string); !ok {
		return 0, fmt.Errorf("unexpected type for sub, expected string, got %T", claims["sub"])
	}

	sub := claims["sub"]
	expect := "CHARACTER:EVE:"
	if !strings.HasPrefix(sub.(string), expect) {
		return 0, fmt.Errorf("invalid sub, expected sub to start with %s, got %s", expect, sub)
	}

	pieces := strings.Split(sub.(string), ":")
	if len(pieces) != 3 {
		return 0, fmt.Errorf("subject malformed. expect 3 pieces, got %d", len(pieces))
	}

	userID, err := strconv.ParseUint(pieces[2], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse userID from sub, got err: %w", err)
	}

	return userID, nil

}

func (s *service) OwnerHashFromToken(ctx context.Context, token *jwt.Token) (string, error) {

	if _, ok := token.Claims.(jwt.MapClaims); !ok {
		return "", fmt.Errorf("invalid structure to claims, expected jwt.MapClaims, got %T", token.Claims)
	}

	claims := token.Claims.(jwt.MapClaims)

	if _, ok := claims["owner"]; !ok {
		return "", fmt.Errorf("claims missing owner")
	}

	if _, ok := claims["owner"].(string); !ok {
		return "", fmt.Errorf("unexpected type for owner, expected string, got %T", claims["owner"])
	}

	return claims["owner"].(string), nil

}

func (s *service) ExpiresFromToken(ctx context.Context, token *jwt.Token) (time.Time, error) {

	if _, ok := token.Claims.(jwt.MapClaims); !ok {
		return time.Time{}, fmt.Errorf("invalid structure to claims, expected jwt.MapClaims, got %T", token.Claims)
	}

	claims := token.Claims.(jwt.MapClaims)

	if _, ok := claims["exp"]; !ok {
		return time.Time{}, fmt.Errorf("claims missing exp")
	}

	if _, ok := claims["exp"].(json.Number); !ok {
		return time.Time{}, fmt.Errorf("unexpected type for exp, expected json.Number, got %T", claims["exp"])
	}

	exp, err := claims["exp"].(json.Number).Int64()
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to convert exp to unix timestamp: %w", err)
	}

	if exp == 0 {
		exp = time.Now().Add(time.Minute * 19).Unix()
	}

	return time.Unix(exp, 0), nil

}

func (s *service) getSignatureKey(token *jwt.Token) (interface{}, error) {

	ctx := context.Background()

	result, err := s.redis.Get(ctx, zrule.REDIS_CCP_JWKS).Bytes()
	if err != nil && err.Error() != "redis: nil" {
		return nil, errors.Wrap(err, "unexpected error looking for jwk in redis")
	}

	if len(result) == 0 {
		res, err := s.client.Get(s.jwksURL)
		if err != nil {
			return nil, errors.Wrap(err, "unable to retrieve jwks from sso")
		}

		if res.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("unexpected status code recieved while fetching jwks. %d", res.StatusCode)
		}

		buf, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, errors.Wrap(err, "faile dto read jwk response body")
		}

		_, err = s.redis.Set(ctx, zrule.REDIS_CCP_JWKS, buf, time.Hour*24).Result()
		if err != nil {
			return nil, errors.Wrap(err, "failed to cache jwks in redis")
		}

		result = buf
	}

	set, err := jwk.ParseBytes(result)
	if err != nil {
		return nil, errors.Wrap(err, "unable to parse jwks bytes")
	}

	keyID, ok := token.Header["kid"].(string)
	if !ok {
		return nil, errors.New("expected jwt header to have string kid")
	}

	webkey := set.LookupKeyID(keyID)
	if len(webkey) == 1 {
		return webkey[0].Materialize()
	}

	return nil, fmt.Errorf("unable to find key with kid of %s", keyID)
}
