package nsqservice

import (
	"fmt"
	"zservice/zservice"

	"github.com/nsqio/go-nsq"
)

type NsqProducerService struct {
	*zservice.ZService
	Producer *nsq.Producer
}

type NsqProducerServiceConfig struct {
	Name string // 服务名
	Addr string // nsq地址

	OnStart func(*nsq.Producer) // 启动的回调
}

// 创建一个新的nsq服务
func NewNsqProducerService(c *NsqProducerServiceConfig) *NsqProducerService {
	if c == nil {
		zservice.LogPanic("NsqProducerServiceConfig is nil")
		return nil
	}
	name := "NsqProducerService"
	if c.Name != "" {
		name = fmt.Sprint(name, "-", zservice.GetServiceName())
	}

	nps := &NsqProducerService{}

	nps.ZService = zservice.NewService(name, func(s *zservice.ZService) {

		producer, e := nsq.NewProducer(c.Addr, nsq.NewConfig())
		if e != nil {
			s.LogPanic(e)
		}
		nps.Producer = producer

		if c.OnStart != nil {
			c.OnStart(nps.Producer)
		}
		s.StartDone()
	})

	return nps
}
