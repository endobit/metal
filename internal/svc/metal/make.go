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
// Make
//

// CreateMake implements the grpc StackServiceServer interface.
func (s Service) CreateMake(ctx context.Context, in *pb.CreateMakeRequest) (*emptypb.Empty, error) {
	if !in.HasName() {
		return nil, errMissingMake
	}

	exists, err := s.store.CreateMake(ctx, in.GetName())
	if err != nil {
		return nil, dberr(err)
	}

	if exists {
		return nil, status.Errorf(codes.AlreadyExists, "make %q", in.GetName())
	}

	return new(emptypb.Empty), nil
}

// ReadMakes implements the grpc StackServiceServer interface.
func (s Service) ReadMakes(in *pb.ReadMakesRequest, out grpc.ServerStreamingServer[pb.ReadMakesResponse]) error {
	var filter data.Filter

	ctx := context.Background()

	switch {
	case in.HasName():
		filter = data.FilterByName(in.GetName())
	case in.HasGlob():
		filter = data.FilterByGlob(in.GetGlob())
	}

	rows, err := s.store.ReadMakes(ctx, filter)
	if err != nil {
		return dberr(err)
	}

	for i := range rows {
		resp := pb.ReadMakesResponse_builder{
			Name: &rows[i].Name,
		}.Build()

		if err := out.Send(resp); err != nil {
			return err
		}
	}

	return nil
}

// UpdateMake implements the grpc StackServiceServer interface.
func (s Service) UpdateMake(ctx context.Context, in *pb.UpdateMakeRequest) (*emptypb.Empty, error) {
	fields := in.GetFields()

	scope := data.MakeScope{Make: in.GetName()}

	if fields.HasName() {
		if err := s.store.UpdateMakeName(ctx, fields.GetName(), scope); err != nil {
			return nil, dberr(err)
		}
	}

	return new(emptypb.Empty), nil
}

// DeleteMakes implements the grpc StackServiceServer interface.
func (s Service) DeleteMakes(ctx context.Context, in *pb.DeleteMakesRequest) (*emptypb.Empty, error) {
	var filter data.Filter

	switch {
	case in.HasName():
		filter = data.FilterByName(in.GetName())
	case in.HasGlob():
		filter = data.FilterByGlob(in.GetGlob())
	default:
		return nil, errNameOrGlob
	}

	err := s.store.DeleteMakes(ctx, filter)
	if err != nil {
		return nil, dberr(err)
	}

	return new(emptypb.Empty), nil
}
