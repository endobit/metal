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
// Network
//

// CreateNetwork implements the grpc StackServiceServer interface.
func (s Service) CreateNetwork(ctx context.Context, in *pb.CreateNetworkRequest) (*emptypb.Empty, error) {
	switch {
	case !in.HasZone():
		return nil, errMissingZone
	case !in.HasName():
		return nil, errMissingNetwork
	}

	exists, err := s.store.CreateNetwork(ctx, in.GetName(), data.ZoneScope{Zone: in.GetZone()})
	if err != nil {
		return nil, dberr(err)
	}

	if exists {
		return nil, status.Errorf(codes.AlreadyExists, "zone %q network %q", in.GetZone(), in.GetName())
	}

	return new(emptypb.Empty), nil
}

// ReadNetworks implements the grpc StackServiceServer interface.
func (s Service) ReadNetworks(in *pb.ReadNetworksRequest, out grpc.ServerStreamingServer[pb.ReadNetworksResponse]) error {
	var filter data.Filter

	ctx := context.Background()

	switch {
	case in.HasName():
		filter = data.FilterByName(in.GetName())
	case in.HasGlob():
		filter = data.FilterByGlob(in.GetGlob())
	}

	rows, err := s.store.ReadNetworks(ctx, filter, data.ZoneScope{Zone: in.GetZone()})
	if err != nil {
		return dberr(err)
	}

	for i := range rows {
		resp := pb.ReadNetworksResponse_builder{
			Zone: &rows[i].Zone,
			Name: &rows[i].Name,
		}.Build()

		if err := out.Send(resp); err != nil {
			return err
		}
	}

	return nil
}

// UpdateNetwork implements the grpc StackServiceServer interface.
func (s Service) UpdateNetwork(ctx context.Context, in *pb.UpdateNetworkRequest) (*emptypb.Empty, error) {
	switch {
	case !in.HasZone():
		return nil, errMissingZone
	case !in.HasName():
		return nil, errMissingNetwork
	}

	fields := in.GetFields()

	scope := data.NetworkScope{
		Zone:    in.GetZone(),
		Network: in.GetName(),
	}

	if fields.HasName() {
		if err := s.store.UpdateNetworkName(ctx, fields.GetName(), scope); err != nil {
			return nil, dberr(err)
		}
	}

	return new(emptypb.Empty), nil
}

// DeleteNetworks implements the grpc StackServiceServer interface.
func (s Service) DeleteNetworks(ctx context.Context, in *pb.DeleteNetworksRequest) (*emptypb.Empty, error) {
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

	err := s.store.DeleteNetworks(ctx, filter, data.ZoneScope{Zone: in.GetZone()})
	if err != nil {
		return nil, dberr(err)
	}

	return new(emptypb.Empty), nil
}
