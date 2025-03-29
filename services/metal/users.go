package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"

	"endobit.io/metal/internal/data/db"
)

func setupAdminPassword(database *sql.DB, username, password string) error {
	var sqlerr *sqlite.Error

	ctx := context.Background()
	q := db.New(database)

	if err := q.CreateUser(ctx, db.CreateUserParams{Name: username}); err != nil {
		if errors.As(err, &sqlerr) && sqlerr.Code() == sqlite3.SQLITE_CONSTRAINT_UNIQUE {
			return nil // admin user already exists
		}

		return fmt.Errorf("failed to create admin user: %w", err)
	}

	err := q.UpdateUserPassword(ctx, db.UpdateUserPasswordParams{
		User:         username,
		PasswordHash: &password,
	})
	if err != nil {
		return fmt.Errorf("failed set admin password: %w", err)
	}

	user, err := q.ReadUser(ctx, db.ReadUserParams{User: username})
	if err != nil {
		return fmt.Errorf("failed to read admin user: %w", err)
	}

	err = q.CreateRole(ctx, db.CreateRoleParams{Name: "admin"})
	if err != nil {
		return fmt.Errorf("failed to create admin role: %w", err)
	}

	role, err := q.ReadRole(ctx, db.ReadRoleParams{Role: "admin"})
	if err != nil {
		return fmt.Errorf("failed to read admin role: %w", err)
	}

	return q.AssignRoleToUser(ctx, db.AssignRoleToUserParams{
		UserID: user.ID,
		RoleID: role.ID,
	})
}
