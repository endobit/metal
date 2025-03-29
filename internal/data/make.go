package data

import (
	"context"
	"database/sql"
	"errors"

	"endobit.io/metal/internal/data/db"
)

// CreateMake creates a new make in the database. It returns true if
// the make already exists.
func (s Store) CreateMake(ctx context.Context, name string) (bool, error) {
	if _, err := s.ReadMakes(ctx, FilterByName(name)); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, s.db.CreateMake(ctx, db.CreateMakeParams{Name: name})
		}

		return false, err
	}

	return true, nil
}

func newMake(src db.Make) Make {
	return Make{
		id:   src.ID,
		Name: src.Name,
	}
}

// ReadMakes reads makes from the database.
func (s Store) ReadMakes(ctx context.Context, filter Filter) ([]Make, error) {
	var (
		makes []Make
		err   error
	)

	switch filter.Type {
	case NameFilter:
		var row db.Make

		row, err = s.db.ReadMake(ctx, db.ReadMakeParams{Name: filter.Pattern})
		makes = append(makes, newMake(row))

	case GlobFilter:
		var rows []db.Make

		rows, err = s.db.ReadMakesByGlob(ctx, db.ReadMakesByGlobParams{Glob: filter.Pattern})
		for i := range rows {
			makes = append(makes, newMake(rows[i]))
		}

	default:
		var rows []db.Make

		rows, err = s.db.ReadMakes(ctx, db.ReadMakesParams{})
		for i := range rows {
			makes = append(makes, newMake(rows[i]))
		}
	}

	if err != nil {
		return nil, err
	}

	return makes, nil
}

// UpdateMakeName updates the name of a make in the database.
func (s Store) UpdateMakeName(ctx context.Context, name string, scope MakeScope) error {
	return s.db.UpdateMakeName(ctx, db.UpdateMakeNameParams{
		Make: scope.Make,
		Name: name,
	})
}

// DeleteMakes deletes makes from the database.
func (s Store) DeleteMakes(ctx context.Context, filter Filter) error {
	rows, err := s.ReadMakes(ctx, filter)
	if err != nil {
		return err
	}

	for i := range rows {
		err := s.db.DeleteMake(ctx, db.DeleteMakeParams{ID: rows[i].id})
		if err != nil {
			return err
		}
	}

	return nil
}
