// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.25.1
// source: zauth.proto

package zauth_pb

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
	Zauth_Logout_FullMethodName              = "/zauth_pb.zauth/Logout"
	Zauth_LoginByPhone_FullMethodName        = "/zauth_pb.zauth/LoginByPhone"
	Zauth_LoginByAccount_FullMethodName      = "/zauth_pb.zauth/LoginByAccount"
	Zauth_PermissionCreate_FullMethodName    = "/zauth_pb.zauth/PermissionCreate"
	Zauth_PermissionListGet_FullMethodName   = "/zauth_pb.zauth/PermissionListGet"
	Zauth_PermissionUpdate_FullMethodName    = "/zauth_pb.zauth/PermissionUpdate"
	Zauth_PermissionBind_FullMethodName      = "/zauth_pb.zauth/PermissionBind"
	Zauth_OrgCreate_FullMethodName           = "/zauth_pb.zauth/OrgCreate"
	Zauth_OrgListGet_FullMethodName          = "/zauth_pb.zauth/OrgListGet"
	Zauth_OrgUpdate_FullMethodName           = "/zauth_pb.zauth/OrgUpdate"
	Zauth_SMSVerifyCodeSend_FullMethodName   = "/zauth_pb.zauth/SMSVerifyCodeSend"
	Zauth_SMSVerifyCodeVerify_FullMethodName = "/zauth_pb.zauth/SMSVerifyCodeVerify"
	Zauth_CheckAuth_FullMethodName           = "/zauth_pb.zauth/CheckAuth"
	Zauth_GetFileConfig_FullMethodName       = "/zauth_pb.zauth/GetFileConfig"
)

// ZauthClient is the client API for Zauth service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ZauthClient interface {
	Logout(ctx context.Context, in *Default_REQ, opts ...grpc.CallOption) (*Default_RES, error)
	LoginByPhone(ctx context.Context, in *LoginByPhone_REQ, opts ...grpc.CallOption) (*Default_RES, error)
	LoginByAccount(ctx context.Context, in *LoginByAccount_REQ, opts ...grpc.CallOption) (*Default_RES, error)
	PermissionCreate(ctx context.Context, in *PermissionInfo, opts ...grpc.CallOption) (*PermissionInfo_RES, error)
	PermissionListGet(ctx context.Context, in *PermissionListGet_REQ, opts ...grpc.CallOption) (*PermissionInfoList_RES, error)
	PermissionUpdate(ctx context.Context, in *PermissionInfo, opts ...grpc.CallOption) (*PermissionInfo_RES, error)
	PermissionBind(ctx context.Context, in *PermissionBind_REQ, opts ...grpc.CallOption) (*Default_RES, error)
	OrgCreate(ctx context.Context, in *OrgInfo, opts ...grpc.CallOption) (*OrgInfo_RES, error)
	OrgListGet(ctx context.Context, in *OrgListGet_REQ, opts ...grpc.CallOption) (*OrgInfoList_RES, error)
	OrgUpdate(ctx context.Context, in *OrgInfo, opts ...grpc.CallOption) (*OrgInfo_RES, error)
	SMSVerifyCodeSend(ctx context.Context, in *SMSVerifyCodeSend_REQ, opts ...grpc.CallOption) (*SMSSendVerifyCode_RES, error)
	SMSVerifyCodeVerify(ctx context.Context, in *SMSVerifyCodeVerify_REQ, opts ...grpc.CallOption) (*Default_RES, error)
	CheckAuth(ctx context.Context, in *CheckAuth_REQ, opts ...grpc.CallOption) (*CheckAuth_RES, error)
	GetFileConfig(ctx context.Context, in *GetFileConfig_REQ, opts ...grpc.CallOption) (*GetFileConfig_RES, error)
}

type zauthClient struct {
	cc grpc.ClientConnInterface
}

func NewZauthClient(cc grpc.ClientConnInterface) ZauthClient {
	return &zauthClient{cc}
}

