package goredisservice

import (
	"context"
	"encoding/json"
	"fmt"
	"runtime"
	"time"

	"zserviceapps/packages/zservice"

	"github.com/redis/go-redis/v9"
)

type GoRedisEX struct {
	*redis.Client
}

var lockKey = "__goredisservice__lock__:"

// 是否是空数据错误
func (r *GoRedisEX) IsNotFoundErr(e error) bool {
	return redis.Nil == e
}

// 加锁 timeout 默认 1分钟, 已经加锁的直接返回错误
func (r *GoRedisEX) Lock(ctx *zservice.Context, key string, timeout ...time.Duration) (func(), *zservice.Error) {
	lockKey := fmt.Sprint(lockKey, key)
	if len(timeout) == 0 {
		timeout = append(timeout, Time_1m)
	}

	ok, e := r.SetNX(ctx, lockKey, "1", timeout[0]).Result()
	if e != nil {
		return nil, zservice.NewError(e).SetCode(zservice.Code_Fatal)
	}
	if !ok {
		return nil, zservice.NewErrorf("lock %s fail", lockKey).SetCode(zservice.Code_Repetition).SetMsg("数据正在处理，请稍后重试")
	}

	return func() {
		_, e = r.Del(ctx, lockKey).Result()
		if e != nil {
			ctx.LogErrorf("unlock %s fail: %s", lockKey, e)
		}
	}, nil
}

type LockOrWaitOption struct {
	RetryCount    uint32        // 重试次数 def:5
	RetryInterval time.Duration // 单次重试间隔 def: 10ms
	Timeout       time.Duration // 上锁超时时间 def: 3s
}

// 默认锁配置
var LockOrWaitOptionDefault = LockOrWaitOption{
	RetryCount:    5,
	RetryInterval: 10 * time.Millisecond,
	Timeout:       3 * time.Second,
}

// 加锁，等待直到超时, timeout[0]等待时间 timeout[1]超时时间
func (r *GoRedisEX) LockOrWait(ctx *zservice.Context, key string, opts ...LockOrWaitOption) (func(), *zservice.Error) {
	var opt LockOrWaitOption
	if len(opts) > 0 {
		opt = opts[0]
	}

	if opt.RetryCount == 0 {
		opt.RetryCount = LockOrWaitOptionDefault.RetryCount
	}
	if opt.RetryInterval.Milliseconds() <= 0 {
		opt.RetryInterval = LockOrWaitOptionDefault.RetryInterval
	}
	if opt.Timeout.Milliseconds() <= 0 {
		opt.Timeout = LockOrWaitOptionDefault.Timeout
	}
	for {
		un, e := r.Lock(ctx, key, opt.Timeout)
		if e != nil {
			if e.GetCode() == zservice.Code_Repetition {
				time.Sleep(opt.RetryInterval)
				continue
			}
			return nil, e
		}
		return un, nil
	}
}

// 重写 Set
func (r *GoRedisEX) Set(ctx context.Context, key string, value interface{}) *redis.StatusCmd {
	return r.Client.Set(ctx, key, value, 0)
}

// 查询到的内容直接转结构体
func (r *GoRedisEX) GetScan(ctx context.Context, key string, v any) *zservice.Error {
	if s, e := r.Get(ctx, key).Result(); e != nil {
		if r.IsNotFoundErr(e) {
			return zservice.NewError(e).SetCode(zservice.Code_NotFound)
		}
		return zservice.NewError(e).SetCode(zservice.Code_Fatal)
	} else if e := json.Unmarshal([]byte(s), v); e != nil {
		return zservice.NewError(e, key).SetCode(zservice.Code_Fatal)
	} else {
		return nil
	}
}

// 将结构体转为json字符串并存储
func (r *GoRedisEX) SetScan(ctx context.Context, key string, v any) *zservice.Error {
	if s, e := json.Marshal(v); e != nil {
		return zservice.NewError(e).SetCode(zservice.Code_Fatal)
	} else if e := r.Set(ctx, key, string(s)).Err(); e != nil {
		return zservice.NewError(e).SetCode(zservice.Code_Fatal)
	} else {
		return nil
	}
}

// 将结构体转为json字符串并存储
func (r *GoRedisEX) SetScanEX(ctx context.Context, key string, v any, expiration time.Duration) *zservice.Error {
	if s, e := json.Marshal(v); e != nil {
		return zservice.NewError(e).SetCode(zservice.Code_Fatal)
	} else if e := r.SetEx(ctx, key, string(s), expiration).Err(); e != nil {
		return zservice.NewError(e).SetCode(zservice.Code_Fatal)
	} else {
		return nil
	}
}

// redis 消息订阅发布参数
type PubsubBody struct {
	S2S string `json:"s2s"` // 上下文
	Val []byte `json:"val"` // 内存数据
}

func (r *GoRedisEX) Pub(channel string, data []byte) *zservice.Error {
	return r.PubCtx(zservice.NewContext(), channel, data)
}

func (r *GoRedisEX) PubCtx(ctx *zservice.Context, channel string, data []byte) *zservice.Error {

	body := &PubsubBody{
		S2S: ctx.ToContextString(),
		Val: data,
	}

	err := r.Publish(ctx, channel, string(zservice.JsonMustMarshal(body))).Err()
	if err != nil {
		return zservice.NewError(err)
	}
	return nil
}

// 订阅消息，返回取消函数
func (r *GoRedisEX) Sub(channel string, cb func(*zservice.Context, string, []byte)) *redis.PubSub {
	return r.SubCtx(zservice.NewContext(), channel, cb)
}

// 订阅消息，返回取消函数，带上下文
func (r *GoRedisEX) SubCtx(ctx *zservice.Context, channel string, cb func(*zservice.Context, string, []byte)) *redis.PubSub {
	sub := r.Subscribe(ctx, channel)

	zservice.Go(func() {
		defer sub.Close()

		for msg := range sub.Channel() {

			channel := msg.Channel
			data := []byte(msg.Payload)

			body := &PubsubBody{}
			if e := json.Unmarshal(data, body); e != nil {
				buf := make([]byte, 1<<12)
				stackSize := runtime.Stack(buf, true)
				ctx.LogErrorf("%v :E %v :ST %v",
					channel, e, string(buf[:stackSize]),
				)
				continue
			}

			cb(zservice.NewContext(body.S2S), channel, body.Val)

		}

	})

	return sub
}
