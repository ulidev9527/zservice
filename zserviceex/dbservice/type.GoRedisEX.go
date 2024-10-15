package dbservice

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/ulidev9527/zservice/zservice"
)

type GoRedisEX struct {
	client          *redis.Client
	keyPrefix       string // 前缀
	ignoreKeyPrefix bool   // 是否忽略前缀, 当前传入缀为空则忽略
	keyLockPrefix   string // 锁前缀
}

func NewGoRedisEX(opt DBServiceOption) *GoRedisEX {
	r := &GoRedisEX{
		client: redis.NewClient(&redis.Options{
			Addr:         opt.RedisAddr,
			Password:     opt.RedisPass,
			MaxIdleConns: opt.MaxIdleConns,
			PoolSize:     opt.MaxOpenConns,
			PoolTimeout:  time.Duration(opt.ConnMaxLifetime) * time.Second,
		}),

		keyLockPrefix: "__zserviceKeyLock:",
	}

	r.ignoreKeyPrefix = opt.RedisPrefix == ""

	if !r.ignoreKeyPrefix {
		r.keyPrefix = fmt.Sprint(opt.RedisPrefix, ":")
		r.keyLockPrefix = fmt.Sprint(r.keyLockPrefix, opt.RedisPrefix, ":")
	}

	_, e := r.client.Info(context.TODO(), "stats").Result()
	if e != nil {
		zservice.LogPanic(e)
	}
	return r
}

// 是否是空数据错误
func (r *GoRedisEX) IsNotFoundErr(e error) bool {
	return redis.Nil == e
}

// 添加前缀
func (r *GoRedisEX) AddKeyPrefix(key string) string {
	if r.ignoreKeyPrefix {
		return key
	}

	if strings.HasPrefix(key, r.keyPrefix) {
		return key
	}

	return r.keyPrefix + key
}
func (r *GoRedisEX) AddkeyPrefixs(key ...string) []string {

	if r.ignoreKeyPrefix {
		return key
	}
	for i := 0; i < len(key); i++ {
		key[i] = r.AddKeyPrefix(key[i])
	}
	return key
}

// 获取原生客户端
func (r *GoRedisEX) Client() *redis.Client {
	return r.client
}

// 加锁 timeout 默认 1分钟, 已经加锁的直接返回错误
func (r *GoRedisEX) Lock(key string, timeout ...time.Duration) (func(), *zservice.Error) {
	return r.LockCtx(context.TODO(), key)
}

func (r *GoRedisEX) LockCtx(ctx context.Context, key string, timeout ...time.Duration) (func(), *zservice.Error) {
	lockKey := fmt.Sprint(r.keyLockPrefix, key)
	if len(timeout) == 0 {
		timeout = append(timeout, zservice.Time_1m)
	}

	ok, e := r.SetNX(lockKey, "1", timeout[0]).Result()
	if e != nil {
		return nil, zservice.NewError(e).SetCode(zservice.Code_Fatal)
	}
	if !ok {
		return nil, zservice.NewErrorf("lock %s fail", lockKey).SetCode(zservice.Code_Repetition).SetMsg("数据正在处理，请稍后重试")
	}

	return func() {
		_, e = r.Del(lockKey).Result()
		if e != nil {
			zservice.LogErrorf("unlock %s fail: %s", lockKey, e)
		}
	}, nil
}

type LockOrWaitOption struct {
	RetryCount    uint32        // 重试次数 def:10
	RetryInterval time.Duration // 单次重试间隔 def: 10ms
	Timeout       time.Duration // 上锁超时时间 def: 1s
}

// 默认锁配置
var LockOrWaitOptionDefault = LockOrWaitOption{}

// 加锁，等待直到超时, timeout[0]等待时间 timeout[1]超时时间
func (r *GoRedisEX) LockOrWait(key string, opt LockOrWaitOption) (func(), *zservice.Error) {
	return r.LockOrWaitCtx(context.TODO(), key, opt)
}
func (r *GoRedisEX) LockOrWaitCtx(ctx context.Context, key string, opt LockOrWaitOption) (func(), *zservice.Error) {
	if opt.RetryCount == 0 {
		opt.RetryCount = 10
	}
	if opt.RetryInterval.Milliseconds() <= 0 {
		opt.RetryInterval = 10 * time.Millisecond
	}
	if opt.Timeout.Milliseconds() <= 0 {
		opt.Timeout = 1 * time.Second
	}
	for {
		un, e := r.LockCtx(ctx, key, opt.Timeout)
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

func (r *GoRedisEX) Get(key string) *redis.StringCmd {
	return r.GetCtx(context.TODO(), key)
}

func (r *GoRedisEX) GetCtx(ctx context.Context, key string) *redis.StringCmd {
	return r.client.Get(ctx, r.AddKeyPrefix(key))
}

// 查询到的内容直接转结构体
func (r *GoRedisEX) GetScan(key string, v any) *zservice.Error {
	if s, e := r.Get(key).Result(); e != nil {
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
func (r *GoRedisEX) SetScan(key string, v any) *zservice.Error {
	if s, e := json.Marshal(v); e != nil {
		return zservice.NewError(e).SetCode(zservice.Code_Fatal)
	} else if e := r.Set(key, string(s)); e != nil {
		return zservice.NewError(e).SetCode(zservice.Code_Fatal)
	} else {
		return nil
	}
}

// 将结构体转为json字符串并存储
func (r *GoRedisEX) SetScanEX(key string, v any, expiration time.Duration) *zservice.Error {
	if s, e := json.Marshal(v); e != nil {
		return zservice.NewError(e).SetCode(zservice.Code_Fatal)
	} else if e := r.SetEX(key, string(s), expiration); e != nil {
		return zservice.NewError(e).SetCode(zservice.Code_Fatal)
	} else {
		return nil
	}
}

func (r *GoRedisEX) Set(key string, value string) *redis.StatusCmd {
	return r.SetCtx(context.TODO(), key, value)
}
func (r *GoRedisEX) SetCtx(ctx context.Context, key string, value string) *redis.StatusCmd {
	return r.client.Set(ctx, r.AddKeyPrefix(key), value, 0)
}

func (r *GoRedisEX) SetEX(key string, value string, expiration time.Duration) *redis.StatusCmd {
	return r.SetExCtx(context.TODO(), key, value, expiration)
}
func (r *GoRedisEX) SetExCtx(ctx context.Context, key string, value string, expiration time.Duration) *redis.StatusCmd {
	return r.client.SetEx(ctx, r.AddKeyPrefix(key), value, expiration)
}

func (r *GoRedisEX) SetNX(key string, value string, expiration time.Duration) *redis.BoolCmd {
	return r.SetNXCtx(context.TODO(), key, value, expiration)
}
func (r *GoRedisEX) SetNXCtx(ctx context.Context, key string, value string, expiration time.Duration) *redis.BoolCmd {
	return r.client.SetNX(ctx, r.AddKeyPrefix(key), value, expiration)
}