func (c *zauthClient) Logout(ctx context.Context, in *Default_REQ, opts ...grpc.CallOption) (*Default_RES, error) {
	out := new(Default_RES)
	err := c.cc.Invoke(ctx, Zauth_Logout_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *zauthClient) LoginByPhone(ctx context.Context, in *LoginByPhone_REQ, opts ...grpc.CallOption) (*Default_RES, error) {
	out := new(Default_RES)
	err := c.cc.Invoke(ctx, Zauth_LoginByPhone_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *zauthClient) LoginByAccount(ctx context.Context, in *LoginByAccount_REQ, opts ...grpc.CallOption) (*Default_RES, error) {
	out := new(Default_RES)
	err := c.cc.Invoke(ctx, Zauth_LoginByAccount_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *zauthClient) PermissionCreate(ctx context.Context, in *PermissionInfo, opts ...grpc.CallOption) (*PermissionInfo_RES, error) {
	out := new(PermissionInfo_RES)
	err := c.cc.Invoke(ctx, Zauth_PermissionCreate_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *zauthClient) PermissionListGet(ctx context.Context, in *PermissionListGet_REQ, opts ...grpc.CallOption) (*PermissionInfoList_RES, error) {
	out := new(PermissionInfoList_RES)
	err := c.cc.Invoke(ctx, Zauth_PermissionListGet_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *zauthClient) PermissionUpdate(ctx context.Context, in *PermissionInfo, opts ...grpc.CallOption) (*PermissionInfo_RES, error) {
	out := new(PermissionInfo_RES)
	err := c.cc.Invoke(ctx, Zauth_PermissionUpdate_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *zauthClient) PermissionBind(ctx context.Context, in *PermissionBind_REQ, opts ...grpc.CallOption) (*Default_RES, error) {
	out := new(Default_RES)
	err := c.cc.Invoke(ctx, Zauth_PermissionBind_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *zauthClient) OrgCreate(ctx context.Context, in *OrgInfo, opts ...grpc.CallOption) (*OrgInfo_RES, error) {
	out := new(OrgInfo_RES)
	err := c.cc.Invoke(ctx, Zauth_OrgCreate_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *zauthClient) OrgListGet(ctx context.Context, in *OrgListGet_REQ, opts ...grpc.CallOption) (*OrgInfoList_RES, error) {
	out := new(OrgInfoList_RES)
	err := c.cc.Invoke(ctx, Zauth_OrgListGet_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *zauthClient) OrgUpdate(ctx context.Context, in *OrgInfo, opts ...grpc.CallOption) (*OrgInfo_RES, error) {
	out := new(OrgInfo_RES)
	err := c.cc.Invoke(ctx, Zauth_OrgUpdate_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *zauthClient) SMSVerifyCodeSend(ctx context.Context, in *SMSVerifyCodeSend_REQ, opts ...grpc.CallOption) (*SMSSendVerifyCode_RES, error) {
	out := new(SMSSendVerifyCode_RES)
	err := c.cc.Invoke(ctx, Zauth_SMSVerifyCodeSend_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *zauthClient) SMSVerifyCodeVerify(ctx context.Context, in *SMSVerifyCodeVerify_REQ, opts ...grpc.CallOption) (*Default_RES, error) {
	out := new(Default_RES)
	err := c.cc.Invoke(ctx, Zauth_SMSVerifyCodeVerify_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *zauthClient) CheckAuth(ctx context.Context, in *CheckAuth_REQ, opts ...grpc.CallOption) (*CheckAuth_RES, error) {
	out := new(CheckAuth_RES)
	err := c.cc.Invoke(ctx, Zauth_CheckAuth_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *zauthClient) GetFileConfig(ctx context.Context, in *GetFileConfig_REQ, opts ...grpc.CallOption) (*GetFileConfig_RES, error) {
	out := new(GetFileConfig_RES)
	err := c.cc.Invoke(ctx, Zauth_GetFileConfig_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ZauthServer is the server API for Zauth service.
// All implementations must embed UnimplementedZauthServer
// for forward compatibility
type ZauthServer interface {
	Logout(context.Context, *Default_REQ) (*Default_RES, error)
	LoginByPhone(context.Context, *LoginByPhone_REQ) (*Default_RES, error)
	LoginByAccount(context.Context, *LoginByAccount_REQ) (*Default_RES, error)
	PermissionCreate(context.Context, *PermissionInfo) (*PermissionInfo_RES, error)
	PermissionListGet(context.Context, *PermissionListGet_REQ) (*PermissionInfoList_RES, error)
	PermissionUpdate(context.Context, *PermissionInfo) (*PermissionInfo_RES, error)
	PermissionBind(context.Context, *PermissionBind_REQ) (*Default_RES, error)
	OrgCreate(context.Context, *OrgInfo) (*OrgInfo_RES, error)
	OrgListGet(context.Context, *OrgListGet_REQ) (*OrgInfoList_RES, error)
	OrgUpdate(context.Context, *OrgInfo) (*OrgInfo_RES, error)
	SMSVerifyCodeSend(context.Context, *SMSVerifyCodeSend_REQ) (*SMSSendVerifyCode_RES, error)
	SMSVerifyCodeVerify(context.Context, *SMSVerifyCodeVerify_REQ) (*Default_RES, error)
	CheckAuth(context.Context, *CheckAuth_REQ) (*CheckAuth_RES, error)
	GetFileConfig(context.Context, *GetFileConfig_REQ) (*GetFileConfig_RES, error)
	mustEmbedUnimplementedZauthServer()
}

// UnimplementedZauthServer must be embedded to have forward compatible implementations.
type UnimplementedZauthServer struct {
}

func (UnimplementedZauthServer) Logout(context.Context, *Default_REQ) (*Default_RES, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Logout not implemented")
}
func (UnimplementedZauthServer) LoginByPhone(context.Context, *LoginByPhone_REQ) (*Default_RES, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LoginByPhone not implemented")
}
func (UnimplementedZauthServer) LoginByAccount(context.Context, *LoginByAccount_REQ) (*Default_RES, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LoginByAccount not implemented")
}
func (UnimplementedZauthServer) PermissionCreate(context.Context, *PermissionInfo) (*PermissionInfo_RES, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PermissionCreate not implemented")
}
func (UnimplementedZauthServer) PermissionListGet(context.Context, *PermissionListGet_REQ) (*PermissionInfoList_RES, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PermissionListGet not implemented")
}
func (UnimplementedZauthServer) PermissionUpdate(context.Context, *PermissionInfo) (*PermissionInfo_RES, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PermissionUpdate not implemented")
}
func (UnimplementedZauthServer) PermissionBind(context.Context, *PermissionBind_REQ) (*Default_RES, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PermissionBind not implemented")
}
func (UnimplementedZauthServer) OrgCreate(context.Context, *OrgInfo) (*OrgInfo_RES, error) {
	return nil, status.Errorf(codes.Unimplemented, "method OrgCreate not implemented")
}
func (UnimplementedZauthServer) OrgListGet(context.Context, *OrgListGet_REQ) (*OrgInfoList_RES, error) {
	return nil, status.Errorf(codes.Unimplemented, "method OrgListGet not implemented")
}
func (UnimplementedZauthServer) OrgUpdate(context.Context, *OrgInfo) (*OrgInfo_RES, error) {
	return nil, status.Errorf(codes.Unimplemented, "method OrgUpdate not implemented")
}
func (UnimplementedZauthServer) SMSVerifyCodeSend(context.Context, *SMSVerifyCodeSend_REQ) (*SMSSendVerifyCode_RES, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SMSVerifyCodeSend not implemented")
}
func (UnimplementedZauthServer) SMSVerifyCodeVerify(context.Context, *SMSVerifyCodeVerify_REQ) (*Default_RES, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SMSVerifyCodeVerify not implemented")
}
func (UnimplementedZauthServer) CheckAuth(context.Context, *CheckAuth_REQ) (*CheckAuth_RES, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CheckAuth not implemented")
}
func (UnimplementedZauthServer) GetFileConfig(context.Context, *GetFileConfig_REQ) (*GetFileConfig_RES, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetFileConfig not implemented")
}
func (UnimplementedZauthServer) mustEmbedUnimplementedZauthServer() {}

// UnsafeZauthServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ZauthServer will
// result in compilation errors.
type UnsafeZauthServer interface {
	mustEmbedUnimplementedZauthServer()
}

func RegisterZauthServer(s grpc.ServiceRegistrar, srv ZauthServer) {
	s.RegisterService(&Zauth_ServiceDesc, srv)
}

func _Zauth_Logout_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Default_REQ)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ZauthServer).Logout(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Zauth_Logout_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ZauthServer).Logout(ctx, req.(*Default_REQ))
	}
	return interceptor(ctx, in, info, handler)
}

func _Zauth_LoginByPhone_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LoginByPhone_REQ)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ZauthServer).LoginByPhone(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Zauth_LoginByPhone_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ZauthServer).LoginByPhone(ctx, req.(*LoginByPhone_REQ))
	}
	return interceptor(ctx, in, info, handler)
}

func _Zauth_LoginByAccount_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LoginByAccount_REQ)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ZauthServer).LoginByAccount(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Zauth_LoginByAccount_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ZauthServer).LoginByAccount(ctx, req.(*LoginByAccount_REQ))
	}
	return interceptor(ctx, in, info, handler)
}

func _Zauth_PermissionCreate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PermissionInfo)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ZauthServer).PermissionCreate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Zauth_PermissionCreate_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ZauthServer).PermissionCreate(ctx, req.(*PermissionInfo))
	}
	return interceptor(ctx, in, info, handler)
}

