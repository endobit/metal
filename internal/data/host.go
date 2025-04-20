package data

import (
	"context"
	"database/sql"
	"errors"
	"math"

	"endobit.io/metal/internal/data/db"
)

// CreateHost creates a new host in the database. It returns true if
// the host already exists.
func (s Store) CreateHost(ctx context.Context, name string, scope ClusterScope) (bool, error) {
	if _, err := s.ReadHosts(ctx, FilterByName(name), scope); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, s.db.CreateHost(ctx, db.CreateHostParams{
				Zone:    scope.Zone,
				Cluster: scope.Cluster,
				Name:    name,
			})
		}

		return false, err
	}

	return true, nil
}

// CreateHostAttr creates a new host attribute in the database. If
// returns true if the attribute already exists.
func (s Store) CreateHostAttr(ctx context.Context, name string, scope HostScope) (bool, error) {
	if _, err := s.ReadHostAttrs(ctx, FilterByName(name), scope); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, s.db.CreateHostAttribute(ctx, db.CreateHostAttributeParams{
				Zone:    scope.Zone,
				Cluster: scope.Cluster,
				Host:    scope.Host,
				Name:    name,
			})
		}

		return false, err
	}

	return true, nil
}

// CreateHostInterface creates a new host interface in the database. If
// returns true if the interface already exists.
func (s Store) CreateHostInterface(ctx context.Context, name string, scope HostScope) (bool, error) {
	if _, err := s.ReadHostInterfaces(ctx, FilterByName(name), scope); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, s.db.CreateHostInterface(ctx, db.CreateHostInterfaceParams{
				Zone:    scope.Zone,
				Cluster: scope.Cluster,
				Host:    scope.Host,
				Name:    name,
			})
		}

		return false, err
	}

	return true, nil
}

func newHost(src *db.ReadHostRow) Host {
	host := Host{
		ClusterScope: ClusterScope{
			Zone:    value(src.Zone),
			Cluster: value(src.Cluster),
		},
		id:          src.ID,
		Name:        src.Name,
		Make:        value(src.Make),
		Model:       value(src.Model),
		Environment: value(src.Environment),
		Appliance:   value(src.Appliance),
		Location:    value(src.Location),
		Rack:        value(src.Rack),
		Type:        value(src.Type),
	}

	if src.Rank != nil && *src.Rank >= 0 && *src.Rank <= math.MaxUint32 {
		x := uint32(*src.Rank)
		host.rank = &x
	}

	if src.Slot != nil && *src.Slot >= 0 && *src.Slot <= math.MaxUint32 {
		x := uint32(*src.Slot)
		host.slot = &x
	}

	return host
}

// ReadHosts reads hosts from the database.
func (s Store) ReadHosts(ctx context.Context, filter Filter, scope ClusterScope) ([]Host, error) {
	type from = *db.ReadHostRow

	var (
		hosts []Host
		err   error
	)

	switch {
	case scope.Zone != "" && filter.Type == NameFilter:
		var row db.ReadHostRow

		row, err = s.db.ReadHost(ctx, db.ReadHostParams{
			Zone:    scope.Zone,
			Cluster: scope.Cluster, // "" means standalone host (same for queries below)
			Name:    filter.Pattern,
		})
		hosts = append(hosts, newHost(&row))

	case scope.Zone != "" && filter.Type == GlobFilter:
		var rows []db.ReadHostsByGlobRow

		rows, err = s.db.ReadHostsByGlob(ctx, db.ReadHostsByGlobParams{
			Zone:    scope.Zone,
			Cluster: scope.Cluster,
			Glob:    filter.Pattern,
		})
		for i := range rows {
			hosts = append(hosts, newHost(from(&rows[i])))
		}

	case scope.Zone != "" && scope.Cluster != "":
		var rows []db.ReadHostsByClusterRow

		rows, err = s.db.ReadHostsByCluster(ctx, db.ReadHostsByClusterParams{
			Zone:    scope.Zone,
			Cluster: scope.Cluster,
		})
		for i := range rows {
			hosts = append(hosts, newHost(from(&rows[i])))
		}

	case scope.Zone != "":
		var rows []db.ReadHostsByZoneRow

		rows, err = s.db.ReadHostsByZone(ctx, db.ReadHostsByZoneParams{
			Zone: scope.Zone,
		})
		for i := range rows {
			hosts = append(hosts, newHost(from(&rows[i])))
		}

	default:
		var rows []db.ReadHostsRow

		rows, err = s.db.ReadHosts(ctx, db.ReadHostsParams{})
		for i := range rows {
			hosts = append(hosts, newHost(from(&rows[i])))
		}
	}

	return hosts, err
}

