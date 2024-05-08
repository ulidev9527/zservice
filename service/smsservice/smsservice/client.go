package smsservice

import (
	"zservice/service/smsservice/smsservice_pb"
	"zservice/zservice"
)

type SmsserviceClient struct {
}

func NewSmsserviceClient() *SmsserviceClient {
	return &SmsserviceClient{}
}

func (s *SmsserviceClient) SendVerifyCode(ctx *zservice.ZContext, req *smsservice_pb.SendVerifyCode_REQ) (*smsservice_pb.Default_RES, error) {
	return &smsservice_pb.Default_RES{}, nil
}

func (s *SmsserviceClient) VerifyCode(ctx *zservice.ZContext, req *smsservice_pb.VerifyCode_REQ) (*smsservice_pb.Default_RES, error) {
	return &smsservice_pb.Default_RES{}, nil
}
