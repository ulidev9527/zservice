package internal

import (
	"context"
	"zservice/service/zauth/zauth_pb"
	"zservice/zservice/ex/grpcservice"
)

type ZauthServer struct {
	zauth_pb.UnimplementedZauthServer
}

// 手机号登陆
func (s *ZauthServer) LoginByPhone(ctx context.Context, in *zauth_pb.LoginByPhone_REQ) (*zauth_pb.Default_RES, error) {
	return LoginByPhone(grpcservice.GetCtxEX(ctx), in), nil
}

func (s *ZauthServer) CheckAuth(ctx context.Context, in *zauth_pb.CheckAuth_REQ) (*zauth_pb.CheckAuth_RES, error) {
	return CheckAuth(grpcservice.GetCtxEX(ctx), in), nil
}
