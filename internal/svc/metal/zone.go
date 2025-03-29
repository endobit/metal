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
// Zone
//

// CreateZone implements the grpc StackServiceServer interface.
func (s Service) CreateZone(ctx context.Context, in *pb.CreateZoneRequest) (*emptypb.Empty, error) {
	if !in.HasName() {
		return nil, errMissingZone
	}

	exists, err := s.store.CreateZone(ctx, in.GetName())
	if err != nil {
		return nil, dberr(err)
	}

	if exists {
		return nil, status.Errorf(codes.AlreadyExists, "zone %q", in.GetName())
	}

	return new(emptypb.Empty), nil
}

// ReadZones implements the grpc StackServiceServer interface.
func (s Service) ReadZones(in *pb.ReadZonesRequest, out grpc.ServerStreamingServer[pb.ReadZonesResponse]) error {
	var filter data.Filter

	ctx := context.Background()

	switch {
	case in.HasName():
		filter = data.FilterByName(in.GetName())
	case in.HasGlob():
		filter = data.FilterByGlob(in.GetGlob())
	}

	rows, err := s.store.ReadZones(ctx, filter)
	if err != nil {
		return dberr(err)
	}

	for i := range rows {
		resp := pb.ReadZonesResponse_builder{
			Name:     &rows[i].Name,
			TimeZone: &rows[i].TimeZone,
		}.Build()

		if err := out.Send(resp); err != nil {
			return err
		}
	}

	return nil
}

// UpdateZone implements the grpc StackServiceServer interface.
func (s Service) UpdateZone(ctx context.Context, in *pb.UpdateZoneRequest) (*emptypb.Empty, error) {
	fields := in.GetFields()

	scope := data.ZoneScope{Zone: in.GetName()}

	if fields.HasName() {
		if err := s.store.UpdateZoneName(ctx, fields.GetName(), scope); err != nil {
			return nil, dberr(err)
		}
	}

	if fields.HasTimeZone() {
		if err := s.store.UpdateZoneTimeZone(ctx, fields.GetTimeZone(), scope); err != nil {
			return nil, dberr(err)
		}

	}

	return new(emptypb.Empty), nil
}

// DeleteZones implements the grpc StackServiceServer interface.
func (s Service) DeleteZones(ctx context.Context, in *pb.DeleteZonesRequest) (*emptypb.Empty, error) {
	var filter data.Filter

	switch {
	case in.HasName():
		filter = data.FilterByName(in.GetName())
	case in.HasGlob():
		filter = data.FilterByGlob(in.GetGlob())
	default:
		return nil, errNameOrGlob
	}

	err := s.store.DeleteZones(ctx, filter)
	if err != nil {
		return nil, dberr(err)
	}

	return new(emptypb.Empty), nil
}

//
// Zone Attributes
//

// CreateZoneAttr implements the grpc StackServiceServer interface.
func (s *Service) CreateZoneAttr(ctx context.Context, in *pb.CreateZoneAttrRequest) (*emptypb.Empty, error) {
	switch {
	case !in.HasZone():
		return nil, errMissingZone
	case !in.HasName():
		return nil, errMissingAttribute
	}

	exists, err := s.store.CreateZoneAttr(ctx, in.GetName(), data.ZoneScope{Zone: in.GetZone()})
	if err != nil {
		return nil, dberr(err)
	}

	if exists {
		return nil, status.Errorf(codes.AlreadyExists, "zone %q attribute %q", in.GetZone(), in.GetName())
	}

	return new(emptypb.Empty), nil
}

// ReadZoneAttrs implements the grpc StackServiceServer interface.
func (s *Service) ReadZoneAttrs(in *pb.ReadZoneAttrsRequest, out grpc.ServerStreamingServer[pb.ReadZoneAttrsResponse]) error {
	var filter data.Filter

	ctx := context.Background()

	switch {
	case in.HasName():
		filter = data.FilterByName(in.GetName())
	case in.HasGlob():
		filter = data.FilterByGlob(in.GetGlob())
	}

	rows, err := s.store.ReadZoneAttrs(ctx, filter, data.ZoneScope{Zone: in.GetZone()})
	if err != nil {
		return err
	}

	for i := range rows {
		resp := pb.ReadZoneAttrsResponse_builder{
			Zone:      &rows[i].Zone,
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

// UpdateZoneAttr implements the grpc StackServiceServer interface.
func (s *Service) UpdateZoneAttr(ctx context.Context, in *pb.UpdateZoneAttrRequest) (*emptypb.Empty, error) {
	fields := in.GetFields()

	scope := data.ZoneAttrScope{
		Zone: in.GetZone(),
		Attr: in.GetName(),
	}

	if fields.HasName() {
		if err := s.store.UpdateZoneAttrName(ctx, fields.GetName(), scope); err != nil {
			return nil, dberr(err)
		}
	}

	if fields.HasValue() {
		if err := s.store.UpdateZoneAttrValue(ctx, fields.GetValue(), scope); err != nil {
			return nil, dberr(err)
		}
	}

	if fields.HasProtected() {
		if err := s.store.UpdateZoneAttrProtection(ctx, fields.GetProtected(), scope); err != nil {
			return nil, dberr(err)
		}
	}

	return new(emptypb.Empty), nil
}

// DeleteZoneAttrs implements the grpc StackServiceServer interface.
func (s *Service) DeleteZoneAttrs(ctx context.Context, in *pb.DeleteZoneAttrsRequest) (*emptypb.Empty, error) {
	var filter data.Filter

	switch {
	case in.HasName():
		filter = data.FilterByName(in.GetName())
	case in.HasGlob():
		filter = data.FilterByGlob(in.GetGlob())
	default:
		return nil, errNameOrGlob
	}

	err := s.store.DeleteZoneAttrs(ctx, filter, data.ZoneScope{Zone: in.GetZone()})
	if err != nil {
		return nil, dberr(err)
	}

	return new(emptypb.Empty), nil
}
