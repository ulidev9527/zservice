package internal

import (
	"context"
	"zservice/service/zsms/zsms_pb"
	"zservice/zservice/ex/grpcservice"
)

type ZsmsServer struct {
	zsms_pb.UnimplementedZsmsServer
}

func NewZsmsServer() *ZsmsServer {
	return &ZsmsServer{}
}

func (s *ZsmsServer) SendVerifyCode(ctx context.Context, in *zsms_pb.SendVerifyCode_REQ) (*zsms_pb.Default_RES, error) {
	return &zsms_pb.Default_RES{
		Code: SendVerifyCode(grpcservice.GetCtxEX(ctx), in),
	}, nil
}

func (s *ZsmsServer) VerifyCode(ctx context.Context, in *zsms_pb.VerifyCode_REQ) (*zsms_pb.Default_RES, error) {
	return &zsms_pb.Default_RES{
		Code: VerifyCode(grpcservice.GetCtxEX(ctx), in),
	}, nil
}