func _Zauth_PermissionListGet_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PermissionListGet_REQ)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ZauthServer).PermissionListGet(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Zauth_PermissionListGet_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ZauthServer).PermissionListGet(ctx, req.(*PermissionListGet_REQ))
	}
	return interceptor(ctx, in, info, handler)
}

func _Zauth_PermissionUpdate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PermissionInfo)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ZauthServer).PermissionUpdate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Zauth_PermissionUpdate_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ZauthServer).PermissionUpdate(ctx, req.(*PermissionInfo))
	}
	return interceptor(ctx, in, info, handler)
}

func _Zauth_PermissionBind_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PermissionBind_REQ)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ZauthServer).PermissionBind(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Zauth_PermissionBind_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ZauthServer).PermissionBind(ctx, req.(*PermissionBind_REQ))
	}
	return interceptor(ctx, in, info, handler)
}

func _Zauth_OrgCreate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(OrgInfo)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ZauthServer).OrgCreate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Zauth_OrgCreate_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ZauthServer).OrgCreate(ctx, req.(*OrgInfo))
	}
	return interceptor(ctx, in, info, handler)
}

func _Zauth_OrgListGet_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(OrgListGet_REQ)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ZauthServer).OrgListGet(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Zauth_OrgListGet_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ZauthServer).OrgListGet(ctx, req.(*OrgListGet_REQ))
	}
	return interceptor(ctx, in, info, handler)
}

