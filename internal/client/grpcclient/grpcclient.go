package grpcclient

import (
	"context"
	"log"

	"github.com/GrebenschikovDI/metalsys.git/internal/common/models"
	pb "github.com/GrebenschikovDI/metalsys.git/internal/proto"
	structpb "github.com/golang/protobuf/ptypes/struct"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GrpcClient struct {
	conn    *grpc.ClientConn
	service pb.MetricsServiceClient
}

func NewGrpcClient(serverAddress string) (*GrpcClient, error) {
	conn, err := grpc.Dial(serverAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
		return nil, err
	}

	client := pb.NewMetricsServiceClient(conn)

	return &GrpcClient{
		conn:    conn,
		service: client,
	}, nil
}

func (c *GrpcClient) Close() {
	c.conn.Close()
}

func (c *GrpcClient) UpdateMetrics(ctx context.Context, metrics []models.Metric) error {
	var grpcMetrics []*pb.Metric

	for _, metric := range metrics {
		grpcMetric := ConvertToGRPC(metric)
		grpcMetrics = append(grpcMetrics, grpcMetric)
	}

	request := &pb.MetricsRequest{
		Metrics: grpcMetrics,
	}

	_, err := c.service.UpdateMetrics(ctx, request)
	if err != nil {
		log.Printf("Failed to update metrics: %v", err)
		return err
	}

	return nil
}

func ConvertToGRPC(originalMetric models.Metric) *pb.Metric {
	grpcMetric := &pb.Metric{
		Id:    originalMetric.ID,
		Mtype: originalMetric.Mtype,
	}

	if originalMetric.Delta != nil {
		delta := *originalMetric.Delta
		grpcMetric.Delta = &structpb.Value{
			Kind: &structpb.Value_NumberValue{
				NumberValue: float64(delta),
			},
		}

	}

	if originalMetric.Value != nil {
		value := *originalMetric.Value
		grpcMetric.Value = &structpb.Value{
			Kind: &structpb.Value_NumberValue{
				NumberValue: value,
			},
		}
	}

	return grpcMetric
}
