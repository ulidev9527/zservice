package nsqservice

import (
	"runtime"
	"zservice/zservice"

	"github.com/nsqio/go-nsq"
)

type NsqProducerService struct {
	*zservice.ZService
	Producer *nsq.Producer
}

type NsqProducerServiceConfig struct {
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

func (nps *NsqProducerService) Publish(ctx *zservice.Context, topic string, body []byte) *zservice.Error {
	bex := zservice.JsonMustMarshal(&BodyEx{
		S2S:  ctx.GetS2S(),
		Body: body,
	})

	defer func() {
		if e := recover(); e != nil {
			buf := make([]byte, 1<<12)
			stackSize := runtime.Stack(buf, true)
			ctx.LogErrorf("NSQ SEND :T %s :E %s :ST %s", topic, e, string(buf[:stackSize]))
		}
	}()

	if e := nps.Producer.Publish(topic, bex); e != nil {
		ctx.LogErrorf("NSQ SEND :T %s :E %s", topic, e)
		return zservice.NewError(e)
	}
	ctx.LogInfof("NSQ SEND :T %s :Q %s", topic, string(body))
	return nil
}

func (nps *NsqProducerService) Stop() {

	nps.Producer.Stop()
}
