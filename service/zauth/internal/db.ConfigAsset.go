package internal

import "zservice/zservice/ex/gormservice"

// 配置资源
type ConfigAsset struct {
	gormservice.Model
	Service  string
	FileName string
	FileMd5  string
}
