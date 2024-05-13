package zconfig

import (
	"fmt"
	"zservice/service/zconfig/internal"
	"zservice/zservice"
	"zservice/zservice/ex/nsqservice"

	"github.com/nsqio/go-nsq"
)

type NsqConsumerConfig struct {
	Addr      string
	IsNsqd    bool // 是否是nsqlookupd地址
	OnMessage func(*nsq.Message) error
}

// 监听配置文件改变
func NewNsqConsumer_FileConfigChange(c *NsqConsumerConfig) {
	consumer, e := nsq.NewConsumer(internal.NSQ_FileConfig_Change, fmt.Sprintf("%s-%s", zservice.GetServiceName(), zservice.RandomXID()), nsq.NewConfig())
	if e != nil {
		zservice.LogPanic(e)
	}

	consumer.AddHandler(nsq.HandlerFunc(c.OnMessage))
	consumer.SetLogger(&nsqservice.LogEx{}, nsq.LogLevelInfo)

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
