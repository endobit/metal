package data

import (
	"context"
	"database/sql"
	"errors"

	"endobit.io/metal/internal/data/db"
)

// CreateGlobalAttr creates a new global attribute in the database. It returns true if
// the attribute already exists.
func (s Store) CreateGlobalAttr(ctx context.Context, name string) (bool, error) {
	if _, err := s.ReadGlobalAttrs(ctx, FilterByName(name)); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, s.db.CreateGlobalAttribute(ctx, db.CreateGlobalAttributeParams{Name: name})
		}

		return false, err
	}

	return true, nil
}

func newGlobalAttr(src db.ReadGlobalAttributeRow) *GlobalAttr {
	return &GlobalAttr{
		Attr: Attr{
			id:          src.ID,
			Name:        src.Name,
			Value:       value(src.Value),
			IsProtected: value(src.IsProtected) == 1,
		},
	}
}

// ReadGlobalAttrs reads global attributes from the database.
func (s Store) ReadGlobalAttrs(ctx context.Context, filter Filter) ([]GlobalAttr, error) {
	var (
		attrs []GlobalAttr
		err   error
	)

	switch filter.Type {
	case NameFilter:
		var row db.ReadGlobalAttributeRow

		row, err = s.db.ReadGlobalAttribute(ctx, db.ReadGlobalAttributeParams{
			Name: filter.Pattern,
		})
		attrs = append(attrs, *newGlobalAttr(row))

	case GlobFilter:
		var rows []db.ReadGlobalAttributesByGlobRow

		rows, err = s.db.ReadGlobalAttributesByGlob(ctx, db.ReadGlobalAttributesByGlobParams{
			Glob: filter.Pattern,
		})
		for i := range rows {
			attrs = append(attrs, *newGlobalAttr(db.ReadGlobalAttributeRow(rows[i])))
		}

	default:
		var rows []db.ReadGlobalAttributesRow

		rows, err = s.db.ReadGlobalAttributes(ctx, db.ReadGlobalAttributesParams{})
		for i := range rows {
			attrs = append(attrs, *newGlobalAttr(db.ReadGlobalAttributeRow(rows[i])))
		}
	}

	return attrs, err
}

// UpdateGlobalAttrName updates the name of a global attribute in the database.
func (s Store) UpdateGlobalAttrName(ctx context.Context, name string, scope GlobalAttrScope) error {
	return s.db.UpdateGlobalAttributeName(ctx, db.UpdateGlobalAttributeNameParams{
		Attr: scope.Attr,
		Name: name,
	})
}

// UpdateGlobalAttrValue updates the value of a global attribute in the database.
func (s Store) UpdateGlobalAttrValue(ctx context.Context, value string, scope GlobalAttrScope) error {
	return s.db.UpdateGlobalAttributeValue(ctx, db.UpdateGlobalAttributeValueParams{
		Attr:  scope.Attr,
		Value: Optional(value),
	})
}

// UpdateGlobalAttrProtection updates the protection of a global attribute in the database.
func (s Store) UpdateGlobalAttrProtection(ctx context.Context, protected bool, scope GlobalAttrScope) error {
	return s.db.UpdateGlobalAttributeProtection(ctx, db.UpdateGlobalAttributeProtectionParams{
		Attr:        scope.Attr,
		IsProtected: Optional(boolean(protected)),
	})
}

// DeleteGlobalAttrs deletes global attribute from the database.
func (s Store) DeleteGlobalAttrs(ctx context.Context, filter Filter) error {
	rows, err := s.ReadGlobalAttrs(ctx, filter)
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
