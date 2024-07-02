package zauth

import (
	"fmt"
	"zservice/service/zauth/zauth_pb"
	"zservice/zservice"
	"zservice/zservice/ex/etcdservice"
	"zservice/zservice/ex/grpcservice"
	"zservice/zservice/zglobal"
)

var grpcClient zauth_pb.ZauthClient
var zauthInitConfig *ZAuthInitOption

type ZAuthInitOption struct {
	ServiceName string // 权限服务名称
	EtcdService *etcdservice.EtcdService
	GrpcAddr    string // rpc addr
}

func Init(opt *ZAuthInitOption) {
	zauthInitConfig = opt
	if conn, e := grpcservice.NewGrpcClient(&grpcservice.GrpcClientOption{
		ServiceName: opt.ServiceName,
		EtcdClient:  opt.EtcdService.EtcdClient,
		Addr:        opt.GrpcAddr,
	}); e != nil {
		zservice.LogPanic(e)
	} else {
		grpcClient = zauth_pb.NewZauthClient(conn)
	}

	// 服务配置改变监听
	opt.EtcdService.WatchEvent(fmt.Sprintf(zglobal.EV_Config_ServiceFileConfigChange, zservice.GetServiceName()), func(ctx *zservice.Context, val string) {
		fileConfigCache.Delete(val)
	})
}
func GetGrpcClient() zauth_pb.ZauthClient {
	return grpcClient
}
func GetGrpcServiceName() string {
	return zauthInitConfig.ServiceName
}
