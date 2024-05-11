package internal

import (
	"context"
	"zservice/service/zconfig/zconfig_pb"
	"zservice/zservice/ex/grpcservice"
)

type ZconfigServer struct {
	zconfig_pb.UnimplementedZconfigServer
}

func NewZconfigServer() *ZconfigServer {
	return &ZconfigServer{}
}

func (s *ZconfigServer) GetFileConfig(ctx context.Context, in *zconfig_pb.GetFileConfig_REQ) (*zconfig_pb.GetFileConfig_RES, error) {
	code, val := GetFileConfig(grpcservice.GetCtxEX(ctx), in)
	return &zconfig_pb.GetFileConfig_RES{
		Code:  code,
		Value: val,
	}, nil
}
