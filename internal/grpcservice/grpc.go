package grpcservice

import (
	"context"
	"zservice/zservice"
)

func GetCtxEX(ctx context.Context) *zservice.Context {
	v := ctx.Value(GRPC_contextEX_Middleware_Key)
	if v != nil {
		return v.(*zservice.Context)
	}
	return nil
}
