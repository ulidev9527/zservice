package zservice

import "time"

type ZService struct {
	isAlreadyStarting bool // 是否启动
	childServices     []*ZService
	ChanLock          chan any
	startTime         time.Time
	onDepend          func(*ZService)
	onStart           func(*ZService)
}

func NewService(name string, onDepend func(*ZService), onStart func(*ZService)) *ZService {
	return &ZService{
		onDepend: onDepend,
		onStart:  onStart,
	}
}

// 添加子服务
func (s *ZService) AddService(service *ZService) {
	s.childServices = append(s.childServices, service)
}

// 启动服务
func (s *ZService) Start() {
	if s.isAlreadyStarting {
		s.LogError("service is already starting")
		return
	}
	s.startTime = time.Now()
	for _, service := range s.childServices {
		if !service.isAlreadyStarting {
			continue
		}
		go service.Start()
	}
	for _, v := range s.childServices {
		<-v.ChanLock
	}

}

// 启动完成
func (s *ZService) StartDone() {
	close(s.ChanLock)
}

// 停止服务
func (s *ZService) Stop() error {
	return nil
}

func (s *ZService) LogError(v ...any) {

}
