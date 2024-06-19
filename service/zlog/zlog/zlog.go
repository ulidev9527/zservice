package zlog

import (
	"zservice/zservice/ex/nsqservice"
)

var nsqPService *nsqservice.NsqProducerService

type ZlogInitConfig struct {
	NsqProducerService *nsqservice.NsqProducerService
}

func Init(c *ZlogInitConfig) {
	nsqPService = c.NsqProducerService
}
