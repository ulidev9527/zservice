package zservice

import (
	"encoding/json"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type ContextTrace struct {
	TraceTime time.Time `json:"tt"`  // 链路初始化时间
	TraceID   string    `json:"tid"` // 链路ID
	SpanID    int       `json:"sid"` // 链路 , 自增处理
}

// 集成链路、日志、错误功能
type Context struct {
	ContextTrace
	StartTime  time.Time // 当前上下文启动时间
	Service    *ZService // 服务
	CTX_mu     sync.Mutex
	CTX_done   atomic.Value
	CTX_err    error
	CTX_values sync.Map
}

// 创建上下文
func NewContext(traceJsonStr string) *Context {
	ctx := &Context{
		StartTime:  time.Now(),
		Service:    mainService,
		CTX_mu:     sync.Mutex{},
		CTX_values: sync.Map{},
	}
	if traceJsonStr != "" {
		e := json.Unmarshal([]byte(traceJsonStr), &ctx.ContextTrace)
		if e != nil {
			mainService.LogError(e, "[zservice.NewContext] => fail, traceJsonStr: %v", traceJsonStr)
		}

		ctx.StartTime = time.Now()
		ctx.SpanID++

		return ctx
	} else {
		ctx.ContextTrace = ContextTrace{
			TraceTime: ctx.StartTime,
			TraceID:   RandomXID(),
			SpanID:    0,
		}
	}
	return ctx
}

// 创建一个空的上下文
func NewEmptyContext() *Context {
	return NewContext("")
}

// -------- 打印消息
// 获取日志的打印信息
func (ctx *Context) logCtxStr() string {
	return fmt.Sprintf("[%v %v-%v %v]", ctx.Service.tranceName, ctx.TraceID, ctx.SpanID, ctx.SinceTrace())
}
func (ctx *Context) LogInfo(v ...any) {
	LogInfoCaller(2, ctx.logCtxStr(), Sprint(v...))
}
func (ctx *Context) LogInfof(f string, v ...any) {
	LogInfoCaller(2, ctx.logCtxStr(), fmt.Sprintf(f, v...))
}
func (ctx *Context) LogWarn(v ...any) {
	LogWarnCaller(2, ctx.logCtxStr(), Sprint(v...))
}
func (ctx *Context) LogWarnf(f string, v ...any) {
	LogWarnCaller(2, ctx.logCtxStr(), fmt.Sprintf(f, v...))
}
func (ctx *Context) LogError(v ...any) {
	LogErrorCaller(2, ctx.logCtxStr(), Sprint(v...))
}
func (ctx *Context) LogErrorf(f string, v ...any) {
	LogErrorCaller(2, ctx.logCtxStr(), fmt.Sprintf(f, v...))
}

func (ctx *Context) LogPanic(v ...any) {
	LogPanicCaller(2, ctx.logCtxStr(), Sprint(v...))
}
func (ctx *Context) LogPanicf(f string, v ...any) {
	LogPanicCaller(2, ctx.logCtxStr(), fmt.Sprintf(f, v...))
}

// 获取上下文创建到现在的时间
func (ctx *Context) Since() time.Duration {
	return time.Since(ctx.StartTime)
}

// 获取链路创建到现在的时间
func (ctx *Context) SinceTrace() time.Duration {
	return time.Since(ctx.TraceTime)
}

// Deadline implements context.Context.
func (ctx *Context) Deadline() (deadline time.Time, ok bool) {
	return
}

// Done implements context.Context.
func (ctx *Context) Done() <-chan struct{} {
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
func (ctx *Context) Err() error {
	return ctx.CTX_err
}

// Value implements context.Context.
func (ctx *Context) Value(key any) any {
	v, _ := ctx.CTX_values.Load(key)
	return v
}
