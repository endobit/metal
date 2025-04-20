package metal

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"google.golang.org/protobuf/types/known/emptypb"

	"endobit.io/metal"
	pb "endobit.io/metal/gen/go/proto/metal/v1"
	"endobit.io/metal/internal/data"
)

// CreateSchema implements the grpc StackServiceServer interface.
func (s Service) CreateSchema(ctx context.Context, in *pb.CreateSchemaRequest) (*emptypb.Empty, error) {
	l := loader{
		store:  s.store,
		logger: s.logger,
	}

	if err := l.Load(ctx, in.GetSchema()); err != nil {
		return nil, err
	}

	return new(emptypb.Empty), nil
}

// ReadSchema implements the grpc StackServiceServer interface. It uses the
// metal.Dumper to dump a metal.Report and then loads into the protobuf schema.
// The metal.Report is unique to the Go code only, and the protobuf format is
// cross-platform.
func (s Service) ReadSchema(ctx context.Context, in *pb.ReadSchemaRequest) (*pb.ReadSchemaResponse, error) {
	var resp pb.ReadSchemaResponse

	d := metal.Dumper{Store: s.store}

	if in.HasZone() {
		d.Filter.Zone = data.FilterByName(in.GetZone())
	}

	if in.HasCluster() {
		d.Filter.Cluster = data.FilterByName(in.GetCluster())
	}

	if in.HasHost() {
		d.Filter.Host = data.FilterByName(in.GetHost())
	}

	doc, err := d.Dump(ctx)
	if err != nil {
		return nil, err
	}

	resp.SetSchema(newSchema(doc))

	return &resp, nil
}

func newMake(src *metal.Make) *pb.Make {
	models := make([]*pb.Model, 0, len(src.Models))

	for i := range src.Models {
		model := &src.Models[i]
		attrs := make(map[string]string)

		for key, value := range model.Attrs {
			attrs[key] = value.Value
		}

		models = append(models, &pb.Model{
			Name:         &model.Name,
			Architecture: data.Optional(model.Architecture),
			Attributes:   attrs,
		})
	}

	return &pb.Make{
		Name:   &src.Name,
		Models: models,
	}
}

func newCluster(src *metal.Cluster) *pb.Cluster {
	attrs := make(map[string]string)
	hosts := make([]*pb.Host, 0, len(src.Hosts))

	for key, value := range src.Attrs {
		attrs[key] = value.Value
	}

	for i := range src.Hosts {
		hosts = append(hosts, newHost(&src.Hosts[i]))
	}

	return &pb.Cluster{
		Name:       &src.Name,
		Attributes: attrs,
		Hosts:      hosts,
	}
}

func newHost(src *metal.Host) *pb.Host {
	attrs := make(map[string]string)
	ifaces := make([]*pb.Host_Interface, 0, len(src.Interfaces))

	for key, value := range src.Attrs {
		attrs[key] = value.Value
	}

	for i := range src.Interfaces {
		ifaces = append(ifaces, newHostInterface(&src.Interfaces[i]))
	}

	host := pb.Host{
		Attributes:  attrs,
		Name:        &src.Name,
		Make:        data.Optional(src.Make),
		Model:       data.Optional(src.Model),
		Environment: data.Optional(src.Environment),
		Location:    data.Optional(src.Location),
		Rack:        data.Optional(src.Rack),
		Appliance:   data.Optional(src.Appliance),
		Type:        data.Optional(src.Type),
		Interfaces:  ifaces,
	}

	if rank, ok := src.Host.Rank(); ok {
		host.Rank = &rank
	}

	if slot, ok := src.Host.Slot(); ok {
		host.Slot = &slot
	}

	return &host
}

func newAppliance(src *metal.Appliance) *pb.Appliance {
	attrs := make(map[string]string)

	for key, value := range src.Attrs {
		attrs[key] = value.Value
	}

	return &pb.Appliance{
		Name:       &src.Name,
		Attributes: attrs,
	}
}

func newEnvironment(src *metal.Environment) *pb.Environment {
	attrs := make(map[string]string)

	for key, value := range src.Attrs {
		attrs[key] = value.Value
	}

	return &pb.Environment{
		Name:       &src.Name,
		Attributes: attrs,
	}
}

func newRack(src *metal.Rack) *pb.Rack {
	attrs := make(map[string]string)

	for key, value := range src.Attrs {
		attrs[key] = value.Value
	}

	return &pb.Rack{
		Name:       &src.Name,
		Attributes: attrs,
	}
}

