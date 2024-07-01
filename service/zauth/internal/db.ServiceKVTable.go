package internal

import (
	"encoding/json"
	"fmt"
	"zservice/zservice"
	"zservice/zservice/ex/gormservice"
	"zservice/zservice/ex/redisservice"
	"zservice/zservice/zglobal"

	"gorm.io/gorm"
)

type ServiceKVTable struct {
	gorm.Model
	Key     string
	Value   string
	Service string // 所属服务
}

func GetOrCreateServiceKVTable(ctx *zservice.Context, service string, key string) (*ServiceKVTable, *zservice.Error) {
	if tab, e := GetServiceKVTable(ctx, service, key); e != nil {
		if e.GetCode() != zglobal.Code_NotFound {
			return nil, e
		}
	} else {
		return tab, nil
	}

	// 创建
	tab := &ServiceKVTable{
		Key:     key,
		Service: service,
	}
	if e := tab.Save(ctx); e != nil {
		return nil, e

	}
	return tab, nil
}

func GetServiceKVTable(ctx *zservice.Context, service string, key string) (*ServiceKVTable, *zservice.Error) {
	rk_info := fmt.Sprintf(RK_ServiceKVInfo, service, key)
	tab := &ServiceKVTable{}

	// 查缓存
	if s, e := Redis.Get(rk_info).Result(); e != nil {
		if !redisservice.IsNilErr(e) {
			return nil, zservice.NewError(e)
		}

	} else {
		if e := json.Unmarshal([]byte(s), tab); e != nil {
			ctx.LogError(e)
			zservice.Go(func() {
				if e := Redis.Del(rk_info).Err(); e != nil {
					ctx.LogError(e)
				}
			})
		} else {
			return tab, nil
		}
	}

	// 查库
	if e := Mysql.Where("service = ? AND `key` = ?", service, key).First(&tab).Error; e != nil {
		if gormservice.IsNotFound(e) {
			return nil, zservice.NewError(e).SetCode(zglobal.Code_NotFound)
		}
		return nil, zservice.NewError(e)
	}

	// 更新缓存
	zservice.Go(func() {
		if e := Redis.Set(rk_info, zservice.JsonMustMarshalString(tab)).Err(); e != nil {
			ctx.LogError(e)
		}
	})

	return tab, nil
}

func (tab *ServiceKVTable) Save(ctx *zservice.Context) *zservice.Error {
	rk_info := fmt.Sprintf(RK_ServiceKVInfo, tab.Service, tab.Key)

	un, e := Redis.Lock(rk_info)
	if e != nil {
		return e
	}
	defer un()

	if e := Mysql.Save(tab).Error; e != nil {
		return zservice.NewError(e)
	}

	zservice.Go(func() {
		if e := Redis.Del(rk_info).Err(); e != nil {
			ctx.LogError(e)
		}
	})
	return nil
}
