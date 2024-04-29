package zservice

import (
	"encoding/json"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// 集成链路、日志、错误功能
type ZContext struct {
	StartTime time.Time // 当前上下文启动时间
	TraceTime time.Time // 链路初始化时间
	TraceID   string    // 链路ID
	SpanID    int       // 链路 , 自增处理
	service   *ZService // 服务
	mu        sync.Mutex
	done      atomic.Value
	err       error
	values    sync.Map
}

// 创建上下文
func NewContext(s *ZService, traceJsonStr string) *ZContext {
	ctx := &ZContext{
		service: s,
		mu:      sync.Mutex{},
		values:  sync.Map{},
	}
	if traceJsonStr != "" {
		e := json.Unmarshal([]byte(traceJsonStr), &ctx)
		if e != nil {
			s.LogError(e, "[zserver.NewContext] => fail, traceJsonStr: %v", traceJsonStr)
		}

		ctx.StartTime = time.Now()
		ctx.SpanID++

		return ctx
	}
	t := time.Now()
	return &ZContext{
		StartTime: t,
		TraceTime: t,
		TraceID:   RandomXID(),
		SpanID:    0,
	}
}

// -------- 打印消息
// 获取日志的打印信息
func (ctx *ZContext) logCtxStr() string {
	return fmt.Sprintf("[%v %v-%v %v]", ctx.service.tranceName, ctx.TraceID, ctx.SpanID, ctx.SinceTrace())
}
func (ctx *ZContext) LogInfo(v ...any) {
	LogInfoCaller(2, ctx.logCtxStr(), Sprint(v...))
}
func (ctx *ZContext) LogInfof(f string, v ...any) {
	LogInfoCaller(2, ctx.logCtxStr(), fmt.Sprintf(f, v...))
}
func (ctx *ZContext) LogWarn(v ...any) {
	LogWarnCaller(2, ctx.logCtxStr(), Sprint(v...))
}
func (ctx *ZContext) LogWarnf(f string, v ...any) {
	LogWarnCaller(2, ctx.logCtxStr(), fmt.Sprintf(f, v...))
}
func (ctx *ZContext) LogError(v ...any) {
	LogErrorCaller(2, ctx.logCtxStr(), Sprint(v...))
}
func (ctx *ZContext) LogErrorf(f string, v ...any) {
	LogErrorCaller(2, ctx.logCtxStr(), fmt.Sprintf(f, v...))
}

// 获取上下文创建到现在的时间
func (ctx *ZContext) Since() time.Duration {
	return time.Since(ctx.StartTime)
}

// 获取链路创建到现在的时间
func (ctx *ZContext) SinceTrace() time.Duration {
	return time.Since(ctx.TraceTime)
}

// Deadline implements context.Context.
func (ctx *ZContext) Deadline() (deadline time.Time, ok bool) {
	return
}

// Done implements context.Context.
func (ctx *ZContext) Done() <-chan struct{} {
	d := ctx.done.Load()
	if d != nil {
		return d.(chan struct{})
	}
	ctx.mu.Lock()
	defer ctx.mu.Unlock()
	d = ctx.done.Load()
	if d == nil {
		d = make(chan struct{})
		ctx.done.Store(d)
	}
	return d.(chan struct{})
}

// Err implements context.Context.
func (ctx *ZContext) Err() error {
	return ctx.err
}

// Value implements context.Context.
func (ctx *ZContext) Value(key any) any {
	v, _ := ctx.values.Load(key)
	return v
}
