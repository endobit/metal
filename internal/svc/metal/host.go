package metal

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	pb "endobit.io/metal/gen/go/proto/metal/v1"
	"endobit.io/metal/internal/data"
)

//
// Host
//

// CreateHost implements the grpc StackServiceServer interface.
func (s Service) CreateHost(ctx context.Context, in *pb.CreateHostRequest) (*emptypb.Empty, error) {
	switch {
	case !in.HasZone():
		return nil, errMissingZone
	case !in.HasName():
		return nil, errMissingHost
	}

	exists, err := s.store.CreateHost(ctx, in.GetName(), data.ClusterScope{
		Zone:    in.GetZone(),
		Cluster: in.GetCluster(),
	})
	if err != nil {
		return nil, dberr(err)
	}

	if exists {
		return nil, status.Errorf(codes.AlreadyExists, "zone %q cluster %q host %q",
			in.GetZone(), in.GetCluster(), in.GetName())
	}

	return new(emptypb.Empty), nil
}

// ReadHosts implements the grpc StackServiceServer interface.
func (s Service) ReadHosts(in *pb.ReadHostsRequest, out grpc.ServerStreamingServer[pb.ReadHostsResponse]) error {
	var filter data.Filter

	ctx := context.Background()

	switch {
	case in.HasName():
		filter = data.FilterByName(in.GetName())
	case in.HasGlob():
		filter = data.FilterByGlob(in.GetGlob())
	}

	rows, err := s.store.ReadHosts(ctx, filter, data.ClusterScope{
		Zone:    in.GetZone(),
		Cluster: in.GetCluster(),
	})
	if err != nil {
		return err
	}

	for i := range rows {
		resp := pb.ReadHostsResponse_builder{
			Zone:        &rows[i].Zone,
			Cluster:     &rows[i].Cluster,
			Name:        &rows[i].Name,
			Make:        &rows[i].Make,
			Model:       &rows[i].Model,
			Environment: &rows[i].Environment,
			Appliance:   &rows[i].Appliance,
			Rack:        &rows[i].Rack,
			Location:    &rows[i].Location,
		}.Build()

		if rank, ok := rows[i].Rank(); ok {
			resp.SetRank(rank)
		}

		if slot, ok := rows[i].Slot(); ok {
			resp.SetSlot(slot)
		}

		if rows[i].Type != "" {
			resp.SetType(pb.HostType(pb.HostType_value[rows[i].Type]))
		}

		if err := out.Send(resp); err != nil {
			return err
		}
	}

	return nil
}

