// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.12.4
// source: metrics.proto

package __

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	MetricsService_UpdateMetrics_FullMethodName = "/metalsys.MetricsService/UpdateMetrics"
)

// MetricsServiceClient is the client API for MetricsService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MetricsServiceClient interface {
	UpdateMetrics(ctx context.Context, in *MetricsRequest, opts ...grpc.CallOption) (*MetricResponse, error)
}

type metricsServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewMetricsServiceClient(cc grpc.ClientConnInterface) MetricsServiceClient {
	return &metricsServiceClient{cc}
}

func (c *metricsServiceClient) UpdateMetrics(ctx context.Context, in *MetricsRequest, opts ...grpc.CallOption) (*MetricResponse, error) {
	out := new(MetricResponse)
	err := c.cc.Invoke(ctx, MetricsService_UpdateMetrics_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MetricsServiceServer is the server API for MetricsService service.
// All implementations must embed UnimplementedMetricsServiceServer
// for forward compatibility
type MetricsServiceServer interface {
	UpdateMetrics(context.Context, *MetricsRequest) (*MetricResponse, error)
	mustEmbedUnimplementedMetricsServiceServer()
}

// UnimplementedMetricsServiceServer must be embedded to have forward compatible implementations.
type UnimplementedMetricsServiceServer struct {
}

func (UnimplementedMetricsServiceServer) UpdateMetrics(context.Context, *MetricsRequest) (*MetricResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateMetrics not implemented")
}
func (UnimplementedMetricsServiceServer) mustEmbedUnimplementedMetricsServiceServer() {}

// UnsafeMetricsServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MetricsServiceServer will
// result in compilation errors.
type UnsafeMetricsServiceServer interface {
	mustEmbedUnimplementedMetricsServiceServer()
}

func RegisterMetricsServiceServer(s grpc.ServiceRegistrar, srv MetricsServiceServer) {
	s.RegisterService(&MetricsService_ServiceDesc, srv)
}

func _MetricsService_UpdateMetrics_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MetricsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MetricsServiceServer).UpdateMetrics(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MetricsService_UpdateMetrics_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MetricsServiceServer).UpdateMetrics(ctx, req.(*MetricsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// MetricsService_ServiceDesc is the grpc.ServiceDesc for MetricsService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var MetricsService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "metalsys.MetricsService",
	HandlerType: (*MetricsServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "UpdateMetrics",
			Handler:    _MetricsService_UpdateMetrics_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "metrics.proto",
}