package nsqservice

import (
	"encoding/json"
	"zservice/zservice"

	"github.com/nsqio/go-nsq"
)

type NsqConsumerConfig struct {
	Addrs      string // 地址 多个地址用 , 隔开
	IsNsqdAddr bool   // 是否是 nsqd 地址
	Topic      string // 主题
	Channel    string // 频道
	OnMessage  func(ctx *zservice.Context, body []byte)
}

// nsq consumer
func NewNsqConsumer(c *NsqConsumerConfig) {
	consumer, e := nsq.NewConsumer(c.Topic, c.Channel, nsq.NewConfig())
	if e != nil {
		zservice.LogPanic(e)
	}

	consumer.AddHandler(nsq.HandlerFunc(func(message *nsq.Message) error {

		bex := &BodyEx{}

		if e := json.Unmarshal(message.Body, bex); e != nil {
			return e
		}
		ctx := zservice.NewContext(bex.S2S)
		c.OnMessage(ctx, bex.Body)
		return nil
	}))
	consumer.SetLogger(&LogEx{}, nsq.LogLevelInfo)

	startChan := make(chan any, 1)

	addrs := zservice.StringSplit(c.Addrs, ",", true)
	go func() {
		e := func() error {
			if c.IsNsqdAddr {
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
