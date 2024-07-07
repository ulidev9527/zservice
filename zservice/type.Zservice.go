package zservice

import (
	"fmt"
	"sync"
	"time"
)

type ZService struct {
	name                 string          // 服务名称
	tranceName           string          // 链路名称
	dependService        []*ZService     // 等待的依赖服务
	chanServiceStopLock  chan any        // 服务完成锁，表示服务已经执行结束
	chanServiceStartLock chan any        // 服务启动锁，表示服务已经执行启动
	createTime           time.Time       // 创建时间
	startTime            time.Time       // 启动时间
	onStart              func(*ZService) // 等待启动
	state                uint32          // 服务状态 0已创建 1等待启动 2已启动 3已停止
	mu                   sync.Mutex      // 互斥锁
}

// 创建一个服务
func createService(name string, onStart func(*ZService)) *ZService {

	tName := name
	if mainService != nil {
		tName = fmt.Sprintf("%s/%s", mainService.tranceName, name)
	}

	return &ZService{
		name:                 name,
		tranceName:           tName,
		dependService:        []*ZService{},
		chanServiceStopLock:  make(chan any, 1),
		chanServiceStartLock: make(chan any, 1),
		mu:                   sync.Mutex{},
		createTime:           time.Now(),
		onStart:              onStart,
	}
}

// 外部创建服务入口
func NewService(name string, onStart func(*ZService)) *ZService {
	if mainService == nil {
		LogPanic("you need use zservice.Init first")
	}
	return createService(name, onStart)
}

// 启动服务
func (z *ZService) Start() *ZService {
	z.mu.Lock()
	if z.state != 0 {
		return z
	}
	z.state = 1
	z.startTime = time.Now()
	z.mu.Unlock()

	z.LogInfo("start service")
	// 启动依赖
	if len(z.dependService) > 0 {
		z.LogInfo("waiting depend service")
		if z == mainService {
			for i := 0; i < len(z.dependService); i++ {
				go z.dependService[i].Start()
			}
		}
		// 等待依赖启动
		for i := 0; i < len(z.dependService); i++ {
			z.dependService[i].WaitStart()
		}
	}
	// 启动自己
	if z.onStart != nil {
		z.LogInfo("waiting start service")
		z.onStart(z)
		z.WaitStart()
	}
	if z == mainService {
		z.LogInfof("[[[[[[ %s start done %s]]]]]]", GetServiceName(), time.Since(z.createTime))
		z.StartDone()
	} else {
		z.LogInfo("start service done", time.Since(z.createTime))
	}
	return z
}

// 启动完成
func (z *ZService) StartDone() {
	close(z.chanServiceStartLock)
}

// 等待启动
func (z *ZService) WaitStart() *ZService {
	<-z.chanServiceStartLock
	z.state = 2
	return z
}

// 等待停止
func (z *ZService) WaitStop() {
	<-z.chanServiceStopLock
	z.state = 3
}

// 添加依赖
func (z *ZService) AddDependService(sArr ...*ZService) *ZService {
	z.dependService = append(z.dependService, sArr...)
	if z != mainService {
		mainService.AddDependService(sArr...)
	}
	return z
}
func (z *ZService) GetState() uint32 {
	return z.state
}

// 停止服务
func (z *ZService) Stop() error {
	LogInfo("stop service", z.tranceName)
	close(z.chanServiceStopLock)
	return nil
}

// -------- 打印消息
// 获取日志的打印信息
func (z *ZService) logCtxStr() string {
	return fmt.Sprintf("[%s]", z.tranceName)
}

func (z *ZService) LogInfo(v ...any) {
	LogInfoCaller(2, z.logCtxStr(), Sprint(v...))
}
func (z *ZService) LogInfof(f string, v ...any) {
	LogInfoCaller(2, z.logCtxStr(), fmt.Sprintf(f, v...))
}
func (z *ZService) LogInfoCaller(caller int, v ...any) {
	LogInfoCaller(2+caller, z.logCtxStr(), Sprint(v...))
}
func (z *ZService) LogWarn(v ...any) {
	LogWarnCaller(2, z.logCtxStr(), Sprint(v...))
}
func (z *ZService) LogWarnf(f string, v ...any) {
	LogWarnCaller(2, z.logCtxStr(), fmt.Sprintf(f, v...))
}
func (z *ZService) LogWarnCaller(caller int, v ...any) {
	LogWarnCaller(2+caller, z.logCtxStr(), Sprint(v...))
}
func (z *ZService) LogError(v ...any) {
	LogErrorCaller(2, z.logCtxStr(), Sprint(v...))
}
func (z *ZService) LogErrorf(f string, v ...any) {
	LogErrorCaller(2, z.logCtxStr(), fmt.Sprintf(f, v...))
}
func (z *ZService) LogErrorCaller(caller int, v ...any) {
	LogErrorCaller(2+caller, z.logCtxStr(), Sprint(v...))
}
func (z *ZService) LogPanic(v ...any) {
	LogPanicCaller(2, z.logCtxStr(), Sprint(v...))
}
func (z *ZService) LogPanicf(f string, v ...any) {
	LogPanicCallerf(2, z.logCtxStr(), fmt.Sprintf(f, v...))
}
func (z *ZService) LogPanicCaller(caller int, v ...any) {
	LogPanicCaller(2+caller, z.logCtxStr(), Sprint(v...))
}
