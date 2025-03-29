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
// Environment
//

// CreateEnvironment implements the grpc StackServiceServer interface.
func (s Service) CreateEnvironment(ctx context.Context, in *pb.CreateEnvironmentRequest) (*emptypb.Empty, error) {
	switch {
	case !in.HasZone():
		return nil, errMissingZone
	case !in.HasName():
		return nil, errMissingEnvironment
	}

	exists, err := s.store.CreateEnvironment(ctx, in.GetName(), data.ZoneScope{Zone: in.GetZone()})
	if err != nil {
		return nil, dberr(err)
	}

	if exists {
		return nil, status.Errorf(codes.AlreadyExists, "zone %q environment %q", in.GetZone(), in.GetName())
	}

	return new(emptypb.Empty), nil
}

// ReadEnvironments implements the grpc StackServiceServer interface.
func (s Service) ReadEnvironments(in *pb.ReadEnvironmentsRequest, out grpc.ServerStreamingServer[pb.ReadEnvironmentsResponse]) error {
	var filter data.Filter

	ctx := context.Background()

	switch {
	case in.HasName():
		filter = data.FilterByName(in.GetName())
	case in.HasGlob():
		filter = data.FilterByGlob(in.GetGlob())
	}

	rows, err := s.store.ReadEnvironments(ctx, filter, data.ZoneScope{Zone: in.GetZone()})
	if err != nil {
		return dberr(err)
	}

	for i := range rows {
		resp := pb.ReadEnvironmentsResponse_builder{
			Zone: &rows[i].Zone,
			Name: &rows[i].Name,
		}.Build()

		if err := out.Send(resp); err != nil {
			return err
		}
	}

	return nil
}

// UpdateEnvironment implements the grpc StackServiceServer interface.
func (s Service) UpdateEnvironment(ctx context.Context, in *pb.UpdateEnvironmentRequest) (*emptypb.Empty, error) {
	switch {
	case !in.HasZone():
		return nil, errMissingZone
	case !in.HasName():
		return nil, errMissingEnvironment
	}

	fields := in.GetFields()

	scope := data.EnvironmentScope{
		Zone:        in.GetZone(),
		Environment: in.GetName(),
	}

	if fields.HasName() {
		if err := s.store.UpdateEnvironmentName(ctx, fields.GetName(), scope); err != nil {
			return nil, dberr(err)
		}
	}

	return new(emptypb.Empty), nil
}

// DeleteEnvironments implements the grpc StackServiceServer interface.
func (s Service) DeleteEnvironments(ctx context.Context, in *pb.DeleteEnvironmentsRequest) (*emptypb.Empty, error) {
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

	err := s.store.DeleteEnvironments(ctx, filter, data.ZoneScope{Zone: in.GetZone()})
	if err != nil {
		return nil, dberr(err)
	}

	return new(emptypb.Empty), nil
}

//
// Environment Attributes
//

// CreateEnvironmentAttr implements the grpc StackServiceServer interface.
func (s *Service) CreateEnvironmentAttr(ctx context.Context, in *pb.CreateEnvironmentAttrRequest) (*emptypb.Empty, error) {
	switch {
	case !in.HasZone():
		return nil, errMissingZone
	case !in.HasEnvironment():
		return nil, errMissingEnvironment
	case !in.HasName():
		return nil, errMissingAttribute
	}

	exists, err := s.store.CreateEnvironmentAttr(ctx, in.GetName(), data.EnvironmentScope{
		Zone:        in.GetZone(),
		Environment: in.GetEnvironment(),
	})
	if err != nil {
		return nil, dberr(err)
	}

	if exists {
		return nil, status.Errorf(codes.AlreadyExists, "zone %q environment %q attr %q",
			in.GetZone(), in.GetName(), in.GetName())
	}

	return new(emptypb.Empty), nil
}

// ReadEnvironmentAttrs implements the grpc StackServiceServer interface.
func (s *Service) ReadEnvironmentAttrs(in *pb.ReadEnvironmentAttrsRequest, out grpc.ServerStreamingServer[pb.ReadEnvironmentAttrsResponse]) error {
	var filter data.Filter

	ctx := context.Background()

	switch {
	case in.HasName():
		filter = data.FilterByName(in.GetName())
	case in.HasGlob():
		filter = data.FilterByGlob(in.GetGlob())
	}

	rows, err := s.store.ReadEnvironmentAttrs(ctx, filter, data.EnvironmentScope{
		Zone:        in.GetZone(),
		Environment: in.GetEnvironment(),
	})
	if err != nil {
		return dberr(err)
	}

	for i := range rows {
		resp := pb.ReadEnvironmentAttrsResponse_builder{
			Environment: &rows[i].Environment,
			Name:        &rows[i].Name,
			Value:       &rows[i].Value,
			Protected:   &rows[i].IsProtected,
		}.Build()

		if err := out.Send(resp); err != nil {
			return err
		}
	}

	return nil
}

// UpdateEnvironmentAttr implements the grpc StackServiceServer interface.
func (s *Service) UpdateEnvironmentAttr(ctx context.Context, in *pb.UpdateEnvironmentAttrRequest) (*emptypb.Empty, error) {
	switch {
	case !in.HasZone():
		return nil, errMissingZone
	case !in.HasEnvironment():
		return nil, errMissingEnvironment
	case !in.HasName():
		return nil, errMissingAttribute
	}

	fields := in.GetFields()

	scope := data.EnvironmentAttrScope{
		Zone:        in.GetZone(),
		Environment: in.GetEnvironment(),
		Attr:        in.GetName(),
	}

	if fields.HasName() {
		err := s.store.UpdateEnvironmentAttrName(ctx, fields.GetName(), scope)
		if err != nil {
			return nil, dberr(err)
		}
	}

	if fields.HasValue() {
		err := s.store.UpdateEnvironmentAttrValue(ctx, fields.GetValue(), scope)
		if err != nil {
			return nil, dberr(err)
		}
	}

	if fields.HasProtected() {
		err := s.store.UpdateEnvironmentAttrProtection(ctx, fields.GetProtected(), scope)
		if err != nil {
			return nil, dberr(err)
		}
	}

	return new(emptypb.Empty), nil
}

// DeleteEnvironmentAttrs implements the grpc StackServiceServer interface.
func (s *Service) DeleteEnvironmentAttrs(ctx context.Context, in *pb.DeleteEnvironmentAttrsRequest) (*emptypb.Empty, error) {
	var filter data.Filter

	switch {
	case !in.HasZone():
		return nil, errMissingZone
	case !in.HasEnvironment():
		return nil, errMissingEnvironment
	case in.HasName():
		filter = data.FilterByName(in.GetName())
	case in.HasGlob():
		filter = data.FilterByGlob(in.GetGlob())
	default:
		return nil, errNameOrGlob
	}

	err := s.store.DeleteEnvironmentAttrs(ctx, filter, data.EnvironmentScope{
		Zone:        in.GetZone(),
		Environment: in.GetEnvironment(),
	})
	if err != nil {
		return nil, dberr(err)
	}

	return new(emptypb.Empty), nil
}
