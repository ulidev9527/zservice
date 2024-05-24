package gormservice

import (
	"zservice/zservice"
)

type GormLoggerWrite struct {
	Service *zservice.ZService
}

func (cw *GormLoggerWrite) Write(p []byte) (n int, err error) {
	cw.Service.LogError(string(p))
	return len(p), nil
}
