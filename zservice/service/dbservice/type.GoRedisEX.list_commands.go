package dbservice

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

func (r *GoRedisEX) BLPop(ctx context.Context, timeout time.Duration, keys ...string) *redis.StringSliceCmd {
	return r.BLPopCtx(ctx, timeout, keys...)
}

func (r *GoRedisEX) BLPopCtx(ctx context.Context, timeout time.Duration, keys ...string) *redis.StringSliceCmd {
	return r.client.BLPop(ctx, timeout, r.AddkeyPrefixs(keys...)...)
}
func (r *GoRedisEX) BLMPop(ctx context.Context, timeout time.Duration, direction string, count int64, keys ...string) *redis.KeyValuesCmd {
	return r.BLMPopCtx(ctx, timeout, direction, count, keys...)
}
func (r *GoRedisEX) BLMPopCtx(ctx context.Context, timeout time.Duration, direction string, count int64, keys ...string) *redis.KeyValuesCmd {
	return r.client.BLMPop(ctx, timeout, direction, count, r.AddkeyPrefixs(keys...)...)
}

func (r *GoRedisEX) BRPop(ctx context.Context, timeout time.Duration, keys ...string) *redis.StringSliceCmd {
	return r.BRPopCtx(ctx, timeout, keys...)
}
func (r *GoRedisEX) BRPopCtx(ctx context.Context, timeout time.Duration, keys ...string) *redis.StringSliceCmd {
	return r.client.BRPop(ctx, timeout, r.AddkeyPrefixs(keys...)...)
}

func (r *GoRedisEX) BRPopLPush(ctx context.Context, source, destination string, timeout time.Duration) *redis.StringCmd {
	return r.BRPopLPushCtx(ctx, source, destination, timeout)
}
func (r *GoRedisEX) BRPopLPushCtx(ctx context.Context, source, destination string, timeout time.Duration) *redis.StringCmd {
	return r.client.BRPopLPush(ctx, r.AddKeyPrefix(source), r.AddKeyPrefix(destination), timeout)
}

func (r *GoRedisEX) LIndex(key string, index int64) *redis.StringCmd {
	return r.LIndexCtx(context.TODO(), key, index)
}
func (r *GoRedisEX) LIndexCtx(ctx context.Context, key string, index int64) *redis.StringCmd {
	return r.client.LIndex(ctx, r.AddKeyPrefix(key), index)
}

func (r *GoRedisEX) LMPop(ctx context.Context, direction string, count int64, keys ...string) *redis.KeyValuesCmd {
	return r.LMPopCtx(ctx, direction, count, keys...)
}
func (r *GoRedisEX) LMPopCtx(ctx context.Context, direction string, count int64, keys ...string) *redis.KeyValuesCmd {
	return r.client.LMPop(ctx, direction, count, r.AddkeyPrefixs(keys...)...)
}

func (r *GoRedisEX) LInsert(key string, op string, pivot any, value any) *redis.IntCmd {
	return r.LInsertCtx(context.TODO(), key, op, pivot, value)
}
func (r *GoRedisEX) LInsertCtx(ctx context.Context, key string, op string, pivot any, value any) *redis.IntCmd {
	return r.client.LInsert(ctx, r.AddKeyPrefix(key), op, pivot, value)
}

func (r *GoRedisEX) LInsertBefore(key string, pivot any, value any) *redis.IntCmd {
	return r.LInsertBeforeCtx(context.TODO(), key, pivot, value)
}
func (r *GoRedisEX) LInsertBeforeCtx(ctx context.Context, key string, pivot any, value any) *redis.IntCmd {
	return r.client.LInsertBefore(ctx, r.AddKeyPrefix(key), pivot, value)
}

func (r *GoRedisEX) LInsertAfter(key string, pivot any, value any) *redis.IntCmd {
	return r.LInsertAfterCtx(context.TODO(), key, pivot, value)
}
func (r *GoRedisEX) LInsertAfterCtx(ctx context.Context, key string, pivot any, value any) *redis.IntCmd {
	return r.client.LInsertAfter(ctx, r.AddKeyPrefix(key), pivot, value)
}

func (r *GoRedisEX) LLen(key string) *redis.IntCmd {
	return r.LLenCtx(context.TODO(), key)
}
func (r *GoRedisEX) LLenCtx(ctx context.Context, key string) *redis.IntCmd {
	return r.client.LLen(ctx, r.AddKeyPrefix(key))
}

func (r *GoRedisEX) LPop(key string) *redis.StringCmd {
	return r.LPopCtx(context.TODO(), key)
}
func (r *GoRedisEX) LPopCtx(ctx context.Context, key string) *redis.StringCmd {
	return r.client.LPop(ctx, r.AddKeyPrefix(key))
}

func (r *GoRedisEX) LPopCount(key string, count int) *redis.StringSliceCmd {
	return r.LPopCountCtx(context.TODO(), key, count)
}
func (r *GoRedisEX) LPopCountCtx(ctx context.Context, key string, count int) *redis.StringSliceCmd {
	return r.client.LPopCount(ctx, r.AddKeyPrefix(key), count)
}

