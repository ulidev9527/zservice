package dbservice

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

func (r *GoRedisEX) Del(keys ...string) *redis.IntCmd {
	return r.DelCtx(context.TODO(), keys...)
}
func (r *GoRedisEX) DelCtx(ctx context.Context, keys ...string) *redis.IntCmd {
	return r.client.Del(ctx, r.AddkeyPrefixs(keys...)...)
}

func (r *GoRedisEX) Dump(key string) *redis.StringCmd {
	return r.DumpCtx(context.TODO(), key)
}
func (r *GoRedisEX) DumpCtx(ctx context.Context, key string) *redis.StringCmd {
	return r.client.Dump(ctx, r.AddKeyPrefix(key))
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

func (r *GoRedisEX) ExpireAt(key string, expiration time.Time) *redis.BoolCmd {
	return r.ExpireAtCtx(context.TODO(), key, expiration)
}
func (r *GoRedisEX) ExpireAtCtx(ctx context.Context, key string, expiration time.Time) *redis.BoolCmd {
	return r.client.ExpireAt(ctx, r.AddKeyPrefix(key), expiration)
}

func (r *GoRedisEX) ExpireTime(key string) *redis.DurationCmd {
	return r.ExpireTimeCtx(context.TODO(), key)
}
func (r *GoRedisEX) ExpireTimeCtx(ctx context.Context, key string) *redis.DurationCmd {
	return r.client.ExpireTime(ctx, r.AddKeyPrefix(key))
}

func (r *GoRedisEX) ExpireNX(key string, expiration time.Duration) *redis.BoolCmd {
	return r.ExpireNXCtx(context.TODO(), key, expiration)
}
func (r *GoRedisEX) ExpireNXCtx(ctx context.Context, key string, expiration time.Duration) *redis.BoolCmd {
	return r.client.ExpireNX(ctx, r.AddKeyPrefix(key), expiration)
}

func (r *GoRedisEX) ExpireXX(key string, expiration time.Duration) *redis.BoolCmd {
	return r.ExpireXXCtx(context.TODO(), key, expiration)
}
func (r *GoRedisEX) ExpireXXCtx(ctx context.Context, key string, expiration time.Duration) *redis.BoolCmd {
	return r.client.ExpireXX(ctx, r.AddKeyPrefix(key), expiration)
}

func (r *GoRedisEX) ExpireGT(key string, expiration time.Duration) *redis.BoolCmd {
	return r.ExpireGTCtx(context.TODO(), key, expiration)
}
func (r *GoRedisEX) ExpireGTCtx(ctx context.Context, key string, expiration time.Duration) *redis.BoolCmd {
	return r.client.ExpireGT(ctx, r.AddKeyPrefix(key), expiration)
}

func (r *GoRedisEX) ExpireLT(key string, expiration time.Duration) *redis.BoolCmd {
	return r.ExpireLTCtx(context.TODO(), key, expiration)
}
func (r *GoRedisEX) ExpireLTCtx(ctx context.Context, key string, expiration time.Duration) *redis.BoolCmd {
	return r.client.ExpireLT(ctx, r.AddKeyPrefix(key), expiration)
}

// 注意，返回数据包含 key 前缀
func (r *GoRedisEX) Keys(pattern string) *redis.StringSliceCmd {
	return r.KeysCtx(context.TODO(), pattern)
}

// 注意，返回数据包含 key 前缀
func (r *GoRedisEX) KeysCtx(ctx context.Context, pattern string) *redis.StringSliceCmd {
	return r.client.Keys(ctx, r.AddKeyPrefix(pattern))
}

func (r *GoRedisEX) Migrate(host, port, key string, db int, timeout time.Duration) *redis.StatusCmd {
	return r.MigrateCtx(context.TODO(), host, port, key, db, timeout)
}
func (r *GoRedisEX) MigrateCtx(ctx context.Context, host, port, key string, db int, timeout time.Duration) *redis.StatusCmd {
	return r.client.Migrate(ctx, host, port, r.AddKeyPrefix(key), db, timeout)
}

func (r *GoRedisEX) Move(key string, db int) *redis.BoolCmd {
	return r.MoveCtx(context.TODO(), key, db)
}
func (r *GoRedisEX) MoveCtx(ctx context.Context, key string, db int) *redis.BoolCmd {
	return r.client.Move(ctx, r.AddKeyPrefix(key), db)
}

func (r *GoRedisEX) ObjectFreq(key string) *redis.IntCmd {
	return r.ObjectFreqCtx(context.TODO(), key)
}
func (r *GoRedisEX) ObjectFreqCtx(ctx context.Context, key string) *redis.IntCmd {
	return r.client.ObjectFreq(ctx, r.AddKeyPrefix(key))
}

func (r *GoRedisEX) ObjectRefCount(key string) *redis.IntCmd {
	return r.ObjectRefCountCtx(context.TODO(), key)
}
func (r *GoRedisEX) ObjectRefCountCtx(ctx context.Context, key string) *redis.IntCmd {
	return r.client.ObjectRefCount(ctx, r.AddKeyPrefix(key))
}

func (r *GoRedisEX) ObjectEncoding(key string) *redis.StringCmd {
	return r.ObjectEncodingCtx(context.TODO(), key)
}
func (r *GoRedisEX) ObjectEncodingCtx(ctx context.Context, key string) *redis.StringCmd {
	return r.client.ObjectEncoding(ctx, r.AddKeyPrefix(key))
}

func (r *GoRedisEX) ObjectIdleTime(key string) *redis.DurationCmd {
	return r.ObjectIdleTimeCtx(context.TODO(), key)
}
func (r *GoRedisEX) ObjectIdleTimeCtx(ctx context.Context, key string) *redis.DurationCmd {
	return r.client.ObjectIdleTime(ctx, r.AddKeyPrefix(key))
}

func (r *GoRedisEX) Persist(key string) *redis.BoolCmd {
	return r.PersistCtx(context.TODO(), key)
}
func (r *GoRedisEX) PersistCtx(ctx context.Context, key string) *redis.BoolCmd {
	return r.client.Persist(ctx, r.AddKeyPrefix(key))
}

func (r *GoRedisEX) PExpire(key string, expiration time.Duration) *redis.BoolCmd {
	return r.PExpireCtx(context.TODO(), key, expiration)
}
func (r *GoRedisEX) PExpireCtx(ctx context.Context, key string, expiration time.Duration) *redis.BoolCmd {
	return r.client.PExpire(ctx, r.AddKeyPrefix(key), expiration)
}

func (r *GoRedisEX) PExpireAt(key string, t time.Time) *redis.BoolCmd {
	return r.PExpireAtCtx(context.TODO(), key, t)
}
func (r *GoRedisEX) PExpireAtCtx(ctx context.Context, key string, t time.Time) *redis.BoolCmd {
	return r.client.PExpireAt(ctx, r.AddKeyPrefix(key), t)
}

func (r *GoRedisEX) PExpireTime(key string) *redis.DurationCmd {
	return r.PExpireTimeCtx(context.TODO(), key)
}
func (r *GoRedisEX) PExpireTimeCtx(ctx context.Context, key string) *redis.DurationCmd {
	return r.client.PExpireTime(ctx, r.AddKeyPrefix(key))
}

func (r *GoRedisEX) PTTL(key string) *redis.DurationCmd {
	return r.PTTLCtx(context.TODO(), key)
}
func (r *GoRedisEX) PTTLCtx(ctx context.Context, key string) *redis.DurationCmd {
	return r.client.PTTL(ctx, r.AddKeyPrefix(key))
}

func (r *GoRedisEX) RandomKey() *redis.StringCmd {
	return r.RandomKeyCtx(context.TODO())
}
func (r *GoRedisEX) RandomKeyCtx(ctx context.Context) *redis.StringCmd {
	return r.client.RandomKey(ctx)
}

func (r *GoRedisEX) Rename(key, newkey string) *redis.StatusCmd {
	return r.RenameCtx(context.TODO(), key, newkey)
}
func (r *GoRedisEX) RenameCtx(ctx context.Context, key, newkey string) *redis.StatusCmd {
	return r.client.Rename(ctx, r.AddKeyPrefix(key), r.AddKeyPrefix(newkey))
}

func (r *GoRedisEX) RenameNX(key, newkey string) *redis.BoolCmd {
	return r.RenameNXCtx(context.TODO(), key, newkey)
}
func (r *GoRedisEX) RenameNXCtx(ctx context.Context, key, newkey string) *redis.BoolCmd {
	return r.client.RenameNX(ctx, r.AddKeyPrefix(key), r.AddKeyPrefix(newkey))
}

func (r *GoRedisEX) Restore(key string, ttl time.Duration, value string) *redis.StatusCmd {
	return r.RestoreCtx(context.TODO(), key, ttl, value)
}
func (r *GoRedisEX) RestoreCtx(ctx context.Context, key string, ttl time.Duration, value string) *redis.StatusCmd {
	return r.client.Restore(ctx, r.AddKeyPrefix(key), ttl, value)
}

func (r *GoRedisEX) RestoreReplace(key string, ttl time.Duration, value string) *redis.StatusCmd {
	return r.RestoreReplaceCtx(context.TODO(), key, ttl, value)
}
func (r *GoRedisEX) RestoreReplaceCtx(ctx context.Context, key string, ttl time.Duration, value string) *redis.StatusCmd {
	return r.client.RestoreReplace(ctx, r.AddKeyPrefix(key), ttl, value)
}

func (r *GoRedisEX) Sort(ctx context.Context, key string, sort *redis.Sort) *redis.StringSliceCmd {
	return r.SortCtx(ctx, key, sort)
}
func (r *GoRedisEX) SortCtx(ctx context.Context, key string, sort *redis.Sort) *redis.StringSliceCmd {
	return r.client.Sort(ctx, r.AddKeyPrefix(key), sort)
}

func (r *GoRedisEX) SortRO(ctx context.Context, key string, sort *redis.Sort) *redis.StringSliceCmd {
	return r.SortROCtx(ctx, key, sort)
}
func (r *GoRedisEX) SortROCtx(ctx context.Context, key string, sort *redis.Sort) *redis.StringSliceCmd {
	return r.client.SortRO(ctx, r.AddKeyPrefix(key), sort)
}

func (r *GoRedisEX) SortStore(ctx context.Context, key, store string, sort *redis.Sort) *redis.IntCmd {
	return r.SortStoreCtx(ctx, key, store, sort)
}
func (r *GoRedisEX) SortStoreCtx(ctx context.Context, key, store string, sort *redis.Sort) *redis.IntCmd {
	return r.client.SortStore(ctx, r.AddKeyPrefix(key), r.AddKeyPrefix(store), sort)
}

func (r *GoRedisEX) SortInterfaces(ctx context.Context, key string, sort *redis.Sort) *redis.SliceCmd {
	return r.SortInterfacesCtx(ctx, key, sort)
}
func (r *GoRedisEX) SortInterfacesCtx(ctx context.Context, key string, sort *redis.Sort) *redis.SliceCmd {
	return r.client.SortInterfaces(ctx, r.AddKeyPrefix(key), sort)
}

func (r *GoRedisEX) Touch(keys ...string) *redis.IntCmd {
	return r.TouchCtx(context.TODO(), keys...)
}
func (r *GoRedisEX) TouchCtx(ctx context.Context, keys ...string) *redis.IntCmd {
	return r.client.Touch(ctx, r.AddkeyPrefixs(keys...)...)
}

func (r *GoRedisEX) TTL(key string) *redis.DurationCmd {
	return r.TTLCtx(context.TODO(), key)
}
func (r *GoRedisEX) TTLCtx(ctx context.Context, key string) *redis.DurationCmd {
	return r.client.TTL(ctx, r.AddKeyPrefix(key))
}

func (r *GoRedisEX) Type(key string) *redis.StatusCmd {
	return r.TypeCtx(context.TODO(), key)
}
func (r *GoRedisEX) TypeCtx(ctx context.Context, key string) *redis.StatusCmd {
	return r.client.Type(ctx, r.AddKeyPrefix(key))
}

func (r *GoRedisEX) Copy(key, newkey string, db int, replace bool) *redis.IntCmd {
	return r.CopyCtx(context.TODO(), key, newkey, db, replace)
}
func (r *GoRedisEX) CopyCtx(ctx context.Context, key, newkey string, db int, replace bool) *redis.IntCmd {
	return r.client.Copy(ctx, r.AddKeyPrefix(key), r.AddKeyPrefix(newkey), db, replace)
}

// 注意，返回数据包含 key 前缀
func (r *GoRedisEX) Scan(cursor uint64, match string, count int64) *redis.ScanCmd {
	return r.ScanCtx(context.TODO(), cursor, match, count)
}
func (r *GoRedisEX) ScanCtx(ctx context.Context, cursor uint64, match string, count int64) *redis.ScanCmd {
	return r.client.Scan(ctx, cursor, r.AddKeyPrefix(match), count)
}

// 注意，返回数据包含 key 前缀
func (r *GoRedisEX) ScanType(cursor uint64, match string, count int64, keyType string) *redis.ScanCmd {
	return r.ScanTypeCtx(context.TODO(), cursor, match, count, keyType)
}
func (r *GoRedisEX) ScanTypeCtx(ctx context.Context, cursor uint64, match string, count int64, keyType string) *redis.ScanCmd {
	return r.client.ScanType(ctx, cursor, r.AddKeyPrefix(match), count, keyType)
}