// UpdateHost implements the grpc StackServiceServer interface.
func (s Service) UpdateHost(ctx context.Context, in *pb.UpdateHostRequest) (*emptypb.Empty, error) {
	switch {
	case !in.HasZone():
		return nil, errMissingZone
	case !in.HasName():
		return nil, errMissingHost
	}

	scope := data.HostScope{
		Zone:    in.GetZone(),
		Cluster: in.GetCluster(),
		Host:    in.GetName(),
	}

	switch {
	case in.HasSet():
		fields := in.GetSet()

		switch {
		case fields.HasMake() && !fields.HasModel():
			return nil, errMissingModel
		case fields.HasModel() && !fields.HasMake():
			return nil, errMissingMake
		}

		if fields.HasName() {
			err := s.store.UpdateHostName(ctx, fields.GetName(), scope)
			if err != nil {
				return nil, dberr(err)
			}
		}

		if fields.HasModel() {
			_, _ = s.store.CreateModel(ctx, fields.GetModel(), data.MakeScope{Make: fields.GetMake()}) // try to create if doesn't exists

			err := s.store.UpdateHostModel(ctx, fields.GetMake(), fields.GetModel(), scope)
			if err != nil {
				return nil, dberr(err)
			}
		}

		if fields.HasEnvironment() {
			_, _ = s.store.CreateEnvironment(ctx, fields.GetEnvironment(), data.ZoneScope{Zone: in.GetZone()}) // try to create if doesn't exists

			err := s.store.UpdateHostEnvironment(ctx, fields.GetEnvironment(), scope)
			if err != nil {
				return nil, dberr(err)
			}
		}

		if fields.HasAppliance() {
			_, _ = s.store.CreateAppliance(ctx, fields.GetAppliance(), data.ZoneScope{Zone: in.GetZone()}) // try to create if doesn't exists

			err := s.store.UpdateHostAppliance(ctx, fields.GetAppliance(), scope)
			if err != nil {
				return nil, dberr(err)
			}
		}

		if fields.HasLocation() {
			err := s.store.UpdateHostLocation(ctx, fields.GetLocation(), scope)
			if err != nil {
				return nil, dberr(err)
			}
		}

		if fields.HasRack() {
			_, _ = s.store.CreateRack(ctx, fields.GetRack(), data.ZoneScope{Zone: in.GetZone()}) // try to create if doesn't exists

			err := s.store.UpdateHostRack(ctx, fields.GetRack(), scope)
			if err != nil {
				return nil, dberr(err)
			}
		}

		if fields.HasRank() {
			err := s.store.UpdateHostRank(ctx, int64(fields.GetRank()), scope)
			if err != nil {
				return nil, dberr(err)
			}
		}

		if fields.HasSlot() {
			err := s.store.UpdateHostSlot(ctx, int64(fields.GetSlot()), scope)
			if err != nil {
				return nil, dberr(err)
			}
		}

		if fields.HasType() {
			t := pb.HostType_name[int32(fields.GetType())]
			if err := s.store.UpdateHostType(ctx, t, scope); err != nil {
				return nil, dberr(err)
			}
		}

	case in.HasUnset():
		fields := in.GetUnset()

		if fields.HasModel() || fields.HasMake() {
			if err := s.store.UpdateHostModel(ctx, "", "", scope); err != nil {
				return nil, dberr(err)
			}
		}

		if fields.HasEnvironment() {
			if err := s.store.UpdateHostEnvironment(ctx, "", scope); err != nil {
				return nil, dberr(err)
			}
		}

		if fields.HasAppliance() {
			if err := s.store.UpdateHostAppliance(ctx, "", scope); err != nil {
				return nil, dberr(err)
			}
		}

		if fields.HasLocation() {
			if err := s.store.UpdateHostLocation(ctx, "", scope); err != nil {
				return nil, dberr(err)
			}
		}

		if fields.HasRack() {
			if err := s.store.UpdateHostRack(ctx, "", scope); err != nil {
				return nil, dberr(err)
			}
		}

		if fields.HasRank() {
			if err := s.store.UpdateHostRank(ctx, -1, scope); err != nil {
				return nil, dberr(err)
			}
		}

		if fields.HasSlot() {
			if err := s.store.UpdateHostSlot(ctx, -1, scope); err != nil {
				return nil, dberr(err)
			}
		}
	}

	return new(emptypb.Empty), nil
}

// DeleteHosts implements the grpc StackServiceServer interface.
func (s Service) DeleteHosts(ctx context.Context, in *pb.DeleteHostsRequest) (*emptypb.Empty, error) {
	var filter data.Filter

	switch {
	case !in.HasZone():
		return nil, errMissingZone
	case in.HasName():
		filter = data.FilterByName(in.GetName())
	case in.HasGlob():
		filter = data.FilterByGlob(in.GetGlob())
	default:
		return nil, errNameOrGlob
	}

	err := s.store.DeleteHosts(ctx, filter, data.ClusterScope{
		Zone:    in.GetZone(),
		Cluster: in.GetCluster(),
	})
	if err != nil {
		return nil, dberr(err)
	}

	return new(emptypb.Empty), nil
}

//
// Host Attributes
//

// CreateHostAttr implements the grpc StackServiceServer interface.
func (s *Service) CreateHostAttr(ctx context.Context, in *pb.CreateHostAttrRequest) (*emptypb.Empty, error) {
	switch {
	case !in.HasZone():
		return nil, errMissingZone
	case !in.HasHost():
		return nil, errMissingHost
	case !in.HasName():
		return nil, errMissingAttribute
	}

	exists, err := s.store.CreateHostAttr(ctx, in.GetName(), data.HostScope{
		Zone:    in.GetZone(),
		Cluster: in.GetCluster(),
		Host:    in.GetHost(),
	})

	if err != nil {
		return nil, dberr(err)
	}

	if exists {
		return nil, status.Errorf(codes.AlreadyExists, "zone %q cluster %q host %q attribute %q",
			in.GetZone(), in.GetCluster(), in.GetHost(), in.GetName())
	}

	return new(emptypb.Empty), nil
}

