package data

import (
	"context"

	"endobit.io/metal/internal/data/db"
)

// ReadUser reads a user from the database.
func (s Store) ReadUser(ctx context.Context, username string) (*User, error) {
	user, err := s.db.ReadUser(ctx, db.ReadUserParams{User: username})
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s Store) ReadRolesForUser(ctx context.Context, username string) ([]Role, error) {
	user, err := s.ReadUser(ctx, username)
	if err != nil {
		return nil, err
	}

	roles, err := s.db.ReadRolesForUser(ctx, db.ReadRolesForUserParams{UserID: user.ID})
	if err != nil {
		return nil, err
	}

	return roles, nil
}
