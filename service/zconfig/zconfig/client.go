package zconfig

import (
	"context"
	"zservice/service/zconfig/zconfig_pb"
	"zservice/zservice"
	"zservice/zservice/ex/grpcservice"

	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
)

type ZconfigClient struct {
	conn   *grpc.ClientConn
	client zconfig_pb.ZconfigClient
}

func NewZconfigClient(etcdClient *clientv3.Client) *ZconfigClient {
	conn, e := grpcservice.NewGrpcClient(&grpcservice.GrpcClientConfig{
		EtcdServiceName: "zconfig",
		EtcdServer:      etcdClient,
	})
	if e != nil {
		zservice.LogPanic(e)
		return nil
	}
	return &ZconfigClient{
		conn:   conn,
		client: zconfig_pb.NewZconfigClient(conn),
	}
}

func (s *ZconfigClient) GetFileConfig(ctx *zservice.Context, req *zconfig_pb.GetFileConfig_REQ) (*zconfig_pb.GetFileConfig_RES, error) {
	return s.client.GetFileConfig(context.WithValue(context.Background(), grpcservice.GRPC_contextEX_Middleware_Key, ctx), req)
}