// ReadHostAttrs implements the grpc StackServiceServer interface.
func (s *Service) ReadHostAttrs(in *pb.ReadHostAttrsRequest, out grpc.ServerStreamingServer[pb.ReadHostAttrsResponse]) error {
	var filter data.Filter

	ctx := context.Background()

	switch {
	case in.HasName():
		filter = data.FilterByName(in.GetName())
	case in.HasGlob():
		filter = data.FilterByGlob(in.GetGlob())
	}

	rows, err := s.store.ReadHostAttrs(ctx, filter, data.HostScope{
		Zone:    in.GetZone(),
		Cluster: in.GetCluster(),
		Host:    in.GetHost(),
	})
	if err != nil {
		return err
	}

	for i := range rows {
		resp := pb.ReadHostAttrsResponse_builder{
			Zone:      &rows[i].Zone,
			Cluster:   &rows[i].Cluster,
			Host:      &rows[i].Host,
			Name:      &rows[i].Name,
			Value:     &rows[i].Value,
			Protected: &rows[i].IsProtected,
		}.Build()

		if err := out.Send(resp); err != nil {
			return err
		}
	}

	return nil
}

// UpdateHostAttr implements the grpc StackServiceServer interface.
func (s *Service) UpdateHostAttr(ctx context.Context, in *pb.UpdateHostAttrRequest) (*emptypb.Empty, error) {
	switch {
	case !in.HasZone():
		return nil, errMissingZone
	case !in.HasHost():
		return nil, errMissingHost
	case !in.HasName():
		return nil, errMissingAttribute
	}

	fields := in.GetFields()

	scope := data.HostAttrScope{
		Zone:    in.GetZone(),
		Cluster: in.GetCluster(),
		Host:    in.GetHost(),
		Attr:    in.GetName(),
	}

	if fields.HasName() {
		if err := s.store.UpdateHostAttrName(ctx, fields.GetName(), scope); err != nil {
			return nil, dberr(err)
		}
	}

	if fields.HasValue() {
		if err := s.store.UpdateHostAttrValue(ctx, fields.GetValue(), scope); err != nil {
			return nil, dberr(err)
		}
	}

	if fields.HasProtected() {
		if err := s.store.UpdateHostAttrProtection(ctx, fields.GetProtected(), scope); err != nil {
			return nil, dberr(err)
		}
	}

	return new(emptypb.Empty), nil
}

// DeleteHostAttrs implements the grpc StackServiceServer interface.
func (s *Service) DeleteHostAttrs(ctx context.Context, in *pb.DeleteHostAttrsRequest) (*emptypb.Empty, error) {
	var filter data.Filter

	switch {
	case !in.HasZone():
		return nil, errMissingZone
	case !in.HasHost():
		return nil, errMissingHost
	case in.HasName():
		filter = data.FilterByName(in.GetName())
	case in.HasGlob():
		filter = data.FilterByGlob(in.GetGlob())
	default:
		return nil, errNameOrGlob
	}

	err := s.store.DeleteHostAttrs(ctx, filter, data.HostScope{
		Zone:    in.GetZone(),
		Cluster: in.GetCluster(),
		Host:    in.GetHost(),
	})

	if err != nil {
		return nil, dberr(err)
	}

	return new(emptypb.Empty), nil
}

//
// Host Interfaces
//

// CreateHostInterface implements the grpc StackServiceServer interface.
func (s Service) CreateHostInterface(ctx context.Context, in *pb.CreateHostInterfaceRequest) (*emptypb.Empty, error) {
	switch {
	case !in.HasZone():
		return nil, errMissingZone
	case !in.HasHost():
		return nil, errMissingHost
	case !in.HasName():
		return nil, errMissingHost
	}

	exists, err := s.store.CreateHostInterface(ctx, in.GetZone(), data.HostScope{
		Zone:    in.GetZone(),
		Cluster: in.GetCluster(),
		Host:    in.GetHost(),
	})
	if err != nil {
		return nil, dberr(err)
	}

	if exists {
		return nil, status.Errorf(codes.AlreadyExists, "zone %q cluster %q host %q interface %q",
			in.GetZone(), in.GetCluster(), in.GetHost(), in.GetName())
	}

	return new(emptypb.Empty), nil
}

// ReadHostInterfaces implements the grpc StackServiceServer interface.
func (s *Service) ReadHostInterfaces(in *pb.ReadHostInterfacesRequest, out grpc.ServerStreamingServer[pb.ReadHostInterfacesResponse]) error {
	var filter data.Filter

	ctx := context.Background()

	switch {
	case in.HasName():
		filter = data.FilterByName(in.GetName())
	case in.HasGlob():
		filter = data.FilterByGlob(in.GetGlob())
	}

	rows, err := s.store.ReadHostInterfaces(ctx, filter, data.HostScope{
		Zone:    in.GetZone(),
		Cluster: in.GetCluster(),
		Host:    in.GetHost(),
	})
	if err != nil {
		return dberr(err)
	}

	for i := range rows {
		resp := pb.ReadHostInterfacesResponse_builder{
			Zone:       &rows[i].Zone,
			Cluster:    &rows[i].Cluster,
			Host:       &rows[i].Host,
			Name:       &rows[i].Name,
			Ip:         &rows[i].IP,
			Mac:        &rows[i].MAC,
			Netmask:    &rows[i].Netmask,
			Dhcp:       &rows[i].IsDHCP,
			Pxe:        &rows[i].IsPXE,
			Management: &rows[i].IsManagement,
			Master:     &rows[i].MasterInterface,
			Network:    &rows[i].Network,
		}.Build()

		if rows[i].Type != "" {
			resp.SetType(pb.InterfaceType(pb.InterfaceType_value[rows[i].Type]))
		}

		if rows[i].BondMode != "" {
			resp.SetBondMode(pb.BondMode(pb.BondMode_value[rows[i].BondMode]))
		}

		if err := out.Send(resp); err != nil {
			return err
		}
	}

	return nil
}

