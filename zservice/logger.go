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

	isUseTime := Convert_StringToBoolean(os.Getenv("LOGGER_USE_TIME"))

	if !isUseTime {
		logCW.FormatTimestamp = func(i interface{}) string {
			return ""
		}
	}

	logger = zerolog.New(logCW)
	logCtx := logger.With()

	if isUseTime {

		logger = logCtx.Timestamp().Logger()
	} else {
		logger = logCtx.Logger()
	}

}

// info 消息
func LogInfo(v ...any)                 { LogInfoCaller(2, v...) }
func LogInfof(format string, v ...any) { LogInfoCaller(2, fmt.Sprintf(format, v...)) }

func LogInfoCaller(caller int, v ...any) { logger.Info().Caller(caller).Msg(Sprint(v...)) }
func LogInfoCallerf(caller int, format string, v ...any) {
	LogInfoCaller(caller+1, fmt.Sprintf(format, v...))
}

// Debugf
func LogDebug(v ...any)                   { LogDebugCaller(2, v...) }
func LogDebugf(format string, v ...any)   { LogDebugCallerf(2, format, v...) }
func LogDebugCaller(caller int, v ...any) { logger.Debug().Caller(caller).Msg(SprintQuote(v...)) }
func LogDebugCallerf(caller int, format string, v ...any) {
	LogDebugCaller(caller+1, fmt.Sprintf(format, v...))
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

// fatal
func LogFatal(v ...any)                   { LogFatalCaller(2, v...) }
func LogFatalf(format string, v ...any)   { LogFatalCallerf(2, format, v...) }
func LogFatalCaller(caller int, v ...any) { logger.Fatal().Caller(caller).Msg(SprintQuote(v...)) }
func LogFatalCallerf(caller int, format string, v ...any) {
	LogFatalCaller(caller+1, fmt.Sprintf(format, v...))
}