func _Zauth_OrgUpdate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(OrgInfo)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ZauthServer).OrgUpdate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Zauth_OrgUpdate_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ZauthServer).OrgUpdate(ctx, req.(*OrgInfo))
	}
	return interceptor(ctx, in, info, handler)
}

func _Zauth_SMSVerifyCodeSend_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SMSVerifyCodeSend_REQ)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ZauthServer).SMSVerifyCodeSend(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Zauth_SMSVerifyCodeSend_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ZauthServer).SMSVerifyCodeSend(ctx, req.(*SMSVerifyCodeSend_REQ))
	}
	return interceptor(ctx, in, info, handler)
}

func _Zauth_SMSVerifyCodeVerify_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SMSVerifyCodeVerify_REQ)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ZauthServer).SMSVerifyCodeVerify(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Zauth_SMSVerifyCodeVerify_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ZauthServer).SMSVerifyCodeVerify(ctx, req.(*SMSVerifyCodeVerify_REQ))
	}
	return interceptor(ctx, in, info, handler)
}

func _Zauth_CheckAuth_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CheckAuth_REQ)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ZauthServer).CheckAuth(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Zauth_CheckAuth_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ZauthServer).CheckAuth(ctx, req.(*CheckAuth_REQ))
	}
	return interceptor(ctx, in, info, handler)
}

func _Zauth_GetFileConfig_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetFileConfig_REQ)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ZauthServer).GetFileConfig(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Zauth_GetFileConfig_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ZauthServer).GetFileConfig(ctx, req.(*GetFileConfig_REQ))
	}
	return interceptor(ctx, in, info, handler)
}

// Zauth_ServiceDesc is the grpc.ServiceDesc for Zauth service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Zauth_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "zauth_pb.zauth",
	HandlerType: (*ZauthServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Logout",
			Handler:    _Zauth_Logout_Handler,
		},
		{
			MethodName: "LoginByPhone",
			Handler:    _Zauth_LoginByPhone_Handler,
		},
		{
			MethodName: "LoginByAccount",
			Handler:    _Zauth_LoginByAccount_Handler,
		},
		{
			MethodName: "PermissionCreate",
			Handler:    _Zauth_PermissionCreate_Handler,
		},
		{
			MethodName: "PermissionListGet",
			Handler:    _Zauth_PermissionListGet_Handler,
		},
		{
			MethodName: "PermissionUpdate",
			Handler:    _Zauth_PermissionUpdate_Handler,
		},
		{
			MethodName: "PermissionBind",
			Handler:    _Zauth_PermissionBind_Handler,
		},
		{
			MethodName: "OrgCreate",
			Handler:    _Zauth_OrgCreate_Handler,
		},
		{
			MethodName: "OrgListGet",
			Handler:    _Zauth_OrgListGet_Handler,
		},
		{
			MethodName: "OrgUpdate",
			Handler:    _Zauth_OrgUpdate_Handler,
		},
		{
			MethodName: "SMSVerifyCodeSend",
			Handler:    _Zauth_SMSVerifyCodeSend_Handler,
		},
		{
			MethodName: "SMSVerifyCodeVerify",
			Handler:    _Zauth_SMSVerifyCodeVerify_Handler,
		},
		{
			MethodName: "CheckAuth",
			Handler:    _Zauth_CheckAuth_Handler,
		},
		{
			MethodName: "GetFileConfig",
			Handler:    _Zauth_GetFileConfig_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "zauth.proto",
}
