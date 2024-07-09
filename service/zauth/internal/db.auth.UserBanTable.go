package internal

import (
	"fmt"
	"zservice/zservice"
	"zservice/zserviceex/dbservice"

	"gorm.io/gorm"
)

// 用户禁用表
type UserBanTable struct {
	gorm.Model

	TraceID string        // 链路ID, 最好一次更新的链路ID
	UID     uint32        // 用户ID
	Service string        // 拒绝服务
	Expire  zservice.Time // 封禁时间

}

// 获取或者创建用户禁用信息
func GetOrCreateUserBanTable(ctx *zservice.Context, uid uint32, service string) (*UserBanTable, *zservice.Error) {

	tab := &UserBanTable{}
	if e := DBService.GetTableValue(ctx, dbservice.GetTableValueOption{
		Tab:      tab,
		RK:       fmt.Sprintf(RK_UserBan, uid, service),
		SQLConds: []any{"uid = ? AND service = ?", uid, service},
	}); e != nil {
		if e.GetCode() != zservice.Code_NotFound {
			return nil, e.AddCaller()
		}
	}

	// 创建
	tab.UID = uid
	tab.Service = service
	tab.TraceID = ctx.TraceID

	if e := tab.Save(ctx); e != nil {
		return nil, e.AddCaller()
	}

	return tab, nil
}

func (tab *UserBanTable) Save(ctx *zservice.Context) *zservice.Error {

	if e := DBService.SaveTableValue(ctx, tab, fmt.Sprintf(RK_UserBan, tab.UID, tab.Service)); e != nil {
		return e.AddCaller()
	}
	return nil

}
