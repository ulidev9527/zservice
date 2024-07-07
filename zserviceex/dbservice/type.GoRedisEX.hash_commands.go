package dbservice

import (
	"context"

	"github.com/redis/go-redis/v9"
)

func (r *GoRedisEX) HDel(key string, fields ...string) *redis.IntCmd {
	return r.HDelCtx(context.TODO(), key, fields...)
}
func (r *GoRedisEX) HDelCtx(ctx context.Context, key string, fields ...string) *redis.IntCmd {
	return r.client.HDel(ctx, r.AddKeyPrefix(key), fields...)
}

func (r *GoRedisEX) HExists(key, field string) *redis.BoolCmd {
	return r.HExistsCtx(context.TODO(), key, field)
}
func (r *GoRedisEX) HExistsCtx(ctx context.Context, key, field string) *redis.BoolCmd {
	return r.client.HExists(ctx, r.AddKeyPrefix(key), field)
}

func (r *GoRedisEX) HGet(key, field string) *redis.StringCmd {
	return r.HGetCtx(context.TODO(), key, field)
}
func (r *GoRedisEX) HGetCtx(ctx context.Context, key, field string) *redis.StringCmd {
	return r.client.HGet(ctx, r.AddKeyPrefix(key), field)
}

func (r *GoRedisEX) HGetAll(key string) *redis.MapStringStringCmd {
	return r.HGetAllCtx(context.TODO(), key)
}
func (r *GoRedisEX) HGetAllCtx(ctx context.Context, key string) *redis.MapStringStringCmd {
	return r.client.HGetAll(ctx, r.AddKeyPrefix(key))
}

func (r *GoRedisEX) HIncrBy(key, field string, incr int64) *redis.IntCmd {
	return r.HIncrByCtx(context.TODO(), key, field, incr)
}
func (r *GoRedisEX) HIncrByCtx(ctx context.Context, key, field string, incr int64) *redis.IntCmd {
	return r.client.HIncrBy(ctx, r.AddKeyPrefix(key), field, incr)
}

func (r *GoRedisEX) HIncrByFloat(key, field string, incr float64) *redis.FloatCmd {
	return r.HIncrByFloatCtx(context.TODO(), key, field, incr)
}
func (r *GoRedisEX) HIncrByFloatCtx(ctx context.Context, key, field string, incr float64) *redis.FloatCmd {
	return r.client.HIncrByFloat(ctx, r.AddKeyPrefix(key), field, incr)
}

func (r *GoRedisEX) HKeys(key string) *redis.StringSliceCmd {
	return r.HKeysCtx(context.TODO(), key)
}
func (r *GoRedisEX) HKeysCtx(ctx context.Context, key string) *redis.StringSliceCmd {
	return r.client.HKeys(ctx, r.AddKeyPrefix(key))
}

func (r *GoRedisEX) HLen(key string) *redis.IntCmd {
	return r.HLenCtx(context.TODO(), key)
}
func (r *GoRedisEX) HLenCtx(ctx context.Context, key string) *redis.IntCmd {
	return r.client.HLen(ctx, r.AddKeyPrefix(key))
}

func (r *GoRedisEX) HMGet(key string, fields ...string) *redis.SliceCmd {
	return r.HMGetCtx(context.TODO(), key, fields...)
}
func (r *GoRedisEX) HMGetCtx(ctx context.Context, key string, fields ...string) *redis.SliceCmd {
	return r.client.HMGet(ctx, r.AddKeyPrefix(key), fields...)
}

func (r *GoRedisEX) HSet(key string, values ...interface{}) *redis.IntCmd {
	return r.HSetCtx(context.TODO(), key, values...)
}
func (r *GoRedisEX) HSetCtx(ctx context.Context, key string, values ...interface{}) *redis.IntCmd {
	return r.client.HSet(ctx, r.AddKeyPrefix(key), values...)
}

func (r *GoRedisEX) HMSet(key string, values ...interface{}) *redis.BoolCmd {
	return r.HMSetCtx(context.TODO(), key, values...)
}
func (r *GoRedisEX) HMSetCtx(ctx context.Context, key string, values ...interface{}) *redis.BoolCmd {
	return r.client.HMSet(ctx, r.AddKeyPrefix(key), values...)
}

func (r *GoRedisEX) HSetNX(key, field string, value interface{}) *redis.BoolCmd {
	return r.HSetNxCtx(context.TODO(), key, field, value)
}
func (r *GoRedisEX) HSetNxCtx(ctx context.Context, key, field string, value interface{}) *redis.BoolCmd {
	return r.client.HSetNX(ctx, r.AddKeyPrefix(key), field, value)
}

func (r *GoRedisEX) HScan(key string, cursor uint64, match string, count int64) *redis.ScanCmd {
	return r.HScanCtx(context.TODO(), key, cursor, match, count)
}
func (r *GoRedisEX) HScanCtx(ctx context.Context, key string, cursor uint64, match string, count int64) *redis.ScanCmd {
	return r.client.HScan(ctx, r.AddKeyPrefix(key), cursor, match, count)
}

func (r *GoRedisEX) HVals(key string) *redis.StringSliceCmd {
	return r.HValsCtx(context.TODO(), key)
}
func (r *GoRedisEX) HValsCtx(ctx context.Context, key string) *redis.StringSliceCmd {
	return r.client.HVals(ctx, r.AddKeyPrefix(key))
}

func (r *GoRedisEX) HRandField(key string, count int) *redis.StringSliceCmd {
	return r.HRandFieldCtx(context.TODO(), key, count)
}
func (r *GoRedisEX) HRandFieldCtx(ctx context.Context, key string, count int) *redis.StringSliceCmd {
	return r.client.HRandField(ctx, r.AddKeyPrefix(key), count)
}

func (r *GoRedisEX) HRandFieldWithValues(key string, count int) *redis.KeyValueSliceCmd {
	return r.HRandFieldWithValuesCtx(context.TODO(), key, count)
}
func (r *GoRedisEX) HRandFieldWithValuesCtx(ctx context.Context, key string, count int) *redis.KeyValueSliceCmd {
	return r.client.HRandFieldWithValues(ctx, r.AddKeyPrefix(key), count)
}
