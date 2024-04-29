package main

import (
	"zservice/zservice"

	"github.com/gin-gonic/gin"
)

var DBService = zservice.NewService(&zservice.ZServiceConfig{
	Name: "DBService",
	OnBeforeStart: func(s *zservice.ZService) {
		zservice.LogDebug("DBService")
	},
	OnStart: func(s *zservice.ZService) {
		s.StartDone()
	},
})

var Gineng *gin.Engine
var GinService = zservice.NewService(&zservice.ZServiceConfig{
	Name: "GinService",
	OnBeforeStart: func(s *zservice.ZService) {

	},
	OnStart: func(s *zservice.ZService) {

		Gineng = gin.New()
		Gineng.Use(zservice.GinCORSMiddleware())
		Gineng.Use(zservice.GinContextEXTMiddleware(s))

		Gineng.GET("/ping", func(c *gin.Context) {
			c.String(200, "pong")
		})

		Gineng.Run(":3000")
		s.StartDone()
	},
})

func main() {

	service := zservice.NewService(&zservice.ZServiceConfig{Name: "TestService"})

	service.AddService(DBService)
	service.AddService(GinService)

	service.Start()
	service.WaitingDone()

}
