package main

import (
	"net/http"
	"strings"
	"zservice/service/zauth/zauth"
	"zservice/zservice"
	"zservice/zservice/ex/etcdservice"
	"zservice/zservice/ex/ginservice"

	"github.com/gin-gonic/gin"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func init() {

	zservice.Init("zauth.fileConfig", "0.1.0")
}

func main() {

	etcdS := etcdservice.NewEtcdService(&etcdservice.EtcdServiceConfig{

		Addr: zservice.Getenv("ETCD_ADDR"),
		OnStart: func(etcd *clientv3.Client) {
			// do something
		},
	})

	grpcClient := zservice.NewService("zauth.grpc", func(z *zservice.ZService) {

		zauth.Init(&zauth.ZAuthInitConfig{
			ZauthServiceName: "zauth",
			Etcd:             etcdS.Etcd,
		})
		z.StartDone()
	})

	ginS := ginservice.NewGinService(&ginservice.GinServiceConfig{
		ListenAddr: zservice.Getenv("GIN_ADDR"),
		OnStart: func(engine *gin.Engine) {
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

					e := zauth.GetFileConfig(zctx, "test.xlsx", &arr)
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

					e := zauth.GetFileConfig(zctx, "test.xlsx", &arr, zservice.StringSplit(id, ",")...)
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
					e := zauth.GetFileConfig(zctx, "test.xlsx", &m, id)
					if e != nil {
						zctx.LogError(e)
					}
					ctx.String(http.StatusOK, string(zservice.JsonMustMarshal(m)))
				}
			})
		},
	})

	zservice.AddDependService(etcdS.ZService)
	zservice.AddDependService(grpcClient)
	zservice.AddDependService(ginS.ZService)

	grpcClient.AddDependService(etcdS.ZService)
	ginS.ZService.AddDependService(etcdS.ZService)

	zservice.Start()

	zservice.WaitStop()

}
