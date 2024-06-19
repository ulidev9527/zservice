package internal

import (
	"encoding/json"
	"zservice/service/zlog/zlog_pb"
	"zservice/zservice"
	"zservice/zservice/ex/nsqservice"
	"zservice/zservice/zglobal"
)

func NsqInit() {

	// kv 日志
	go nsqservice.NewNsqConsumer(&nsqservice.NsqConsumerConfig{
		Addrs:         zservice.Getenv("NSQ_ADDRS"),
		UseNsqlookupd: zservice.GetenvBool("USE_NSQLOOKUPD"),
		Topic:         zglobal.NSQ_Topic_zlog_AddKV,
		Channel:       zservice.GetServiceName(),
		OnMessage: func(ctx *zservice.Context, body []byte) {

			bd := &zlog_pb.LogKV_REQ{}
			if e := json.Unmarshal(body, bd); e != nil {
				ctx.LogError(e)
				return
			}

			Logic_AddLogKV(ctx, bd)
		}})

}
