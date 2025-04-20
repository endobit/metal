package metal

import (
	"context"
	"fmt"
	"maps"

	"endobit.io/metal/internal/data"
)

//go:generate go tool "github.com/dmarkham/enumer" -type AttrType -linecomment -text

type AttrType int

const (
	GlobalAttr      AttrType = iota // G
	ModelAttr                       // M
	ZoneAttr                        // Z
	ApplianceAttr                   // A
	RackAttr                        // R
	ClusterAttr                     // C
	EnvironmentAttr                 // E
	HostAttr                        // H
	SwitchAttr                      // S
	BMCAttr                         // B
)

type (
	// Report is the data returned by the ReadReportData gRPC metal call.
	// This model is used by the metal-ops template engine to render reports.
	Report struct {
		Makes []Make
		Zones []Zone
		Attrs map[string]Attr
	}

	Attr struct {
		Type  AttrType
		Value string
	}

	Make struct {
		*data.Make
		Models []Model
	}

	Model struct {
		*data.Model
		Attrs map[string]Attr
	}

	Zone struct {
		*data.Zone
		Attrs        map[string]Attr
		Networks     []Network
		Appliances   []Appliance
		Environments []Environment
		Racks        []Rack
		Hosts        []Host
		Clusters     []Cluster
	}

	Appliance struct {
		*data.Appliance
		Attrs map[string]Attr
	}

	Environment struct {
		*data.Environment
		Attrs map[string]Attr
	}

	Network struct {
		*data.Network
	}

	Rack struct {
		*data.Rack
		Attrs map[string]Attr
	}

	Cluster struct {
		*data.Cluster
		Attrs map[string]Attr
		Hosts []Host
	}

	Host struct {
		*data.Host
		Rank         string
		Slot         string
		Architecture string
		Attrs        map[string]Attr
		Interfaces   []HostInterface
	}

	HostInterface struct {
		*data.HostInterface
		Network *data.Network
	}
)

type Dumper struct {
	Store        *data.Store
	Filter       Filters
	ResolveAttrs bool
	attrCache    attrs
	networkCache map[data.ZoneScope]map[string]*data.Network
}

type attrs struct {
	Model       map[data.ModelScope]map[string]Attr
	Appliance   map[data.ApplianceScope]map[string]Attr
	Environment map[data.EnvironmentScope]map[string]Attr
	Rack        map[data.RackScope]map[string]Attr
}

type Filters struct {
	Zone    data.Filter
	Cluster data.Filter
	Host    data.Filter
}

func (d *Dumper) Dump(ctx context.Context) (*Report, error) {
	d.networkCache = make(map[data.ZoneScope]map[string]*data.Network)

	attrs, err := d.GlobalAttrs(ctx)
	if err != nil {
		return nil, err
	}
	makes, err := d.Makes(ctx, attrs)
	if err != nil {
		return nil, err
	}
	zones, err := d.Zones(ctx, attrs)
	if err != nil {
		return nil, err
	}

	report := Report{
		Attrs: attrs,
		Makes: makes,
		Zones: zones,
	}

	return &report, nil
}

func (d *Dumper) GlobalAttrs(ctx context.Context) (map[string]Attr, error) {
	attrs := make(map[string]Attr)

	rows, err := d.Store.ReadGlobalAttrs(ctx, data.Filter{})
	if err != nil {
		return nil, err
	}

	for _, row := range rows {
		attrs[row.Name] = Attr{Type: GlobalAttr, Value: row.Value}
	}

	return attrs, nil
}

func (d *Dumper) ModelAttrs(ctx context.Context, scope data.ModelScope) (map[string]Attr, error) {
	attrs := make(map[string]Attr)

	rows, err := d.Store.ReadModelAttrs(ctx, data.Filter{}, scope)
	if err != nil {
		return nil, err
	}

	for _, row := range rows {
		attrs[row.Name] = Attr{Type: ModelAttr, Value: row.Value}
	}

	return attrs, nil
}

func (d *Dumper) ZoneAttrs(ctx context.Context, scope data.ZoneScope) (map[string]Attr, error) {
	attrs := make(map[string]Attr)

	rows, err := d.Store.ReadZoneAttrs(ctx, data.Filter{}, scope)
	if err != nil {
		return nil, err
	}

	for _, row := range rows {
		attrs[row.Name] = Attr{Type: ZoneAttr, Value: row.Value}
	}

	return attrs, nil
}