func newHostAttr(src db.ReadHostAttributeRow) HostAttr {
	return HostAttr{
		HostScope: HostScope{
			Zone:    value(src.Zone),
			Cluster: value(src.Cluster),
			Host:    src.Host,
		},
		Attr: Attr{
			id:          src.ID,
			Name:        src.Name,
			Value:       value(src.Value),
			IsProtected: src.IsProtected == 1,
		},
	}
}

// ReadHostAttrs reads host attributes from the database.
func (s Store) ReadHostAttrs(ctx context.Context, filter Filter, scope HostScope) ([]HostAttr, error) {
	type from = db.ReadHostAttributeRow

	var (
		attrs []HostAttr
		err   error
	)

	switch {
	case scope.Zone != "" && scope.Host != "" && filter.Type == NameFilter:
		var row db.ReadHostAttributeRow

		row, err = s.db.ReadHostAttribute(ctx, db.ReadHostAttributeParams{
			Zone:    scope.Zone,
			Cluster: scope.Cluster, // "" means standalone (same for queries below)
			Host:    scope.Host,
			Name:    filter.Pattern,
		})
		attrs = append(attrs, newHostAttr(row))

	case scope.Zone != "" && scope.Host != "" && filter.Type == GlobFilter:
		var rows []db.ReadHostAttributesByGlobRow

		rows, err = s.db.ReadHostAttributesByGlob(ctx, db.ReadHostAttributesByGlobParams{
			Zone:    scope.Zone,
			Cluster: scope.Cluster,
			Host:    scope.Host,
			Glob:    filter.Pattern,
		})
		for i := range rows {
			attrs = append(attrs, newHostAttr(from(rows[i])))
		}

	case scope.Zone != "" && scope.Host != "":
		var rows []db.ReadHostAttributesByHostRow

		rows, err = s.db.ReadHostAttributesByHost(ctx, db.ReadHostAttributesByHostParams{
			Zone:    scope.Zone,
			Cluster: scope.Cluster,
			Host:    scope.Host,
		})
		for i := range rows {
			attrs = append(attrs, newHostAttr(from(rows[i])))
		}

	case scope.Zone != "" && scope.Cluster != "":
		var rows []db.ReadHostAttributesByClusterRow

		rows, err = s.db.ReadHostAttributesByCluster(ctx, db.ReadHostAttributesByClusterParams{
			Zone:    scope.Zone,
			Cluster: scope.Cluster,
		})
		for i := range rows {
			attrs = append(attrs, newHostAttr(from(rows[i])))
		}

	case scope.Zone != "":
		var rows []db.ReadHostAttributesByZoneRow

		rows, err = s.db.ReadHostAttributesByZone(ctx, db.ReadHostAttributesByZoneParams{
			Zone: scope.Zone,
		})
		for i := range rows {
			attrs = append(attrs, newHostAttr(from(rows[i])))
		}

	default:
		var rows []db.ReadHostAttributesRow

		rows, err = s.db.ReadHostAttributes(ctx, db.ReadHostAttributesParams{})
		for i := range rows {
			attrs = append(attrs, newHostAttr(from(rows[i])))
		}
	}

	return attrs, err
}

func newHostInterface(src *db.ReadHostInterfaceRow) HostInterface {
	return HostInterface{
		HostScope: HostScope{
			Zone:    value(src.Zone),
			Cluster: value(src.Cluster),
			Host:    src.Host,
		},
		id:              src.ID,
		Name:            src.Name,
		IP:              value(src.IP),
		MAC:             value(src.MAC),
		IsDHCP:          src.IsDHCP == 1,
		IsPXE:           src.IsPXE == 1,
		IsManagement:    src.IsManagement == 1,
		BondMode:        value(src.BondMode),
		MasterInterface: value(src.MasterInterface),
		Network:         value(src.Network),
	}
}

