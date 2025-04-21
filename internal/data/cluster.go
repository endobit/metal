package data

import (
	"context"
	"database/sql"
	"errors"

	"endobit.io/metal/internal/data/db"
)

// CreateCluster creates a new cluster in the database. It returns true if
// the cluster already exists.
func (s Store) CreateCluster(ctx context.Context, name string, scope ZoneScope) (bool, error) {
	filter := FilterByName(name)

	if _, err := s.ReadClusters(ctx, filter, scope); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, s.db.CreateCluster(ctx, db.CreateClusterParams{
				Zone: scope.Zone,
				Name: name,
			})
		}

		return false, err
	}

	return true, nil
}

// CreateClusterAttr creates a new cluster attribute in the database. If
// returns true if the attribute already exists.
func (s Store) CreateClusterAttr(ctx context.Context, name string, scope ClusterScope) (bool, error) {
	if _, err := s.ReadClusterAttrs(ctx, FilterByName(name), scope); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, s.db.CreateClusterAttribute(ctx, db.CreateClusterAttributeParams{
				Zone:    scope.Zone,
				Cluster: scope.Cluster,
				Name:    name,
			})
		}

		return false, err
	}

	return true, nil
}

func newCluster(src db.ReadClusterRow) Cluster {
	return Cluster{
		ZoneScope: ZoneScope{
			Zone: src.Zone,
		},
		id:   src.ID,
		Name: src.Name,
	}
}

// ReadClusters reads clusters from the database.
func (s Store) ReadClusters(ctx context.Context, filter Filter, scope ZoneScope) ([]Cluster, error) {
	type from = db.ReadClusterRow

	var (
		clusters []Cluster
		err      error
	)

	switch {
	case scope.Zone != "" && filter.Type == NameFilter:
		var row from

		row, err = s.db.ReadCluster(ctx, db.ReadClusterParams{
			Zone: scope.Zone,
			Name: filter.Pattern,
		})
		clusters = append(clusters, newCluster(row))

	case scope.Zone != "" && filter.Type == GlobFilter:
		var rows []db.ReadClustersByGlobRow

		rows, err = s.db.ReadClustersByGlob(ctx, db.ReadClustersByGlobParams{
			Zone: scope.Zone,
			Glob: filter.Pattern,
		})
		for i := range rows {
			clusters = append(clusters, newCluster(from(rows[i])))
		}

	case scope.Zone != "":
		var rows []db.ReadClustersByZoneRow

		rows, err = s.db.ReadClustersByZone(ctx, db.ReadClustersByZoneParams{
			Zone: scope.Zone,
		})
		for i := range rows {
			clusters = append(clusters, newCluster(from(rows[i])))
		}

	default:
		var rows []db.ReadClustersRow

		rows, err = s.db.ReadClusters(ctx, db.ReadClustersParams{})
		for i := range rows {
			clusters = append(clusters, newCluster(from(rows[i])))
		}
	}

	if err != nil {
		return nil, err
	}

	return clusters, nil
}

func newClusterAttr(src db.ReadClusterAttributeRow) ClusterAttr {
	return ClusterAttr{
		ClusterScope: ClusterScope{
			Zone:    src.Zone,
			Cluster: src.Cluster,
		},
		Attr: Attr{
			id:          src.ID,
			Name:        src.Name,
			Value:       value(src.Value),
			IsProtected: value(src.IsProtected) == 1,
		},
	}
}

