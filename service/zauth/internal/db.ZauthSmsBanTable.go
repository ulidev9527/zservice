package internal

import (
	"errors"
	"fmt"
	"time"
	"zservice/zservice"
	"zservice/zservice/ex/gormservice"
	"zservice/zservice/zglobal"

	"gorm.io/gorm"
)

// 短信封禁
type ZauthSmsBanTable struct {
	gormservice.AllModel
	Phone   string // 手机号
	Expires uint64 // 过期时间
	BanMsg  string // 封禁原因
}

// 账号是否封禁
func IsSmsBan(ctx *zservice.Context, phone string) (bool, *zservice.Error) {
	if phone == "" {
		return false, nil
	}
	if phone[0] != '+' {
		return false, nil
	}

	// 查缓存
	rk_phoneBan := fmt.Sprintf(RK_Sms_PhoneBan, phone)
	has, e := Redis.Exists(rk_phoneBan).Result()
	if e != nil {
		return true, zservice.NewError(e).SetCode(zglobal.Code_ErrorBreakoff)
	}
	if has > 0 {
		return true, nil
	}

	// 查数据库
	tb := &ZauthSmsBanTable{}
	if e := Mysql.Model(&ZauthSmsBanTable{}).Where("phone = ? and expires > now()", phone).Order("expires DESC").Limit(1).First(&tb).Error; e != nil {
		if errors.Is(e, gorm.ErrRecordNotFound) {
			return false, nil
		} else {
			return false, zservice.NewError(e).SetCode(zglobal.Code_ErrorBreakoff)
		}
	} else {
		// 更新缓存
		_, e = Redis.SetEX(rk_phoneBan, "1", time.Second*time.Duration(tb.Expires)).Result()
		if e != nil {
			zservice.LogError(e)
		}
		return true, nil
	}
}
