package zlog

import (
	"time"
	"zservice/service/zlog/zlog_pb"
	"zservice/zservice"
	"zservice/zservice/zglobal"
)

func LogKV(ctx *zservice.Context, in *zlog_pb.LogKV_REQ) *zlog_pb.Default_RES {
	in.SaveTime = time.Now().UnixMilli()
	if e := nsqPService.Publish(ctx, zglobal.NSQ_Topic_zlog_AddKV, zservice.JsonMustMarshal(in)); e != nil {
		ctx.LogError(e)
		return &zlog_pb.Default_RES{Code: e.GetCode()}
	}
	return &zlog_pb.Default_RES{Code: zglobal.Code_SUCC}
}
