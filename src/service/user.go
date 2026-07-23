package service

import (
	"context"
	"database/sql"
	"errors"

	"goploy/src/core/throw"
	"goploy/src/db/repos"
)

type UserService struct {
	queries *repos.Queries
}

func NewUserService(queries *repos.Queries) *UserService {
	return &UserService{queries: queries}
}

// GetByID fetches a user by their unique ID.
func (s *UserService) GetByID(
	ctx context.Context,
	id int64,
) (*repos.User, error) {
	user, err := s.queries.GetUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, throw.NotFoundError(
				"User not found",
				"USER_NOT_FOUND",
			)
		}
		return nil, throw.InternalServerError(
			"Failed to fetch user",
			"USER_FETCH_ERROR",
			throw.WithCause(err),
		)
	}
	return &user, nil
}
