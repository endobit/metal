package data

import (
	"context"
	"database/sql"
	"errors"

	"endobit.io/metal/internal/data/db"
)

// CreateZone creates a new zone in the database. It returns true if
// the zone already exists.
func (s Store) CreateZone(ctx context.Context, name string) (bool, error) {
	if _, err := s.ReadZones(ctx, FilterByName(name)); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, s.db.CreateZone(ctx, db.CreateZoneParams{Name: name})
		}

		return false, err
	}

	return true, nil
}

// CreateZoneAttr creates a new zone attribute in the database. If
// returns true if the attribute already exists.
func (s Store) CreateZoneAttr(ctx context.Context, name string, scope ZoneScope) (bool, error) {
	if _, err := s.ReadZoneAttrs(ctx, FilterByName(name), scope); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, s.db.CreateZoneAttribute(ctx, db.CreateZoneAttributeParams{
				Zone: scope.Zone,
				Name: name,
			})
		}

		return false, err
	}

	return true, nil
}

func newZone(src db.Zone) Zone {
	return Zone{
		id:   src.ID,
		Name: src.Name,
	}
}

// ReadZones reads zones from the database.
func (s Store) ReadZones(ctx context.Context, filter Filter) ([]Zone, error) {
	var (
		zones []Zone
		err   error
	)

	switch filter.Type {
	case NameFilter:
		var row db.Zone

		row, err = s.db.ReadZone(ctx, db.ReadZoneParams{Name: filter.Pattern})
		zones = append(zones, newZone(row))

	case GlobFilter:
		var rows []db.Zone

		rows, err = s.db.ReadZonesByGlob(ctx, db.ReadZonesByGlobParams{Glob: filter.Pattern})
		for i := range rows {
			zones = append(zones, newZone(rows[i]))
		}

	default:
		var rows []db.Zone

		rows, err = s.db.ReadZones(ctx, db.ReadZonesParams{})
		for i := range rows {
			zones = append(zones, newZone(rows[i]))
		}
	}

	if err != nil {
		return nil, err
	}

	return zones, nil
}

func newZoneAttr(src db.ReadZoneAttributeRow) ZoneAttr {
	return ZoneAttr{
		ZoneScope: ZoneScope{
			Zone: src.Zone,
		},
		Attr: Attr{
			id:          src.ID,
			Name:        src.Name,
			Value:       value(src.Value),
			IsProtected: value(src.IsProtected) == 1,
		},
	}
}

// ReadZoneAttrs reads zone attributes from the database.
func (s Store) ReadZoneAttrs(ctx context.Context, filter Filter, scope ZoneScope) ([]ZoneAttr, error) {
	type from = db.ReadZoneAttributeRow

	var (
		zones []ZoneAttr
		err   error
	)

	switch {
	case scope.Zone != "" && filter.Type == NameFilter:
		var row from

		row, err = s.db.ReadZoneAttribute(ctx, db.ReadZoneAttributeParams{
			Zone: scope.Zone,
			Name: filter.Pattern,
		})
		zones = append(zones, newZoneAttr(row))

	case scope.Zone != "" && filter.Type == GlobFilter:
		var rows []db.ReadZoneAttributesByGlobRow

		rows, err = s.db.ReadZoneAttributesByGlob(ctx, db.ReadZoneAttributesByGlobParams{
			Zone: scope.Zone,
			Glob: filter.Pattern,
		})
		for i := range rows {
			zones = append(zones, newZoneAttr(from(rows[i])))
		}

	case scope.Zone != "":
		var rows []db.ReadZoneAttributesByZoneRow

		rows, err = s.db.ReadZoneAttributesByZone(ctx, db.ReadZoneAttributesByZoneParams{
			Zone: scope.Zone,
		})
		for i := range rows {
			zones = append(zones, newZoneAttr(from(rows[i])))
		}

	default:
		var rows []db.ReadZoneAttributesRow

		rows, err = s.db.ReadZoneAttributes(ctx, db.ReadZoneAttributesParams{})
		for i := range rows {
			zones = append(zones, newZoneAttr(from(rows[i])))
		}
	}

	return zones, err
}

// UpdateZoneName updates the name of a zone in the database.
func (s Store) UpdateZoneName(ctx context.Context, name string, scope ZoneScope) error {
	return s.db.UpdateZoneName(ctx, db.UpdateZoneNameParams{
		Zone: scope.Zone,
		Name: name,
	})
}

// UpdateZoneTimeZone updates the timezone of a zone in the database.
func (s Store) UpdateZoneTimeZone(ctx context.Context, timezone string, scope ZoneScope) error {
	return s.db.UpdateZoneTimeZone(ctx, db.UpdateZoneTimeZoneParams{
		Zone:     scope.Zone,
		TimeZone: Optional(timezone),
	})
}

// UpdateZoneAttrName updates the name of a zone attribute in the database.
func (s Store) UpdateZoneAttrName(ctx context.Context, name string, scope ZoneAttrScope) error {
	return s.db.UpdateZoneAttributeName(ctx, db.UpdateZoneAttributeNameParams{
		Zone: scope.Zone,
		Attr: scope.Attr,
		Name: name,
	})
}

// UpdateZoneAttrValue updates the value of a zone attribute in the database.
func (s Store) UpdateZoneAttrValue(ctx context.Context, value string, scope ZoneAttrScope) error {
	return s.db.UpdateZoneAttributeValue(ctx, db.UpdateZoneAttributeValueParams{
		Zone:  scope.Zone,
		Attr:  scope.Attr,
		Value: Optional(value),
	})
}

// UpdateZoneAttrProtection updates the protection of a zone attribute in the database.
func (s Store) UpdateZoneAttrProtection(ctx context.Context, protected bool, scope ZoneAttrScope) error {
	return s.db.UpdateZoneAttributeProtection(ctx, db.UpdateZoneAttributeProtectionParams{
		Zone:        scope.Zone,
		Attr:        scope.Attr,
		IsProtected: Optional(boolean(protected)),
	})
}

// DeleteZones deletes zones from the database.
func (s Store) DeleteZones(ctx context.Context, filter Filter) error {
	rows, err := s.ReadZones(ctx, filter)
	if err != nil {
		return err
	}

	for i := range rows {
		err := s.db.DeleteZone(ctx, db.DeleteZoneParams{ID: rows[i].id})
		if err != nil {
			return err
		}
	}

	return nil
}

// DeleteZoneAttrs deletes zones attributes from the database.
func (s Store) DeleteZoneAttrs(ctx context.Context, filter Filter, scope ZoneScope) error {
	rows, err := s.ReadZoneAttrs(ctx, filter, scope)
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
