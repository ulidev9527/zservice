package internal

import (
	"zservice/service/zauth/zauth"
	"zservice/zservice"
)

// 道具配置
type Config_Item struct {
	ID     uint32 `json:"id"`     // 道具 ID
	Name   string `json:"name"`   // 道具名称
	Desc   string `json:"desc"`   // 道具描述
	Type   uint32 `json:"type"`   // 道具类型
	Expire int64  `json:"expire"` // 过期时间（秒）
}

// 获取道具配置
func GetConfig_Item(ctx *zservice.Context, id uint32) (*Config_Item, *zservice.Error) {
	item := Config_Item{}

	if e := zauth.DownloadConfigAsset(ctx, "Config_Item.xlsx", &item, zservice.Uint32ToString(id)); e != nil {
		return nil, e
	}
	return &item, nil
}
