package zauth

import (
	"context"
	"zservice/service/zauth/zauth_pb"
	"zservice/zservice"
	"zservice/zservice/ex/grpcservice"
	"zservice/zservice/ex/redisservice"
	"zservice/zservice/zglobal"

	clientv3 "go.etcd.io/etcd/client/v3"
)

var grpcClient zauth_pb.ZauthClient

type ZAuthConfig struct {
	Etcd            *clientv3.Client
	Redis           *redisservice.GoRedisEX
	NsqConsumerAddr string // nsq consumer addr
	IsNsqd          bool
}

func Init(c *ZAuthConfig) {
	func() {
		conn, e := grpcservice.NewGrpcClient(&grpcservice.GrpcClientConfig{
			EtcdServiceName: "zauth",
			EtcdServer:      c.Etcd,
		})
		if e != nil {
			zservice.LogPanic(e)
			return
		}

		grpcClient = zauth_pb.NewZauthClient(conn)
	}()

}

// 检查权限, 没返回错误表示检查成功
func CheckAuth(ctx *zservice.Context, req *zauth_pb.CheckAuth_REQ) *zservice.Error {
	if res, e := grpcClient.CheckAuth(context.WithValue(context.Background(), grpcservice.GRPC_contextEX_Middleware_Key, ctx.ContextS2S), req); e != nil {
		return zservice.NewError(e).SetCode(zglobal.Code_ErrorBreakoff)
	} else if res.Code != zglobal.Code_SUCC {
		return zservice.NewError("check auth fail").SetCode(res.Code)
	} else {
		return nil
	}
}
