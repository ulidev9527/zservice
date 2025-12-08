package zservice

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

var zLogger = zerolog.New(zerolog.ConsoleWriter{
	Out:        os.Stdout,
	TimeFormat: time.RFC3339,
}).With().Timestamp().Logger()

// info 消息
func LogInfo(v ...any)                   { zLogger.Info().Caller(1).Msg(StringSprint(v...)) }
func LogInfof(f string, v ...any)        { zLogger.Info().Caller(1).Msgf(f, v...) }
func LogInfoCaller(caller int, v ...any) { zLogger.Info().Caller(caller).Msg(StringSprint(v...)) }
func LogInfoCallerf(caller int, f string, v ...any) {
	zLogger.Info().Caller(caller).Msgf(f, v...)
}

// warn 消息
func LogWarn(v ...any)                   { zLogger.Warn().Caller(1).Msg(StringSprint(v...)) }
func LogWarnf(f string, v ...any)        { zLogger.Warn().Caller(1).Msgf(f, v...) }
func LogWarnCaller(caller int, v ...any) { zLogger.Warn().Caller(caller).Msg(StringSprint(v...)) }
func LogWarnCallerf(caller int, f string, v ...any) {
	zLogger.Warn().Caller(caller).Msgf(f, v...)
}

// err 消息
func LogError(v ...any)                   { zLogger.Error().Caller(1).Msg(StringSprint(v...)) }
func LogErrorf(f string, v ...any)        { zLogger.Error().Caller(1).Msgf(f, v...) }
func LogErrorCaller(caller int, v ...any) { zLogger.Error().Caller(caller).Msg(StringSprint(v...)) }
func LogErrorCallerf(caller int, f string, v ...any) {
	zLogger.Error().Caller(caller).Msgf(f, v...)
}

// panic 消息
func LogPanic(v ...any)                   { zLogger.Panic().Caller(1).Msg(StringSprint(v...)) }
func LogPanicf(f string, v ...any)        { zLogger.Panic().Caller(1).Msgf(f, v...) }
func LogPanicCaller(caller int, v ...any) { zLogger.Panic().Caller(caller).Msg(StringSprint(v...)) }
func LogPanicCallerf(caller int, f string, v ...any) {
	zLogger.Panic().Caller(caller).Msgf(f, v...)
}

// debug 消息
func LogDebug(v ...any)                   { zLogger.Debug().Caller(1).Msg(StringSprint(v...)) }
func LogDebugf(f string, v ...any)        { zLogger.Debug().Caller(1).Msgf(f, v...) }
func LogDebugCaller(caller int, v ...any) { zLogger.Debug().Caller(caller).Msg(StringSprint(v...)) }
func LogDebugCallerf(caller int, f string, v ...any) {
	zLogger.Debug().Caller(caller).Msgf(f, v...)
}
