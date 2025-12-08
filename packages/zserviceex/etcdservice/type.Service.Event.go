package etcdservice

import (
	"context"
	"zserviceapps/packages/zservice"
	"zserviceapps/packages/zserviceex/etcdservice/pb"

	"google.golang.org/protobuf/proto"
)

// 发送事件
func (ser *Service) SendEvent(ctx *zservice.Context, key string, body []byte) *zservice.Error {
	eventBody := pb.Get_EventBody().SetCtx(ctx.ToContextString()).SetBody(body)

	if _, e := ser.Client.Put(ctx, key, string(eventBody.Put_Bytes())); e != nil {
		ctx.LogInfof("ETCD K:%s E:%s", key, e)
		return zservice.NewError(e)
	}
	return nil
}

// 监听事件，支持取消功能
func (ser *Service) WatchEvent(key string, cb func(ctx *zservice.Context, val []byte)) (cancelFunc context.CancelFunc) {
	// 创建带取消功能的上下文
	ctx, cancel := context.WithCancel(zservice.ContextTODO())

	zservice.Go(func() {
		watcher := ser.Client.Watch(ctx, key)
		for resp := range watcher {
			for _, event := range resp.Events {
				eventBody := pb.Get_EventBody()

				if e := proto.Unmarshal(event.Kv.Value, eventBody); e != nil {
					ser.LogErrorf("ETCD K:%s EK:%s E:%s", key, event.Kv.Key, e)
					eventBody.Put()
					continue
				}

				cb(zservice.NewContext(eventBody.Ctx), eventBody.Body)
				eventBody.Put()
			}
		}
	})

	// 返回 cancel 函数，以便外部可以调用取消
	return cancel
}
