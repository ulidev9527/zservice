package internal

import (
	"zservice/zservice"

	"gorm.io/gorm"
)

type UserBanLogTable struct {
	gorm.Model

	TraceID  string        // 链路ID
	UID      uint32        // 用户ID
	Services string        // 封禁哪些服务，逗号隔开
	Msg      string        // 禁言原因
	Expire   zservice.Time // 解禁时间

}

func NewUserBanLogTable(ctx *zservice.Context, in UserBanLogTable) (*UserBanLogTable, *zservice.Error) {

	tab := &UserBanLogTable{
		TraceID:  ctx.TraceID,
		UID:      in.UID,
		Services: in.Services,
		Msg:      in.Msg,
		Expire:   in.Expire,
	}
	if e := tab.Save(ctx); e != nil {
		return nil, e.AddCaller()
	}
	return tab, nil
}

func (tab *UserBanLogTable) Save(ctx *zservice.Context) *zservice.Error {
	if e := DBService.SaveTableValue(ctx, tab, ""); e != nil {
		return e.AddCaller()
	}
	return nil
}
