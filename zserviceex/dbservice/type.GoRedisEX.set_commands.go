package dbservice

import (
	"context"

	"github.com/redis/go-redis/v9"
)

func (r *GoRedisEX) SAdd(key string, members ...interface{}) *redis.IntCmd {
	return r.SAddCtx(context.TODO(), key, members...)
}
func (r *GoRedisEX) SAddCtx(ctx context.Context, key string, members ...interface{}) *redis.IntCmd {
	return r.client.SAdd(ctx, r.AddKeyPrefix(key), members...)
}

func (r *GoRedisEX) SCard(key string) *redis.IntCmd {
	return r.SCardCtx(context.TODO(), key)
}
func (r *GoRedisEX) SCardCtx(ctx context.Context, key string) *redis.IntCmd {
	return r.client.SCard(ctx, r.AddKeyPrefix(key))
}

func (r *GoRedisEX) SDiff(keys ...string) *redis.StringSliceCmd {
	return r.SDiffCtx(context.TODO(), keys...)
}
func (r *GoRedisEX) SDiffCtx(ctx context.Context, keys ...string) *redis.StringSliceCmd {
	return r.client.SDiff(ctx, r.AddkeyPrefixs(keys...)...)
}

func (r *GoRedisEX) SDiffStore(destination string, keys ...string) *redis.IntCmd {
	return r.SDiffStoreCtx(context.TODO(), destination, keys...)
}
func (r *GoRedisEX) SDiffStoreCtx(ctx context.Context, destination string, keys ...string) *redis.IntCmd {
	return r.client.SDiffStore(ctx, r.AddKeyPrefix(destination), r.AddkeyPrefixs(keys...)...)
}

func (r *GoRedisEX) SInter(keys ...string) *redis.StringSliceCmd {
	return r.SInterCtx(context.TODO(), keys...)
}
func (r *GoRedisEX) SInterCtx(ctx context.Context, keys ...string) *redis.StringSliceCmd {
	return r.client.SInter(ctx, r.AddkeyPrefixs(keys...)...)
}

func (r *GoRedisEX) SInterStore(destination string, keys ...string) *redis.IntCmd {
	return r.SInterStoreCtx(context.TODO(), destination, keys...)
}
func (r *GoRedisEX) SInterStoreCtx(ctx context.Context, destination string, keys ...string) *redis.IntCmd {
	return r.client.SInterStore(ctx, r.AddKeyPrefix(destination), r.AddkeyPrefixs(keys...)...)
}

func (r *GoRedisEX) SIsMember(key string, member interface{}) *redis.BoolCmd {
	return r.SIsMemberCtx(context.TODO(), key, member)
}
func (r *GoRedisEX) SIsMemberCtx(ctx context.Context, key string, member interface{}) *redis.BoolCmd {
	return r.client.SIsMember(ctx, r.AddKeyPrefix(key), member)
}

func (r *GoRedisEX) SMIsMember(key string, members ...interface{}) *redis.BoolSliceCmd {
	return r.SMIsMemberCtx(context.TODO(), key, members...)
}
func (r *GoRedisEX) SMIsMemberCtx(ctx context.Context, key string, members ...interface{}) *redis.BoolSliceCmd {
	return r.client.SMIsMember(ctx, r.AddKeyPrefix(key), members...)
}

func (r *GoRedisEX) SMembers(key string) *redis.StringSliceCmd {
	return r.SMembersCtx(context.TODO(), key)
}
func (r *GoRedisEX) SMembersCtx(ctx context.Context, key string) *redis.StringSliceCmd {
	return r.client.SMembers(ctx, r.AddKeyPrefix(key))
}

func (r *GoRedisEX) SMembersMap(key string) *redis.StringStructMapCmd {
	return r.SMembersMapCtx(context.TODO(), key)
}
func (r *GoRedisEX) SMembersMapCtx(ctx context.Context, key string) *redis.StringStructMapCmd {
	return r.client.SMembersMap(ctx, r.AddKeyPrefix(key))
}

func (r *GoRedisEX) SMove(source, destination string, member interface{}) *redis.BoolCmd {
	return r.SMoveCtx(context.TODO(), source, destination, member)
}
func (r *GoRedisEX) SMoveCtx(ctx context.Context, source, destination string, member interface{}) *redis.BoolCmd {
	return r.client.SMove(ctx, r.AddKeyPrefix(source), r.AddKeyPrefix(destination), member)
}

func (r *GoRedisEX) SPop(key string) *redis.StringCmd {
	return r.SPopCtx(context.TODO(), key)
}
func (r *GoRedisEX) SPopCtx(ctx context.Context, key string) *redis.StringCmd {
	return r.client.SPop(ctx, r.AddKeyPrefix(key))
}

func (r *GoRedisEX) SPopN(key string, count int64) *redis.StringSliceCmd {
	return r.SPopNCtx(context.TODO(), key, count)
}
func (r *GoRedisEX) SPopNCtx(ctx context.Context, key string, count int64) *redis.StringSliceCmd {
	return r.client.SPopN(ctx, r.AddKeyPrefix(key), count)
}

func (r *GoRedisEX) SRandMember(key string) *redis.StringCmd {
	return r.SRandMemberCtx(context.TODO(), key)
}
func (r *GoRedisEX) SRandMemberCtx(ctx context.Context, key string) *redis.StringCmd {
	return r.client.SRandMember(ctx, r.AddKeyPrefix(key))
}

func (r *GoRedisEX) SRandMemberN(key string, count int64) *redis.StringSliceCmd {
	return r.SRandMemberNCtx(context.TODO(), key, count)
}
func (r *GoRedisEX) SRandMemberNCtx(ctx context.Context, key string, count int64) *redis.StringSliceCmd {
	return r.client.SRandMemberN(ctx, r.AddKeyPrefix(key), count)
}

func (r *GoRedisEX) SRem(key string, members ...interface{}) *redis.IntCmd {
	return r.SRemCtx(context.TODO(), key, members...)
}
func (r *GoRedisEX) SRemCtx(ctx context.Context, key string, members ...interface{}) *redis.IntCmd {
	return r.client.SRem(ctx, r.AddKeyPrefix(key), members...)
}

func (r *GoRedisEX) SScan(key string, cursor uint64, match string, count int64) *redis.ScanCmd {
	return r.SScanCtx(context.TODO(), key, cursor, match, count)
}
func (r *GoRedisEX) SScanCtx(ctx context.Context, key string, cursor uint64, match string, count int64) *redis.ScanCmd {
	return r.client.SScan(ctx, r.AddKeyPrefix(key), cursor, match, count)
}

func (r *GoRedisEX) SUnion(keys ...string) *redis.StringSliceCmd {
	return r.SUnionCtx(context.TODO(), keys...)
}
func (r *GoRedisEX) SUnionCtx(ctx context.Context, keys ...string) *redis.StringSliceCmd {
	return r.client.SUnion(ctx, r.AddkeyPrefixs(keys...)...)
}
