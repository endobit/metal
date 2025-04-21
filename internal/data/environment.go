package data

import (
	"context"
	"database/sql"
	"errors"

	"endobit.io/metal/internal/data/db"
)

// CreateEnvironment creates a new environment in the database. It returns true if
// the environment already exists.
func (s Store) CreateEnvironment(ctx context.Context, name string, scope ZoneScope) (bool, error) {
	filter := FilterByName(name)

	if _, err := s.ReadEnvironments(ctx, filter, scope); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, s.db.CreateEnvironment(ctx, db.CreateEnvironmentParams{
				Zone: scope.Zone,
				Name: name,
			})
		}

		return false, err
	}

	return true, nil
}

// CreateEnvironmentAttr creates a new environment attribute in the database. If
// returns true if the attribute already exists.
func (s Store) CreateEnvironmentAttr(ctx context.Context, name string, scope EnvironmentScope) (bool, error) {
	if _, err := s.ReadEnvironmentAttrs(ctx, FilterByName(name), scope); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, s.db.CreateEnvironmentAttribute(ctx, db.CreateEnvironmentAttributeParams{
				Zone:        scope.Zone,
				Environment: scope.Environment,
				Name:        name,
			})
		}

		return false, err
	}

	return true, nil
}

func newEnvironment(src db.ReadEnvironmentRow) Environment {
	return Environment{
		ZoneScope: ZoneScope{
			Zone: src.Zone,
		},
		id:   src.ID,
		Name: src.Name,
	}
}

// ReadEnvironments reads environments from the database.
func (s Store) ReadEnvironments(ctx context.Context, filter Filter, scope ZoneScope) ([]Environment, error) {
	type from = db.ReadEnvironmentRow

	var (
		environments []Environment
		err          error
	)

	switch {
	case scope.Zone != "" && filter.Type == NameFilter:
		var row from

		row, err = s.db.ReadEnvironment(ctx, db.ReadEnvironmentParams{
			Zone: scope.Zone,
			Name: filter.Pattern,
		})
		environments = append(environments, newEnvironment(row))

	case scope.Zone != "" && filter.Type == GlobFilter:
		var rows []db.ReadEnvironmentsByGlobRow

		rows, err = s.db.ReadEnvironmentsByGlob(ctx, db.ReadEnvironmentsByGlobParams{
			Zone: scope.Zone,
			Glob: filter.Pattern,
		})
		for i := range rows {
			environments = append(environments, newEnvironment(from(rows[i])))
		}

	case scope.Zone != "":
		var rows []db.ReadEnvironmentsByZoneRow

		rows, err = s.db.ReadEnvironmentsByZone(ctx, db.ReadEnvironmentsByZoneParams{
			Zone: scope.Zone,
		})
		for i := range rows {
			environments = append(environments, newEnvironment(from(rows[i])))
		}

	default:
		var rows []db.ReadEnvironmentsRow

		rows, err = s.db.ReadEnvironments(ctx, db.ReadEnvironmentsParams{})
		for i := range rows {
			environments = append(environments, newEnvironment(from(rows[i])))
		}
	}

	if err != nil {
		return nil, err
	}

	return environments, nil
}

func newEnvironmentAttr(src db.ReadEnvironmentAttributeRow) EnvironmentAttr {
	return EnvironmentAttr{
		EnvironmentScope: EnvironmentScope{
			Zone:        src.Zone,
			Environment: src.Environment,
		},
		Attr: Attr{
			id:          src.ID,
			Name:        src.Name,
			Value:       value(src.Value),
			IsProtected: value(src.IsProtected) == 1,
		},
	}
}

