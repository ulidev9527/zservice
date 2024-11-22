package ginservice

import (
	"fmt"
	"mime/multipart"

	"github.com/ulidev9527/zservice/zservice"

	"github.com/gin-gonic/gin"
)

type GinService struct {
	*zservice.ZService
	Engine *gin.Engine
}

type GinServiceConfig struct {
	Port    int               // 监听端口
	OnStart func(*GinService) // 启动的回调
}

func init() {
	gin.SetMode(gin.ReleaseMode)
}

// gin 服务扩展
func NewGinService(c GinServiceConfig) *GinService {

	name := fmt.Sprint("GinService-", c.Port)
	gs := &GinService{}
	gs.Engine = gin.New()

	// 服务
	gs.ZService = zservice.NewService(zservice.ZserviceOption{
		Name: name,
		OnStart: func(s *zservice.ZService) {

			if c.OnStart != nil {
				c.OnStart(gs)
			}

			go func() {
				e := gs.Engine.Run(fmt.Sprint(":", c.Port))
				if e != nil {
					s.LogPanic(e)
				}
			}()

			gs.LogInfof("ginService listen on :%v", c.Port)
			s.StartDone()

		},
	})

	// 中间件
	gs.Engine.Use(GinMiddlewareCORS(gs.ZService))
	gs.Engine.Use(GinMiddlewareContext(gs.ZService))

	return gs
}

// 读取上传的文件信息
func ReadUploadFile(file *multipart.FileHeader) ([]byte, *zservice.Error) {

	fs, e := file.Open()
	if e != nil {
		return nil, zservice.NewError(e)
	}

	data := make([]byte, file.Size)
	if i, e := fs.Read(data); e != nil {
		return nil, zservice.NewError(i, e)
	} else {
		return data, nil
	}
}

// 默认响应Json
func DefResJson(code uint32, msg ...any) map[string]any {
	return gin.H{"code": code, "msg": zservice.NewError(code).SetMsg(msg...).GetMsg()}
}
