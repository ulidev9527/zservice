package redisservice

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	"zservice/zservice"
	"zservice/zservice/zglobal"

	"github.com/redis/go-redis/v9"
)

type GoRedisEX struct {
	client          *redis.Client
	keyPrefix       string // 前缀
	ignoreKeyPrefix bool   // 是否忽略前缀
	keyLockPrefix   string // 锁前缀
}

// 添加前缀
func (r *GoRedisEX) AddKeyPrefix(key string) string {
	if r.ignoreKeyPrefix {
		return key
	}

	return r.keyPrefix + key
}
func (r *GoRedisEX) AddkeyPrefixs(key ...string) []string {

	if r.ignoreKeyPrefix {
		return key
	}
	for i := 0; i < len(key); i++ {
		key[i] = r.keyPrefix + key[i]
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
		timeout = append(timeout, zglobal.Time_1m)
	}

	has, e := r.SetNX(lockKey, "1", timeout[0]).Result()
	if e != nil {
		return nil, zservice.NewError(e).SetCode(zglobal.Code_ErrorBreakoff)
	}
	if !has {
		return nil, zservice.NewErrorf("lock %s fail", lockKey).SetCode(zglobal.Code_RedisKeyLockFail)
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
		return zservice.NewError(e)
	} else if e := json.Unmarshal([]byte(s), v); e != nil {
		return zservice.NewError(e)
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

func (r *GoRedisEX) Del(keys ...string) *redis.IntCmd {
	return r.DelCtx(context.TODO(), keys...)
}
func (r *GoRedisEX) DelCtx(ctx context.Context, keys ...string) *redis.IntCmd {
	return r.client.Del(ctx, r.AddkeyPrefixs(keys...)...)
}

func (r *GoRedisEX) Exists(keys ...string) *redis.IntCmd {
	return r.ExistsCtx(context.TODO(), keys...)
}
func (r *GoRedisEX) ExistsCtx(ctx context.Context, keys ...string) *redis.IntCmd {
	return r.client.Exists(ctx, r.AddkeyPrefixs(keys...)...)
}

func (r *GoRedisEX) Expire(key string, expiration time.Duration) *redis.BoolCmd {
	return r.ExpireCtx(context.TODO(), key, expiration)
}
func (r *GoRedisEX) ExpireCtx(ctx context.Context, key string, expiration time.Duration) *redis.BoolCmd {
	return r.client.Expire(ctx, r.AddKeyPrefix(key), expiration)
}

func (r *GoRedisEX) HMGet(key string, fields ...string) *redis.SliceCmd {
	return r.HMGetCtx(context.TODO(), key, fields...)
}
func (r *GoRedisEX) HMGetCtx(ctx context.Context, key string, fields ...string) *redis.SliceCmd {
	return r.client.HMGet(ctx, r.AddKeyPrefix(key), fields...)
}

func (r *GoRedisEX) HGetAll(key string) *redis.MapStringStringCmd {
	return r.HGetAllCtx(context.TODO(), key)
}
func (r *GoRedisEX) HGetAllCtx(ctx context.Context, key string) *redis.MapStringStringCmd {
	return r.client.HGetAll(ctx, r.AddKeyPrefix(key))
}

func (r *GoRedisEX) HGet(key, field string) *redis.StringCmd {
	return r.HGetCtx(context.TODO(), key, field)
}
func (r *GoRedisEX) HGetCtx(ctx context.Context, key, field string) *redis.StringCmd {
	return r.client.HGet(ctx, r.AddKeyPrefix(key), field)
}
func (r *GoRedisEX) HSet(key string, values ...any) *redis.IntCmd {
	return r.HSetCtx(context.TODO(), key, values...)
}
func (r *GoRedisEX) HSetCtx(ctx context.Context, key string, values ...any) *redis.IntCmd {
	return r.client.HSet(ctx, r.AddKeyPrefix(key), values...)
}
