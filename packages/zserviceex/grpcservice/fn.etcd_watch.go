package grpcservice

import (
	"context"
	"fmt"
	"os"
	"time"
	"zserviceapps/packages/zservice"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/endpoints"
)

// 使用 etcd 在线保持
func Etcd_watch(ser *zservice.ZService, etcdClient *clientv3.Client, etcdKey string, port int32) {

	if etcdClient == nil {
		ser.LogError("can`t find etcdClient")
		os.Exit(1)
		return
	}

	if etcdKey == "" {
		ser.LogError("EtcdKey is nil")
		os.Exit(1)
		return
	}

	if etcdKey[0] != '/' {
		etcdKey = "/zserviceapps/" + etcdKey
	}

	// 创建 etcd 客户端
	mgr, e := endpoints.NewManager(etcdClient, etcdKey)
	if e != nil {
		ser.LogError(e)
		os.Exit(1)
	}

	ctx, cancel := context.WithCancel(zservice.NewContext())

	isConnected := false // 是否连接到 etcd

	// 在线保持
	zservice.Go(func() {
		for {
			// 创建一个租约，每隔 10s 需要向 etcd 汇报一次心跳，证明当前节点仍然存活
			lease, e := etcdClient.Grant(ctx, 10)
			if e != nil {
				ser.LogError("wait 1s reconnect", e)
				time.Sleep(time.Second) // 等待1秒重连
				continue
			}

			hostName, e := os.Hostname()
			if e != nil {
				ser.LogError("can`t find Hostname", e)
				os.Exit(1)
			}

			addr := fmt.Sprint(hostName, ":", port)
			endpointKey := fmt.Sprintf("%s/%s", etcdKey, addr)
			ser.LogInfof("grpc endpointKey: %s", endpointKey)
			// 添加注册节点到 etcd 中，并且携带上租约 id
			// 以 serverName/serverAddr 为 key，serverAddr 为 value
			// serverName/serverAddr 中的 serverAddr 可以自定义，只要能够区分同一个 grpc 服务器功能的不同机器即可

			e = mgr.AddEndpoint(ctx, endpointKey,
				endpoints.Endpoint{
					Addr:     addr,
					Metadata: map[string]string{"ccc": "212"},
				}, clientv3.WithLease(lease.ID))
			if e != nil {
				ser.LogError("add endpoint fail,wait 3s again", e)
				time.Sleep(time.Second * 3) // 等待1秒重连
				continue
			}

			isConnected = true
			// 处理租约续期，如果续租失败或者租约过期则退出
			for {
				isTimeout := false
				select {
				case <-time.After(5 * time.Second):
					// 租约
					_, err := etcdClient.KeepAliveOnce(context.Background(), lease.ID)
					if err != nil {
						ser.LogErrorf("Failed to keep lease alive: %s\n", err.Error())
						isTimeout = true
					}
				case <-ctx.Done():
					ser.LogError("grpc in etcd timeout", etcdClient.Ctx().Err())
					isTimeout = true
				}
				if isTimeout { // 重连
					break
				}
			}
			ser.LogWarn("GRPC In Etcd Reconnecting...")
		}
	})
	zservice.Go(func() {
		time.Sleep(time.Second * 5)
		if !isConnected {
			cancel()
		}
	})
	for {
		if isConnected {
			break
		}
	}
}
