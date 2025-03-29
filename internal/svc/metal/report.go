package metal

import (
	"context"
	"encoding/json"

	"endobit.io/metal"
	pb "endobit.io/metal/gen/go/proto/metal/v1"
	"endobit.io/metal/internal/data"
)

// ReadSchema implements the grpc StackServiceServer interface.
func (s Service) ReadReportData(ctx context.Context, in *pb.ReadReportDataRequest) (*pb.ReadReportDataResponse, error) {
	var resp pb.ReadReportDataResponse

	d := metal.Dumper{Store: s.store}

	if in.HasZone() {
		d.Filter.Zone = data.FilterByName(in.GetZone())
	}

	if in.HasCluster() {
		d.Filter.Cluster = data.FilterByName(in.GetCluster())
	}

	if in.HasHost() {
		d.Filter.Host = data.FilterByName(in.GetHost())
	}

	report, err := d.Dump(ctx)
	if err != nil {
		return nil, err
	}

	doc, err := json.Marshal(report)
	if err != nil {
		return nil, err
	}

	resp.SetData(doc)

	return &resp, nil
}
