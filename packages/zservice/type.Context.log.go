package zservice

import "fmt"

// -------- 打印消息
// 获取日志的打印信息
func (ctx *Context) logCtxStr() string {
	str := fmt.Sprintf("[%v %v-%v %v]", ctx.Trace.TraceService, ctx.Trace.TraceID, ctx.Trace.TraceSpanID, ctx.Since())
	if ctx.IsFirst_userLogCtx {
		ctx.IsFirst_userLogCtx = false
		str += fmt.Sprint("New Trace By:", ctx.Create_Path, " ")
	}
	return str
}
func (ctx *Context) LogInfo(v ...any) {
	LogInfoCaller(2, ctx.logCtxStr(), StringSprint(v...))
}
func (ctx *Context) LogInfof(f string, v ...any) {
	LogInfoCaller(2, ctx.logCtxStr(), fmt.Sprintf(f, v...))
}

func (ctx *Context) LogInfoCaller(caller int, v ...any) {
	LogInfoCaller(caller+1, ctx.logCtxStr(), StringSprint(v...))
}
func (ctx *Context) LogInfoCallerf(caller int, f string, v ...any) {
	LogInfoCaller(caller+1, ctx.logCtxStr(), fmt.Sprintf(f, v...))
}

func (ctx *Context) LogWarn(v ...any) {
	LogWarnCaller(2, ctx.logCtxStr(), StringSprint(v...))
}
func (ctx *Context) LogWarnf(f string, v ...any) {
	LogWarnCaller(2, ctx.logCtxStr(), fmt.Sprintf(f, v...))
}
func (ctx *Context) LogWarnCaller(caller int, v ...any) {
	LogWarnCaller(caller+1, ctx.logCtxStr(), StringSprint(v...))
}
func (ctx *Context) LogWarnCallerf(caller int, f string, v ...any) {
	LogWarnCaller(caller+1, ctx.logCtxStr(), fmt.Sprintf(f, v...))
}

func (ctx *Context) LogDebug(v ...any) {
	LogDebugCaller(2, ctx.logCtxStr(), StringSprint(v...))
}
func (ctx *Context) LogDebugf(f string, v ...any) {
	LogDebugCaller(2, ctx.logCtxStr(), fmt.Sprintf(f, v...))
}
func (ctx *Context) LogDebugCaller(caller int, v ...any) {
	LogDebugCaller(caller+1, ctx.logCtxStr(), StringSprint(v...))
}
func (ctx *Context) LogDebugCallerf(caller int, f string, v ...any) {
	LogDebugCaller(caller+1, ctx.logCtxStr(), fmt.Sprintf(f, v...))
}

func (ctx *Context) LogError(v ...any) {
	LogErrorCaller(2, ctx.logCtxStr(), StringSprint(v...))
}
func (ctx *Context) LogErrorf(f string, v ...any) {
	LogErrorCaller(2, ctx.logCtxStr(), fmt.Sprintf(f, v...))
}
func (ctx *Context) LogErrorCaller(caller int, v ...any) {
	LogErrorCaller(caller+1, ctx.logCtxStr(), StringSprint(v...))
}
func (ctx *Context) LogErrorCallerf(caller int, f string, v ...any) {
	LogErrorCaller(caller+1, ctx.logCtxStr(), fmt.Sprintf(f, v...))
}

func (ctx *Context) LogPanic(v ...any) {
	LogPanicCaller(2, ctx.logCtxStr(), StringSprint(v...))
}
func (ctx *Context) LogPanicf(f string, v ...any) {
	LogPanicCaller(2, ctx.logCtxStr(), fmt.Sprintf(f, v...))
}
func (ctx *Context) LogPanicCaller(caller int, v ...any) {
	LogPanicCaller(caller+1, ctx.logCtxStr(), StringSprint(v...))
}
func (ctx *Context) LogPanicCallerf(caller int, f string, v ...any) {
	LogPanicCaller(caller+1, ctx.logCtxStr(), fmt.Sprintf(f, v...))
}
