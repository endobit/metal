package metal

import (
	"errors"
	"io"
	"iter"

	"google.golang.org/grpc"

	pb "endobit.io/metal/gen/go/proto/metal/v1"
)

//
// Global
//

// NewGlobalAttrReader creates a new stream reader for global attributes.
func (c *Client) NewGlobalAttrReader(attr string) *streamReader[pb.ReadGlobalAttrsResponse] {
	var req pb.ReadGlobalAttrsRequest

	if attr != "" {
		req.SetGlob(attr)
	}

	return newStreamReader(c.ReadGlobalAttrs(c.Context(), &req))
}

//
// Makes and Models
//

// NewMakeReader creates a new stream reader for makes.
func (c *Client) NewMakeReader(mk string) *streamReader[pb.ReadMakesResponse] {
	var req pb.ReadMakesRequest

	if mk != "" {
		req.SetGlob(mk)
	}

	return newStreamReader(c.ReadMakes(c.Context(), &req))
}

// NewModelReader creates a new stream reader for models.
func (c *Client) NewModelReader(mk, model string) *streamReader[pb.ReadModelsResponse] {
	var req pb.ReadModelsRequest

	if mk != "" {
		req.SetMake(mk)
	}
	if model != "" {
		req.SetGlob(model)
	}

	return newStreamReader(c.ReadModels(c.Context(), &req))
}

// NewModelAttrReader creates a new stream reader for model attributes.
func (c *Client) NewModelAttrReader(model, attr string) *streamReader[pb.ReadModelAttrsResponse] {
	var req pb.ReadModelAttrsRequest

	if model != "" {
		req.SetModel(model)
	}
	if attr != "" {
		req.SetGlob(attr)
	}

	return newStreamReader(c.ReadModelAttrs(c.Context(), &req))
}

//
// Zones
//

// NewZoneReader creates a new stream reader for zones.
func (c *Client) NewZoneReader(zone string) *streamReader[pb.ReadZonesResponse] {
	var req pb.ReadZonesRequest

	if zone != "" {
		req.SetGlob(zone)
	}

	return newStreamReader(c.ReadZones(c.Context(), &req))
}

// NewZoneAttrReader creates a new stream reader for zone attributes.
func (c *Client) NewZoneAttrReader(zone, attr string) *streamReader[pb.ReadZoneAttrsResponse] {
	var req pb.ReadZoneAttrsRequest

	if zone != "" {
		req.SetZone(zone)
	}
	if attr != "" {
		req.SetGlob(attr)
	}

	return newStreamReader(c.ReadZoneAttrs(c.Context(), &req))
}

//
// Appliances
//

// NewApplianceReader creates a new stream reader for appliances.
func (c *Client) NewApplianceReader(zone, appliance string) *streamReader[pb.ReadAppliancesResponse] {
	var req pb.ReadAppliancesRequest

	if zone != "" {
		req.SetZone(zone)
	}
	if appliance != "" {
		req.SetGlob(appliance)
	}

	return newStreamReader(c.ReadAppliances(c.Context(), &req))
}

// NewApplianceAttrReader creates a new stream reader for appliance attributes.
func (c *Client) NewApplianceAttrReader(zone, appliance, attr string) *streamReader[pb.ReadApplianceAttrsResponse] {
	var req pb.ReadApplianceAttrsRequest

	if zone != "" {
		req.SetZone(zone)
	}
	if appliance != "" {
		req.SetAppliance(appliance)
	}
	if attr != "" {
		req.SetGlob(attr)
	}

	return newStreamReader(c.ReadApplianceAttrs(c.Context(), &req))
}

//
// Environments
//

// NewEnvironmentReader creates a new stream reader for environments.
func (c *Client) NewEnvironmentReader(zone, environment string) *streamReader[pb.ReadEnvironmentsResponse] {
	var req pb.ReadEnvironmentsRequest

	if zone != "" {
		req.SetZone(zone)
	}
	if environment != "" {
		req.SetGlob(environment)
	}

	return newStreamReader(c.ReadEnvironments(c.Context(), &req))
}

// NewEnvironmentAttrReader creates a new stream reader for environment attributes.
func (c *Client) NewEnvironmentAttrReader(zone, environment, attr string) *streamReader[pb.ReadEnvironmentAttrsResponse] {
	var req pb.ReadEnvironmentAttrsRequest

	if zone != "" {
		req.SetZone(zone)
	}
	if environment != "" {
		req.SetEnvironment(environment)
	}
	if attr != "" {
		req.SetGlob(attr)
	}

	return newStreamReader(c.ReadEnvironmentAttrs(c.Context(), &req))
}

//
// Networks
//

// NewNetworkReader creates a new stream reader for networks.
func (c *Client) NewNetworkReader(zone, network string) *streamReader[pb.ReadNetworksResponse] {
	var req pb.ReadNetworksRequest

	if zone != "" {
		req.SetZone(zone)
	}
	if network != "" {
		req.SetGlob(network)
	}

	return newStreamReader(c.ReadNetworks(c.Context(), &req))
}

