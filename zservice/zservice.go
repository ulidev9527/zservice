package zservice

import (
	"fmt"
	"time"
)

type ZService struct {
	name              string          // 服务名称
	tranceName        string          // 链路名称
	isAlreadyStarting bool            // 是否启动
	childServices     []*ZService     // 子服务
	parentService     *ZService       // 父服务
	chanLock          chan any        // 锁
	startTime         time.Time       // 启动时间
	onBeforeStart     func(*ZService) // 等待依赖
	onStart           func(*ZService) // 等待启动
}

// 服务配置
type ZServiceConfig struct {
	Name          string          // 服务名称
	OnBeforeStart func(*ZService) // 启动前的回调
	OnStart       func(*ZService) // 等待启动
}

func NewService(c *ZServiceConfig) *ZService {
	if c == nil {
		LogError("service is nil")
		return nil
	}
	return &ZService{
		name:              c.Name,
		tranceName:        c.Name,
		isAlreadyStarting: false,
		childServices:     []*ZService{},
		chanLock:          make(chan any, 1),
		startTime:         time.Time{},
		onBeforeStart:     c.OnBeforeStart,
		onStart:           c.OnStart,
	}
}

// 获取服务名称
func (s *ZService) GetName() string { return s.name }

// 添加子服务
func (s *ZService) AddService(service *ZService) {
	if s.parentService != nil {
		s.LogErrorf("service {%v} already has parent service", service.name)
		return
	}
	s.childServices = append(s.childServices, service)
	service.parentService = s
	service.tranceName = s.tranceName + "/" + service.name
}

// 启动服务
func (z *ZService) Start() {
	if z.isAlreadyStarting {
		z.LogError("service is already starting")
		return
	}
	z.startTime = time.Now()
	z.LogInfo("start service")
	// 子服务启动
	if len(z.childServices) > 0 {
		for i := 0; i < len(z.childServices); i++ {
			go func(item *ZService) {
				item.Start()
			}(z.childServices[i])
		}
		z.LogInfo("waiting child service")
		for i := 0; i < len(z.childServices); i++ {
			z.childServices[i].WaitingDone()
		}
	}
	// 启动自己
	if z.onBeforeStart != nil {
		z.LogInfo("waiting depend service")
		z.onBeforeStart(z)
	}
	if z.onStart != nil {
		z.LogInfo("waiting start service")
		z.onStart(z)
		z.WaitingDone()
	}
	z.LogInfo("start service done")
}

// 启动完成
func (z *ZService) StartDone() {
	close(z.chanLock)
}

// 等待完成
func (z *ZService) WaitingDone() {
	<-z.chanLock
}

// 停止服务
func (z *ZService) Stop() error {
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
func (z *ZService) LogWarn(v ...any) {
	LogWarnCaller(2, z.logCtxStr(), Sprint(v...))
}
func (z *ZService) LogWarnf(f string, v ...any) {
	LogWarnCaller(2, z.logCtxStr(), fmt.Sprintf(f, v...))
}
func (z *ZService) LogError(v ...any) {
	LogErrorCaller(2, z.logCtxStr(), Sprint(v...))
}
func (z *ZService) LogErrorf(f string, v ...any) {
	LogErrorCaller(2, z.logCtxStr(), fmt.Sprintf(f, v...))
}
