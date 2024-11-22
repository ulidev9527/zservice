package main

import (
	"github.com/ulidev9527/zservice/zserviceex/ginservice"

	"github.com/ulidev9527/zservice/zservice"

	"github.com/gin-gonic/gin"
)

func main() {

	zservice.Init(zservice.ZserviceOption{

		Name:    "ginservice.test",
		Version: "1.0.0",
	})

	ginS := ginservice.NewGinService(ginservice.GinServiceConfig{
		Port: 8811,
		OnStart: func(s *ginservice.GinService) {

			s.Engine.GET("/test_auth", func(ctx *gin.Context) {

				ctx.String(200, "ok")

			})

			s.Engine.GET("/test", func(ctx *gin.Context) {

				ctx.String(200, "ok")
			})
			s.Engine.GET("/Test", func(ctx *gin.Context) {
				ctx.String(200, "OK___")
			})
		},
	})

	zservice.AddDependService(ginS.ZService)

	zservice.Start()

	zservice.WaitStart()
	zservice.WaitStop()
}
