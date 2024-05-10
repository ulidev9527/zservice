package main

import (
	"time"
	"zservice/internal/etcdservice"
	"zservice/service/smsservice/smsservice"
	"zservice/service/smsservice/smsservice_pb"
	"zservice/zservice"

	clientv3 "go.etcd.io/etcd/client/v3"
)

func init() {
	zservice.Init(&zservice.ZServiceConfig{
		Name:    "smsservice.test",
		Version: "1.0.0",
	})
}
func main() {

	// func() { // 直连测试 成功
	// 	conn, err := grpc.Dial("0.0.0.0:3002", grpc.WithTransportCredentials(insecure.NewCredentials()))
	// 	if err != nil {
	// 		log.Fatalf("did not connect: %v", err)
	// 	}
	// 	defer conn.Close()

	// 	c := smsservice_pb.NewSmsserviceClient(conn)

	// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	// 	defer cancel()
	// 	res, e := c.SendVerifyCode(ctx, &smsservice_pb.SendVerifyCode_REQ{

	// 		Phone: "18888888888",
	// 	})
	// 	if e != nil {
	// 		zservice.LogError(e)
	// 	}
	// 	zservice.LogDebug(res)

	// }()

	func() { // grpc.client

		var smsserviceClient *smsservice.SmsserviceClient

		etcdS := etcdservice.NewEtcdService(&etcdservice.EtcdServiceConfig{

			Addr: zservice.Getenv("ETCD_ADDR"),
			OnStart: func(etcd *clientv3.Client) {
				// do something
			},
		})

		grcpClientS := zservice.NewService("grcpClientS", func(s *zservice.ZService) {
			smsserviceClient = smsservice.NewSmsserviceClient(etcdS.Etcd)
			s.StartDone()
		})

		zservice.AddDependService(etcdS.ZService)
		zservice.AddDependService(grcpClientS)
		grcpClientS.AddDependService(etcdS.ZService)

		zservice.Start()
		zservice.WaitStart()

		for {
			time.Sleep(time.Second * 3)

			res, e := smsserviceClient.SendVerifyCode(zservice.NewEmptyContext(), &smsservice_pb.SendVerifyCode_REQ{
				Phone: "18888888888",
			})
			if e != nil {
				zservice.LogError(e)
			}

			zservice.LogDebug(res.String())
		}

	}()

	// func() {

	// 	cli, e := clientv3.NewFromURL("10.223.223.100:8103")
	// 	if e != nil {
	// 		zservice.LogError(e)
	// 	}
	// 	etcdResolver, e := resolver.NewBuilder(cli)
	// 	if e != nil {
	// 		zservice.LogError(e)
	// 	}
	// 	conn, e := grpc.Dial("etcd:////zservice/services/smsservice",
	// 		grpc.WithTransportCredentials(insecure.NewCredentials()),
	// 		grpc.WithResolvers(etcdResolver))
	// 	if e != nil {
	// 		zservice.LogError(e)
	// 	}
	// 	zservice.LogDebug(conn)
	// 	ce := smsservice_pb.NewSmsserviceClient(conn)
	// 	cannctx := context.WithoutCancel(context.Background())

	// 	res, e := ce.SendVerifyCode(cannctx, &smsservice_pb.SendVerifyCode_REQ{
	// 		Phone: "18888888888",
	// 	})
	// 	if e != nil {
	// 		zservice.LogError(e)

	// 	}
	// 	zservice.LogDebug(res)
	// }()

	zservice.WaitStop()

}
