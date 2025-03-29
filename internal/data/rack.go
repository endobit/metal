package data

import (
	"context"
	"database/sql"
	"errors"

	"endobit.io/metal/internal/data/db"
)

// CreateRack creates a new rack in the database. It returns true if
// the rack already exists.
func (s Store) CreateRack(ctx context.Context, name string, scope ZoneScope) (bool, error) {
	filter := FilterByName(name)

	if _, err := s.ReadRacks(ctx, filter, scope); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, s.db.CreateRack(ctx, db.CreateRackParams{
				Zone: scope.Zone,
				Name: name,
			})
		}

		return false, err
	}

	return true, nil
}

// CreateRackAttr creates a new rack attribute in the database. If
// returns true if the attribute already exists.
func (s Store) CreateRackAttr(ctx context.Context, name string, scope RackScope) (bool, error) {
	if _, err := s.ReadRackAttrs(ctx, FilterByName(name), scope); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, s.db.CreateRackAttribute(ctx, db.CreateRackAttributeParams{
				Zone: scope.Zone,
				Rack: scope.Rack,
				Name: name,
			})
		}

		return false, err
	}

	return true, nil
}

func newRack(src db.ReadRackRow) Rack {
	return Rack{
		ZoneScope: ZoneScope{
			Zone: src.Zone,
		},
		id:   src.ID,
		Name: src.Name,
	}
}

// ReadRacks reads racks from the database.
func (s Store) ReadRacks(ctx context.Context, filter Filter, scope ZoneScope) ([]Rack, error) {
	type from = db.ReadRackRow

	var (
		racks []Rack
		err   error
	)

	switch {
	case scope.Zone != "" && filter.Type == NameFilter:
		var row db.ReadRackRow

		row, err = s.db.ReadRack(ctx, db.ReadRackParams{
			Zone: scope.Zone,
			Name: filter.Pattern,
		})
		racks = append(racks, newRack(row))

	case scope.Zone != "" && filter.Type == GlobFilter:
		var rows []db.ReadRacksByGlobRow

		rows, err = s.db.ReadRacksByGlob(ctx, db.ReadRacksByGlobParams{
			Zone: scope.Zone,
			Glob: filter.Pattern,
		})
		for i := range rows {
			racks = append(racks, newRack(from(rows[i])))
		}

	case scope.Zone != "":
		var rows []db.ReadRacksByZoneRow

		rows, err = s.db.ReadRacksByZone(ctx, db.ReadRacksByZoneParams{
			Zone: scope.Zone,
		})
		for i := range rows {
			racks = append(racks, newRack(from(rows[i])))
		}

	default:
		var rows []db.ReadRacksRow

		rows, err = s.db.ReadRacks(ctx, db.ReadRacksParams{})
		for i := range rows {
			racks = append(racks, newRack(from(rows[i])))
		}
	}

	if err != nil {
		return nil, err
	}

	return racks, nil
}

func newRackAttr(src db.ReadRackAttributeRow) RackAttr {
	return RackAttr{
		RackScope: RackScope{
			Zone: src.Zone,
			Rack: src.Rack,
		},
		Attr: Attr{
			id:          src.ID,
			Name:        src.Name,
			Value:       value(src.Value),
			IsProtected: src.IsProtected == 1,
		},
	}
}

func (s Store) ReadRackAttrs(ctx context.Context, filter Filter, scope RackScope) ([]RackAttr, error) {
	type from = db.ReadRackAttributeRow

	var (
		attrs []RackAttr
		err   error
	)

	switch {
	case scope.Zone != "" && scope.Rack != "" && filter.Type == NameFilter:
		var row from

		row, err = s.db.ReadRackAttribute(ctx, db.ReadRackAttributeParams{
			Zone: scope.Zone,
			Rack: scope.Rack,
			Name: filter.Pattern,
		})
		attrs = append(attrs, newRackAttr(row))

	case scope.Zone != "" && scope.Rack != "" && filter.Type == GlobFilter:
		var rows []db.ReadRackAttributesByGlobRow

		rows, err = s.db.ReadRackAttributesByGlob(ctx, db.ReadRackAttributesByGlobParams{
			Zone: scope.Zone,
			Rack: scope.Rack,
			Glob: filter.Pattern,
		})
		for i := range rows {
			attrs = append(attrs, newRackAttr(from(rows[i])))
		}

	case scope.Zone != "" && scope.Rack != "":
		var rows []db.ReadRackAttributesByRackRow

		rows, err = s.db.ReadRackAttributesByRack(ctx, db.ReadRackAttributesByRackParams{
			Zone: scope.Zone,
			Rack: scope.Rack,
		})
		for i := range rows {
			attrs = append(attrs, newRackAttr(from(rows[i])))
		}

	case scope.Zone != "":
		var rows []db.ReadRackAttributesByZoneRow

		rows, err = s.db.ReadRackAttributesByZone(ctx, db.ReadRackAttributesByZoneParams{
			Zone: scope.Zone,
		})
		for i := range rows {
			attrs = append(attrs, newRackAttr(from(rows[i])))
		}

	default:
		var rows []db.ReadRackAttributesRow

		rows, err = s.db.ReadRackAttributes(ctx, db.ReadRackAttributesParams{})
		for i := range rows {
			attrs = append(attrs, newRackAttr(from(rows[i])))
		}
	}

	return attrs, err
}

// UpdateRackName updates the name of a rack in the database.
func (s Store) UpdateRackName(ctx context.Context, name string, scope RackScope) error {
	return s.db.UpdateRackName(ctx, db.UpdateRackNameParams{
		Zone: scope.Zone,
		Rack: scope.Rack,
		Name: name,
	})
}

// UpdateRackAttrName updates the name of a rack attribute in the database.
func (s Store) UpdateRackAttrName(ctx context.Context, name string, scope RackAttrScope) error {
	return s.db.UpdateRackAttributeName(ctx, db.UpdateRackAttributeNameParams{
		Zone: scope.Zone,
		Rack: scope.Rack,
		Attr: scope.Attr,
		Name: name,
	})
}

// UpdateRackAttrValue updates the value of a rack attribute in the database.
func (s Store) UpdateRackAttrValue(ctx context.Context, value string, scope RackAttrScope) error {
	return s.db.UpdateRackAttributeValue(ctx, db.UpdateRackAttributeValueParams{
		Zone:  scope.Zone,
		Rack:  scope.Rack,
		Attr:  scope.Attr,
		Value: Optional(value),
	})
}

// UpdateRackAttrProtection updates the protection of a rack attribute in the database.
func (s Store) UpdateRackAttrProtection(ctx context.Context, protected bool, scope RackAttrScope) error {
	return s.db.UpdateRackAttributeProtection(ctx, db.UpdateRackAttributeProtectionParams{
		Zone:        scope.Zone,
		Rack:        scope.Rack,
		Attr:        scope.Attr,
		IsProtected: boolean(protected),
	})
}

// DeleteRacks deletes racks from the database.
func (s Store) DeleteRacks(ctx context.Context, filter Filter, scope ZoneScope) error {
	rows, err := s.ReadRacks(ctx, filter, scope)
	if err != nil {
		return err
	}

	for i := range rows {
		err := s.db.DeleteRack(ctx, db.DeleteRackParams{ID: rows[i].id})
		if err != nil {
			return err
		}
	}

	return nil
}

// DeleteRackAttrs deletes rack attributes from the database.
func (s Store) DeleteRackAttrs(ctx context.Context, filter Filter, scope RackScope) error {
	rows, err := s.ReadRackAttrs(ctx, filter, scope)
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
