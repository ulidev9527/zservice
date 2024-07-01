package main

import (
	"net/http"
	"strings"
	"zservice/service/zauth/zauth"
	"zservice/service/zauth/zauth_ex"
	"zservice/service/zauth/zauth_pb"
	"zservice/zservice"
	"zservice/zservice/ex/etcdservice"
	"zservice/zservice/ex/ginservice"

	"github.com/gin-gonic/gin"
)

func init() {

	zservice.Init("zauth.test", "0.1.0")
}

func main() {

	etcdS := etcdservice.NewEtcdService(&etcdservice.EtcdServiceConfig{

		Addr: zservice.Getenv("ETCD_ADDR"),
		OnStart: func(s *etcdservice.EtcdService) {
			// do something
		},
	})

	grpcClient := zservice.NewService("zauth.grpc", func(z *zservice.ZService) {
		ctx := zservice.NewContext()
		zauth.Init(&zauth.ZAuthInitConfig{
			ServiceName: zservice.Getenv("ZAUTH_SERVICE_NAME"),
			Etcd:        etcdS.Etcd,
			GrpcAddr:    zservice.Getenv("zauth_grpc_addr"),
			UseGrpcEtcd: zservice.GetenvBool("USE_GRPC_ETCD"),
		})

		zauth_ex.ServiceInfo.Regist(ctx, &zauth_pb.ServiceRegist_REQ{})
		z.StartDone()
	})
	grpcClient.AddDependService(etcdS.ZService)

	ginS := ginservice.NewGinService(&ginservice.GinServiceConfig{
		ListenPort: zservice.Getenv("GIN_PORT"),
		OnStart: func(s *ginservice.GinService) {
			engine := s.Engine
			engine.GET("/", func(ctx *gin.Context) {
				zctx := ginservice.GetCtxEX(ctx)
				id := ctx.Query("id")
				if id == "" {

					arr := []struct {
						ID         string `json:"id"`
						Name       string `json:"name"`
						Desc       string `json:"desc"`
						Icon       string `json:"icon"`
						LimitCount uint32 `json:"limit_count"`
					}{}

					e := zauth.ConfigGetFileConfig(zctx, "test.xlsx", &arr)
					if e != nil {
						zctx.LogError(e)
					}
					ctx.String(http.StatusOK, string(zservice.JsonMustMarshal(arr)))
				} else if strings.Contains(id, ",") {
					arr := []struct {
						ID         string `json:"id"`
						Name       string `json:"name"`
						Desc       string `json:"desc"`
						Icon       string `json:"icon"`
						LimitCount uint32 `json:"limit_count"`
					}{}

					e := zauth.ConfigGetFileConfig(zctx, "test.xlsx", &arr, zservice.StringSplit(id, ",")...)
					if e != nil {
						zctx.LogError(e)
					}
					ctx.String(http.StatusOK, string(zservice.JsonMustMarshal(arr)))
				} else {

					m := struct {
						ID         int    `json:"id"`
						Name       string `json:"name"`
						Desc       string `json:"desc"`
						Icon       string `json:"icon"`
						LimitCount uint32 `json:"limit_count"`
					}{}
					e := zauth.ConfigGetFileConfig(zctx, "test.xlsx", &m, id)
					if e != nil {
						zctx.LogError(e)
					}
					ctx.String(http.StatusOK, string(zservice.JsonMustMarshal(m)))
				}
			})
		},
	})

	ginS.ZService.AddDependService(grpcClient, zservice.NewService("init", func(z *zservice.ZService) {

		z.StartDone()

	}))

	zservice.AddDependService(ginS.ZService)

	zservice.Start().WaitStart().WaitStop()

}
