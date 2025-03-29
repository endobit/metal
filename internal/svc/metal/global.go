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
// Attributes
//

// CreateGlobalAttr implements the grpc StackServiceServer interface.
func (s Service) CreateGlobalAttr(ctx context.Context, in *pb.CreateGlobalAttrRequest) (*emptypb.Empty, error) {
	if !in.HasName() {
		return nil, errMissingAttribute
	}

	exists, err := s.store.CreateGlobalAttr(ctx, in.GetName())
	if err != nil {
		return nil, dberr(err)
	}

	if exists {
		return nil, status.Errorf(codes.AlreadyExists, "attribute %q", in.GetName())
	}

	return new(emptypb.Empty), nil
}

// ReadGlobalAttrs implements the grpc StackServiceServer interface.
func (s Service) ReadGlobalAttrs(in *pb.ReadGlobalAttrsRequest, out grpc.ServerStreamingServer[pb.ReadGlobalAttrsResponse]) error {
	var filter data.Filter

	ctx := context.Background()

	switch {
	case in.HasName():
		filter = data.FilterByName(in.GetName())
	case in.HasGlob():
		filter = data.FilterByGlob(in.GetGlob())
	}

	rows, err := s.store.ReadGlobalAttrs(ctx, filter)
	if err != nil {
		return err
	}

	for i := range rows {
		resp := pb.ReadGlobalAttrsResponse_builder{
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

// UpdateGlobalAttr implements the grpc StackServiceServer interface.
func (s Service) UpdateGlobalAttr(ctx context.Context, in *pb.UpdateGlobalAttrRequest) (*emptypb.Empty, error) {
	if !in.HasName() {
		return nil, errMissingAttribute
	}

	fields := in.GetFields()

	scope := data.GlobalAttrScope{
		Attr: in.GetName(),
	}

	if fields.HasName() {
		if err := s.store.UpdateGlobalAttrName(ctx, fields.GetName(), scope); err != nil {
			return nil, dberr(err)
		}
	}

	if fields.HasValue() {
		if err := s.store.UpdateGlobalAttrValue(ctx, fields.GetValue(), scope); err != nil {
			return nil, dberr(err)
		}
	}

	if fields.HasProtected() {
		if err := s.store.UpdateGlobalAttrProtection(ctx, fields.GetProtected(), scope); err != nil {
			return nil, dberr(err)
		}
	}

	return new(emptypb.Empty), nil
}

// DeleteGlobalAttrs implements the grpc StackServiceServer interface.
func (s Service) DeleteGlobalAttrs(ctx context.Context, in *pb.DeleteGlobalAttrsRequest) (*emptypb.Empty, error) {
	var filter data.Filter

	switch {
	case in.HasName():
		filter = data.FilterByName(in.GetName())
	case in.HasGlob():
		filter = data.FilterByGlob(in.GetGlob())
	default:
		return nil, errNameOrGlob
	}

	err := s.store.DeleteGlobalAttrs(ctx, filter)
	if err != nil {
		return nil, dberr(err)
	}

	return new(emptypb.Empty), nil
}
