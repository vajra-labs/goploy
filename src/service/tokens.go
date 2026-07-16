package service

import (
	"context"
	"strconv"
	"time"

	"goploy/src/conf"
	"goploy/src/core/throw"
	"goploy/src/db/repos"
	"goploy/src/pkg/jwt"
	"goploy/src/types"
)

type AuthTokens struct {
	Access  *jwt.Token
	Refresh *jwt.Token
}

type TokenService struct {
	cfg     *conf.Config
	jwt     *jwt.JwtToken
	queries *repos.Queries
}

func NewTokenService(
	cfg *conf.Config,
	jwtToken *jwt.JwtToken,
	queries *repos.Queries,
) *TokenService {
	return &TokenService{cfg: cfg, jwt: jwtToken, queries: queries}
}

// create generates a signed JWT and optionally persists it in the database.
func (s *TokenService) create(
	ctx context.Context,
	userID int64,
	tokenType types.TOKEN,
	exp time.Duration,
	save bool,
) (*jwt.Token, error) {
	sub := strconv.FormatInt(userID, 10)
	payload := s.jwt.Payload(sub, tokenType)
	token, err := s.jwt.Sign(payload, exp)
	if err != nil {
		return nil, err
	}
	if save {
		expiredAt := time.Now().Add(exp).Unix()
		_, err = s.queries.CreateJwtToken(ctx, repos.CreateJwtTokenParams{
			Jti:       payload.ID,
			UserID:    userID,
			ExpiredAt: &expiredAt,
		})
		if err != nil {
			return nil, throw.InternalServerError(
				"Failed to save token", "TOKEN_SAVE_ERROR",
				throw.WithCause(err),
			)
		}
	}
	return &token, nil
}

// Generate creates a new access + refresh token pair for the given user.
func (s *TokenService) Generate(
	ctx context.Context,
	userID int64,
) (*AuthTokens, error) {
	access, err := s.create(
		ctx,
		userID,
		types.ACC_TOKEN,
		s.cfg.JWT_ACCESS_EXP,
		false,
	)
	if err != nil {
		return nil, err
	}
	refresh, err := s.create(
		ctx,
		userID,
		types.REF_TOKEN,
		s.cfg.JWT_REFRESH_EXP,
		true,
	)
	if err != nil {
		return nil, err
	}
	return &AuthTokens{Access: access, Refresh: refresh}, nil
}

// Verify parses and validates a JWT string, and checks the token type.
func (s *TokenService) Verify(
	tokenStr string,
	tokenType types.TOKEN,
) (*jwt.Payload, error) {
	payload, err := s.jwt.Verify(tokenStr)
	if err != nil {
		return nil, err
	}
	if payload.TokenType != tokenType {
		return nil, throw.UnauthorizedError(
			"Token type is invalid",
			"INVALID_TOKEN_TYPE",
		)
	}
	return payload, nil
}

// AddBlacklist blacklists a refresh token.
// If many is true, all tokens for that user are blacklisted.
func (s *TokenService) AddBlacklist(
	ctx context.Context,
	tokenStr string,
	many bool,
) error {
	payload, err := s.Verify(tokenStr, types.REF_TOKEN)
	if err != nil {
		return err
	}
	record, err := s.queries.GetJwtTokenByJti(ctx, payload.ID)
	if err != nil {
		return throw.UnauthorizedError(
			"Token not found",
			"TOKEN_NOT_FOUND",
			throw.WithCause(err),
		)
	}
	isBlacklisted := record.IsBlacklist != nil && *record.IsBlacklist == 1
	if isBlacklisted {
		return throw.BadRequestError(
			"Token is already blacklisted",
			"TOKEN_ALREADY_BLACKLISTED",
		)
	}
	now := time.Now().Unix()
	blacklisted := int64(1)
	if many {
		userID, parseErr := strconv.ParseInt(payload.Subject, 10, 64)
		if parseErr != nil {
			return throw.InternalServerError(
				"Invalid subject in token",
				"INVALID_TOKEN_SUB",
			)
		}
		err = s.queries.UpdateJwtTokensByUserID(
			ctx,
			repos.UpdateJwtTokensByUserIDParams{
				IsBlacklist: &blacklisted,
				BlacklistAt: &now,
				UserID:      userID,
			},
		)
	} else {
		err = s.queries.UpdateJwtTokenByJti(
			ctx,
			repos.UpdateJwtTokenByJtiParams{
				IsBlacklist: &blacklisted,
				BlacklistAt: &now,
				Jti:         payload.ID,
			},
		)
	}
	if err != nil {
		return throw.InternalServerError(
			"Failed to blacklist token",
			"TOKEN_BLACKLIST_ERROR",
			throw.WithCause(err),
		)
	}
	return nil
}

// RefreshAccess issues a new access token using a valid non-blacklisted refresh token.
func (s *TokenService) RefreshAccess(
	ctx context.Context,
	refreshToken string,
) (*jwt.Token, error) {
	payload, err := s.Verify(refreshToken, types.REF_TOKEN)
	if err != nil {
		return nil, err
	}
	notBlacklisted := int64(0)
	_, err = s.queries.GetJwtTokenByJtiAndBlacklist(
		ctx,
		repos.GetJwtTokenByJtiAndBlacklistParams{
			Jti:         payload.ID,
			IsBlacklist: &notBlacklisted,
		},
	)
	if err != nil {
		return nil, throw.UnauthorizedError(
			"Refresh token is blacklisted or invalid",
			"TOKEN_BLACKLISTED",
			throw.WithCause(err),
		)
	}
	userID, err := strconv.ParseInt(payload.Subject, 10, 64)
	if err != nil {
		return nil, throw.InternalServerError(
			"Invalid subject in token",
			"INVALID_TOKEN_SUB",
		)
	}
	return s.create(ctx, userID, types.ACC_TOKEN, s.cfg.JWT_ACCESS_EXP, false)
}
