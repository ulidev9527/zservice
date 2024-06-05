package zservice

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// 上下文内部交互信息
type ContextS2S struct {
	TraceTime    time.Time `json:"tt"`  // 链路初始化时间
	TraceID      string    `json:"ti"`  // 链路ID
	TraceSpanID  int       `json:"tsi"` // 链路 , 自增处理
	TraceService string    `json:"ts"`  // 链路服务
	RequestIP    string    `json:"ip"`  // 请求IP
	Service      string    `json:"s"`   // 服务
	AuthToken    string    `json:"at"`  // token
	AuthSign     string    `json:"as"`  // 授权的签名
	ClientSign   string    `json:"cs"`  // 客户端签名
}

// 集成链路、日志、错误功能
type Context struct {
	ContextS2S
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
	// 链路记录
	if traceJsonStr != "" {
		if len(traceJsonStr) > 200 {
			traceJsonStr = traceJsonStr[:200]
		}
		e := json.Unmarshal([]byte(traceJsonStr), &ctx.ContextS2S)
		if e != nil {
			mainService.LogError(e, "[zservice.NewContext] => fail, traceJsonStr: %v", traceJsonStr)
		}
	}

	// 链路数据更新
	if ctx.ContextS2S.TraceID == "" || len(ctx.ContextS2S.TraceID) != 20 {
		ctx.ContextS2S.TraceTime = ctx.StartTime
		ctx.ContextS2S.TraceID = RandomXID()
		ctx.ContextS2S.TraceSpanID = 0
	} else {
		ctx.ContextS2S.TraceSpanID++
	}

	// 链路服务名更新
	if ctx.ContextS2S.Service == "" {
		ctx.ContextS2S.TraceService = mainService.name
	} else {
		ctx.ContextS2S.TraceService = ctx.ContextS2S.Service
	}
	ctx.ContextS2S.Service = mainService.name

	return ctx
}

// 创建一个空的上下文
func NewEmptyContext() *Context {
	return NewContext("")
}

func ContextTODO() context.Context {
	return context.TODO()
}

// -------- 打印消息
// 获取日志的打印信息
func (ctx *Context) logCtxStr() string {
	return fmt.Sprintf("[%v %v-%v %v]", ctx.Service.tranceName, ctx.TraceID, ctx.TraceSpanID, ctx.SinceTrace())
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
