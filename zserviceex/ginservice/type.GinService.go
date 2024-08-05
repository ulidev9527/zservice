package ginservice

import (
	"fmt"
	"mime/multipart"
	"time"

	"github.com/ulidev9527/zservice/zservice"

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

// 获取携带的 上下文
func (s *GinService) GetCtx(ctx *gin.Context) *zservice.Context {
	z, has := ctx.Get(GIN_contextEX_Middleware_Key)
	if !has {
		return nil
	}
	zctx := z.(*zservice.Context)
	zctx.GinCtx = ctx
	return zctx
}

// 同步 header 信息
func (s *GinService) SyncHeader(ctx *gin.Context) {
	zctx := s.GetCtx(ctx)

	if zctx.ClientSign != "" {
		ctx.Header(zservice.S_C2S_Sign, zctx.ClientSign)
		ctx.Header(zservice.S_C2S_Time, zservice.Int64ToString(time.Now().UnixMilli()))
	}

	ctx.Header(zservice.S_C2S_Token, zctx.AuthToken)
}

// 默认响应Json
func DefResJson(code uint32, msg ...any) map[string]any {
	return gin.H{"code": code, "msg": zservice.NewError(code).SetMsg(msg...).GetMsg()}
}
