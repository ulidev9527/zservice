package main

import (
	"zservice/service/zconfig/zconfig"
	"zservice/service/zconfig/zconfig_pb"
	"zservice/zglobal"
	"zservice/zservice"
	"zservice/zservice/ex/etcdservice"

	clientv3 "go.etcd.io/etcd/client/v3"
)

func init() {

	zservice.Init(&zservice.ZServiceConfig{
		Name:    "zconfig.fileConfig",
		Version: "0.1.0",
	})
}

func main() {

	etcdS := etcdservice.NewEtcdService(&etcdservice.EtcdServiceConfig{

		Addr: zservice.Getenv("ETCD_ADDR"),
		OnStart: func(etcd *clientv3.Client) {
			// do something
		},
	})

	grpcClient := zservice.NewService("zconfig.grpc", func(z *zservice.ZService) {
		zc := zconfig.NewZconfigClient(etcdS.Etcd)

		zservice.TestAction("GetFileConfig-all", func() {

			res, e := zc.GetFileConfig(zservice.NewEmptyContext(), &zconfig_pb.GetFileConfig_REQ{
				Parser:   zglobal.E_ZConfig_Parser_Excel,
				FileName: "test.xlsx",
			})
			if e != nil {
				z.LogError(e)
			} else {
				z.LogInfo(res)
			}
		})

		zservice.TestAction("getFileConfig-byID", func() {
			res, e := zc.GetFileConfig(zservice.NewEmptyContext(), &zconfig_pb.GetFileConfig_REQ{
				Parser:   zglobal.E_ZConfig_Parser_Excel,
				FileName: "test.xlsx",
				Keys:     "1,3,5,",
			})
			if e != nil {
				z.LogError(e)
			} else {
				z.LogInfo(res)
			}
		})

		zservice.TestAction("getFileConfig-one", func() {
			res, e := zc.GetFileConfig(zservice.NewEmptyContext(), &zconfig_pb.GetFileConfig_REQ{
				Parser:   zglobal.E_ZConfig_Parser_Excel,
				FileName: "test.xlsx",
				Keys:     "1",
			})
			if e != nil {
				z.LogError(e)
			} else {
				z.LogInfo(res)
			}
		})

		zservice.TestAction("getFileConfig-one empty", func() {
			res, e := zc.GetFileConfig(zservice.NewEmptyContext(), &zconfig_pb.GetFileConfig_REQ{
				Parser:   zglobal.E_ZConfig_Parser_Excel,
				FileName: "test.xlsx",
				Keys:     "18",
			})
			if e != nil {
				z.LogError(e)
			} else {
				z.LogInfo(res)
			}
		})

	})

	zservice.AddDependService(etcdS.ZService)
	zservice.AddDependService(grpcClient)

	grpcClient.AddDependService(etcdS.ZService)

	zservice.Start()

	zservice.WaitStop()

}
