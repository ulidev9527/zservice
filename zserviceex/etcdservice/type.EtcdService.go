package etcdservice

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/ulidev9527/zservice/zservice"

	clientv3 "go.etcd.io/etcd/client/v3"
)

type EtcdService struct {
	*zservice.ZService
	EtcdClient *clientv3.Client
}

type EtcdServiceConfig struct {
	Addr    string             // ETCD 服务地址
	OnStart func(*EtcdService) // 启动的回调
}

func NewEtcdService(c *EtcdServiceConfig) *EtcdService {

	if c == nil {
		zservice.LogPanic("EtcdServiceConfig is nil")
		return nil
	}

	name := fmt.Sprint("EtcdService-", c.Addr)

	es := &EtcdService{}
	es.ZService = zservice.NewService(name, func(s *zservice.ZService) {

		s.LogInfof("etcdService connect on %v", c.Addr)

		timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if status, e := es.EtcdClient.Status(timeoutCtx, c.Addr); e != nil {
			s.LogPanic(e)
		} else {
			s.LogInfo("ETCD Status:", string(zservice.JsonMustMarshal(status)))
		}
		if c.OnStart != nil {
			c.OnStart(es)
		}
		s.StartDone()

	})

	etcd, e := clientv3.New(clientv3.Config{
		Endpoints:   []string{c.Addr},
		DialTimeout: 5 * time.Second,
	})

	if e != nil {
		es.LogPanic(e)
	}

	es.EtcdClient = etcd
	return es
}

// 发送事件
func (es *EtcdService) SendEvent(ctx *zservice.Context, key string, val string) *zservice.Error {
	eb := &EventBody{
		S2S: ctx.GetS2S(),
		Val: val,
	}
	ctx.LogDebug(zservice.JsonMustMarshalString(eb))
	if _, e := es.EtcdClient.Put(ctx, key, zservice.JsonMustMarshalString(eb)); e != nil {
		ctx.LogInfof("ETCD K:%s V:%s E:%s", key, val, e)
		return zservice.NewError(e)
	} else {
		ctx.LogInfof("ETCD K:%s V:%s", key, val)
	}
	return nil
}

// 监听事件，支持取消功能
func (es *EtcdService) WatchEvent(key string, cb func(ctx *zservice.Context, val string)) (cancelFunc context.CancelFunc) {
	// 创建带取消功能的上下文
	ctx, cancel := context.WithCancel(zservice.ContextTODO())

	zservice.Go(func() {
		watcher := es.EtcdClient.Watch(ctx, key)
		for resp := range watcher {
			for _, event := range resp.Events {
				eb := &EventBody{}
				if e := json.Unmarshal([]byte(event.Kv.Value), eb); e != nil {
					es.LogErrorf("ETCD K:%s E:%s", key, e)
					continue
				}
				ctx := zservice.NewContext(eb.S2S)

				ctx.LogInfof("ETCD K:%s V:%s", key, eb.Val)
				cb(ctx, eb.Val)
			}
		}
	})

	// 返回 cancel 函数，以便外部可以调用取消
	return cancel
}