//
// Update
//

// UpdateHostInterface implements the grpc StackServiceServer interface.
func (s *Service) UpdateHostInterface(ctx context.Context, in *pb.UpdateHostInterfaceRequest) (*emptypb.Empty, error) {
	switch {
	case !in.HasZone():
		return nil, errMissingZone
	case !in.HasHost():
		return nil, errMissingHost
	case !in.HasName():
		return nil, errMissingInterface
	}

	fields := in.GetFields()

	scope := data.HostInterfaceScope{
		Zone:      in.GetZone(),
		Cluster:   in.GetCluster(),
		Host:      in.GetHost(),
		Interface: in.GetName(),
	}

	if fields.HasName() {
		if err := s.store.UpdateHostInterfaceName(ctx, fields.GetName(), scope); err != nil {
			return nil, dberr(err)
		}
	}

	if fields.HasIp() {
		if err := s.store.UpdateHostInterfaceIP(ctx, fields.GetIp(), scope); err != nil {
			return nil, dberr(err)
		}
	}

	if fields.HasMac() {
		if err := s.store.UpdateHostInterfaceMAC(ctx, fields.GetMac(), scope); err != nil {
			return nil, dberr(err)
		}
	}

	if fields.HasNetmask() {
		if err := s.store.UpdateHostInterfaceNetmask(ctx, fields.GetNetmask(), scope); err != nil {
			return nil, dberr(err)
		}
	}

	if fields.HasDhcp() {
		if err := s.store.UpdateHostInterfaceDHCP(ctx, fields.GetDhcp(), scope); err != nil {
			return nil, dberr(err)
		}
	}

	if fields.HasPxe() {
		if err := s.store.UpdateHostInterfacePXE(ctx, fields.GetPxe(), scope); err != nil {
			return nil, dberr(err)
		}
	}

	if fields.HasManagement() {
		if err := s.store.UpdateHostInterfaceManagement(ctx, fields.GetManagement(), scope); err != nil {
			return nil, dberr(err)
		}
	}

	if fields.HasType() {
		if err := s.store.UpdateHostInterfaceType(ctx, fields.GetType().String(), scope); err != nil {
			return nil, dberr(err)
		}
	}

	if fields.HasBondMode() {
		if err := s.store.UpdateHostInterfaceBondMode(ctx, fields.GetBondMode().String(), scope); err != nil {
			return nil, dberr(err)
		}
	}

	if fields.HasMaster() {
		if err := s.store.UpdateHostInterfaceMaster(ctx, fields.GetMaster(), scope); err != nil {
			return nil, dberr(err)
		}
	}

	if fields.HasNetwork() {
		if err := s.store.UpdateHostInterfaceNetwork(ctx, fields.GetNetwork(), scope); err != nil {
			return nil, dberr(err)
		}
	}

	return new(emptypb.Empty), nil
}

// DeleteHostInterfaces implements the grpc StackServiceServer interface.
func (s *Service) DeleteHostInterfaces(ctx context.Context, in *pb.DeleteHostInterfacesRequest) (*emptypb.Empty, error) {
	var filter data.Filter

	switch {
	case !in.HasZone():
		return nil, errMissingZone
	case !in.HasHost():
		return nil, errMissingHost
	case in.HasName():
		filter = data.FilterByName(in.GetName())
	case in.HasGlob():
		filter = data.FilterByGlob(in.GetGlob())
	default:
		return nil, errNameOrGlob
	}

	err := s.store.DeleteHostInterfaces(ctx, filter, data.HostScope{
		Zone:    in.GetZone(),
		Cluster: in.GetCluster(),
		Host:    in.GetHost(),
	})

	if err != nil {
		return nil, dberr(err)
	}

	return new(emptypb.Empty), nil
}
