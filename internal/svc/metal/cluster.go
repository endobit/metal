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
// Cluster
//

// CreateCluster implements the grpc StackServiceServer interface.
func (s Service) CreateCluster(ctx context.Context, in *pb.CreateClusterRequest) (*emptypb.Empty, error) {
	switch {
	case !in.HasZone():
		return nil, errMissingZone
	case !in.HasName():
		return nil, errMissingCluster
	}

	exists, err := s.store.CreateCluster(ctx, in.GetName(), data.ZoneScope{Zone: in.GetZone()})
	if err != nil {
		return nil, dberr(err)
	}

	if exists {
		return nil, status.Errorf(codes.AlreadyExists, "zone %q cluster %q", in.GetZone(), in.GetName())
	}

	return new(emptypb.Empty), nil
}

// ReadClusters implements the grpc StackServiceServer interface.
func (s Service) ReadClusters(in *pb.ReadClustersRequest, out grpc.ServerStreamingServer[pb.ReadClustersResponse]) error {
	var filter data.Filter

	ctx := context.Background()

	switch {
	case in.HasName():
		filter = data.FilterByName(in.GetName())
	case in.HasGlob():
		filter = data.FilterByGlob(in.GetGlob())
	}

	rows, err := s.store.ReadClusters(ctx, filter, data.ZoneScope{Zone: in.GetZone()})
	if err != nil {
		return dberr(err)
	}

	for i := range rows {
		resp := pb.ReadClustersResponse_builder{
			Zone: &rows[i].Zone,
			Name: &rows[i].Name,
		}.Build()

		if err := out.Send(resp); err != nil {
			return err
		}
	}

	return nil
}

// UpdateCluster implements the grpc StackServiceServer interface.
func (s Service) UpdateCluster(ctx context.Context, in *pb.UpdateClusterRequest) (*emptypb.Empty, error) {
	switch {
	case !in.HasZone():
		return nil, errMissingZone
	case !in.HasName():
		return nil, errMissingCluster
	}

	fields := in.GetFields()

	scope := data.ClusterScope{
		Zone:    in.GetZone(),
		Cluster: in.GetName(),
	}

	if fields.HasName() {
		if err := s.store.UpdateClusterName(ctx, fields.GetName(), scope); err != nil {
			return nil, dberr(err)
		}
	}

	return new(emptypb.Empty), nil
}

// DeleteClusters implements the grpc StackServiceServer interface.
func (s Service) DeleteClusters(ctx context.Context, in *pb.DeleteClustersRequest) (*emptypb.Empty, error) {
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

	err := s.store.DeleteClusters(ctx, filter, data.ZoneScope{Zone: in.GetZone()})
	if err != nil {
		return nil, dberr(err)
	}

	return new(emptypb.Empty), nil
}

//
// Cluster Attributes
//

// CreateClusterAttr implements the grpc StackServiceServer interface.
func (s *Service) CreateClusterAttr(ctx context.Context, in *pb.CreateClusterAttrRequest) (*emptypb.Empty, error) {
	switch {
	case !in.HasZone():
		return nil, errMissingZone
	case !in.HasCluster():
		return nil, errMissingCluster
	case !in.HasName():
		return nil, errMissingAttribute
	}

	exists, err := s.store.CreateClusterAttr(ctx, in.GetName(), data.ClusterScope{
		Zone:    in.GetZone(),
		Cluster: in.GetCluster(),
	})
	if err != nil {
		return nil, dberr(err)
	}

	if exists {
		return nil, status.Errorf(codes.AlreadyExists, "zone %q cluster %q attr %q",
			in.GetZone(), in.GetName(), in.GetName())
	}

	return new(emptypb.Empty), nil
}

// ReadClusterAttrs implements the grpc StackServiceServer interface.
func (s *Service) ReadClusterAttrs(in *pb.ReadClusterAttrsRequest, out grpc.ServerStreamingServer[pb.ReadClusterAttrsResponse]) error {
	var filter data.Filter

	ctx := context.Background()

	switch {
	case in.HasName():
		filter = data.FilterByName(in.GetName())
	case in.HasGlob():
		filter = data.FilterByGlob(in.GetGlob())
	}

	rows, err := s.store.ReadClusterAttrs(ctx, filter, data.ClusterScope{
		Zone:    in.GetZone(),
		Cluster: in.GetCluster(),
	})
	if err != nil {
		return dberr(err)
	}

	for i := range rows {
		resp := pb.ReadClusterAttrsResponse_builder{
			Cluster:   &rows[i].Cluster,
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

// UpdateClusterAttr implements the grpc StackServiceServer interface.
func (s *Service) UpdateClusterAttr(ctx context.Context, in *pb.UpdateClusterAttrRequest) (*emptypb.Empty, error) {
	switch {
	case !in.HasZone():
		return nil, errMissingZone
	case !in.HasCluster():
		return nil, errMissingCluster
	case !in.HasName():
		return nil, errMissingAttribute
	}

	fields := in.GetFields()

	scope := data.ClusterAttrScope{
		Zone:    in.GetZone(),
		Cluster: in.GetCluster(),
		Attr:    in.GetName(),
	}

	if fields.HasName() {
		err := s.store.UpdateClusterAttrName(ctx, fields.GetName(), scope)
		if err != nil {
			return nil, dberr(err)
		}
	}

	if fields.HasValue() {
		err := s.store.UpdateClusterAttrValue(ctx, fields.GetValue(), scope)
		if err != nil {
			return nil, dberr(err)
		}
	}

	if fields.HasProtected() {
		err := s.store.UpdateClusterAttrProtection(ctx, fields.GetProtected(), scope)
		if err != nil {
			return nil, dberr(err)
		}
	}

	return new(emptypb.Empty), nil
}

// DeleteClusterAttrs implements the grpc StackServiceServer interface.
func (s *Service) DeleteClusterAttrs(ctx context.Context, in *pb.DeleteClusterAttrsRequest) (*emptypb.Empty, error) {
	var filter data.Filter

	switch {
	case !in.HasZone():
		return nil, errMissingZone
	case !in.HasCluster():
		return nil, errMissingCluster
	case in.HasName():
		filter = data.FilterByName(in.GetName())
	case in.HasGlob():
		filter = data.FilterByGlob(in.GetGlob())
	default:
		return nil, errNameOrGlob
	}

	err := s.store.DeleteClusterAttrs(ctx, filter, data.ClusterScope{
		Zone:    in.GetZone(),
		Cluster: in.GetCluster(),
	})
	if err != nil {
		return nil, dberr(err)
	}

	return new(emptypb.Empty), nil
}
