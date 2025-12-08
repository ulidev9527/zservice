package etcdservice

import (
	"context"
	"time"
	"zserviceapps/packages/zservice"
)

// 测试所有节点连通性
// @return 失败的节点
func (ser *Service) Ping(ctx *zservice.Context) []string {

	fail := []string{}

	c, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	for _, endpoint := range ser.Client.Endpoints() {
		if res, e := ser.Client.Status(c, endpoint); e != nil {
			ctx.LogError("ping fail", endpoint, e)
			fail = append(fail, endpoint)
		} else {
			ctx.LogInfo("ping succ", endpoint, zservice.JsonMustMarshalString(res))
		}
	}

	return fail

}
