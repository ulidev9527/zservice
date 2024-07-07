package internal

import (
	"fmt"
	"time"
	"zservice/zservice"

	"gorm.io/gorm"
)

// 短信封禁
type SmsBanTable struct {
	gorm.Model
	Phone   string        // 手机号
	Expires zservice.Time // 过期时间
	BanMsg  string        // 封禁原因
}

// 账号是否封禁
func IsSmsBan(ctx *zservice.Context, phone string) (bool, *zservice.Error) {
	if phone == "" || phone[0] != '+' {
		return true, zservice.NewError("phone error")
	}

	// 查缓存
	rk_phoneBan := fmt.Sprintf(RK_Sms_PhoneBan, phone)
	has, e := Redis.Exists(rk_phoneBan).Result()
	if e != nil {
		return true, zservice.NewError(e)
	}
	if has > 0 {
		return true, nil
	}

	// 查数据库
	tb := &SmsBanTable{}
	if e := Gorm.Order("expires DESC").Limit(1).First(&tb, "phone = ? and expires > ?", phone, time.Now()).Error; e != nil {
		if DBService.IsNotFoundErr(e) {
			return false, nil
		} else {
			return false, zservice.NewError(e)
		}
	} else {
		// 更新缓存
		_, e = Redis.SetEX(rk_phoneBan, "1", time.Until(tb.Expires.Time)).Result()
		if e != nil {
			zservice.LogError(e)
		}
		return true, nil
	}
}
