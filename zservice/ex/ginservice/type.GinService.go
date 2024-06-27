package ginservice

import (
	"fmt"
	"zservice/zservice"

	"github.com/gin-gonic/gin"
)

type GinService struct {
	*zservice.ZService
	Engine *gin.Engine
}

type GinServiceConfig struct {
	ListenPort string            // 监听地址
	OnStart    func(*GinService) // 启动的回调
}

func init() {
	gin.SetMode(gin.ReleaseMode)
}

// gin 服务扩展
func NewGinService(c *GinServiceConfig) *GinService {

	if c == nil {
		zservice.LogPanic("GinServiceConfig is nil")
		return nil
	}
	name := fmt.Sprint("GinService-", c.ListenPort)
	gs := &GinService{}
	gs.Engine = gin.New()

	// 服务
	gs.ZService = zservice.NewService(name, func(s *zservice.ZService) {

		if c.OnStart != nil {
			c.OnStart(gs)
		}

		go func() {
			e := gs.Engine.Run(fmt.Sprint(":", c.ListenPort))
			if e != nil {
				s.LogPanic(e)
			}
		}()

		gs.LogInfof("ginService listen on :%v", c.ListenPort)
		s.StartDone()

	})

	// 中间件
	gs.Engine.Use(GinMiddlewareCORS(gs.ZService))
	gs.Engine.Use(GinMiddlewareContext(gs.ZService))

	return gs
}
