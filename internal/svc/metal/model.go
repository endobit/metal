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
// Model
//

// CreateModel implements the grpc StackServiceServer interface.
func (s Service) CreateModel(ctx context.Context, in *pb.CreateModelRequest) (*emptypb.Empty, error) {
	switch {
	case !in.HasMake():
		return nil, errMissingMake
	case !in.HasName():
		return nil, errMissingModel
	}

	exists, err := s.store.CreateModel(ctx, in.GetName(), data.MakeScope{Make: in.GetMake()})
	if err != nil {
		return nil, dberr(err)
	}

	if exists {
		return nil, status.Errorf(codes.AlreadyExists, "make %q model %q", in.GetMake(), in.GetName())
	}

	return new(emptypb.Empty), nil
}

// ReadModels implements the grpc StackServiceServer interface.
func (s Service) ReadModels(in *pb.ReadModelsRequest, out grpc.ServerStreamingServer[pb.ReadModelsResponse]) error {
	var filter data.Filter

	switch {
	case in.HasName():
		filter = data.FilterByName(in.GetName())
	case in.HasGlob():
		filter = data.FilterByGlob(in.GetGlob())
	}

	rows, err := s.store.ReadModels(context.Background(), filter, data.MakeScope{Make: in.GetMake()})
	if err != nil {
		return dberr(err)
	}

	for i := range rows {
		resp := pb.ReadModelsResponse_builder{
			Make: &rows[i].Make,
			Name: &rows[i].Name,
		}.Build()
		if rows[i].Architecture != "" {
			resp.SetArchitecture(pb.Architecture(pb.Architecture_value[rows[i].Architecture]))
		}

		if err := out.Send(resp); err != nil {
			return err
		}
	}

	return nil
}

// UpdateModel implements the grpc StackServiceServer interface.
func (s Service) UpdateModel(ctx context.Context, in *pb.UpdateModelRequest) (*emptypb.Empty, error) {
	switch {
	case !in.HasMake():
		return nil, errMissingMake
	case !in.HasName():
		return nil, errMissingModel
	}

	fields := in.GetFields()

	scope := data.ModelScope{
		Make:  in.GetMake(),
		Model: in.GetName(),
	}

	if fields.HasName() {
		if err := s.store.UpdateModelName(ctx, fields.GetName(), scope); err != nil {
			return nil, dberr(err)
		}
	}

	if fields.HasArchitecture() {
		err := s.store.UpdateModelArchitecture(ctx, fields.GetArchitecture().String(), scope)
		if err != nil {
			return nil, dberr(err)
		}
	}

	return new(emptypb.Empty), nil
}

// DeleteModels implements the grpc StackServiceServer interface.
func (s Service) DeleteModels(ctx context.Context, in *pb.DeleteModelsRequest) (*emptypb.Empty, error) {
	var filter data.Filter

	switch {
	case !in.HasMake():
		return nil, errMissingMake
	case in.HasName():
		filter = data.FilterByName(in.GetName())
	case in.HasGlob():
		filter = data.FilterByGlob(in.GetGlob())
	default:
		return nil, errNameOrGlob
	}

	err := s.store.DeleteModels(ctx, filter, data.MakeScope{Make: in.GetMake()})
	if err != nil {
		return nil, dberr(err)
	}

	return new(emptypb.Empty), nil
}

//
// Model Attributes
//

// CreateModelAttr implements the grpc StackServiceServer interface.
func (s *Service) CreateModelAttr(ctx context.Context, in *pb.CreateModelAttrRequest) (*emptypb.Empty, error) {
	switch {
	case !in.HasMake():
		return nil, errMissingMake
	case !in.HasModel():
		return nil, errMissingModel
	case !in.HasName():
		return nil, errMissingAttribute
	}

	exists, err := s.store.CreateModelAttr(ctx, in.GetName(), data.ModelScope{
		Make:  in.GetMake(),
		Model: in.GetModel(),
	})
	if err != nil {
		return nil, dberr(err)
	}

	if exists {
		return nil, status.Errorf(codes.AlreadyExists, "make %q model %q attribute %q",
			in.GetMake(), in.GetModel(), in.GetName())
	}

	return new(emptypb.Empty), nil
}

// ReadModelAttrs implements the grpc StackServiceServer interface.
func (s *Service) ReadModelAttrs(in *pb.ReadModelAttrsRequest, out grpc.ServerStreamingServer[pb.ReadModelAttrsResponse]) error {
	var filter data.Filter

	ctx := context.Background()

	switch {
	case in.HasName():
		filter = data.FilterByName(in.GetName())
	case in.HasGlob():
		filter = data.FilterByGlob(in.GetGlob())
	}

	rows, err := s.store.ReadModelAttrs(ctx, filter, data.ModelScope{
		Make:  in.GetMake(),
		Model: in.GetModel(),
	})
	if err != nil {
		return err
	}

	for i := range rows {
		resp := pb.ReadModelAttrsResponse_builder{
			Make:      &rows[i].Make,
			Model:     &rows[i].Model,
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

// UpdateModelAttr implements the grpc StackServiceServer interface.
func (s *Service) UpdateModelAttr(ctx context.Context, in *pb.UpdateModelAttrRequest) (*emptypb.Empty, error) {
	switch {
	case !in.HasMake():
		return nil, errMissingMake
	case !in.HasModel():
		return nil, errMissingModel
	case !in.HasName():
		return nil, errMissingAttribute
	}

	fields := in.GetFields()

	scope := data.ModelAttrScope{
		Make:  in.GetMake(),
		Model: in.GetModel(),
		Attr:  in.GetName(),
	}

	if fields.HasName() {
		err := s.store.UpdateModelAttrName(ctx, fields.GetName(), scope)
		if err != nil {
			return nil, dberr(err)
		}
	}

	if fields.HasValue() {
		err := s.store.UpdateModelAttrValue(ctx, fields.GetValue(), scope)
		if err != nil {
			return nil, dberr(err)
		}
	}

	if fields.HasProtected() {
		if err := s.store.UpdateModelAttrProtection(ctx, fields.GetProtected(), scope); err != nil {
			return nil, dberr(err)
		}
	}

	return new(emptypb.Empty), nil
}

// DeleteModelAttrs implements the grpc StackServiceServer interface.
func (s *Service) DeleteModelAttrs(ctx context.Context, in *pb.DeleteModelAttrsRequest) (*emptypb.Empty, error) {
	var filter data.Filter

	switch {
	case !in.HasMake():
		return nil, errMissingMake
	case !in.HasModel():
		return nil, errMissingModel
	case in.HasName():
		filter = data.FilterByName(in.GetName())
	case in.HasGlob():
		filter = data.FilterByGlob(in.GetGlob())
	default:
		return nil, errNameOrGlob
	}

	err := s.store.DeleteModelAttrs(ctx, filter, data.ModelScope{
		Make:  in.GetMake(),
		Model: in.GetModel(),
	})
	if err != nil {
		return nil, dberr(err)
	}

	return new(emptypb.Empty), nil
}
