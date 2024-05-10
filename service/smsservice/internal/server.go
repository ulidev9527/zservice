package internal

import (
	"context"
	"zservice/service/smsservice/smsservice_pb"
)

type SmsserviceServer struct {
	smsservice_pb.UnimplementedSmsserviceServer
}

func NewSmsserviceServer() *SmsserviceServer {
	return &SmsserviceServer{}
}

func (s *SmsserviceServer) SendVerifyCode(ctx context.Context, in *smsservice_pb.SendVerifyCode_REQ) (*smsservice_pb.Default_RES, error) {
	// zctx := grpcservice.GetCtxEX(ctx)
	return &smsservice_pb.Default_RES{
		Code: 200,
	}, nil
}

func (s *SmsserviceServer) VerifyCode(ctx context.Context, in *smsservice_pb.VerifyCode_REQ) (*smsservice_pb.Default_RES, error) {
	return &smsservice_pb.Default_RES{
		Code: 200,
	}, nil
}
