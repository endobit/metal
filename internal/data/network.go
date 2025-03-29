package data

import (
	"context"
	"database/sql"
	"errors"
	"math"

	"endobit.io/metal/internal/data/db"
)

// CreateNetwork creates a new network in the database. It returns true if
// the network already exists.
func (s Store) CreateNetwork(ctx context.Context, name string, scope ZoneScope) (bool, error) {
	filter := FilterByName(name)

	if _, err := s.ReadNetworks(ctx, filter, scope); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, s.db.CreateNetwork(ctx, db.CreateNetworkParams{
				Zone: scope.Zone,
				Name: name,
			})
		}

		return false, err
	}

	return true, nil
}

func newNetwork(src db.ReadNetworkRow) Network {
	var mtu uint32

	if src.MTU > 0 && src.MTU < math.MaxUint32 {
		mtu = uint32(src.MTU)
	}

	return Network{
		ZoneScope: ZoneScope{
			Zone: src.Zone,
		},
		id:      src.ID,
		Name:    src.Name,
		Address: value(src.Address),
		Gateway: value(src.Gateway),
		IsPXE:   src.IsPXE != 0,
		MTU:     mtu,
	}
}

// ReadNetworks reads networks from the database.
func (s Store) ReadNetworks(ctx context.Context, filter Filter, scope ZoneScope) ([]Network, error) {
	type from = db.ReadNetworkRow

	var (
		networks []Network
		err      error
	)

	switch {
	case scope.Zone != "" && filter.Type == NameFilter:
		var row from

		row, err = s.db.ReadNetwork(ctx, db.ReadNetworkParams{
			Zone: scope.Zone,
			Name: filter.Pattern,
		})
		networks = append(networks, newNetwork(row))

	case scope.Zone != "" && filter.Type == GlobFilter:
		var rows []db.ReadNetworksByGlobRow

		rows, err = s.db.ReadNetworksByGlob(ctx, db.ReadNetworksByGlobParams{
			Zone: scope.Zone,
			Glob: filter.Pattern,
		})
		for i := range rows {
			networks = append(networks, newNetwork(from(rows[i])))
		}

	case scope.Zone != "":
		var rows []db.ReadNetworksByZoneRow

		rows, err = s.db.ReadNetworksByZone(ctx, db.ReadNetworksByZoneParams{
			Zone: scope.Zone,
		})
		for i := range rows {
			networks = append(networks, newNetwork(from(rows[i])))
		}

	default:
		var rows []db.ReadNetworksRow

		rows, err = s.db.ReadNetworks(ctx, db.ReadNetworksParams{})
		for i := range rows {
			networks = append(networks, newNetwork(from(rows[i])))
		}
	}

	if err != nil {
		return nil, err
	}

	return networks, nil
}

// UpdateNetworkName updates the name of a network in the database.
func (s Store) UpdateNetworkName(ctx context.Context, name string, scope NetworkScope) error {
	return s.db.UpdateNetworkName(ctx, db.UpdateNetworkNameParams{
		Zone:    scope.Zone,
		Network: scope.Network,
		Name:    name,
	})
}

// UpdateNetworkAddress updates the address of a network in the database.
func (s Store) UpdateNetworkAddress(ctx context.Context, address string, scope NetworkScope) error {
	return s.db.UpdateNetworkAddress(ctx, db.UpdateNetworkAddressParams{
		Zone:    scope.Zone,
		Network: scope.Network,
		Address: Optional(address),
	})
}

// UpdateNetworkGateway updates the gateway of a network in the database.
func (s Store) UpdateNetworkGateway(ctx context.Context, gateway string, scope NetworkScope) error {
	return s.db.UpdateNetworkGateway(ctx, db.UpdateNetworkGatewayParams{
		Zone:    scope.Zone,
		Network: scope.Network,
		Gateway: Optional(gateway),
	})
}

// UpdateNetworkPXE updates the PXE of a network in the database. True if the network supports PXE.
func (s Store) UpdateNetworkPXE(ctx context.Context, pxe bool, scope NetworkScope) error {
	var x int64

	if pxe {
		x = 1
	}

	return s.db.UpdateNetworkPXE(ctx, db.UpdateNetworkPXEParams{
		Zone:    scope.Zone,
		Network: scope.Network,
		IsPXE:   x,
	})
}

// UpdateNetworkMTU updates the mtu of a network in the database.
func (s Store) UpdateNetworkMTU(ctx context.Context, mtu uint32, scope NetworkScope) error {
	return s.db.UpdateNetworkMTU(ctx, db.UpdateNetworkMTUParams{
		Zone:    scope.Zone,
		Network: scope.Network,
		MTU:     int64(mtu),
	})
}

// DeleteNetworks deletes networks from the database.
func (s Store) DeleteNetworks(ctx context.Context, filter Filter, scope ZoneScope) error {
	rows, err := s.ReadNetworks(ctx, filter, scope)
	if err != nil {
		return err
	}

	for i := range rows {
		err := s.db.DeleteNetwork(ctx, db.DeleteNetworkParams{ID: rows[i].id})
		if err != nil {
			return err
		}
	}

	return nil
}
