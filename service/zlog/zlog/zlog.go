package zlog

import (
	"zservice/service/zlog/zlog_pb"
	"zservice/zservice"
	"zservice/zservice/ex/nsqservice"
	"zservice/zservice/zglobal"
)

var nsqPService *nsqservice.NsqProducerService

type ZlogInitConfig struct {
	NsqProducerService *nsqservice.NsqProducerService
}

func Init(c *ZlogInitConfig) {
	nsqPService = c.NsqProducerService
}

func LogKV(ctx *zservice.Context, in *zlog_pb.LogKV_REQ) *zlog_pb.Default_RES {
	if e := nsqPService.Publish(ctx, zglobal.NSQ_Topic_zlog_AddKV, zservice.JsonMustMarshal(in)); e != nil {
		ctx.LogError(e)
		return &zlog_pb.Default_RES{Code: e.GetCode()}
	}
	return &zlog_pb.Default_RES{Code: zglobal.Code_SUCC}
}
