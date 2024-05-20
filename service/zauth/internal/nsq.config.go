package internal

import (
	"zservice/zservice"
	"zservice/zservice/zglobal"
)

// 文件配置改变
func NsqFileConfigChange(fileName string) *zservice.Error {
	if e := Nsq.Publish(NSQ_FileConfig_Change, []byte(fileName)); e != nil {
		return zservice.NewError(e).SetCode(zglobal.Code_ErrorBreakoff)
	}
	return nil
}