func (s Store) ReadEnvironmentAttrs(ctx context.Context, filter Filter, scope EnvironmentScope) ([]EnvironmentAttr, error) {
	type from = db.ReadEnvironmentAttributeRow

	var (
		attrs []EnvironmentAttr
		err   error
	)

	switch {
	case scope.Zone != "" && scope.Environment != "" && filter.Type == NameFilter:
		var row from

		row, err = s.db.ReadEnvironmentAttribute(ctx, db.ReadEnvironmentAttributeParams{
			Zone:        scope.Zone,
			Environment: scope.Environment,
			Attr:        filter.Pattern,
		})
		attrs = append(attrs, newEnvironmentAttr(row))

	case scope.Zone != "" && scope.Environment != "" && filter.Type == GlobFilter:
		var rows []db.ReadEnvironmentAttributesByGlobRow

		rows, err = s.db.ReadEnvironmentAttributesByGlob(ctx, db.ReadEnvironmentAttributesByGlobParams{
			Zone:        scope.Zone,
			Environment: scope.Environment,
			Glob:        filter.Pattern,
		})
		for i := range rows {
			attrs = append(attrs, newEnvironmentAttr(from(rows[i])))
		}

	case scope.Zone != "" && scope.Environment != "":
		var rows []db.ReadEnvironmentAttributesByEnvironmentRow

		rows, err = s.db.ReadEnvironmentAttributesByEnvironment(ctx, db.ReadEnvironmentAttributesByEnvironmentParams{
			Zone:        scope.Zone,
			Environment: scope.Environment,
		})
		for i := range rows {
			attrs = append(attrs, newEnvironmentAttr(from(rows[i])))
		}

	case scope.Zone != "":
		var rows []db.ReadEnvironmentAttributesByZoneRow

		rows, err = s.db.ReadEnvironmentAttributesByZone(ctx, db.ReadEnvironmentAttributesByZoneParams{
			Zone: scope.Zone,
		})
		for i := range rows {
			attrs = append(attrs, newEnvironmentAttr(from(rows[i])))
		}

	default:
		var rows []db.ReadEnvironmentAttributesRow

		rows, err = s.db.ReadEnvironmentAttributes(ctx, db.ReadEnvironmentAttributesParams{})
		for i := range rows {
			attrs = append(attrs, newEnvironmentAttr(from(rows[i])))
		}
	}

	return attrs, err
}

// UpdateEnvironmentName updates the name of an environment in the database.
func (s Store) UpdateEnvironmentName(ctx context.Context, name string, scope EnvironmentScope) error {
	return s.db.UpdateEnvironmentName(ctx, db.UpdateEnvironmentNameParams{
		Zone:        scope.Zone,
		Environment: scope.Environment,
		Name:        name,
	})
}

// UpdateEnvironmentAttrName updates the name of an environment attribute in the database.
func (s Store) UpdateEnvironmentAttrName(ctx context.Context, name string, scope EnvironmentAttrScope) error {
	return s.db.UpdateEnvironmentAttributeName(ctx, db.UpdateEnvironmentAttributeNameParams{
		Zone:        scope.Zone,
		Environment: scope.Environment,
		Attr:        scope.Attr,
		Name:        name,
	})
}

// UpdateEnvironmentAttrValue updates the value of an environment attribute in the database.
func (s Store) UpdateEnvironmentAttrValue(ctx context.Context, value string, scope EnvironmentAttrScope) error {
	return s.db.UpdateEnvironmentAttributeValue(ctx, db.UpdateEnvironmentAttributeValueParams{
		Zone:        scope.Zone,
		Environment: scope.Environment,
		Attr:        scope.Attr,
		Value:       Optional(value),
	})
}

// UpdateEnvironmentAttrProtection updates the protection of an environment attribute in the database.
func (s Store) UpdateEnvironmentAttrProtection(ctx context.Context, protected bool, scope EnvironmentAttrScope) error {
	return s.db.UpdateEnvironmentAttributeProtection(ctx, db.UpdateEnvironmentAttributeProtectionParams{
		Zone:        scope.Zone,
		Environment: scope.Environment,
		Attr:        scope.Attr,
		IsProtected: Optional(boolean(protected)),
	})
}

// DeleteEnvironments deletes environments from the database.
func (s Store) DeleteEnvironments(ctx context.Context, filter Filter, scope ZoneScope) error {
	rows, err := s.ReadEnvironments(ctx, filter, scope)
	if err != nil {
		return err
	}

	for i := range rows {
		err := s.db.DeleteEnvironment(ctx, db.DeleteEnvironmentParams{ID: rows[i].id})
		if err != nil {
			return err
		}
	}

	return nil
}

// DeleteEnvironmentAttrs deletes environment attributes from the database.
func (s Store) DeleteEnvironmentAttrs(ctx context.Context, filter Filter, scope EnvironmentScope) error {
	rows, err := s.ReadEnvironmentAttrs(ctx, filter, scope)
	if err != nil {
		return err
	}

	for i := range rows {
		err := s.db.DeleteAttribute(ctx, db.DeleteAttributeParams{ID: rows[i].id})
		if err != nil {
			return err
		}
	}

	return nil
}
