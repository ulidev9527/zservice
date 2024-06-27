package zservice

import (
	"fmt"
	"runtime"
	"zservice/zservice/zglobal"
)

// ------------- Error -------------
type Error struct {
	code uint32 // 错误码
	msg  string // 错误消息
}

func NewError(v ...any) *Error {
	return NewErrorCaller(2, Sprint(v...), nil)
}
func NewErrorCaller(skip int, str string, e error) *Error {
	code := uint32(zglobal.Code_Fail)
	s := str
	if e != nil {
		s = fmt.Sprint(s, " | ", e.Error())
	}
	err := &Error{
		code: code,
		msg:  s,
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
		return fmt.Sprintf("%d:%T:\n%s", e.code, &Error{}, e.msg)
	} else {
		return fmt.Sprintf("%T:\n%s", &Error{}, e.msg)
	}
}
func (e *Error) String() string {
	return e.msg
}

// 添加路径记录
// skip 一般不需要填写，涉及到 caller 跳层问题
func (e *Error) AddCaller(skips ...int) *Error {
	skip := 1
	if len(skips) > 0 {
		skip = skips[0] + 1
	}
	_, file, line, _ := runtime.Caller(skip)
	e.msg = fmt.Sprint(file, ":", line, " > ", e.msg)
	return e
}

// 设置错误码
func (e *Error) SetCode(code uint32) *Error {
	e.code = code
	return e
}

// 获取错误码
func (e *Error) GetCode() uint32 {
	return e.code
}
