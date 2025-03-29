package data

import (
	"context"
	"database/sql"
	"errors"

	"endobit.io/metal/internal/data/db"
)

// CreateAppliance creates a new appliance in the database. It returns true if
// the appliance already exists.
func (s Store) CreateAppliance(ctx context.Context, name string, scope ZoneScope) (bool, error) {
	filter := FilterByName(name)

	if _, err := s.ReadAppliances(ctx, filter, scope); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, s.db.CreateAppliance(ctx, db.CreateApplianceParams{
				Zone: scope.Zone,
				Name: name,
			})
		}

		return false, err
	}

	return true, nil
}

// CreateApplianceAttr creates a new appliance attribute in the database. If
// returns true if the attribute already exists.
func (s Store) CreateApplianceAttr(ctx context.Context, name string, scope ApplianceScope) (bool, error) {
	if _, err := s.ReadApplianceAttrs(ctx, FilterByName(name), scope); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, s.db.CreateApplianceAttribute(ctx, db.CreateApplianceAttributeParams{
				Zone:      scope.Zone,
				Appliance: scope.Appliance,
				Name:      name,
			})
		}

		return false, err
	}

	return true, nil
}

func newAppliance(src db.ReadApplianceRow) Appliance {
	return Appliance{
		ZoneScope: ZoneScope{
			Zone: src.Zone,
		},
		Name: src.Name,
		id:   src.ID,
	}
}

// ReadAppliances reads appliances from the database.
func (s Store) ReadAppliances(ctx context.Context, filter Filter, scope ZoneScope) ([]Appliance, error) {
	type from = db.ReadApplianceRow

	var (
		appliances []Appliance
		err        error
	)

	switch {
	case scope.Zone != "" && filter.Type == NameFilter:
		var row from

		row, err = s.db.ReadAppliance(ctx, db.ReadApplianceParams{
			Zone: scope.Zone,
			Name: filter.Pattern,
		})
		appliances = append(appliances, newAppliance(row))

	case scope.Zone != "" && filter.Type == GlobFilter:
		var rows []db.ReadAppliancesByGlobRow

		rows, err = s.db.ReadAppliancesByGlob(ctx, db.ReadAppliancesByGlobParams{
			Zone: scope.Zone,
			Glob: filter.Pattern,
		})
		for i := range rows {
			appliances = append(appliances, newAppliance(from(rows[i])))
		}

	case scope.Zone != "":
		var rows []db.ReadAppliancesByZoneRow

		rows, err = s.db.ReadAppliancesByZone(ctx, db.ReadAppliancesByZoneParams{
			Zone: scope.Zone,
		})
		for i := range rows {
			appliances = append(appliances, newAppliance(from(rows[i])))
		}

	default:
		var rows []db.ReadAppliancesRow

		rows, err = s.db.ReadAppliances(ctx, db.ReadAppliancesParams{})
		for i := range rows {
			appliances = append(appliances, newAppliance(from(rows[i])))
		}
	}

	if err != nil {
		return nil, err
	}

	return appliances, nil
}

func newApplianceAttr(src db.ReadApplianceAttributeRow) ApplianceAttr {
	return ApplianceAttr{
		ApplianceScope: ApplianceScope{
			Zone:      src.Zone,
			Appliance: src.Appliance,
		},
		Attr: Attr{
			Name:        src.Name,
			Value:       value(src.Value),
			IsProtected: src.IsProtected == 1,
			id:          src.ID,
		},
	}
}