func (s Store) ReadClusterAttrs(ctx context.Context, filter Filter, scope ClusterScope) ([]ClusterAttr, error) {
	type from = db.ReadClusterAttributeRow

	var (
		attrs []ClusterAttr
		err   error
	)

	switch {
	case scope.Zone != "" && scope.Cluster != "" && filter.Type == NameFilter:
		var row from

		row, err = s.db.ReadClusterAttribute(ctx, db.ReadClusterAttributeParams{
			Zone:    scope.Zone,
			Cluster: scope.Cluster,
			Name:    filter.Pattern,
		})
		attrs = append(attrs, newClusterAttr(row))

	case scope.Zone != "" && scope.Cluster != "" && filter.Type == GlobFilter:
		var rows []db.ReadClusterAttributesByGlobRow

		rows, err = s.db.ReadClusterAttributesByGlob(ctx, db.ReadClusterAttributesByGlobParams{
			Zone:    scope.Zone,
			Cluster: scope.Cluster,
			Glob:    filter.Pattern,
		})
		for i := range rows {
			attrs = append(attrs, newClusterAttr(from(rows[i])))
		}

	case scope.Zone != "" && scope.Cluster != "":
		var rows []db.ReadClusterAttributesByClusterRow

		rows, err = s.db.ReadClusterAttributesByCluster(ctx, db.ReadClusterAttributesByClusterParams{
			Zone:    scope.Zone,
			Cluster: scope.Cluster,
		})
		for i := range rows {
			attrs = append(attrs, newClusterAttr(from(rows[i])))
		}

	case scope.Zone != "":
		var rows []db.ReadClusterAttributesByZoneRow

		rows, err = s.db.ReadClusterAttributesByZone(ctx, db.ReadClusterAttributesByZoneParams{
			Zone: scope.Zone,
		})
		for i := range rows {
			attrs = append(attrs, newClusterAttr(from(rows[i])))
		}

	default:
		var rows []db.ReadClusterAttributesRow

		rows, err = s.db.ReadClusterAttributes(ctx, db.ReadClusterAttributesParams{})
		for i := range rows {
			attrs = append(attrs, newClusterAttr(from(rows[i])))
		}
	}

	return attrs, err
}

// UpdateClusterName updates the name of a cluster in the database.
func (s Store) UpdateClusterName(ctx context.Context, name string, scope ClusterScope) error {
	return s.db.UpdateClusterName(ctx, db.UpdateClusterNameParams{
		Zone:    scope.Zone,
		Cluster: scope.Cluster,
		Name:    name,
	})
}

// UpdateClusterAttrName updates the name of a cluster attribute in the database.
func (s Store) UpdateClusterAttrName(ctx context.Context, name string, scope ClusterAttrScope) error {
	return s.db.UpdateClusterAttributeName(ctx, db.UpdateClusterAttributeNameParams{
		Zone:    scope.Zone,
		Cluster: scope.Cluster,
		Attr:    scope.Attr,
		Name:    name,
	})
}

// UpdateClusterAttrValue updates the value of a cluster attribute in the database.
func (s Store) UpdateClusterAttrValue(ctx context.Context, value string, scope ClusterAttrScope) error {
	return s.db.UpdateClusterAttributeValue(ctx, db.UpdateClusterAttributeValueParams{
		Zone:    scope.Zone,
		Cluster: scope.Cluster,
		Attr:    scope.Attr,
		Value:   Optional(value),
	})
}

// UpdateClusterAttrProtection updates the protection of a cluster attribute in the database.
func (s Store) UpdateClusterAttrProtection(ctx context.Context, protected bool, scope ClusterAttrScope) error {
	return s.db.UpdateClusterAttributeProtection(ctx, db.UpdateClusterAttributeProtectionParams{
		Zone:        scope.Zone,
		Cluster:     scope.Cluster,
		Attr:        scope.Attr,
		IsProtected: Optional(boolean(protected)),
	})
}

// DeleteClusters deletes clusters from the database.
func (s Store) DeleteClusters(ctx context.Context, filter Filter, scope ZoneScope) error {
	rows, err := s.ReadClusters(ctx, filter, scope)
	if err != nil {
		return err
	}

	for i := range rows {
		err := s.db.DeleteCluster(ctx, db.DeleteClusterParams{ID: rows[i].id})
		if err != nil {
			return err
		}
	}

	return nil
}

// DeleteClusterAttrs deletes cluster attributes from the database.
func (s Store) DeleteClusterAttrs(ctx context.Context, filter Filter, scope ClusterScope) error {
	rows, err := s.ReadClusterAttrs(ctx, filter, scope)
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
