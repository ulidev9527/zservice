package nsqservice

import (
	"zservice/zservice"

	"github.com/nsqio/go-nsq"
)

type NsqConsumerConfig struct {
	Addr      string // 地址 多个地址用 , 隔开
	IsNsqd    bool   // 是否是 nsqd 地址
	Topic     string // 主题
	Channel   string // 频道
	OnMessage func(*nsq.Message) error
}

// nsq consumer
func NewNsqConsumer(c *NsqConsumerConfig) {
	consumer, e := nsq.NewConsumer(c.Topic, c.Channel, nsq.NewConfig())
	if e != nil {
		zservice.LogPanic(e)
	}

	consumer.AddHandler(nsq.HandlerFunc(c.OnMessage))
	consumer.SetLogger(&LogEx{}, nsq.LogLevelInfo)

	startChan := make(chan any, 1)

	addrs := zservice.StringSplit(c.Addr, ",", true)
	go func() {
		e := func() error {
			if c.IsNsqd {
				return consumer.ConnectToNSQDs(addrs)
			} else {
				return consumer.ConnectToNSQLookupds(addrs)
			}
		}()
		if e != nil {
			zservice.LogPanic(e)
		}
		close(startChan)
		<-consumer.StopChan
	}()

	<-startChan
}
