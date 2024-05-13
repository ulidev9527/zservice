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

	if c.Addr == "" {
		zservice.LogPanic("NsqProducerServiceConfig.Addr is nil")
	}

	nps := &NsqProducerService{}

	nps.ZService = zservice.NewService(name, func(s *zservice.ZService) {

		producer, e := nsq.NewProducer(c.Addr, nsq.NewConfig())
		if e != nil {
			s.LogPanic(e)
		}
		s.LogInfo("start nsq producer", c.Addr)
		nps.Producer = producer

		producer.SetLogger(&LogEx{s}, nsq.LogLevelInfo)

		e = nps.Producer.Ping()
		if e != nil {
			s.LogPanic(e)
		}

		if c.OnStart != nil {
			c.OnStart(nps.Producer)
		}
		s.StartDone()
	})

	return nps
}
