package zservice

import (
	"context"
	"encoding/json"
	"fmt"
	"runtime"
	"sync"
	"time"
)

// 上下文内部链路信息
type ContextTrace struct {
	TraceService string    `json:"ts,omitempty"`  // 链路服务
	TraceTime    time.Time `json:"tt,omitempty"`  // 链路初始化时间
	TraceID      string    `json:"ti,omitempty"`  // 链路ID
	TraceSpanID  int       `json:"tsi,omitempty"` // 链路 , 自增处理
}

// 集成链路、日志、错误功能
type Context struct {
	context.Context
	isNew              bool               // 是否是新上下文
	IsFirst_userLogCtx bool               `json:"if"` // 是否首次使用日志
	Create_Path        string             `json:"cp"` // 创建路径
	ctx_mu             sync.Mutex         `json:"-"`
	ctx_err            error              `json:"-"`
	ctx_values         sync.Map           `json:"-"`
	ctx_parent         context.Context    `json:"-"`
	ctx_cancel         context.CancelFunc `json:"-"`

	Authorization string        `json:"at,omitempty"` // 授权
	Trace         *ContextTrace `json:"tr,omitempty"` // 链路信息
	IsDebug       bool          `json:"id,omitempty"` // 是否是调试模式
}

// 创建上下文
func newContext_ByDefault() *Context {
	ctx := &Context{
		isNew:              true,
		IsFirst_userLogCtx: true,
		ctx_mu:             sync.Mutex{},
		ctx_values:         sync.Map{},
		Trace: &ContextTrace{
			TraceService: GetMainService().tranceName,
			TraceTime:    time.Now(),
			TraceID:      RandomXID(),
			TraceSpanID:  0,
		},
		IsDebug: GetMainService().IsDebug,
	}
	return ctx
}

// 创建上下文
func NewContext(ctxStrIn ...string) *Context {

	ctx := newContext_ByDefault()
	ctxStr := ""
	if len(ctxStrIn) > 0 {
		ctxStr = ctxStrIn[0]
	}

	if ctxStr != "" {
		// 长度验证
		if len(ctxStr) > 1024 {
			ctxStr = ctxStr[:1024]
			ctx.LogError("ctxStr length error:", ctxStr)
		} else if e := json.Unmarshal([]byte(ctxStr), ctx); e != nil {
			ctx.LogError("ctx json unmarshal error:", e)
		} else {
			ctx.Trace.TraceService = GetMainService().tranceName
			ctx.Trace.TraceSpanID++
		}
	} else {
		_, file, line, _ := runtime.Caller(1)
		ctx.Create_Path = fmt.Sprint(file, ":", line)
	}
	return ctx
}

func ContextTODO() context.Context {
	return context.TODO()
}
func (ctx *Context) SetTranceName(name string) *Context {

	ctx.Trace.TraceService = fmt.Sprintf("%s/%s", GetMainService().tranceName, "name")

	return ctx
}

// 创建错误
func (ctx *Context) NewError(format string, a ...any) error {
	return NewErrorf("[%s] %s", ctx.logCtxStr(), fmt.Sprintf(format, a...))
}

// 获取上下文创建到现在的时间
func (ctx *Context) Since() time.Duration {
	return time.Since(ctx.Trace.TraceTime)
}

// 获取链路创建到现在的时间
func (ctx *Context) SinceTrace() time.Duration {
	return time.Since(ctx.Trace.TraceTime)
}

// Deadline implements context.Context.
func (ctx *Context) Deadline() (deadline time.Time, ok bool) {
	return
}

// Done implements context.Context.
func (ctx *Context) Done() <-chan struct{} {
	if ctx.ctx_parent == nil {
		return nil
	}

	return ctx.ctx_parent.Done()
}

func (ctx *Context) Cancel() {
	if ctx.ctx_cancel != nil {
		ctx.ctx_cancel()
	}
}

// Err implements context.Context.
func (ctx *Context) Err() error {
	return ctx.ctx_err
}

// Value implements context.Context.
func (ctx *Context) Value(key any) any {
	v, _ := ctx.ctx_values.Load(key)
	return v
}

// 获取上下文字符串
func (ctx *Context) ToContextString() string {
	return JsonMustMarshalString(ctx)
}

// 克隆
func (ctx *Context) Clone() *Context {
	clone := NewContext(ctx.ToContextString())
	clone.Trace.TraceService = ctx.Trace.TraceService
	clone.Trace.TraceSpanID = ctx.Trace.TraceSpanID
	return clone
}