func (r *GoRedisEX) LPos(key string, element string, args redis.LPosArgs) *redis.IntCmd {
	return r.LPosCtx(context.TODO(), key, element, args)
}
func (r *GoRedisEX) LPosCtx(ctx context.Context, key string, element string, args redis.LPosArgs) *redis.IntCmd {
	return r.client.LPos(ctx, r.AddKeyPrefix(key), element, args)
}

func (r *GoRedisEX) LPosCount(key string, element string, count int64, args redis.LPosArgs) *redis.IntSliceCmd {
	return r.LPosCountCtx(context.TODO(), key, element, count, args)
}
func (r *GoRedisEX) LPosCountCtx(ctx context.Context, key string, element string, count int64, args redis.LPosArgs) *redis.IntSliceCmd {
	return r.client.LPosCount(ctx, r.AddKeyPrefix(key), element, count, args)
}

func (r *GoRedisEX) LPush(key string, values ...any) *redis.IntCmd {
	return r.LPushCtx(context.TODO(), key, values...)
}
func (r *GoRedisEX) LPushCtx(ctx context.Context, key string, values ...any) *redis.IntCmd {
	return r.client.LPush(ctx, r.AddKeyPrefix(key), values...)
}

func (r *GoRedisEX) LPushX(key string, values ...any) *redis.IntCmd {
	return r.LPushXCtx(context.TODO(), key, values...)
}
func (r *GoRedisEX) LPushXCtx(ctx context.Context, key string, values ...any) *redis.IntCmd {
	return r.client.LPushX(ctx, r.AddKeyPrefix(key), values...)
}
func (r *GoRedisEX) LRange(key string, start, stop int64) *redis.StringSliceCmd {
	return r.LRangeCtx(context.TODO(), key, start, stop)
}
func (r *GoRedisEX) LRangeCtx(ctx context.Context, key string, start, stop int64) *redis.StringSliceCmd {
	return r.client.LRange(ctx, r.AddKeyPrefix(key), start, stop)
}
func (r *GoRedisEX) LRem(key string, count int64, value interface{}) *redis.IntCmd {
	return r.LRemCtx(context.TODO(), key, count, value)
}
func (r *GoRedisEX) LRemCtx(ctx context.Context, key string, count int64, value interface{}) *redis.IntCmd {
	return r.client.LRem(ctx, r.AddKeyPrefix(key), count, value)
}
func (r *GoRedisEX) LSet(key string, index int64, value interface{}) *redis.StatusCmd {
	return r.LSetCtx(context.TODO(), key, index, value)
}
func (r *GoRedisEX) LSetCtx(ctx context.Context, key string, index int64, value interface{}) *redis.StatusCmd {
	return r.client.LSet(ctx, r.AddKeyPrefix(key), index, value)
}
func (r *GoRedisEX) LTrim(key string, start, stop int64) *redis.StatusCmd {
	return r.LTrimCtx(context.TODO(), key, start, stop)
}
func (r *GoRedisEX) LTrimCtx(ctx context.Context, key string, start, stop int64) *redis.StatusCmd {
	return r.client.LTrim(ctx, r.AddKeyPrefix(key), start, stop)
}
func (r *GoRedisEX) RPop(key string) *redis.StringCmd {
	return r.RPopCtx(context.TODO(), key)
}
func (r *GoRedisEX) RPopCtx(ctx context.Context, key string) *redis.StringCmd {
	return r.client.RPop(ctx, r.AddKeyPrefix(key))
}

func (r *GoRedisEX) RPopCount(key string, count int) *redis.StringSliceCmd {
	return r.RPopCountCtx(context.TODO(), key, count)
}
func (r *GoRedisEX) RPopCountCtx(ctx context.Context, key string, count int) *redis.StringSliceCmd {
	return r.client.RPopCount(ctx, r.AddKeyPrefix(key), count)
}

func (r *GoRedisEX) RPopLPush(source, destination string) *redis.StringCmd {
	return r.RPopLPushCtx(context.TODO(), source, destination)
}
func (r *GoRedisEX) RPopLPushCtx(ctx context.Context, source, destination string) *redis.StringCmd {
	return r.client.RPopLPush(ctx, r.AddKeyPrefix(source), r.AddKeyPrefix(destination))
}

func (r *GoRedisEX) RPush(key string, values ...any) *redis.IntCmd {
	return r.RPushCtx(context.TODO(), key, values...)
}
func (r *GoRedisEX) RPushCtx(ctx context.Context, key string, values ...any) *redis.IntCmd {
	return r.client.RPush(ctx, r.AddKeyPrefix(key), values...)
}

func (r *GoRedisEX) RPushX(key string, values ...any) *redis.IntCmd {
	return r.RPushXCtx(context.TODO(), key, values...)
}
func (r *GoRedisEX) RPushXCtx(ctx context.Context, key string, values ...any) *redis.IntCmd {
	return r.client.RPushX(ctx, r.AddKeyPrefix(key), values...)
}