func (s Store) ReadHostInterfaces(ctx context.Context, filter Filter, scope HostScope) ([]HostInterface, error) {
	type from = *db.ReadHostInterfaceRow

	var (
		ifaces []HostInterface
		err    error
	)

	switch {
	case scope.Zone != "" && scope.Host != "" && filter.Type == NameFilter:
		var row db.ReadHostInterfaceRow

		row, err = s.db.ReadHostInterface(ctx, db.ReadHostInterfaceParams{
			Zone:    scope.Zone,
			Cluster: scope.Cluster, // "" means standalone host (same for queries below)
			Host:    scope.Host,
			Name:    filter.Pattern,
		})
		ifaces = append(ifaces, newHostInterface(&row))

	case scope.Zone != "" && scope.Host != "" && filter.Type == GlobFilter:
		var rows []db.ReadHostInterfacesByGlobRow

		rows, err = s.db.ReadHostInterfacesByGlob(ctx, db.ReadHostInterfacesByGlobParams{
			Zone:    scope.Zone,
			Cluster: scope.Cluster,
			Host:    scope.Host,
			Glob:    filter.Pattern,
		})
		for i := range rows {
			ifaces = append(ifaces, newHostInterface(from(&rows[i])))
		}

	case scope.Zone != "" && scope.Host != "":
		var rows []db.ReadHostInterfacesByHostRow

		rows, err = s.db.ReadHostInterfacesByHost(ctx, db.ReadHostInterfacesByHostParams{
			Zone:    scope.Zone,
			Cluster: scope.Cluster,
			Host:    scope.Host,
		})
		for i := range rows {
			ifaces = append(ifaces, newHostInterface(from(&rows[i])))
		}

	case scope.Zone != "" && scope.Cluster != "":
		var rows []db.ReadHostInterfacesByClusterRow

		rows, err = s.db.ReadHostInterfacesByCluster(ctx, db.ReadHostInterfacesByClusterParams{
			Zone:    scope.Zone,
			Cluster: scope.Cluster,
		})
		for i := range rows {
			ifaces = append(ifaces, newHostInterface(from(&rows[i])))
		}

	case scope.Zone != "":
		var rows []db.ReadHostInterfacesByZoneRow

		rows, err = s.db.ReadHostInterfacesByZone(ctx, db.ReadHostInterfacesByZoneParams{
			Zone: scope.Zone,
		})
		for i := range rows {
			ifaces = append(ifaces, newHostInterface(from(&rows[i])))
		}

	default:
		var rows []db.ReadHostInterfacesRow

		rows, err = s.db.ReadHostInterfaces(ctx, db.ReadHostInterfacesParams{})
		for i := range rows {
			ifaces = append(ifaces, newHostInterface(from(&rows[i])))
		}
	}

	if err != nil {
		return nil, err
	}

	return ifaces, nil
}

// UpdateHostName updates the name of a host in the database.
func (s Store) UpdateHostName(ctx context.Context, name string, scope HostScope) error {
	return s.db.UpdateHostName(ctx, db.UpdateHostNameParams{
		Zone:    scope.Zone,
		Cluster: scope.Cluster,
		Host:    scope.Host,
		Name:    name,
	})
}

// UpdateHostType updates the type of a host in the database.
func (s Store) UpdateHostType(ctx context.Context, t string, scope HostScope) error {
	return s.db.UpdateHostType(ctx, db.UpdateHostTypeParams{
		Zone:    scope.Zone,
		Cluster: scope.Cluster,
		Host:    scope.Host,
		Type:    Optional(t),
	})
}

// UpdateHostModel updates the model of a host in the database.
func (s Store) UpdateHostModel(ctx context.Context, vendor, model string, scope HostScope) error {
	return s.db.UpdateHostModel(ctx, db.UpdateHostModelParams{
		Zone:    scope.Zone,
		Cluster: scope.Cluster,
		Host:    scope.Host,
		Make:    vendor,
		Model:   model,
	})
}

// UpdateHostEnvironment updates the environment of a host in the database.
func (s Store) UpdateHostEnvironment(ctx context.Context, environment string, scope HostScope) error {
	return s.db.UpdateHostEnvironment(ctx, db.UpdateHostEnvironmentParams{
		Zone:        scope.Zone,
		Cluster:     scope.Cluster,
		Host:        scope.Host,
		Environment: environment,
	})
}

// UpdateHostAppliance updates the appliance of a host in the database.
func (s Store) UpdateHostAppliance(ctx context.Context, appliance string, scope HostScope) error {
	return s.db.UpdateHostAppliance(ctx, db.UpdateHostApplianceParams{
		Zone:      scope.Zone,
		Cluster:   scope.Cluster,
		Host:      scope.Host,
		Appliance: appliance,
	})
}

// UpdateHostLocation updates the location of a host in the database.
func (s Store) UpdateHostLocation(ctx context.Context, location string, scope HostScope) error {
	return s.db.UpdateHostLocation(ctx, db.UpdateHostLocationParams{
		Zone:     scope.Zone,
		Cluster:  scope.Cluster,
		Host:     scope.Host,
		Location: Optional(location),
	})
}