func newNetwork(src *metal.Network) *pb.Network {
	return &pb.Network{
		Name:    &src.Name,
		Address: data.Optional(src.CIDR),
		Gateway: data.Optional(src.Gateway),
		Pxe:     data.Optional(src.IsPXE),
		Mtu:     data.Optional(src.MTU),
	}
}

func newHostInterface(src *metal.HostInterface) *pb.Host_Interface {
	iface := pb.Host_Interface{
		Name:       &src.Name,
		Ip:         data.Optional(src.IP),
		Mac:        data.Optional(src.MAC),
		Dhcp:       data.Optional(src.IsDHCP),
		Pxe:        data.Optional(src.IsPXE),
		Management: data.Optional(src.IsManagement),
		Type:       data.Optional(src.Type),
		BondMode:   data.Optional(src.BondMode),
	}

	if src.Network != nil {
		iface.Network = data.Optional(src.Network.CIDR)
		iface.Gateway = data.Optional(src.Network.Gateway)
		iface.Mtu = data.Optional(src.Network.MTU)
	}

	return &iface
}

//
// dump
//

func newSchema(report *metal.Report) *pb.Schema {
	attrs := make(map[string]string)
	makes := make([]*pb.Make, 0, len(report.Makes))
	zones := make([]*pb.Zone, 0, len(report.Zones))

	for key, value := range report.Attrs {
		attrs[key] = value.Value
	}

	for _, value := range report.Makes {
		makes = append(makes, newMake(&value))
	}

	for i := range report.Zones {
		zone := &report.Zones[i]

		attrs := make(map[string]string)
		networks := make([]*pb.Network, 0, len(zone.Networks))
		appliances := make([]*pb.Appliance, 0, len(zone.Appliances))
		environments := make([]*pb.Environment, 0, len(zone.Environments))
		racks := make([]*pb.Rack, 0, len(zone.Racks))
		hosts := make([]*pb.Host, 0, len(zone.Hosts))
		clusters := make([]*pb.Cluster, 0, len(zone.Clusters))

		for key, value := range zone.Attrs {
			attrs[key] = value.Value
		}

		for j := range zone.Networks {
			networks = append(networks, newNetwork(&zone.Networks[j]))
		}

		for j := range zone.Appliances {
			appliances = append(appliances, newAppliance(&zone.Appliances[j]))
		}

		for j := range zone.Environments {
			environments = append(environments, newEnvironment(&zone.Environments[j]))
		}

		for j := range zone.Racks {
			racks = append(racks, newRack(&zone.Racks[j]))
		}

		for j := range zone.Hosts {
			hosts = append(hosts, newHost(&zone.Hosts[j]))
		}

		for j := range zone.Clusters {
			clusters = append(clusters, newCluster(&zone.Clusters[j]))
		}

		zones = append(zones, &pb.Zone{
			Name:         &zone.Name,
			TimeZone:     data.Optional(zone.TimeZone),
			Attributes:   attrs,
			Networks:     networks,
			Appliances:   appliances,
			Environments: environments,
			Racks:        racks,
			Hosts:        hosts,
			Clusters:     clusters,
		})

	}

	return &pb.Schema{
		Attributes: attrs,
		Makes:      makes,
		Zones:      zones,
	}
}

//
// loading
//

type loader struct {
	store  *data.Store
	logger *slog.Logger
}

func (l loader) Load(ctx context.Context, doc *pb.Schema) error {
	if err := l.LoadAttrs(ctx, doc.Attributes); err != nil {
		return err
	}
	if err := l.LoadMakes(ctx, doc.Makes); err != nil {
		return err
	}
	if err := l.LoadZones(ctx, doc.Zones); err != nil {
		return err
	}

	return nil
}

