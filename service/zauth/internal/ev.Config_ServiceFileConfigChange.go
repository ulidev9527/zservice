package internal

import (
	"fmt"
	"zservice/zservice"

	clientv3 "go.etcd.io/etcd/client/v3"
)

// 发送配置变更事件
func EV_Send_Config_serviceFileConfigChange(ctx *zservice.Context, etcd *clientv3.Client, service, fileName string) *zservice.Error {
	ev := fmt.Sprintf(EV_Config_ServiceFileConfigChange, service)

	if _, e := etcd.Put(ctx, ev, fileName); e != nil {
		return zservice.NewError(e)
	}
	return nil
}

// 监听配置变更
func EV_Watch_Config_ServiceFileConfigChange(etcd *clientv3.Client, service string, callback func(string)) {

	ev := fmt.Sprintf(EV_Config_ServiceFileConfigChange, service)
	watcher := etcd.Watch(zservice.ContextTODO(), ev)
	for resp := range watcher {
		for _, event := range resp.Events {
			fmt.Printf("Key: %s, Value: %s, Type: %s\n", event.Kv.Key, event.Kv.Value, event.Type)
			callback(string(event.Kv.Value))
		}
	}
}