// ReadApplianceAttrs reads appliance attributes from the database.
func (s Store) ReadApplianceAttrs(ctx context.Context, filter Filter, scope ApplianceScope) ([]ApplianceAttr, error) {
	type from = db.ReadApplianceAttributeRow

	var (
		attrs []ApplianceAttr
		err   error
	)

	switch {
	case scope.Zone != "" && scope.Appliance != "" && filter.Type == NameFilter:
		var row from

		row, err = s.db.ReadApplianceAttribute(ctx, db.ReadApplianceAttributeParams{
			Zone:      scope.Zone,
			Appliance: scope.Appliance,
			Name:      filter.Pattern,
		})
		attrs = append(attrs, newApplianceAttr(row))

	case scope.Zone != "" && scope.Appliance != "" && filter.Type == GlobFilter:
		var rows []db.ReadApplianceAttributesByGlobRow

		rows, err = s.db.ReadApplianceAttributesByGlob(ctx, db.ReadApplianceAttributesByGlobParams{
			Zone:      scope.Zone,
			Appliance: scope.Appliance,
			Glob:      filter.Pattern,
		})
		for i := range rows {
			attrs = append(attrs, newApplianceAttr(from(rows[i])))
		}

	case scope.Zone != "" && scope.Appliance != "":
		var rows []db.ReadApplianceAttributesByApplianceRow

		rows, err = s.db.ReadApplianceAttributesByAppliance(ctx, db.ReadApplianceAttributesByApplianceParams{
			Zone:      scope.Zone,
			Appliance: scope.Appliance,
		})
		for i := range rows {
			attrs = append(attrs, newApplianceAttr(from(rows[i])))
		}

	case scope.Zone != "":
		var rows []db.ReadApplianceAttributesByZoneRow

		rows, err = s.db.ReadApplianceAttributesByZone(ctx, db.ReadApplianceAttributesByZoneParams{
			Zone: scope.Zone,
		})
		for i := range rows {
			attrs = append(attrs, newApplianceAttr(from(rows[i])))
		}

	default:
		var rows []db.ReadApplianceAttributesRow

		rows, err = s.db.ReadApplianceAttributes(ctx, db.ReadApplianceAttributesParams{})
		for i := range rows {
			attrs = append(attrs, newApplianceAttr(from(rows[i])))
		}
	}

	return attrs, err
}

// UpdateApplianceName updates the name of an appliance in the database.
func (s Store) UpdateApplianceName(ctx context.Context, name string, scope ApplianceScope) error {
	return s.db.UpdateApplianceName(ctx, db.UpdateApplianceNameParams{
		Zone:      scope.Zone,
		Appliance: scope.Appliance,
		Name:      name,
	})
}

// UpdateApplianceAttrName updates the name of an appliance attribute in the database.
func (s Store) UpdateApplianceAttrName(ctx context.Context, name string, scope ApplianceAttrScope) error {
	return s.db.UpdateApplianceAttributeName(ctx, db.UpdateApplianceAttributeNameParams{
		Zone:      scope.Zone,
		Appliance: scope.Appliance,
		Attr:      scope.Attr,
		Name:      name,
	})
}

// UpdateApplianceAttrValue updates the value of an appliance attribute in the database.
func (s Store) UpdateApplianceAttrValue(ctx context.Context, value string, scope ApplianceAttrScope) error {
	return s.db.UpdateApplianceAttributeValue(ctx, db.UpdateApplianceAttributeValueParams{
		Zone:      scope.Zone,
		Appliance: scope.Appliance,
		Attr:      scope.Attr,
		Value:     Optional(value),
	})
}

// UpdateApplianceAttrProtection updates the protection of an appliance attribute in the database.
func (s Store) UpdateApplianceAttrProtection(ctx context.Context, protected bool, scope ApplianceAttrScope) error {
	return s.db.UpdateApplianceAttributeProtection(ctx, db.UpdateApplianceAttributeProtectionParams{
		Zone:        scope.Zone,
		Appliance:   scope.Appliance,
		Attr:        scope.Attr,
		IsProtected: boolean(protected),
	})
}

// DeleteAppliances deletes appliances from the database.
func (s Store) DeleteAppliances(ctx context.Context, filter Filter, scope ZoneScope) error {
	rows, err := s.ReadAppliances(ctx, filter, scope)
	if err != nil {
		return err
	}

	for i := range rows {
		err := s.db.DeleteAppliance(ctx, db.DeleteApplianceParams{ID: rows[i].id})
		if err != nil {
			return err
		}
	}

	return nil
}

// DeleteApplianceAttrs deletes appliance attributes from the database.
func (s Store) DeleteApplianceAttrs(ctx context.Context, filter Filter, scope ApplianceScope) error {
	rows, err := s.ReadApplianceAttrs(ctx, filter, scope)
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
