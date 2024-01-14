package grpcserver

import (
	"context"

	"github.com/GrebenschikovDI/metalsys.git/internal/common/models"
	"github.com/GrebenschikovDI/metalsys.git/internal/common/repository"
	pb "github.com/GrebenschikovDI/metalsys.git/internal/proto"
)

type GrpsServer struct {
	pb.UnimplementedMetricsServiceServer
	storage repository.Repository
}

func CreateGrpcServer(storage repository.Repository) *GrpsServer {
	return &GrpsServer{
		storage: storage,
	}
}

func (s *GrpsServer) UpdateMetrics(ctx context.Context, req *pb.MetricsRequest) (*pb.MetricResponse, error) {
	metrics := make([]models.Metric, len(req.Metrics))

	for i, grpcMetric := range req.Metrics {
		metric := ConvertGRPCMetricToOriginal(grpcMetric)
		metrics[i] = metric
	}

	s.storage.AddMetrics(ctx, metrics)

	return &pb.MetricResponse{
		Status: "OK",
	}, nil
}

// ConvertGRPCMetricToOriginal преобразует структуру gRPC в оригинальную структуру Metric
func ConvertGRPCMetricToOriginal(grpcMetric *pb.Metric) models.Metric {
	originalMetric := models.Metric{
		ID:    grpcMetric.Id,
		Mtype: grpcMetric.Mtype,
	}

	if grpcMetric.Delta != nil {
		deltaValue := int64(grpcMetric.Delta.GetNumberValue())
		originalMetric.Delta = &deltaValue
	}

	if grpcMetric.Value != nil {
		valueValue := grpcMetric.Value.GetNumberValue()
		originalMetric.Value = &valueValue
	}

	return originalMetric
}