//
// Racks
//

// NewRackReader creates a new stream reader for racks.
func (c *Client) NewRackReader(zone, rack string) *streamReader[pb.ReadRacksResponse] {
	var req pb.ReadRacksRequest

	if zone != "" {
		req.SetZone(zone)
	}
	if rack != "" {
		req.SetGlob(rack)
	}

	return newStreamReader(c.ReadRacks(c.Context(), &req))
}

// NewRackAttrReader creates a new stream reader for rack attributes.
func (c *Client) NewRackAttrReader(zone, rack, attr string) *streamReader[pb.ReadRackAttrsResponse] {
	var req pb.ReadRackAttrsRequest

	if zone != "" {
		req.SetZone(zone)
	}
	if rack != "" {
		req.SetRack(rack)
	}
	if attr != "" {
		req.SetGlob(attr)
	}

	return newStreamReader(c.ReadRackAttrs(c.Context(), &req))
}

//
// Clusters
//

// NewClusterReader creates a new stream reader for clusters.
func (c *Client) NewClusterReader(zone, cluster string) *streamReader[pb.ReadClustersResponse] {
	var req pb.ReadClustersRequest

	if zone != "" {
		req.SetZone(zone)
	}
	if cluster != "" {
		req.SetGlob(cluster)
	}

	return newStreamReader(c.ReadClusters(c.Context(), &req))
}

// NewClusterAttrReader creates a new stream reader for cluster attributes.
func (c *Client) NewClusterAttrReader(zone, cluster, attr string) *streamReader[pb.ReadClusterAttrsResponse] {
	var req pb.ReadClusterAttrsRequest

	if zone != "" {
		req.SetZone(zone)
	}
	if cluster != "" {
		req.SetCluster(cluster)
	}
	if attr != "" {
		req.SetGlob(attr)
	}

	return newStreamReader(c.ReadClusterAttrs(c.Context(), &req))
}

//
// Hosts
//

// NewHostReader creates a new stream reader for hosts.
func (c *Client) NewHostReader(zone, cluster, host string) *streamReader[pb.ReadHostsResponse] {
	var req pb.ReadHostsRequest

	if zone != "" {
		req.SetZone(zone)
	}
	if cluster != "" {
		req.SetCluster(cluster)
	}
	if host != "" {
		req.SetGlob(host)
	}

	return newStreamReader(c.ReadHosts(c.Context(), &req))
}

// NewHostAttrReader creates a new stream reader for host attributes.
func (c *Client) NewHostAttrReader(zone, cluster, host, attr string) *streamReader[pb.ReadHostAttrsResponse] {
	var req pb.ReadHostAttrsRequest

	if zone != "" {
		req.SetZone(zone)
	}
	if cluster != "" {
		req.SetCluster(cluster)
	}
	if host != "" {
		req.SetHost(host)
	}
	if attr != "" {
		req.SetGlob(attr)
	}

	return newStreamReader(c.ReadHostAttrs(c.Context(), &req))
}

// NewHostInterfaceReader creates a new stream reader for host interfaces.
func (c *Client) NewHostInterfaceReader(zone, cluster, host, iface string) *streamReader[pb.ReadHostInterfacesResponse] {
	var req pb.ReadHostInterfacesRequest

	if zone != "" {
		req.SetZone(zone)
	}
	if cluster != "" {
		req.SetCluster(cluster)
	}
	if host != "" {
		req.SetHost(host)
	}
	if iface != "" {
		req.SetGlob(iface)
	}

	return newStreamReader(c.ReadHostInterfaces(c.Context(), &req))
}

// streamReader is a gRPC streaming response Reader. It's only purpose is to provide
// an iterator for streaming responses. See newStreamReader for more details.
type streamReader[T any] struct {
	grpc.ServerStreamingClient[T]
	err error
}

// newStreamReader returns a new Reader. It takes as parameters the return values from
// a gRPC streaming client method. For example:
//
//	r := newStreamReader(client.ReadZones(ctx, req))
//	for zone, err := range r.Responses() {
//	...
//	}
func newStreamReader[T any](client grpc.ServerStreamingClient[T], err error) *streamReader[T] {
	return &streamReader[T]{
		ServerStreamingClient: client,
		err:                   err,
	}
}

// Responses returns an iterator for processing gRPC streaming responses from the Reader.
func (r *streamReader[T]) Responses() iter.Seq2[*T, error] {
	return func(yield func(*T, error) bool) {
		// constructors are allowed to set and not return an error
		if r.err != nil {
			if !yield(nil, r.err) {
				return
			}
		}

		for {
			resp, err := r.Recv()
			if errors.Is(err, io.EOF) {
				return
			}

			if !yield(resp, err) {
				return
			}
		}
	}
}
