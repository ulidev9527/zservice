package ginservice

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"os"
	"runtime"
	"zserviceapps/packages/zservice"

	"github.com/gin-gonic/gin"
)

type Service struct {
	zservice      *zservice.ZService
	ginEngine     *gin.Engine
	name          string                 // 默认 httpService-host:port
	port          int                    // 默认 随机
	enableContext bool                   // 是否启用上下文
	onStart       func(service *Service) // 启动回调, 会在监听端口之前运行
}

func (s *Service) Head(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	return s.ginEngine.Handle(http.MethodHead, relativePath, handlers...)
}

func (s *Service) GET(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	return s.ginEngine.GET(relativePath, handlers...)
}

func (s *Service) POST(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	return s.ginEngine.POST(relativePath, handlers...)
}

func (s *Service) HEAD(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	return s.ginEngine.HEAD(relativePath, handlers...)
}
func (s *Service) Static(relativePath string, root string) gin.IRoutes {
	return s.ginEngine.Static(relativePath, root)
}

func (s *Service) Group(relativePath string, handlers ...gin.HandlerFunc) *gin.RouterGroup {
	return s.ginEngine.Group(relativePath, handlers...)
}

func (ser *Service) GetZService() *zservice.ZService { return ser.zservice }

// 设置端口
func (s *Service) SetPort(port int) {
	if port == 0 {
		port = zservice.GetFreePort()
	}
	s.port = port
}

// 获取监听的端口
func (s *Service) GetPort() int {
	return s.port
}

type WithXXX func(service *Service)

func WithPort(port int) WithXXX    { return func(service *Service) { service.SetPort(port) } }
func WithName(name string) WithXXX { return func(service *Service) { service.name = name } }
func WithOnStart(onStart func(service *Service)) WithXXX {
	return func(service *Service) { service.onStart = onStart }
}
func WithEnableContext() WithXXX { return func(service *Service) { service.enableContext = true } }

func init() {
	gin.SetMode(gin.ReleaseMode)
}

// gin 服务扩展
func NewService(opts ...WithXXX) *Service {

	ser := &Service{}
	ser.ginEngine = gin.New()

	for _, opt := range opts {
		opt(ser)
	}

	if ser.name == "" {
		ser.name = fmt.Sprint("ginService-", zservice.GetHostname())
	}

	// ping
	ser.ginEngine.Any("/ping", gin_ping)

	// 服务
	ser.zservice = zservice.NewService(zservice.ServiceOptions{
		Name: ser.name,
		OnStart: func(_ *zservice.ZService) {

			if ser.port == 0 {
				ser.port = zservice.GetFreePort()
			}

			go func() {
				ser.zservice.LogInfof("%s listen on :%v", ser.name, ser.port)
				e := ser.ginEngine.Run(fmt.Sprint(":", ser.port))
				if e != nil {
					ser.zservice.LogError(e)
					os.Exit(1)
				}
			}()

			// 等待启动完成
			loopCount := 10
			pingAddr := fmt.Sprint("http://localhost:", ser.port, "/ping")
			for {
				// time.Sleep(time.Millisecond * 500)
				if loopCount < 0 {
					ser.zservice.LogError("can`t connect ")
				}
				loopCount--

				if _, e := zservice.HttpGet_Old(zservice.NewContext(), pingAddr, nil, nil); e != nil {
					ser.zservice.LogError(e)
				} else {
					break
				}
			}
			if ser.onStart != nil {
				ser.onStart(ser)
			}
		},
	})

	// 中间件
	ser.ginEngine.Use(
		func(c *gin.Context) { // 上下文准备
			method := c.Request.Method
			origin := c.Request.Header.Get("Origin")
			if origin != "" {
				c.Header("Access-Control-Allow-Origin", "*") // 可将将 * 替换为指定的域名
				c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")

				if ser.enableContext {
					c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization, ZCtx")
					ctx := zservice.NewContext(c.Request.Header.Get("ZCtx"))
					ctx.Authorization = c.Request.Header.Get("Authorization")
					c.Set(GIN_contextEX_Middleware_Key, ctx)
				} else {
					c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
					ctx := zservice.NewContext()
					ctx.Authorization = c.Request.Header.Get("Authorization")
					c.Set(GIN_contextEX_Middleware_Key, ctx)
				}

				c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
				c.Header("Access-Control-Allow-Credentials", "true")
			}
			if method == "OPTIONS" {
				c.AbortWithStatus(http.StatusNoContent)
			}
		},
		func(c *gin.Context) { // 日志输出

			ctx := GetContext(c)

			defer func() {
				//放在匿名函数里,e捕获到错误信息，并且输出
				e := recover()
				if e != nil {
					buf := make([]byte, 1<<12)
					stackSize := runtime.Stack(buf, true)
					ctx.LogErrorf("%v %v :E %v :T %v",
						c.Request.Method, c.Request.URL, e, string(buf[:stackSize]),
					)
				}
			}()

			c.Next()

			ctx.LogInfo(fmt.Sprintf("GIN %v %v %v %v",
				c.Request.Method, c.Request.URL,
				c.Writer.Status(), ctx.Since(),
			))

		},
	)

	return ser
}

// 获取 gin 携带的上下文
func GetContext(c *gin.Context) *zservice.Context {
	z, has := c.Get(GIN_contextEX_Middleware_Key)
	if !has {

		z = zservice.NewContext()
		ctx := z.(*zservice.Context)
		ctx.Authorization = c.Request.Header.Get("Authorization")
		c.Set(GIN_contextEX_Middleware_Key, ctx)
	}

	return z.(*zservice.Context)
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
