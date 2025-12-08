package nsqservice

import (
	"strings"

	"zserviceapps/packages/zservice"
)

// 日志扩展
type LogEx struct{ Zservice *zservice.ZService }

func (l *LogEx) Output(calldepth int, s string) error {

	if l.Zservice == nil {
		l.Zservice = zservice.GetMainService()
	}

	if strings.HasPrefix(s, "ERR") {
		l.Zservice.LogError(s)
	} else if strings.HasPrefix(s, "WRN") {
		l.Zservice.LogWarn(s)
	} else {
		l.Zservice.LogInfo(s)
	}
	return nil
}