func (d *Dumper) ApplianceAttrs(ctx context.Context, scope data.ApplianceScope) (map[string]Attr, error) {
	attrs := make(map[string]Attr)

	rows, err := d.Store.ReadApplianceAttrs(ctx, data.Filter{}, scope)
	if err != nil {
		return nil, err
	}

	for _, row := range rows {
		attrs[row.Name] = Attr{Type: ApplianceAttr, Value: row.Value}
	}

	return attrs, nil
}

func (d *Dumper) EnvironmentAttrs(ctx context.Context, scope data.EnvironmentScope) (map[string]Attr, error) {
	attrs := make(map[string]Attr)

	rows, err := d.Store.ReadEnvironmentAttrs(ctx, data.Filter{}, scope)
	if err != nil {
		return nil, err
	}

	for _, row := range rows {
		attrs[row.Name] = Attr{Type: EnvironmentAttr, Value: row.Value}
	}

	return attrs, nil
}

func (d *Dumper) RackAttrs(ctx context.Context, scope data.RackScope) (map[string]Attr, error) {
	attrs := make(map[string]Attr)

	rows, err := d.Store.ReadRackAttrs(ctx, data.Filter{}, scope)
	if err != nil {
		return nil, err
	}

	for _, row := range rows {
		attrs[row.Name] = Attr{Type: RackAttr, Value: row.Value}
	}

	return attrs, nil
}

func (d *Dumper) ClusterAttrs(ctx context.Context, scope data.ClusterScope) (map[string]Attr, error) {
	attrs := make(map[string]Attr)

	rows, err := d.Store.ReadClusterAttrs(ctx, data.Filter{}, scope)
	if err != nil {
		return nil, err
	}

	for _, row := range rows {
		attrs[row.Name] = Attr{Type: ClusterAttr, Value: row.Value}
	}

	return attrs, nil
}

func (d *Dumper) HostAttrs(ctx context.Context, scope data.HostScope) (map[string]Attr, error) {
	attrs := make(map[string]Attr)

	rows, err := d.Store.ReadHostAttrs(ctx, data.Filter{}, scope)
	if err != nil {
		return nil, err
	}

	for _, row := range rows {
		attrs[row.Name] = Attr{Type: HostAttr, Value: row.Value}
	}

	return attrs, nil
}

func (d *Dumper) Makes(ctx context.Context, global map[string]Attr) ([]Make, error) {
	var makes []Make

	rows, err := d.Store.ReadMakes(ctx, data.Filter{})
	if err != nil {
		return nil, err
	}

	for i := range rows {
		models, err := d.Models(ctx, global, data.MakeScope{Make: rows[i].Name})
		if err != nil {
			return nil, err
		}

		makes = append(makes, Make{
			Make:   &rows[i],
			Models: models,
		})
	}

	return makes, nil
}

func (d *Dumper) Models(ctx context.Context, globalAttrs map[string]Attr, scope data.MakeScope) ([]Model, error) {
	var models []Model

	rows, err := d.Store.ReadModels(ctx, data.Filter{}, scope)
	if err != nil {
		return nil, err
	}

	for i := range rows {
		model := &rows[i]

		if model.Architecture != "" {
			model.Architecture = ShortArchitecture[model.Architecture]
		}

		scope := data.ModelScope{Make: scope.Make, Model: model.Name}
		attrs, err := d.ModelAttrs(ctx, scope)
		if err != nil {
			return nil, err
		}

		// This is the pattern for resolving attributes at all levels. The
		// owner's (e.g. global, zone, cluster) resolved attributes are passed
		// into the Dumper method (here the owner is global). The owner's
		// attributes are cloned and local attributes overwrite the map.
		//
		// For level that are cache the unresolved attributes are cached.
		// Otherwise when resolving multiple copies of attributes could happen.
		// For example if the appliance attribute cache was resolved it could
		// include zone attribute. This would slow down host attribute
		// resolving.

		if d.ResolveAttrs {
			// Global
			// Model

			d.attrCache.Model[scope] = maps.Clone(attrs)

			resolved := maps.Clone(globalAttrs)
			maps.Copy(resolved, attrs)

			attrs = resolved
		}

		models = append(models, Model{
			Model: model,
			Attrs: attrs,
		})
	}

	return models, nil
}

