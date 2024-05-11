package main

import (
	"zservice/service/zsms/zsms"
	"zservice/service/zsms/zsms_pb"
	"zservice/zservice"
	"zservice/zservice/ex/etcdservice"

	clientv3 "go.etcd.io/etcd/client/v3"
)

func init() {
	zservice.Init(&zservice.ZServiceConfig{
		Name:    "zsms.test",
		Version: "1.0.0",
	})
}
func main() {

	var zsmsClient *zsms.ZsmsClient

	etcdS := etcdservice.NewEtcdService(&etcdservice.EtcdServiceConfig{

		Addr: zservice.Getenv("ETCD_ADDR"),
		OnStart: func(etcd *clientv3.Client) {
			// do something
		},
	})

	grcpClientS := zservice.NewService("grcpClientS", func(s *zservice.ZService) {
		zsmsClient = zsms.NewZsmsClient(etcdS.Etcd)
		s.StartDone()
	})

	zservice.AddDependService(etcdS.ZService)
	zservice.AddDependService(grcpClientS)
	grcpClientS.AddDependService(etcdS.ZService)

	zservice.Start()
	zservice.WaitStart()

	zservice.TestAction("send sms", func() {

		res, e := zsmsClient.SendVerifyCode(zservice.NewEmptyContext(), &zsms_pb.SendVerifyCode_REQ{
			Phone: "+8618888888888",
		})
		if e != nil {
			zservice.LogError(e)
		}

		zservice.LogInfo(res.String())
	})

	zservice.WaitStop()

}
