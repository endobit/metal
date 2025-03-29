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
// Rack
//

// CreateRack implements the grpc StackServiceServer interface.
func (s Service) CreateRack(ctx context.Context, in *pb.CreateRackRequest) (*emptypb.Empty, error) {
	switch {
	case !in.HasZone():
		return nil, errMissingZone
	case !in.HasName():
		return nil, errMissingRack
	}

	exists, err := s.store.CreateRack(ctx, in.GetName(), data.ZoneScope{Zone: in.GetZone()})
	if err != nil {
		return nil, dberr(err)
	}

	if exists {
		return nil, status.Errorf(codes.AlreadyExists, "zone %q rack %q", in.GetZone(), in.GetName())
	}

	return new(emptypb.Empty), nil
}

// ReadRacks implements the grpc StackServiceServer interface.
func (s Service) ReadRacks(in *pb.ReadRacksRequest, out grpc.ServerStreamingServer[pb.ReadRacksResponse]) error {
	var filter data.Filter

	ctx := context.Background()

	switch {
	case in.HasName():
		filter = data.FilterByName(in.GetName())
	case in.HasGlob():
		filter = data.FilterByGlob(in.GetGlob())
	}

	rows, err := s.store.ReadRacks(ctx, filter, data.ZoneScope{Zone: in.GetZone()})
	if err != nil {
		return dberr(err)
	}

	for i := range rows {
		resp := pb.ReadRacksResponse_builder{
			Zone: &rows[i].Zone,
			Name: &rows[i].Name,
		}.Build()

		if err := out.Send(resp); err != nil {
			return err
		}
	}

	return nil
}

// UpdateRack implements the grpc StackServiceServer interface.
func (s Service) UpdateRack(ctx context.Context, in *pb.UpdateRackRequest) (*emptypb.Empty, error) {
	switch {
	case !in.HasZone():
		return nil, errMissingZone
	case !in.HasName():
		return nil, errMissingRack
	}

	fields := in.GetFields()

	scope := data.RackScope{
		Zone: in.GetZone(),
		Rack: in.GetName(),
	}

	if fields.HasName() {
		if err := s.store.UpdateRackName(ctx, fields.GetName(), scope); err != nil {
			return nil, dberr(err)
		}
	}

	return new(emptypb.Empty), nil
}

// DeleteRacks implements the grpc StackServiceServer interface.
func (s Service) DeleteRacks(ctx context.Context, in *pb.DeleteRacksRequest) (*emptypb.Empty, error) {
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

	err := s.store.DeleteRacks(ctx, filter, data.ZoneScope{Zone: in.GetZone()})
	if err != nil {
		return nil, dberr(err)
	}

	return new(emptypb.Empty), nil
}

//
// Rack Attributes
//

// CreateRackAttr implements the grpc StackServiceServer interface.
func (s *Service) CreateRackAttr(ctx context.Context, in *pb.CreateRackAttrRequest) (*emptypb.Empty, error) {
	switch {
	case !in.HasZone():
		return nil, errMissingZone
	case !in.HasRack():
		return nil, errMissingRack
	case !in.HasName():
		return nil, errMissingAttribute
	}

	exists, err := s.store.CreateRackAttr(ctx, in.GetName(), data.RackScope{
		Zone: in.GetZone(),
		Rack: in.GetRack(),
	})
	if err != nil {
		return nil, dberr(err)
	}

	if exists {
		return nil, status.Errorf(codes.AlreadyExists, "zone %q rack %q attr %q",
			in.GetZone(), in.GetName(), in.GetName())
	}

	return new(emptypb.Empty), nil
}

// ReadRackAttrs implements the grpc StackServiceServer interface.
func (s *Service) ReadRackAttrs(in *pb.ReadRackAttrsRequest, out grpc.ServerStreamingServer[pb.ReadRackAttrsResponse]) error {
	var filter data.Filter

	ctx := context.Background()

	switch {
	case in.HasName():
		filter = data.FilterByName(in.GetName())
	case in.HasGlob():
		filter = data.FilterByGlob(in.GetGlob())
	}

	rows, err := s.store.ReadRackAttrs(ctx, filter, data.RackScope{
		Zone: in.GetZone(),
		Rack: in.GetRack(),
	})
	if err != nil {
		return dberr(err)
	}

	for i := range rows {
		resp := pb.ReadRackAttrsResponse_builder{
			Rack:      &rows[i].Rack,
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

// UpdateRackAttr implements the grpc StackServiceServer interface.
func (s *Service) UpdateRackAttr(ctx context.Context, in *pb.UpdateRackAttrRequest) (*emptypb.Empty, error) {
	switch {
	case !in.HasZone():
		return nil, errMissingZone
	case !in.HasRack():
		return nil, errMissingRack
	case !in.HasName():
		return nil, errMissingAttribute
	}

	fields := in.GetFields()

	scope := data.RackAttrScope{
		Zone: in.GetZone(),
		Rack: in.GetRack(),
		Attr: in.GetName(),
	}

	if fields.HasName() {
		err := s.store.UpdateRackAttrName(ctx, fields.GetName(), scope)
		if err != nil {
			return nil, dberr(err)
		}
	}

	if fields.HasValue() {
		err := s.store.UpdateRackAttrValue(ctx, fields.GetValue(), scope)
		if err != nil {
			return nil, dberr(err)
		}
	}

	if fields.HasProtected() {
		err := s.store.UpdateRackAttrProtection(ctx, fields.GetProtected(), scope)
		if err != nil {
			return nil, dberr(err)
		}
	}

	return new(emptypb.Empty), nil
}

// DeleteRackAttrs implements the grpc StackServiceServer interface.
func (s *Service) DeleteRackAttrs(ctx context.Context, in *pb.DeleteRackAttrsRequest) (*emptypb.Empty, error) {
	var filter data.Filter

	switch {
	case !in.HasZone():
		return nil, errMissingZone
	case !in.HasRack():
		return nil, errMissingRack
	case in.HasName():
		filter = data.FilterByName(in.GetName())
	case in.HasGlob():
		filter = data.FilterByGlob(in.GetGlob())
	default:
		return nil, errNameOrGlob
	}

	err := s.store.DeleteRackAttrs(ctx, filter, data.RackScope{
		Zone: in.GetZone(),
		Rack: in.GetRack(),
	})
	if err != nil {
		return nil, dberr(err)
	}

	return new(emptypb.Empty), nil
}
