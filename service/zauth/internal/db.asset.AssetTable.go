package internal

import (
	"fmt"
	"time"
	"zservice/service/zauth/zauth_pb"
	"zservice/zservice"
	"zservice/zserviceex/dbservice"

	"gorm.io/gorm"
)

// 资源
type AssetTable struct {
	gorm.Model
	AssetID string        `gorm:"unique"` // 资源ID
	TraceID string        // 链路ID
	Name    string        // 名称
	Size    uint64        // 文件大小
	MD5     string        // md5
	Expires zservice.Time // 过期时间
}

// 创建资源
func CreateAssetInfo(ctx *zservice.Context, in *zauth_pb.AssetInfo) (*AssetTable, *zservice.Error) {

	if in.Md5 == "" && len(in.Md5) != 32 {
		return nil, zservice.NewError("invalid md5", in.Md5).SetCode(zservice.Code_ParamsErr)
	}

	// 准备写入数据
	tab := &AssetTable{
		AssetID: zservice.MD5String(fmt.Sprintf("%s_%s", in.Md5, zservice.RandomMD5())),
		TraceID: ctx.TraceID,
		Name:    in.Name,
		Size:    in.Size,
		MD5:     in.Md5,
		Expires: zservice.TimeUnixMilli(in.Expires),
	}

	if e := tab.Save(ctx); e != nil {
		return nil, e
	}
	return tab, nil
}

// 获取资源
func GetAssetByID(ctx *zservice.Context, token string) (*AssetTable, *zservice.Error) {
	tab := &AssetTable{}
	if e := DBService.GetTableFirst(ctx, dbservice.GetTableValueOption{
		Tab:      tab,
		RK:       fmt.Sprintf(RK_AssetInfo, token),
		SQLConds: []any{"asset_id = ? AND (expires IS NULL OR expires > ?)", token, time.Now()},
	}); e != nil {
		return nil, e.AddCaller()
	}

	return tab, nil
}

func (tab *AssetTable) ToAssetInfo() *zauth_pb.AssetInfo {
	return &zauth_pb.AssetInfo{
		AssetID: tab.AssetID,
		Name:    tab.Name,
		Md5:     tab.MD5,
		Expires: tab.Expires.UnixMilli(),
		Size:    tab.Size,
	}
}

func (tab *AssetTable) IsExpired() bool {
	if tab.Expires.IsZero() {
		return false
	}
	return tab.Expires.BeforeNow()
}

func (tab *AssetTable) Save(ctx *zservice.Context) *zservice.Error {

	if e := DBService.SaveTableValue(ctx, tab, fmt.Sprintf(RK_AssetInfo, tab.AssetID)); e != nil {
		return e.AddCaller()
	}
	return nil
}
