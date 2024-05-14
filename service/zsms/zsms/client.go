package zsms

import (
	"context"
	_ "embed"
	"zservice/service/zsms/zsms_pb"
	"zservice/zservice"
	"zservice/zservice/ex/grpcservice"

	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
)

type ZsmsClient struct {
	conn   *grpc.ClientConn
	client zsms_pb.ZsmsClient
}

func NewZsmsClient(etcdClient *clientv3.Client) *ZsmsClient {
	conn, e := grpcservice.NewGrpcClient(&grpcservice.GrpcClientConfig{
		EtcdServiceName: "zsms",
		EtcdServer:      etcdClient,
	})
	if e != nil {
		zservice.LogPanic(e)
		return nil
	}
	return &ZsmsClient{
		conn:   conn,
		client: zsms_pb.NewZsmsClient(conn),
	}
}

func (s *ZsmsClient) SendVerifyCode(ctx *zservice.Context, req *zsms_pb.SendVerifyCode_REQ) (*zsms_pb.Default_RES, error) {
	return s.client.SendVerifyCode(context.WithValue(context.Background(), grpcservice.GRPC_contextEX_Middleware_Key, ctx.ContextS2S), req)
}

func (s *ZsmsClient) VerifyCode(ctx *zservice.Context, req *zsms_pb.VerifyCode_REQ) (*zsms_pb.Default_RES, error) {
	return s.client.VerifyCode(context.WithValue(context.Background(), grpcservice.GRPC_contextEX_Middleware_Key, ctx.ContextS2S), req)
}
