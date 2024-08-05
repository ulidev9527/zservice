package zservice

import (
	"fmt"
	"runtime"
)

// ------------- Error -------------
type Error struct {
	code uint32 // 错误码
	log  string // 日志, 仅用于日志收集和打印
	msg  string // 消息, 返回给客户端显示
}

func NewError(v ...any) *Error {
	return NewErrorCaller(2, Sprint(v...), nil)
}
func NewErrorCaller(skip int, str string, e error) *Error {
	code := uint32(Code_Fail)
	s := str
	if e != nil {
		s = fmt.Sprint(s, " | ", e.Error())
	}
	err := &Error{
		code: code,
		log:  s,
	}
	return err.AddCaller(skip)
}
func NewErrore(e error) *Error {
	return NewErrorCaller(2, "", e)
}

func NewErrorf(f string, v ...any) *Error {
	return NewErrorCaller(2, fmt.Sprintf(f, v...), nil)
}
func NewErroref(e error, f string, v ...any) *Error {
	return NewErrorCaller(2, fmt.Sprintf(f, v...), e)
}
func (e *Error) Error() string {
	if e.code != 0 {
		return fmt.Sprintf("%d:%T:\n%s", e.code, &Error{}, e.log)
	} else {
		return fmt.Sprintf("%T:\n%s", &Error{}, e.log)
	}
}
func (e *Error) String() string {
	return e.log
}

// 添加路径记录
// skip 一般不需要填写，涉及到 caller 跳层问题
func (e *Error) AddCaller(skips ...int) *Error {
	skip := 1
	if len(skips) > 0 {
		skip = skips[0] + 1
	}
	_, file, line, _ := runtime.Caller(skip)
	e.log = fmt.Sprint(file, ":", line, " > ", e.log)
	return e
}

// 设置错误码
func (e *Error) SetCode(code uint32) *Error {
	e.code = code

	switch e.code {
	case Code_Zero:
		e.msg = "未知业务错误"
	case Code_Succ:
		e.msg = "成功"
	case Code_Fail:
		e.msg = "失败"
	case Code_Limit:
		e.msg = "业务限制"
	case Code_Auth:
		e.msg = "授权失败"
	case Code_NotImplement:
		e.msg = "功能开发中"
	case Code_Params:
		e.msg = "参数错误"
	case Code_Again:
		e.msg = "需重试"
	case Code_NotFound:
		e.msg = "未查询到数据"
	case Code_Repetition:
		e.msg = "数据重复"
	case Code_Reject:
		e.msg = "服务器拒绝处理"
	case Code_Fatal:
		e.msg = "服务器内部错误"
	}

	return e
}

// 获取错误码
func (e *Error) GetCode() uint32 {
	return e.code
}

// 设置客户端消息
func (e *Error) SetMsg(msg ...any) *Error {
	if len(msg) > 0 {
		e.msg = Sprint(msg...)
	}
	return e
}

// 获取客户端消息
func (e *Error) GetMsg() string {
	return e.msg
}

func (e *Error) Is(target error) bool {
	return false
}