func (d *Dumper) Zones(ctx context.Context, global map[string]Attr) ([]Zone, error) {
	var zones []Zone

	rows, err := d.Store.ReadZones(ctx, d.Filter.Zone)
	if err != nil {
		return nil, err
	}

	for i := range rows {
		zone := &rows[i]

		scope := data.ZoneScope{Zone: zone.Name}

		attrs, err := d.ZoneAttrs(ctx, scope)
		if err != nil {
			return nil, err
		}

		if d.ResolveAttrs {
			resolved := maps.Clone(global)
			maps.Copy(resolved, attrs)

			attrs = resolved
		}

		appliances, err := d.Appliances(ctx, attrs, scope)
		if err != nil {
			return nil, err
		}
		environments, err := d.Environments(ctx, attrs, scope)
		if err != nil {
			return nil, err
		}
		networks, err := d.Networks(ctx, scope)
		if err != nil {
			return nil, err
		}
		racks, err := d.Racks(ctx, attrs, scope)
		if err != nil {
			return nil, err
		}
		hosts, err := d.Hosts(ctx, attrs, data.ClusterScope{Zone: scope.Zone})
		if err != nil {
			return nil, err
		}
		clusters, err := d.Clusters(ctx, attrs, scope)
		if err != nil {
			return nil, err
		}

		zones = append(zones, Zone{
			Zone:         zone,
			Attrs:        attrs,
			Appliances:   appliances,
			Environments: environments,
			Networks:     networks,
			Racks:        racks,
			Hosts:        hosts,
			Clusters:     clusters,
		})
	}

	return zones, nil
}

func (d *Dumper) Appliances(ctx context.Context, zone map[string]Attr, scope data.ZoneScope) ([]Appliance, error) {
	var appliances []Appliance

	rows, err := d.Store.ReadAppliances(ctx, data.Filter{}, scope)
	if err != nil {
		return nil, err
	}

	for i := range rows {
		app := &rows[i]

		scope := data.ApplianceScope{Zone: scope.Zone, Appliance: app.Name}

		attrs, err := d.ApplianceAttrs(ctx, scope)
		if err != nil {
			return nil, err
		}

		if d.ResolveAttrs {
			d.attrCache.Appliance[scope] = maps.Clone(attrs)

			resolved := maps.Clone(zone)
			maps.Copy(resolved, attrs)

			attrs = resolved
		}

		appliances = append(appliances, Appliance{
			Appliance: app,
			Attrs:     attrs,
		})
	}

	return appliances, nil
}

func (d *Dumper) Environments(ctx context.Context, zone map[string]Attr, scope data.ZoneScope) ([]Environment, error) {
	var environments []Environment

	rows, err := d.Store.ReadEnvironments(ctx, data.Filter{}, scope)
	if err != nil {
		return nil, err
	}

	for i := range rows {
		env := &rows[i]

		scope := data.EnvironmentScope{Zone: scope.Zone, Environment: env.Name}

		attrs, err := d.EnvironmentAttrs(ctx, scope)
		if err != nil {
			return nil, err
		}

		if d.ResolveAttrs {
			d.attrCache.Environment[scope] = attrs

			resolved := maps.Clone(zone)
			maps.Copy(resolved, attrs)

			attrs = resolved
		}

		environments = append(environments, Environment{
			Environment: env,
			Attrs:       attrs,
		})
	}

	return environments, nil
}

func (d *Dumper) Networks(ctx context.Context, scope data.ZoneScope) ([]Network, error) {
	var networks []Network

	rows, err := d.Store.ReadNetworks(ctx, data.Filter{}, scope)
	if err != nil {
		return nil, err
	}

	if d.networkCache[scope] == nil {
		d.networkCache[scope] = make(map[string]*data.Network)
	}

	for i := range rows {
		network := &rows[i]

		d.networkCache[scope][network.Name] = network
		networks = append(networks, Network{Network: network})
	}

	return networks, nil
}

func (d *Dumper) Racks(ctx context.Context, zone map[string]Attr, scope data.ZoneScope) ([]Rack, error) {
	var racks []Rack

	rows, err := d.Store.ReadRacks(ctx, data.Filter{}, scope)
	if err != nil {
		return nil, err
	}

	for i := range rows {
		rack := &rows[i]

		scope := data.RackScope{Zone: scope.Zone, Rack: rack.Name}

		attrs, err := d.RackAttrs(ctx, scope)
		if err != nil {
			return nil, err
		}

		if d.ResolveAttrs {
			d.attrCache.Rack[scope] = maps.Clone(attrs)

			resolved := maps.Clone(zone)
			maps.Copy(resolved, attrs)

			attrs = resolved
		}

		racks = append(racks, Rack{
			Rack:  rack,
			Attrs: attrs,
		})
	}

	return racks, nil
}

