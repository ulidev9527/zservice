package internal

import (
	"time"
	"zservice/service/zlog/zlog_pb"
	"zservice/zservice"
)

func Logic_AddLogKV(ctx *zservice.Context, in *zlog_pb.LogKV_REQ) *zlog_pb.Default_RES {

	tab := &LogKVTable{
		Key:      in.Key,
		Value:    in.Value,
		SaveTime: time.UnixMilli(in.SaveTime),
		TraceID:  ctx.TraceID,
		Service:  ctx.TraceService,
	}

	if e := Gorm.Create(tab).Error; e != nil {
		zservice.LogError(e)
		return &zlog_pb.Default_RES{Code: zservice.Code_DB_SaveFail}
	}
	return &zlog_pb.Default_RES{Code: zservice.Code_SUCC}
}
