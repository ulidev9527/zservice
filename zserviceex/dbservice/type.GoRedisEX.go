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
	ignoreKeyPrefix bool   // 是否忽略前缀
	keyLockPrefix   string // 锁前缀
}

func NewGoRedisEX(opt DBServiceOption) *GoRedisEX {
	keyPrefix := opt.RedisPrefix
	if keyPrefix == "" {
		keyPrefix = zservice.GetServiceName()
	}
	r := &GoRedisEX{
		client: redis.NewClient(&redis.Options{
			Addr:     opt.RedisAddr,
			Password: opt.RedisPass,
		}),
		keyPrefix:     fmt.Sprint(keyPrefix, ":"),
		keyLockPrefix: fmt.Sprint("__zserviceKeyLock:", keyPrefix, ":"),
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
