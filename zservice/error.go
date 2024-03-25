package zservice

import (
	"fmt"
	"runtime"
)

type Error struct {
	code int    // 错误码
	msg  string // 错误消息
	str  string // 自己的错误消息
}
type RejectError struct {
	Error
}

func NewError(v ...any) *Error {
	return NewErrorCaller(2, Sprint(v...), nil)
}
func NewErrorCaller(skip int, str string, e error) *Error {
	_, file, line, _ := runtime.Caller(skip)
	s := fmt.Sprint(file, ":", line, " > ", str)
	if e != nil {
		s = fmt.Sprint(s, "\n", e.Error())
	}
	return &Error{
		str: str,
		msg: s,
	}
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
	return fmt.Sprint("Error:\n", e.msg)
}
func (e *Error) String() string {
	return e.str
}

// 设置错误码
func (e *Error) Code(code int) *Error {
	e.code = code
	return e
}
