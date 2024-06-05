// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.25.1
// source: zlog.proto

package zlog_pb

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
	Zlog_AddLogKV_FullMethodName = "/zlog_pb.zlog/AddLogKV"
)

// ZlogClient is the client API for Zlog service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ZlogClient interface {
	AddLogKV(ctx context.Context, in *LogKV_REQ, opts ...grpc.CallOption) (*Default_RES, error)
}

type zlogClient struct {
	cc grpc.ClientConnInterface
}

func NewZlogClient(cc grpc.ClientConnInterface) ZlogClient {
	return &zlogClient{cc}
}

func (c *zlogClient) AddLogKV(ctx context.Context, in *LogKV_REQ, opts ...grpc.CallOption) (*Default_RES, error) {
	out := new(Default_RES)
	err := c.cc.Invoke(ctx, Zlog_AddLogKV_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ZlogServer is the server API for Zlog service.
// All implementations must embed UnimplementedZlogServer
// for forward compatibility
type ZlogServer interface {
	AddLogKV(context.Context, *LogKV_REQ) (*Default_RES, error)
	mustEmbedUnimplementedZlogServer()
}

// UnimplementedZlogServer must be embedded to have forward compatible implementations.
type UnimplementedZlogServer struct {
}

func (UnimplementedZlogServer) AddLogKV(context.Context, *LogKV_REQ) (*Default_RES, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddLogKV not implemented")
}
func (UnimplementedZlogServer) mustEmbedUnimplementedZlogServer() {}

// UnsafeZlogServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ZlogServer will
// result in compilation errors.
type UnsafeZlogServer interface {
	mustEmbedUnimplementedZlogServer()
}

func RegisterZlogServer(s grpc.ServiceRegistrar, srv ZlogServer) {
	s.RegisterService(&Zlog_ServiceDesc, srv)
}

func _Zlog_AddLogKV_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LogKV_REQ)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ZlogServer).AddLogKV(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Zlog_AddLogKV_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ZlogServer).AddLogKV(ctx, req.(*LogKV_REQ))
	}
	return interceptor(ctx, in, info, handler)
}

// Zlog_ServiceDesc is the grpc.ServiceDesc for Zlog service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Zlog_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "zlog_pb.zlog",
	HandlerType: (*ZlogServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AddLogKV",
			Handler:    _Zlog_AddLogKV_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "zlog.proto",
}
