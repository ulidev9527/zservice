package zservice

import (
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog"
)

var logger zerolog.Logger

func init() {
	logCW := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
	}
	logger = zerolog.New(logCW)
	logCtx := logger.With()
	logger = logCtx.Timestamp().Logger()
}

// info 消息
func LogInfo(v ...any)                 { LogInfoCaller(2, v...) }
func LogInfof(format string, v ...any) { LogInfoCaller(2, fmt.Sprintf(format, v...)) }

// func LogInfoCaller(caller int, v ...any) { logger.Info().Caller(caller).Msg(Sprint(v...)) }
func LogInfoCaller(caller int, v ...any) {
	s := Sprint(v...)
	// 限制打印长度
	if len(s) > 1024 {
		s = s[:1000]
	}
	logger.Info().Msg(s)
} // 不打印 caller
func LogInfoCallerf(caller int, format string, v ...any) {
	LogInfoCaller(caller+1, fmt.Sprintf(format, v...))
}

// 警告
func LogWarn(v ...any)                   { LogWarnCaller(2, v...) }
func LogWarnf(format string, v ...any)   { LogWarnCallerf(2, format, v...) }
func LogWarnCaller(caller int, v ...any) { logger.Warn().Caller(caller).Msg(SprintQuote(v...)) }
func LogWarnCallerf(caller int, format string, v ...any) {
	LogWarnCaller(caller+1, fmt.Sprintf(format, v...))
}

// 错误
func LogError(v ...any)                   { LogErrorCaller(2, v...) }
func LogErrorf(format string, v ...any)   { LogErrorCallerf(2, format, v...) }
func LogErrorCaller(caller int, v ...any) { logger.Error().Caller(caller).Msg(SprintQuote(v...)) }
func LogErrorCallerf(caller int, format string, v ...any) {
	LogErrorCaller(caller+1, fmt.Sprintf(format, v...))
}
func Erroref(e error, formmat string, v ...any) {
	LogErrorCallerf(2, "%v \n %v", fmt.Sprintf(formmat, v...), e.Error())
}
func LogErrore(e error) { LogErrorCaller(2, e.Error()) }

// panic
func LogPanic(v ...any)                   { LogPanicCaller(2, v...) }
func LogPanicf(format string, v ...any)   { LogPanicCallerf(2, format, v...) }
func LogPanicCaller(caller int, v ...any) { logger.Panic().Caller(caller).Msg(SprintQuote(v...)) }
func LogPanicCallerf(caller int, format string, v ...any) {
	LogPanicCaller(caller+1, fmt.Sprintf(format, v...))
}
