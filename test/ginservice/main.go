package main

import (
	"zservice/zservice"
	"zservice/zservice/ex/ginservice"

	"github.com/gin-gonic/gin"
)

func init() {

	zservice.Init("ginservice.test", "1.0.0")
}

func main() {
	ginS := ginservice.NewGinService(&ginservice.GinServiceConfig{
		ListenAddr: zservice.Getenv("GIN_ADDR"),
		OnStart: func(engine *gin.Engine) {
			engine.GET("/test_auth", func(ctx *gin.Context) {

				ctx.String(200, "ok")

			})
		},
	})

	zservice.AddDependService(ginS.ZService)

	zservice.Start()

	zservice.WaitStart()
	zservice.WaitStop()
}