func (d *Dumper) Clusters(ctx context.Context, zone map[string]Attr, scope data.ZoneScope) ([]Cluster, error) {
	var clusters []Cluster

	rows, err := d.Store.ReadClusters(ctx, d.Filter.Cluster, scope)
	if err != nil {
		return nil, err
	}

	for i := range rows {
		cluster := &rows[i]

		scope := data.ClusterScope{Zone: scope.Zone, Cluster: cluster.Name}

		attrs, err := d.ClusterAttrs(ctx, scope)
		if err != nil {
			return nil, err
		}

		if d.ResolveAttrs {
			resolved := maps.Clone(zone)
			maps.Copy(resolved, attrs)

			attrs = resolved
		}

		hosts, err := d.Hosts(ctx, attrs, scope)
		if err != nil {
			return nil, err
		}

		clusters = append(clusters, Cluster{
			Cluster: cluster,
			Attrs:   attrs,
			Hosts:   hosts,
		})
	}

	return clusters, nil
}

func (d *Dumper) Hosts(ctx context.Context, base map[string]Attr, scope data.ClusterScope) ([]Host, error) {
	var hosts []Host

	rows, err := d.Store.ReadHosts(ctx, d.Filter.Host, scope)
	if err != nil {
		return nil, err
	}

	for i := range rows {
		host := &rows[i]

		if host.Type != "" {
			host.Type = ShortHostType[host.Type]
		}

		// Reading hosts at the zone level will return all hosts in that zone,
		// including those in clusters. So when at the zone level skipped any
		// clustered hosts.
		if scope.Cluster == "" && host.Cluster != "" {
			continue
		}

		scope := data.HostScope{Zone: scope.Zone, Cluster: scope.Cluster, Host: host.Name}

		attrs, err := d.HostAttrs(ctx, scope)
		if err != nil {
			return nil, err
		}

		if d.ResolveAttrs {
			// Global
			// Model
			// Zone
			// Cluster
			// Environment
			// Appliance
			// Rack
			// Host

			// global and model
			resolved := maps.Clone(d.attrCache.Model[data.ModelScope{Make: host.Make, Model: host.Model}])

			maps.Copy(resolved, base) // mix in the zone/cluster attrs

			maps.Copy(resolved, d.attrCache.Environment[data.EnvironmentScope{Environment: host.Environment}])
			maps.Copy(resolved, d.attrCache.Appliance[data.ApplianceScope{Appliance: host.Appliance}])

			// WARNING: This assumes the Zone is set for clustered hosts, which
			// is currently the case. But there are other comments in the
			// schema/code that imply clustered hosts should not set the Zone in
			// the database. If this changes the Store MUST set it.
			maps.Copy(resolved, d.attrCache.Rack[data.RackScope{Zone: host.Zone, Rack: host.Rack}])

			maps.Copy(resolved, attrs) // mix in the host attrs

			attrs = resolved
		}

		ifaces, err := d.HostInterfaces(ctx, scope)
		if err != nil {
			return nil, err
		}

		// additional info about the host

		var arch, rank, slot string

		model, err := d.Store.ReadModel(ctx, host.Model, data.MakeScope{Make: host.Make})
		if err != nil {
			return nil, err
		}
		if model != nil {
			arch = model.Architecture
		}

		if x, ok := host.Rank(); ok {
			slot = fmt.Sprint(x)

		}
		if x, ok := host.Slot(); ok {
			slot = fmt.Sprint(x)
		}

		hosts = append(hosts, Host{
			Host:         host,
			Attrs:        attrs,
			Interfaces:   ifaces,
			Architecture: arch,
			Rank:         rank,
			Slot:         slot,
		})
	}

	return hosts, nil
}

func (d *Dumper) HostInterfaces(ctx context.Context, scope data.HostScope) ([]HostInterface, error) {
	var ifaces []HostInterface

	rows, err := d.Store.ReadHostInterfaces(ctx, data.Filter{}, scope)
	if err != nil {
		return nil, err
	}

	zone := data.ZoneScope{
		Zone: scope.Zone,
	}

	for i := range rows {
		iface := &rows[i]

		if iface.BondMode != "" {
			iface.BondMode = ShortBondMode[iface.BondMode]
		}
		if iface.Type != "" {
			iface.Type = ShortInterfaceType[iface.Type]
		}

		ifaces = append(ifaces, HostInterface{
			HostInterface: iface,
			Network:       d.networkCache[zone][iface.Network],
		})
	}

	return ifaces, nil
}