func (l loader) LoadMakes(ctx context.Context, makes []*pb.Make) error {
	for _, make := range makes {
		if make.Name == nil {
			return errors.New("make name is required")
		}

		if _, err := l.store.CreateMake(ctx, *make.Name); err != nil {
			return err
		}

		for _, model := range make.Models {
			if model.Name == nil {
				return errors.New("model name is required")
			}

			_, err := l.store.CreateModel(ctx, *model.Name, data.MakeScope{Make: *make.Name})
			if err != nil {
				return err
			}

			if model.Architecture != nil {
				arch, ok := metal.LongArchitecture[*model.Architecture]
				if !ok {
					return fmt.Errorf("%q: %w", *model.Architecture, errInvalidArchitecture)
				}
				err := l.store.UpdateModelArchitecture(ctx, arch, data.ModelScope{
					Make:  *make.Name,
					Model: *model.Name,
				})
				if err != nil {
					return err
				}
			}

			for key, value := range model.Attributes {
				_, err := l.store.CreateModelAttr(ctx, key, data.ModelScope{
					Make:  *make.Name,
					Model: *model.Name,
				})
				if err != nil {
					return err
				}
				err = l.store.UpdateModelAttrValue(ctx, value, data.ModelAttrScope{
					Make:  *make.Name,
					Model: *model.Name,
					Attr:  key,
				})
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (l loader) LoadAttrs(ctx context.Context, attrs map[string]string) error {
	for key, value := range attrs {
		if _, err := l.store.CreateGlobalAttr(ctx, key); err != nil {
			return err
		}

		if err := l.store.UpdateGlobalAttrValue(ctx, value, data.GlobalAttrScope{Attr: key}); err != nil {
			return err
		}
	}

	return nil
}

func (l loader) LoadZones(ctx context.Context, zones []*pb.Zone) error {
	for _, zone := range zones {
		if zone.Name == nil {
			return errors.New("zone name is required")
		}

		if _, err := l.store.CreateZone(ctx, *zone.Name); err != nil {
			return err
		}

		scope := data.ZoneScope{Zone: *zone.Name}

		if zone.TimeZone != nil {
			if err := l.store.UpdateZoneTimeZone(ctx, *zone.TimeZone, scope); err != nil {
				return err
			}
		}

		if err := l.LoadZoneAttrs(ctx, zone.Attributes, scope); err != nil {
			return err
		}

		if err := l.LoadNetworks(ctx, zone.Networks, scope); err != nil {
			return err
		}

		if err := l.LoadRacks(ctx, zone.Racks, scope); err != nil {
			return err
		}

		if err := l.LoadAppliances(ctx, zone.Appliances, scope); err != nil {
			return err
		}

		if err := l.LoadEnvironments(ctx, zone.Environments, scope); err != nil {
			return err
		}

		if err := l.LoadClusters(ctx, zone.Clusters, scope); err != nil {
			return err
		}

		if err := l.LoadHosts(ctx, zone.Hosts, data.ClusterScope{Zone: scope.Zone}); err != nil {
			return err
		}
	}

	return nil
}

func (l loader) LoadZoneAttrs(ctx context.Context, attrs map[string]string, scope data.ZoneScope) error {
	for key, value := range attrs {
		if _, err := l.store.CreateZoneAttr(ctx, key, scope); err != nil {
			return err
		}

		err := l.store.UpdateZoneAttrValue(ctx, value, data.ZoneAttrScope{
			Zone: scope.Zone,
			Attr: key,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (l loader) LoadNetworks(ctx context.Context, networks []*pb.Network, scope data.ZoneScope) error {
	for _, network := range networks {
		if network.Name == nil {
			return errors.New("network name is required")
		}

		if _, err := l.store.CreateNetwork(ctx, *network.Name, scope); err != nil {
			return err
		}

		scope := data.NetworkScope{
			Zone:    scope.Zone,
			Network: *network.Name,
		}

		if network.Address != nil {
			err := l.store.UpdateNetworkAddress(ctx, *network.Address, data.NetworkScope{
				Zone:    scope.Zone,
				Network: *network.Name,
			})
			if err != nil {
				return err
			}
		}

		if network.Gateway != nil {
			err := l.store.UpdateNetworkGateway(ctx, *network.Gateway, data.NetworkScope{
				Zone:    scope.Zone,
				Network: *network.Name,
			})
			if err != nil {
				return err
			}
		}

		if network.Pxe != nil {
			err := l.store.UpdateNetworkPXE(ctx, *network.Pxe, data.NetworkScope{
				Zone:    scope.Zone,
				Network: *network.Name,
			})
			if err != nil {
				return err
			}
		}

		if network.Mtu != nil {
			err := l.store.UpdateNetworkMTU(ctx, *network.Mtu, data.NetworkScope{
				Zone:    scope.Zone,
				Network: *network.Name,
			})
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (l loader) LoadRacks(ctx context.Context, racks []*pb.Rack, scope data.ZoneScope) error {
	for _, rack := range racks {
		if rack.Name == nil {
			return errors.New("rack name is required")
		}

		if _, err := l.store.CreateRack(ctx, *rack.Name, scope); err != nil {
			return err
		}

		scope := data.RackScope{
			Zone: scope.Zone,
			Rack: *rack.Name,
		}

		if err := l.LoadRackAttrs(ctx, rack.Attributes, scope); err != nil {
			return err
		}
	}

	return nil
}

func (l loader) LoadRackAttrs(ctx context.Context, attrs map[string]string, scope data.RackScope) error {
	for key, value := range attrs {
		if _, err := l.store.CreateRackAttr(ctx, key, scope); err != nil {
			return err
		}

		err := l.store.UpdateRackAttrValue(ctx, value, data.RackAttrScope{
			Zone: scope.Zone,
			Rack: scope.Rack,
			Attr: key,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (l loader) LoadAppliances(ctx context.Context, appliances []*pb.Appliance, scope data.ZoneScope) error {
	for _, appliance := range appliances {
		if appliance.Name == nil {
			return errors.New("appliance name is required")
		}

		if _, err := l.store.CreateAppliance(ctx, appliance.GetName(), scope); err != nil {
			return err
		}

		scope := data.ApplianceScope{
			Zone:      scope.Zone,
			Appliance: *appliance.Name,
		}

		if err := l.LoadApplianceAttrs(ctx, appliance.Attributes, scope); err != nil {
			return err
		}
	}

	return nil
}

func (l loader) LoadApplianceAttrs(ctx context.Context, attrs map[string]string, scope data.ApplianceScope) error {
	for key, value := range attrs {
		if _, err := l.store.CreateApplianceAttr(ctx, key, scope); err != nil {
			return err
		}

		err := l.store.UpdateApplianceAttrValue(ctx, value, data.ApplianceAttrScope{
			Zone:      scope.Zone,
			Appliance: scope.Appliance,
			Attr:      key,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (l loader) LoadEnvironments(ctx context.Context, environments []*pb.Environment, scope data.ZoneScope) error {
	for _, environment := range environments {
		if environment.Name == nil {
			return errors.New("environment name is required")
		}

		if _, err := l.store.CreateEnvironment(ctx, *environment.Name, scope); err != nil {
			return err
		}

		scope := data.EnvironmentScope{
			Zone:        scope.Zone,
			Environment: *environment.Name,
		}

		if err := l.LoadEnvironmentAttrs(ctx, environment.Attributes, scope); err != nil {
			return err
		}
	}

	return nil
}

func (l loader) LoadEnvironmentAttrs(ctx context.Context, attrs map[string]string, scope data.EnvironmentScope) error {
	for key, value := range attrs {
		if _, err := l.store.CreateEnvironmentAttr(ctx, key, scope); err != nil {
			return err
		}

		err := l.store.UpdateEnvironmentAttrValue(ctx, value, data.EnvironmentAttrScope{
			Zone:        scope.Zone,
			Environment: scope.Environment,
			Attr:        key,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (l loader) LoadClusters(ctx context.Context, clusters []*pb.Cluster, scope data.ZoneScope) error {
	for _, cluster := range clusters {
		if cluster.Name == nil {
			return errors.New("cluster name is required")
		}

		if _, err := l.store.CreateCluster(ctx, *cluster.Name, scope); err != nil {
			return err
		}

		scope := data.ClusterScope{
			Zone:    scope.Zone,
			Cluster: *cluster.Name,
		}

		if err := l.LoadClusterAttrs(ctx, cluster.Attributes, scope); err != nil {
			return err
		}

		if err := l.LoadHosts(ctx, cluster.Hosts, scope); err != nil {
			return err
		}
	}

	return nil
}

func (l loader) LoadClusterAttrs(ctx context.Context, attrs map[string]string, scope data.ClusterScope) error {
	for key, value := range attrs {
		if _, err := l.store.CreateClusterAttr(ctx, key, scope); err != nil {
			return err
		}

		err := l.store.UpdateClusterAttrValue(ctx, value, data.ClusterAttrScope{
			Zone:    scope.Zone,
			Cluster: scope.Cluster,
			Attr:    key,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (l loader) LoadHosts(ctx context.Context, hosts []*pb.Host, scope data.ClusterScope) error {
	for _, host := range hosts {
		if host.Name == nil {
			return errors.New("host name is required")
		}

		if _, err := l.store.CreateHost(ctx, *host.Name, scope); err != nil {
			return err
		}

		scope := data.HostScope{
			Zone:    scope.Zone,
			Cluster: scope.Cluster,
			Host:    *host.Name,
		}

		if err := l.LoadHostAttrs(ctx, host.Attributes, scope); err != nil {
			return err
		}

		if err := l.LoadHostInterfaces(ctx, host.Interfaces, scope); err != nil {
			return err
		}

		if host.Appliance != nil {
			if err := l.store.UpdateHostAppliance(ctx, *host.Appliance, scope); err != nil {
				return err
			}
		}

		if host.Environment != nil {
			if err := l.store.UpdateHostEnvironment(ctx, *host.Environment, scope); err != nil {
				return err
			}
		}

		if host.Location != nil {
			if err := l.store.UpdateHostLocation(ctx, *host.Location, scope); err != nil {
				return err
			}
		}

		if host.Make != nil && host.Model != nil {
			if err := l.store.UpdateHostModel(ctx, *host.Make, *host.Model, scope); err != nil {
				return err
			}
		}

		if host.Rack != nil {
			if err := l.store.UpdateHostRack(ctx, *host.Rack, scope); err != nil {
				return err
			}
		}

		if host.Rank != nil {
			rank := int64(*host.Rank)
			if err := l.store.UpdateHostRank(ctx, rank, scope); err != nil {
				return err
			}
		}

		if host.Slot != nil {
			slot := int64(*host.Slot)
			if err := l.store.UpdateHostSlot(ctx, slot, scope); err != nil {
				return err
			}
		}

		if host.Type != nil {
			hostType, ok := metal.LongHostType[*host.Type]
			if !ok {
				return fmt.Errorf("%q: %w", *host.Type, errInvalidHostType)
			}

			if err := l.store.UpdateHostType(ctx, hostType, scope); err != nil {
				return err
			}
		}

	}

	return nil
}

func (l loader) LoadHostAttrs(ctx context.Context, attrs map[string]string, scope data.HostScope) error {
	for key, value := range attrs {
		if _, err := l.store.CreateHostAttr(ctx, key, scope); err != nil {
			return err
		}

		err := l.store.UpdateHostAttrValue(ctx, value, data.HostAttrScope{
			Zone:    scope.Zone,
			Cluster: scope.Cluster,
			Host:    scope.Host,
			Attr:    key,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (l loader) LoadHostInterfaces(ctx context.Context, ifaces []*pb.Host_Interface, scope data.HostScope) error {
	for _, iface := range ifaces {
		if iface.Name == nil {
			return errors.New("host interface name is required")
		}

		if _, err := l.store.CreateHostInterface(ctx, *iface.Name, scope); err != nil {
			return err
		}

		scope := data.HostInterfaceScope{
			Zone:      scope.Zone,
			Cluster:   scope.Cluster,
			Host:      scope.Host,
			Interface: *iface.Name,
		}

		if iface.Network != nil {
			if err := l.store.UpdateHostInterfaceNetwork(ctx, *iface.Network, scope); err != nil {
				return err
			}
		}

		if iface.Ip != nil {
			if err := l.store.UpdateHostInterfaceIP(ctx, *iface.Ip, scope); err != nil {
				return err
			}
		}

		if iface.Mac != nil {
			if err := l.store.UpdateHostInterfaceMAC(ctx, *iface.Mac, scope); err != nil {
				return err
			}
		}

		if iface.Dhcp != nil {
			if err := l.store.UpdateHostInterfaceDHCP(ctx, *iface.Dhcp, scope); err != nil {
				return err
			}
		}

		if iface.Pxe != nil {
			if err := l.store.UpdateHostInterfacePXE(ctx, *iface.Pxe, scope); err != nil {
				return err
			}
		}

		if iface.Management != nil {
			if err := l.store.UpdateHostInterfaceManagement(ctx, *iface.Management, scope); err != nil {
				return err
			}
		}

		if iface.Type != nil {
			ifaceType, ok := metal.LongInterfaceType[*iface.Type]
			if !ok {
				return fmt.Errorf("%q: %w", *iface.Type, errInvalidInterfaceType)
			}
			if err := l.store.UpdateHostInterfaceType(ctx, ifaceType, scope); err != nil {
				return err
			}
		}

		if iface.BondMode != nil {
			mode, ok := metal.LongBondMode[*iface.BondMode]
			if !ok {
				return fmt.Errorf("%q: %w", *iface.BondMode, errInvalidBondMode)
			}
			if err := l.store.UpdateHostInterfaceBondMode(ctx, mode, scope); err != nil {
				return err
			}
		}
	}

	return nil
}
