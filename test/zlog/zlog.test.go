package main

import (
	"zservice/service/zlog/zlog"
	"zservice/service/zlog/zlog_pb"
	"zservice/zservice"
	"zservice/zservice/ex/nsqservice"

	"github.com/nsqio/go-nsq"
)

func init() {
	zservice.Init("zlog_test", "1.0.0")
}

func main() {

	nsqPS := nsqservice.NewNsqProducerService(&nsqservice.NsqProducerServiceConfig{
		Addr: zservice.Getenv("NSQ_Producer_ADDR"),
		OnStart: func(p *nsq.Producer) {
		},
	})

	zservice.AddDependService(nsqPS.ZService)

	zservice.Start().WaitStart()

	zlog.Init(&zlog.ZlogInitConfig{
		NsqProducerService: nsqPS,
	})

	zlog.LogKV(zservice.NewContext(), &zlog_pb.LogKV_REQ{
		Key:   "zlog_test",
		Value: "run",
	})

	zservice.WaitStop()

}
