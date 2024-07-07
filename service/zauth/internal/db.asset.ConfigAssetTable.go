package internal

import (
	"fmt"
	"zservice/service/zauth/zauth_pb"
	"zservice/zservice"
	"zservice/zserviceex/dbservice"

	"gorm.io/gorm"
)

// 配置资源
type ConfigAssetTable struct {
	gorm.Model

	TraceID       string // 链路ID
	Name          string // 配置文件名称
	Service       string // 服务
	Parser        uint32 // 解析器
	AssetID       string // 资源ID，来自:AssetTable
	ParserAssetID string // 解析器资源ID，来自:AssetTable
}

// 创建配置资源信息
func CreateConfigAssetInfo(ctx *zservice.Context, in *zauth_pb.ConfigAssetInfo) (*ConfigAssetTable, *zservice.Error) {
	tab := &ConfigAssetTable{
		TraceID:       ctx.TraceID,
		Name:          in.Name,
		Service:       in.Service,
		Parser:        in.Parser,
		AssetID:       in.AssetID,
		ParserAssetID: in.ParserAssetID,
	}
	if e := tab.Save(ctx); e != nil {
		return nil, e.AddCaller()
	}
	return tab, nil
}

// 获取最新一个对应服务的配置资源
func GetConfigAssetInfo(ctx *zservice.Context, service string, name string) (*ConfigAssetTable, *zservice.Error) {
	tab := &ConfigAssetTable{}
	if e := DBService.GetTableFirst(ctx, dbservice.GetTableValueOption{
		Tab:      tab,
		RK:       fmt.Sprintf(RK_ConfigAssetInfo, service, name),
		SQLConds: []any{"service = ? AND name = ?", service, name},
		Order:    "created_at desc",
	}); e != nil {
		return nil, e.AddCaller()
	}
	return tab, nil
}

// 转换到 ConfigAssetInfo
func (tab *ConfigAssetTable) ToConfigAssetInfo() *zauth_pb.ConfigAssetInfo {
	return &zauth_pb.ConfigAssetInfo{
		Name:          tab.Name,
		Service:       tab.Service,
		Parser:        tab.Parser,
		AssetID:       tab.AssetID,
		ParserAssetID: tab.ParserAssetID,
	}
}

// 存储
func (tab *ConfigAssetTable) Save(ctx *zservice.Context) *zservice.Error {
	if e := DBService.SaveTableValue(ctx, tab, fmt.Sprintf(RK_ConfigAssetInfo, tab.Service, tab.Name)); e != nil {
		return e.AddCaller()
	}
	return nil
}
