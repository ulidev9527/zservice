package zservice

import (
	"encoding/json"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type zcontextTrace struct {
	TraceTime time.Time `json:"tt"`  // 链路初始化时间
	TraceID   string    `json:"tid"` // 链路ID
	SpanID    int       `json:"sid"` // 链路 , 自增处理
}

// 集成链路、日志、错误功能
type ZContext struct {
	zcontextTrace
	StartTime  time.Time // 当前上下文启动时间
	Service    *ZService // 服务
	CTX_mu     sync.Mutex
	CTX_done   atomic.Value
	CTX_err    error
	CTX_values sync.Map
}

// 创建上下文
func NewContext(s *ZService, traceJsonStr string) *ZContext {
	ctx := &ZContext{
		StartTime:  time.Now(),
		Service:    s,
		CTX_mu:     sync.Mutex{},
		CTX_values: sync.Map{},
	}
	if traceJsonStr != "" {
		e := json.Unmarshal([]byte(traceJsonStr), &ctx.zcontextTrace)
		if e != nil {
			s.LogError(e, "[zservice.NewContext] => fail, traceJsonStr: %v", traceJsonStr)
		}

		ctx.StartTime = time.Now()
		ctx.SpanID++

		return ctx
	} else {
		ctx.zcontextTrace = zcontextTrace{
			TraceTime: ctx.StartTime,
			TraceID:   RandomXID(),
			SpanID:    0,
		}
	}
	return ctx
}

// 创建一个空的上下文
func NewEmptyContext() *ZContext {
	return NewContext(mainService, "")
}

// -------- 打印消息
// 获取日志的打印信息
func (ctx *ZContext) logCtxStr() string {
	return fmt.Sprintf("[%v %v-%v %v]", ctx.Service.tranceName, ctx.TraceID, ctx.SpanID, ctx.SinceTrace())
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

func (ctx *ZContext) LogPanic(v ...any) {
	LogPanicCaller(2, ctx.logCtxStr(), Sprint(v...))
}
func (ctx *ZContext) LogPanicf(f string, v ...any) {
	LogPanicCaller(2, ctx.logCtxStr(), fmt.Sprintf(f, v...))
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
	d := ctx.CTX_done.Load()
	if d != nil {
		return d.(chan struct{})
	}
	ctx.CTX_mu.Lock()
	defer ctx.CTX_mu.Unlock()
	d = ctx.CTX_done.Load()
	if d == nil {
		d = make(chan struct{})
		ctx.CTX_done.Store(d)
	}
	return d.(chan struct{})
}

// Err implements context.Context.
func (ctx *ZContext) Err() error {
	return ctx.CTX_err
}

// Value implements context.Context.
func (ctx *ZContext) Value(key any) any {
	v, _ := ctx.CTX_values.Load(key)
	return v
}