// UpdateHostRack updates the rack of a host in the database.
func (s Store) UpdateHostRack(ctx context.Context, rack string, scope HostScope) error {
	return s.db.UpdateHostRack(ctx, db.UpdateHostRackParams{
		Zone:    scope.Zone,
		Cluster: scope.Cluster,
		Host:    scope.Host,
		Rack:    rack,
	})
}

// UpdateHostRank updates the rank of a host in the database.
func (s Store) UpdateHostRank(ctx context.Context, rank int64, scope HostScope) error {
	var rp *int64
	if rank >= 0 {
		rp = &rank
	}

	return s.db.UpdateHostRank(ctx, db.UpdateHostRankParams{
		Zone:    scope.Zone,
		Cluster: scope.Cluster,
		Host:    scope.Host,
		Rank:    rp,
	})
}

// UpdateHostSlot updates the slot of a host in the database.
func (s Store) UpdateHostSlot(ctx context.Context, slot int64, scope HostScope) error {
	var sp *int64
	if slot >= 0 {
		sp = &slot
	}

	return s.db.UpdateHostSlot(ctx, db.UpdateHostSlotParams{
		Zone:    scope.Zone,
		Cluster: scope.Cluster,
		Host:    scope.Host,
		Slot:    sp,
	})
}

// UpdateHostAttrName updates the name of a host attribute in the database.
func (s Store) UpdateHostAttrName(ctx context.Context, name string, scope HostAttrScope) error {
	return s.db.UpdateHostAttributeName(ctx, db.UpdateHostAttributeNameParams{
		Zone:    scope.Zone,
		Cluster: scope.Cluster,
		Host:    scope.Host,
		Attr:    scope.Attr,
		Name:    name,
	})
}

// UpdateHostAttrValue updates the value of a host attribute in the database.
func (s Store) UpdateHostAttrValue(ctx context.Context, value string, scope HostAttrScope) error {
	return s.db.UpdateHostAttributeValue(ctx, db.UpdateHostAttributeValueParams{
		Zone:    scope.Zone,
		Cluster: scope.Cluster,
		Host:    scope.Host,
		Attr:    scope.Attr,
		Value:   Optional(value),
	})
}

// UpdateHostAttrProtection updates the protection of a host attribute in the database.
func (s Store) UpdateHostAttrProtection(ctx context.Context, protected bool, scope HostAttrScope) error {
	return s.db.UpdateHostAttributeProtection(ctx, db.UpdateHostAttributeProtectionParams{
		Zone:        scope.Zone,
		Cluster:     scope.Cluster,
		Host:        scope.Host,
		Attr:        scope.Attr,
		IsProtected: boolean(protected),
	})
}

// UpdateHostInterfaceName updates the name of a host interface in the database.
func (s Store) UpdateHostInterfaceName(ctx context.Context, name string, scope HostInterfaceScope) error {
	return s.db.UpdateHostInterfaceName(ctx, db.UpdateHostInterfaceNameParams{
		Zone:      scope.Zone,
		Cluster:   scope.Cluster,
		Host:      scope.Host,
		Interface: scope.Interface,
		Name:      name,
	})
}

// UpdateHostInterfaceIP updates the ip of a host interface in the database.
func (s Store) UpdateHostInterfaceIP(ctx context.Context, ip string, scope HostInterfaceScope) error {
	return s.db.UpdateHostInterfaceIP(ctx, db.UpdateHostInterfaceIPParams{
		Zone:      scope.Zone,
		Cluster:   scope.Cluster,
		Host:      scope.Host,
		Interface: scope.Interface,
		IP:        Optional(ip),
	})
}

// UpdateHostInterfaceMAC updates the mac of a host interface in the database.
func (s Store) UpdateHostInterfaceMAC(ctx context.Context, mac string, scope HostInterfaceScope) error {
	return s.db.UpdateHostInterfaceMAC(ctx, db.UpdateHostInterfaceMACParams{
		Zone:      scope.Zone,
		Cluster:   scope.Cluster,
		Host:      scope.Host,
		Interface: scope.Interface,
		MAC:       Optional(mac),
	})
}

// UpdateHostInterfaceNetmask updates the netmask of a host interface in the database.
func (s Store) UpdateHostInterfaceNetmask(ctx context.Context, mask string, scope HostInterfaceScope) error {
	return s.db.UpdateHostInterfaceNetmask(ctx, db.UpdateHostInterfaceNetmaskParams{
		Zone:      scope.Zone,
		Cluster:   scope.Cluster,
		Host:      scope.Host,
		Interface: scope.Interface,
		Mask:      Optional(mask),
	})
}

