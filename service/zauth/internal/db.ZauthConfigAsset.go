package internal

import "zservice/zservice/ex/gormservice"

// 配置资源
type ZauthConfigAsset struct {
	gormservice.AllModel
	Service  string
	FileName string
	FileMd5  string
}
