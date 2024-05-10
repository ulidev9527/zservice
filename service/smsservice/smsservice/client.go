package smsservice

import (
	"context"
	_ "embed"
	"zservice/internal/grpcservice"
	"zservice/service/smsservice/smsservice_pb"
	"zservice/zservice"

	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
)

type SmsserviceClient struct {
	conn   *grpc.ClientConn
	client smsservice_pb.SmsserviceClient
}

func NewSmsserviceClient(etcdClient *clientv3.Client) *SmsserviceClient {
	conn, e := grpcservice.NewGrpcClient(&grpcservice.GrpcClientConfig{
		EtcdServiceName: "smsservice",
		EtcdServer:      etcdClient,
	})
	if e != nil {
		zservice.LogPanic(e)
		return nil
	}
	return &SmsserviceClient{
		conn:   conn,
		client: smsservice_pb.NewSmsserviceClient(conn),
	}
}

func (s *SmsserviceClient) SendVerifyCode(ctx *zservice.Context, req *smsservice_pb.SendVerifyCode_REQ) (*smsservice_pb.Default_RES, error) {
	zctx := context.WithValue(context.Background(), grpcservice.GRPC_contextEX_Middleware_Key, ctx)
	return s.client.SendVerifyCode(zctx, req)
}

func (s *SmsserviceClient) VerifyCode(ctx *zservice.Context, req *smsservice_pb.VerifyCode_REQ) (*smsservice_pb.Default_RES, error) {
	zctx := context.WithValue(context.Background(), grpcservice.GRPC_contextEX_Middleware_Key, ctx)
	return s.client.VerifyCode(zctx, req)
}