// UpdateHostInterfaceDHCP updates the dhcp of a host interface in the database.
func (s Store) UpdateHostInterfaceDHCP(ctx context.Context, dhcp bool, scope HostInterfaceScope) error {
	var x int64

	if dhcp {
		x = 1
	}
	return s.db.UpdateHostInterfaceDHCP(ctx, db.UpdateHostInterfaceDHCPParams{
		Zone:      scope.Zone,
		Cluster:   scope.Cluster,
		Host:      scope.Host,
		Interface: scope.Interface,
		DHCP:      x,
	})
}

// UpdateHostInterfacePXE updates the pxe of a host interface in the database.
func (s Store) UpdateHostInterfacePXE(ctx context.Context, pxe bool, scope HostInterfaceScope) error {
	var x int64

	if pxe {
		x = 1
	}
	return s.db.UpdateHostInterfacePXE(ctx, db.UpdateHostInterfacePXEParams{
		Zone:      scope.Zone,
		Cluster:   scope.Cluster,
		Host:      scope.Host,
		Interface: scope.Interface,
		PXE:       x,
	})
}

// UpdateHostInterfaceManagement updates the management of a host interface in the database.
func (s Store) UpdateHostInterfaceManagement(ctx context.Context, management bool, scope HostInterfaceScope) error {
	var x int64

	if management {
		x = 1
	}

	return s.db.UpdateHostInterfaceManagement(ctx, db.UpdateHostInterfaceManagementParams{
		Zone:       scope.Zone,
		Cluster:    scope.Cluster,
		Host:       scope.Host,
		Interface:  scope.Interface,
		Management: x,
	})
}

// UpdateHostInterfaceType updates the type of a host interface in the database.
func (s Store) UpdateHostInterfaceType(ctx context.Context, t string, scope HostInterfaceScope) error {
	return s.db.UpdateHostInterfaceType(ctx, db.UpdateHostInterfaceTypeParams{
		Zone:      scope.Zone,
		Cluster:   scope.Cluster,
		Host:      scope.Host,
		Interface: scope.Interface,
		Type:      Optional(t),
	})
}

// UpdateHostInterfaceBondMode updates the bond mode of a host interface in the database.
func (s Store) UpdateHostInterfaceBondMode(ctx context.Context, mode string, scope HostInterfaceScope) error {
	return s.db.UpdateHostInterfaceBondMode(ctx, db.UpdateHostInterfaceBondModeParams{
		Zone:      scope.Zone,
		Cluster:   scope.Cluster,
		Host:      scope.Host,
		Interface: scope.Interface,
		Bond:      Optional(mode),
	})
}

// UpdateHostInterfaceMaster updates the master of a host interface in the database.
func (s Store) UpdateHostInterfaceMaster(ctx context.Context, master string, scope HostInterfaceScope) error {
	return s.db.UpdateHostInterfaceMaster(ctx, db.UpdateHostInterfaceMasterParams{
		Zone:      scope.Zone,
		Cluster:   scope.Cluster,
		Host:      scope.Host,
		Interface: scope.Interface,
		Master:    master,
	})
}

// UpdateHostInterfaceNetwork updates the network of a host interface in the database.
func (s Store) UpdateHostInterfaceNetwork(ctx context.Context, network string, scope HostInterfaceScope) error {
	return s.db.UpdateHostInterfaceNetwork(ctx, db.UpdateHostInterfaceNetworkParams{
		Zone:      scope.Zone,
		Cluster:   scope.Cluster,
		Host:      scope.Host,
		Interface: scope.Interface,
		Network:   network,
	})
}

// DeleteHosts deletes hosts from the database.
func (s Store) DeleteHosts(ctx context.Context, filter Filter, scope ClusterScope) error {
	rows, err := s.ReadHosts(ctx, filter, scope)
	if err != nil {
		return err
	}

	for i := range rows {
		err := s.db.DeleteHost(ctx, db.DeleteHostParams{ID: rows[i].id})
		if err != nil {
			return err
		}
	}

	return nil
}

// DeleteHostAttrs deletes host attributes from the database.
func (s Store) DeleteHostAttrs(ctx context.Context, filter Filter, scope HostScope) error {
	rows, err := s.ReadHostAttrs(ctx, filter, scope)
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

// DeleteHostInterfaces deletes host interfaceibutes from the database.
func (s Store) DeleteHostInterfaces(ctx context.Context, filter Filter, scope HostScope) error {
	rows, err := s.ReadHostInterfaces(ctx, filter, scope)
	if err != nil {
		return err
	}

	for i := range rows {
		err := s.db.DeleteHostInterface(ctx, db.DeleteHostInterfaceParams{ID: rows[i].id})
		if err != nil {
			return err
		}
	}

	return nil
}
