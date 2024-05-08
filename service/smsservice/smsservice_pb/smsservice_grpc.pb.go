// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.25.1
// source: smsservice.proto

package smsservice_pb

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
	Smsservice_SendVerifyCode_FullMethodName = "/smsservice_pb.smsservice/SendVerifyCode"
	Smsservice_VerifyCode_FullMethodName     = "/smsservice_pb.smsservice/VerifyCode"
)

// SmsserviceClient is the client API for Smsservice service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SmsserviceClient interface {
	SendVerifyCode(ctx context.Context, in *SendVerifyCode_REQ, opts ...grpc.CallOption) (*Default_RES, error)
	VerifyCode(ctx context.Context, in *VerifyCode_REQ, opts ...grpc.CallOption) (*Default_RES, error)
}

type smsserviceClient struct {
	cc grpc.ClientConnInterface
}

func NewSmsserviceClient(cc grpc.ClientConnInterface) SmsserviceClient {
	return &smsserviceClient{cc}
}

func (c *smsserviceClient) SendVerifyCode(ctx context.Context, in *SendVerifyCode_REQ, opts ...grpc.CallOption) (*Default_RES, error) {
	out := new(Default_RES)
	err := c.cc.Invoke(ctx, Smsservice_SendVerifyCode_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *smsserviceClient) VerifyCode(ctx context.Context, in *VerifyCode_REQ, opts ...grpc.CallOption) (*Default_RES, error) {
	out := new(Default_RES)
	err := c.cc.Invoke(ctx, Smsservice_VerifyCode_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SmsserviceServer is the server API for Smsservice service.
// All implementations must embed UnimplementedSmsserviceServer
// for forward compatibility
type SmsserviceServer interface {
	SendVerifyCode(context.Context, *SendVerifyCode_REQ) (*Default_RES, error)
	VerifyCode(context.Context, *VerifyCode_REQ) (*Default_RES, error)
	mustEmbedUnimplementedSmsserviceServer()
}

// UnimplementedSmsserviceServer must be embedded to have forward compatible implementations.
type UnimplementedSmsserviceServer struct {
}

func (UnimplementedSmsserviceServer) SendVerifyCode(context.Context, *SendVerifyCode_REQ) (*Default_RES, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendVerifyCode not implemented")
}
func (UnimplementedSmsserviceServer) VerifyCode(context.Context, *VerifyCode_REQ) (*Default_RES, error) {
	return nil, status.Errorf(codes.Unimplemented, "method VerifyCode not implemented")
}
func (UnimplementedSmsserviceServer) mustEmbedUnimplementedSmsserviceServer() {}

// UnsafeSmsserviceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SmsserviceServer will
// result in compilation errors.
type UnsafeSmsserviceServer interface {
	mustEmbedUnimplementedSmsserviceServer()
}

func RegisterSmsserviceServer(s grpc.ServiceRegistrar, srv SmsserviceServer) {
	s.RegisterService(&Smsservice_ServiceDesc, srv)
}

func _Smsservice_SendVerifyCode_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SendVerifyCode_REQ)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SmsserviceServer).SendVerifyCode(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Smsservice_SendVerifyCode_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SmsserviceServer).SendVerifyCode(ctx, req.(*SendVerifyCode_REQ))
	}
	return interceptor(ctx, in, info, handler)
}

func _Smsservice_VerifyCode_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(VerifyCode_REQ)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SmsserviceServer).VerifyCode(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Smsservice_VerifyCode_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SmsserviceServer).VerifyCode(ctx, req.(*VerifyCode_REQ))
	}
	return interceptor(ctx, in, info, handler)
}

// Smsservice_ServiceDesc is the grpc.ServiceDesc for Smsservice service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Smsservice_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "smsservice_pb.smsservice",
	HandlerType: (*SmsserviceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SendVerifyCode",
			Handler:    _Smsservice_SendVerifyCode_Handler,
		},
		{
			MethodName: "VerifyCode",
			Handler:    _Smsservice_VerifyCode_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "smsservice.proto",
}