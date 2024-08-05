package nsqservice

import (
	"encoding/json"
	"runtime"

	"github.com/ulidev9527/zservice/zservice"

	"github.com/nsqio/go-nsq"
)

type NsqConsumerConfig struct {
	Addrs         string // 地址 多个地址用 , 隔开
	UseNsqlookupd bool   // 是否是 nsqd 地址
	Topic         string // 主题
	Channel       string // 频道
	OnMessage     func(ctx *zservice.Context, body []byte)
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

		defer func() {
			if e := recover(); e != nil {
				buf := make([]byte, 1<<12)
				stackSize := runtime.Stack(buf, true)
				ctx.LogErrorf("NSQ SEND :T %s :Q %s :E %s :ST %s", c.Topic, string(bex.Body), e, string(buf[:stackSize]))

			}
		}()

		c.OnMessage(ctx, bex.Body)
		ctx.LogInfof("NSQ ON :T %s :Q %s", c.Topic, bex.Body)
		return nil
	}))
	consumer.SetLogger(&LogEx{}, nsq.LogLevelInfo)

	startChan := make(chan any, 1)

	addrs := zservice.StringSplit(c.Addrs, ",", true)
	go func() {
		e := func() error {
			if !c.UseNsqlookupd {
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
