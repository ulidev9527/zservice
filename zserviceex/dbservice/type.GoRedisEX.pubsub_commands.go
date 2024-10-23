package dbservice

import (
	"context"
	"encoding/json"
	"runtime"

	"github.com/redis/go-redis/v9"
	"github.com/ulidev9527/zservice/zservice"
)

func (r *GoRedisEX) Pub(channel string, data []byte) *zservice.Error {
	return r.PubCtx(zservice.NewContext(), channel, data)
}

func (r *GoRedisEX) PubCtx(ctx *zservice.Context, channel string, data []byte) *zservice.Error {

	body := &PubsubBody{
		S2S: ctx.GetS2S(),
		Val: data,
	}

	err := r.PublishCtx(ctx, channel, string(zservice.JsonMustMarshal(body))).Err()
	if err != nil {
		return zservice.NewError(err)
	}
	return nil
}

// 订阅消息，返回取消函数
func (r *GoRedisEX) Sub(channel string, cb func(*zservice.Context, string, []byte)) *redis.PubSub {
	return r.SubCtx(zservice.NewContext(), channel, cb)
}
func (r *GoRedisEX) SubCtx(ctx *zservice.Context, channel string, cb func(*zservice.Context, string, []byte)) *redis.PubSub {
	sub := r.client.Subscribe(ctx, r.AddKeyPrefix(channel))

	zservice.Go(func() {

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

func (r *GoRedisEX) Publish(channel string, message any) *redis.IntCmd {
	return r.PublishCtx(zservice.ContextTODO(), channel, message)
}
func (r *GoRedisEX) PublishCtx(ctx context.Context, channel string, message any) *redis.IntCmd {
	return r.client.Publish(ctx, r.AddKeyPrefix(channel), message)
}
