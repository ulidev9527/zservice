package zservice

import (
	"fmt"
	"net/http"
	"runtime"
	"slices"
	"sync"
	"time"
)

type ZService struct {
	tranceName    string      // 链路名称
	dependService []*ZService // 依赖的服务
	chanStopLock  chan any    // 服务完成锁，表示服务已经执行结束
	chanStartLock chan any    // 服务启动锁，表示服务已经执行启动
	createTime    time.Time   // 创建时间
	startTime     time.Time   // 启动时间
	state         int32       // 服务状态 Service_State_XXXX
	startLock     *sync.Mutex // 互斥锁
	stopLock      *sync.Mutex
	launchCtx     *Context       // 启动上下文
	option        ServiceOptions // 配置

	IsDebug bool // 是否开启调试模式
}

// 服务器启动配置
type ServiceOptions struct {
	Name    string // 显示在日志中的名称
	Version string // 版本号

	// 服务启动回调
	OnStart func(*ZService)
	OnStop  func(*ZService)
}

func newZservice(opt ServiceOptions) *ZService {
	return &ZService{
		tranceName:    opt.Name,
		dependService: []*ZService{},
		chanStopLock:  make(chan any, 1),
		chanStartLock: make(chan any, 1),
		startLock:     &sync.Mutex{},
		stopLock:      &sync.Mutex{},
		createTime:    time.Now(),
		option:        opt,
	}
}

// 创建服务
func NewService(opt ServiceOptions) *ZService {
	if opt.Name == "" {
		opt.Name = "zservice_" + RandomString(8)
	}

	if opt.Version == "" {
		opt.Version = "version_nil"
	}

	ser := newZservice(opt)

	initMainService()

	ser.launchCtx = ser.NewContext()
	ser.launchCtx.IsFirst_userLogCtx = false
	ser.IsDebug = __mainService.IsDebug

	return ser
}

func initMainService() *ZService {

	if __mainService != nil {
		return __mainService
	}

	__mainService = newZservice(ServiceOptions{
		Name:    "Zservice",
		Version: Version,
	})
	ser := __mainService
	ser.launchCtx = ser.NewContext()
	ser.launchCtx.IsFirst_userLogCtx = false
	ser.IsDebug = GetenvInt("zservice_debug_model") == Debug_Model_On

	pprofPort := GetenvInt("zservice_debug_pprof_port")
	if pprofPort > 0 {
		Go(func() {
			http.ListenAndServe(fmt.Sprint("0.0.0.0:", pprofPort), nil)
		})
	}

	ser.LogInfo("main zservice create")
	ser.LogInfof("    at:         zservice v%s", Version)
	ser.LogInfof("    service:    %s v%s", ser.option.Name, ser.option.Version)
	ser.LogInfof("    debug:      %s, set .env zservice_debug_model to change", BoolToString(ser.IsDebug))
	ser.LogInfof("    pprof-port:      %v, set .env zservice_debug_pprof_port to change", pprofPort)
	return ser
}

func (ser *ZService) NewContext(ctxStrIn ...string) *Context {
	ctx := NewContext(ctxStrIn...)
	if len(ctxStrIn) == 0 || ctxStrIn[0] == "" {
		// 重置服务
		ctx.Trace.TraceService = ser.tranceName

		// 重置记录
		_, file, line, _ := runtime.Caller(1)
		ctx.Create_Path = fmt.Sprint(file, ":", line)
	}

	return ctx
}

// 加载到主服务上
func (ser *ZService) AddDependToMain() *ZService {
	__mainService.AddDependService(ser)
	return ser
}

// 获取启动时的上下文
func (ser *ZService) GetLauncherCtx() *Context {
	return ser.launchCtx
}

func (ser *ZService) GetName() string { return ser.option.Name }

// 获取当前服务状态
func (ser *ZService) GetState() int32 { return ser.state }

// 启动服务
func (ser *ZService) Start() *ZService {
	ser.startLock.Lock()

	if ser.state != Service_State_None {
		ser.startLock.Unlock()
		return ser
	}
	ser.state = Service_State_Starting
	ser.startTime = time.Now()
	ser.startLock.Unlock()

	ser.LogInfo("[service start]")
	// 启动依赖
	if len(ser.dependService) > 0 {
		for _, s := range ser.dependService {
			ser.LogInfo("[service start] depend ", s.tranceName)
			Go(func() {
				s.Start()
			})
		}
		// 等待依赖启动
		for _, s := range ser.dependService {
			s.WaitStart()
		}
	}
	// 启动自己
	if ser.option.OnStart != nil {
		ser.LogInfo("[service start] waiting")
		ser.option.OnStart(ser)
	}

	ser.state = Service_State_Running // 标记启动完成
	close(ser.chanStartLock)
	ser.LogInfo("[service start] done", time.Since(ser.createTime))
	return ser
}

