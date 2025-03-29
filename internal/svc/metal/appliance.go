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
// Appliance
//

// CreateAppliance implements the grpc StackServiceServer interface.
func (s Service) CreateAppliance(ctx context.Context, in *pb.CreateApplianceRequest) (*emptypb.Empty, error) {
	switch {
	case !in.HasZone():
		return nil, errMissingZone
	case !in.HasName():
		return nil, errMissingAppliance
	}

	exists, err := s.store.CreateAppliance(ctx, in.GetName(), data.ZoneScope{Zone: in.GetZone()})
	if err != nil {
		return nil, dberr(err)
	}

	if exists {
		return nil, status.Errorf(codes.AlreadyExists, "zone %q appliance %q", in.GetZone(), in.GetName())
	}

	return new(emptypb.Empty), nil
}

// ReadAppliances implements the grpc StackServiceServer interface.
func (s Service) ReadAppliances(in *pb.ReadAppliancesRequest, out grpc.ServerStreamingServer[pb.ReadAppliancesResponse]) error {
	var filter data.Filter

	ctx := context.Background()

	switch {
	case in.HasName():
		filter = data.FilterByName(in.GetName())
	case in.HasGlob():
		filter = data.FilterByGlob(in.GetGlob())
	}

	rows, err := s.store.ReadAppliances(ctx, filter, data.ZoneScope{Zone: in.GetZone()})
	if err != nil {
		return dberr(err)
	}

	for i := range rows {
		resp := pb.ReadAppliancesResponse_builder{
			Zone: &rows[i].Zone,
			Name: &rows[i].Name,
		}.Build()

		if err := out.Send(resp); err != nil {
			return err
		}
	}

	return nil
}

// UpdateAppliance implements the grpc StackServiceServer interface.
func (s Service) UpdateAppliance(ctx context.Context, in *pb.UpdateApplianceRequest) (*emptypb.Empty, error) {
	switch {
	case !in.HasZone():
		return nil, errMissingZone
	case !in.HasName():
		return nil, errMissingAppliance
	}

	fields := in.GetFields()

	scope := data.ApplianceScope{
		Zone:      in.GetZone(),
		Appliance: in.GetName(),
	}

	if fields.HasName() {
		if err := s.store.UpdateApplianceName(ctx, fields.GetName(), scope); err != nil {
			return nil, dberr(err)
		}
	}

	return new(emptypb.Empty), nil
}

// DeleteAppliances implements the grpc StackServiceServer interface.
func (s Service) DeleteAppliances(ctx context.Context, in *pb.DeleteAppliancesRequest) (*emptypb.Empty, error) {
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

	err := s.store.DeleteAppliances(ctx, filter, data.ZoneScope{Zone: in.GetZone()})
	if err != nil {
		return nil, dberr(err)
	}

	return new(emptypb.Empty), nil
}

//
// Appliance Attributes
//

// CreateApplianceAttr implements the grpc StackServiceServer interface.
func (s *Service) CreateApplianceAttr(ctx context.Context, in *pb.CreateApplianceAttrRequest) (*emptypb.Empty, error) {
	switch {
	case !in.HasZone():
		return nil, errMissingZone
	case !in.HasAppliance():
		return nil, errMissingAppliance
	case !in.HasName():
		return nil, errMissingAttribute
	}

	exists, err := s.store.CreateApplianceAttr(ctx, in.GetName(), data.ApplianceScope{
		Zone:      in.GetZone(),
		Appliance: in.GetAppliance(),
	})
	if err != nil {
		return nil, dberr(err)
	}

	if exists {
		return nil, status.Errorf(codes.AlreadyExists, "zone %q appliance %q attr %q",
			in.GetZone(), in.GetName(), in.GetName())
	}

	return new(emptypb.Empty), nil
}

// ReadApplianceAttrs implements the grpc StackServiceServer interface.
func (s *Service) ReadApplianceAttrs(in *pb.ReadApplianceAttrsRequest, out grpc.ServerStreamingServer[pb.ReadApplianceAttrsResponse]) error {
	var filter data.Filter

	ctx := context.Background()

	switch {
	case in.HasName():
		filter = data.FilterByName(in.GetName())
	case in.HasGlob():
		filter = data.FilterByGlob(in.GetGlob())
	}

	rows, err := s.store.ReadApplianceAttrs(ctx, filter, data.ApplianceScope{
		Zone:      in.GetZone(),
		Appliance: in.GetAppliance(),
	})
	if err != nil {
		return dberr(err)
	}

	for i := range rows {
		resp := pb.ReadApplianceAttrsResponse_builder{
			Appliance: &rows[i].Appliance,
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

// UpdateApplianceAttr implements the grpc StackServiceServer interface.
func (s *Service) UpdateApplianceAttr(ctx context.Context, in *pb.UpdateApplianceAttrRequest) (*emptypb.Empty, error) {
	switch {
	case !in.HasZone():
		return nil, errMissingZone
	case !in.HasAppliance():
		return nil, errMissingAppliance
	case !in.HasName():
		return nil, errMissingAttribute
	}

	fields := in.GetFields()

	scope := data.ApplianceAttrScope{
		Zone:      in.GetZone(),
		Appliance: in.GetAppliance(),
		Attr:      in.GetName(),
	}

	if fields.HasName() {
		err := s.store.UpdateApplianceAttrName(ctx, fields.GetName(), scope)
		if err != nil {
			return nil, dberr(err)
		}
	}

	if fields.HasValue() {
		err := s.store.UpdateApplianceAttrValue(ctx, fields.GetValue(), scope)
		if err != nil {
			return nil, dberr(err)
		}
	}

	if fields.HasProtected() {
		err := s.store.UpdateApplianceAttrProtection(ctx, fields.GetProtected(), scope)
		if err != nil {
			return nil, dberr(err)
		}
	}

	return new(emptypb.Empty), nil
}

// DeleteApplianceAttrs implements the grpc StackServiceServer interface.
func (s *Service) DeleteApplianceAttrs(ctx context.Context, in *pb.DeleteApplianceAttrsRequest) (*emptypb.Empty, error) {
	var filter data.Filter

	switch {
	case !in.HasZone():
		return nil, errMissingZone
	case !in.HasAppliance():
		return nil, errMissingAppliance
	case in.HasName():
		filter = data.FilterByName(in.GetName())
	case in.HasGlob():
		filter = data.FilterByGlob(in.GetGlob())
	default:
		return nil, errNameOrGlob
	}

	err := s.store.DeleteApplianceAttrs(ctx, filter, data.ApplianceScope{
		Zone:      in.GetZone(),
		Appliance: in.GetAppliance(),
	})
	if err != nil {
		return nil, dberr(err)
	}

	return new(emptypb.Empty), nil
}
