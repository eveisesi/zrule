package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/eveisesi/zrule"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/oauth2"
)

func (s *service) VerifyUserRegistrationByToken(ctx context.Context, bearer *oauth2.Token) error {

	token, err := s.token.ParseToken(bearer.AccessToken)
	if err != nil {
		msg := "failed to parse token to valid jwt"
		s.logger.WithError(err).Error(msg)
		return fmt.Errorf("%s: %w", msg, err)
	}

	userID, err := s.token.UserIDFromToken(ctx, token)
	if err != nil {
		msg := "failed to retrieve userID from token"
		s.logger.WithError(err).Error(msg)
		return err
	}

	// Verify that we know who the character is
	character, err := s.universe.Character(ctx, userID)
	if err != nil {
		return err
	}

	owner, err := s.token.OwnerHashFromToken(ctx, token)
	if err != nil {
		return err
	}

	expires, err := s.token.ExpiresFromToken(ctx, token)
	if err != nil {
		return err
	}

	user, err := s.User(ctx, userID)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		s.logger.WithError(err).Error("encountered error searching user by token")
		return err
	}

	if errors.Is(err, mongo.ErrNoDocuments) {
		_, err = s.CreateUser(ctx, &zrule.User{
			CharacterID:  character.ID,
			OwnerHash:    owner,
			AccessToken:  token.Raw,
			RefreshToken: bearer.RefreshToken,
			Expires:      expires,
		})

		return err

	}

	if user.OwnerHash != owner {
		return fmt.Errorf("owner hash does not match token")
	}

	user.AccessToken = token.Raw
	user.Expires = expires

	_, err = s.UpdateUser(ctx, user.CharacterID, user)
	return err
}