// 等待启动
func (ser *ZService) WaitStart() *ZService {
	<-ser.chanStartLock
	return ser
}

// 等待停止
func (z *ZService) WaitStop() {
	<-z.chanStopLock
}

// 添加依赖
func (ser *ZService) AddDependService(sArr ...*ZService) *ZService {

	for _, s := range sArr {
		if slices.Contains(ser.dependService, s) {
			continue
		}

		if s == __mainService {
			ser.LogErrorf("can`t add MainService")
			continue
		}

		ser.dependService = append(ser.dependService, s)
	}

	return ser
}

// 停止服务
func (ser *ZService) Stop() {
	ser.WaitStart()
	ser.stopLock.Lock()
	if ser.state == Service_State_Stopping || ser.state == Service_State_None {
		ser.stopLock.Unlock()
		return
	}
	ser.state = Service_State_Stopping
	ser.stopLock.Unlock()

	for _, s := range ser.dependService {
		s.Stop()
		s.WaitStop()
	}
	if ser.option.OnStop != nil {
		ser.option.OnStop(ser)
	}
	ser.LogInfo("stop service")
	ser.state = Service_State_None
	close(ser.chanStopLock)
}

// -------- 打印消息
// info
func (z *ZService) LogInfo(v ...any)                   { z.launchCtx.LogInfoCaller(2, v...) }
func (z *ZService) LogInfof(f string, v ...any)        { z.launchCtx.LogInfoCallerf(2, f, v...) }
func (z *ZService) LogInfoCaller(caller int, v ...any) { z.launchCtx.LogInfoCaller(2+caller, v...) }
func (z *ZService) LogInfoCallerf(caller int, f string, v ...any) {
	z.launchCtx.LogInfoCallerf(2+caller, f, v...)
}

// warn
func (z *ZService) LogWarn(v ...any)                   { z.launchCtx.LogWarnCaller(2, v...) }
func (z *ZService) LogWarnf(f string, v ...any)        { z.launchCtx.LogWarnCallerf(2, f, v...) }
func (z *ZService) LogWarnCaller(caller int, v ...any) { z.launchCtx.LogWarnCaller(2+caller, v...) }
func (z *ZService) LogWarnCallerf(caller int, f string, v ...any) {
	z.launchCtx.LogWarnCallerf(2+caller, f, v...)
}

// error
func (z *ZService) LogError(v ...any)                   { z.launchCtx.LogErrorCaller(2, v...) }
func (z *ZService) LogErrorf(f string, v ...any)        { z.launchCtx.LogErrorCallerf(2, f, v...) }
func (z *ZService) LogErrorCaller(caller int, v ...any) { z.launchCtx.LogErrorCaller(2+caller, v...) }
func (z *ZService) LogErrorCallerf(caller int, f string, v ...any) {
	z.launchCtx.LogErrorCallerf(2+caller, f, v...)
}

// warn
func (z *ZService) LogPanic(v ...any)                   { z.launchCtx.LogPanicCaller(2, v...) }
func (z *ZService) LogPanicf(f string, v ...any)        { z.launchCtx.LogPanicCallerf(2, f, v...) }
func (z *ZService) LogPanicCaller(caller int, v ...any) { z.launchCtx.LogPanicCaller(2+caller, v...) }
func (z *ZService) LogPanicCallerf(caller int, f string, v ...any) {
	z.launchCtx.LogPanicCallerf(2+caller, f, v...)
}

// warn
func (z *ZService) LogDebug(v ...any)                   { z.launchCtx.LogDebugCaller(2, v...) }
func (z *ZService) LogDebugf(f string, v ...any)        { z.launchCtx.LogDebugCallerf(2, f, v...) }
func (z *ZService) LogDebugCaller(caller int, v ...any) { z.launchCtx.LogDebugCaller(2+caller, v...) }
func (z *ZService) LogDebugCallerf(caller int, f string, v ...any) {
	z.launchCtx.LogDebugCallerf(2+caller, f, v...)
}
