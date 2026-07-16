package service

import (
	"context"
	"database/sql"
	"errors"

	"goploy/src/apis/dtos"
	"goploy/src/core/throw"
	"goploy/src/db/repos"
	"goploy/src/utils"
)

type AuthService struct {
	db      *sql.DB
	queries *repos.Queries
	tokens  *TokenService
}

func NewAuthService(
	db *sql.DB,
	queries *repos.Queries,
	tokens *TokenService,
) *AuthService {
	return &AuthService{db: db, queries: queries, tokens: tokens}
}

// IsOwnerPresent checks if the owner has already been registered.
func (s *AuthService) IsOwnerPresent(ctx context.Context) (bool, error) {
	count, err := s.queries.IsOwnerPresent(ctx)
	if err != nil {
		return false, throw.InternalServerError(
			"Failed to check owner",
			"OWNER_CHECK_ERROR",
			throw.WithCause(err),
		)
	}
	return count > 0, nil
}

// RegisterOwner registers the first user as the application owner.
func (s *AuthService) RegisterOwner(
	ctx context.Context,
	dto *dtos.RegisterDto,
) (*repos.User, error) {
	present, err := s.IsOwnerPresent(ctx)
	if err != nil {
		return nil, err
	}
	if present {
		return nil, throw.ConflictError("Owner already exists", "OWNER_EXISTS")
	}
	hash, err := utils.HashPassword(dto.Password)
	if err != nil {
		return nil, throw.InternalServerError(
			"Failed to hash password",
			"HASH_ERROR",
			throw.WithCause(err),
		)
	}
	isOwner := int64(1)
	user, err := s.queries.CreateUser(ctx, repos.CreateUserParams{
		Email:     &dto.Email,
		FirstName: &dto.FirstName,
		LastName:  &dto.LastName,
		Password:  hash,
		IsOwner:   &isOwner,
		AddedBy:   nil,
	})
	if err != nil {
		return nil, throw.InternalServerError(
			"Failed to create owner",
			"CREATE_OWNER_ERROR",
			throw.WithCause(err),
		)
	}
	return &user, nil
}

// RegisterUser registers an invited user.
func (s *AuthService) RegisterUser(
	ctx context.Context,
	dto *dtos.RegisterDto,
	addedBy int64,
) (*repos.User, error) {
	hash, err := utils.HashPassword(dto.Password)
	if err != nil {
		return nil, throw.InternalServerError(
			"Failed to hash password",
			"HASH_ERROR",
			throw.WithCause(err),
		)
	}
	user, err := s.queries.CreateUser(ctx, repos.CreateUserParams{
		Email:     &dto.Email,
		FirstName: &dto.FirstName,
		LastName:  &dto.LastName,
		Password:  hash,
		IsOwner:   nil,
		AddedBy:   &addedBy,
	})
	if err != nil {
		return nil, throw.InternalServerError(
			"Failed to create user",
			"CREATE_USER_ERROR",
			throw.WithCause(err),
		)
	}
	return &user, nil
}

// Login authenticates a user and returns access + refresh tokens.
func (s *AuthService) Login(
	ctx context.Context,
	dto *dtos.LoginDto,
) (*AuthTokens, error) {
	user, err := s.queries.GetUserByEmail(ctx, &dto.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, throw.UnauthorizedError(
				"Invalid email or password",
				"INVALID_CREDENTIALS",
			)
		}
		return nil, throw.InternalServerError(
			"Failed to fetch user",
			"USER_FETCH_ERROR",
			throw.WithCause(err),
		)
	}
	if !utils.VerifyPassword(dto.Password, user.Password) {
		return nil, throw.UnauthorizedError(
			"Invalid email or password",
			"INVALID_CREDENTIALS",
		)
	}
	tokens, err := s.tokens.Generate(ctx, user.ID)
	if err != nil {
		return nil, err
	}
	return tokens, nil
}
