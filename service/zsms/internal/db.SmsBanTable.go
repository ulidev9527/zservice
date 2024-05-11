package internal

import (
	"time"
	"zservice/zglobal"
	"zservice/zservice"
	"zservice/zservice/ex/redisservice"

	"gorm.io/gorm"
)

// 短信封禁
type SmsBanTable struct {
	gorm.Model
	Phone   string    // 手机号
	Expires time.Time // 过期时间
	Msg     string    // 封禁原因
}

// 账号是否封禁
func IsSmsBan(ctx *zservice.Context, phone string) (bool, error) {
	if phone == "" {
		return false, nil
	}
	if phone[0] != '+' {
		return false, nil
	}

	// 查缓存
	rKey := redisservice.FormatKey(RK_PhoneBan, phone)
	has, e := Redis.Exists(ctx, rKey).Result()
	if e != nil {
		return true, zservice.NewError(e).SetCode(zglobal.Code_ErrorBreakoff)
	}

	if has >= 1 {
		// 提取缓存的时间
		tstr, e := Redis.Get(ctx, rKey).Result()
		if e != nil {
			return true, zservice.NewError(e).SetCode(zglobal.Code_ErrorBreakoff)
		}
		t, e := time.Parse(time.RFC3339, tstr)
		if e != nil {
			_, e := Redis.Del(ctx, rKey).Result()
			if e != nil {
				return true, zservice.NewError(e).SetCode(zglobal.Code_ErrorBreakoff)
			}
		} else {
			return t.After(time.Now()), nil
		}
	}

	// 查数据库
	wh := Mysql.Model(&SmsBanTable{}).Where("phone = ? and expires > now()", phone).Order("expires DESC")
	count := int64(0)
	e = wh.Count(&count).Error
	if e != nil {
		return true, zservice.NewError(e).SetCode(zglobal.Code_ErrorBreakoff)
	}
	// 封禁时间
	banTime := time.Now()
	banCache := zservice.GetenvInt("SMS_BAN_CACHE_DEF")
	isBan := false
	if count > 0 {
		tb := &SmsBanTable{}
		e := wh.First(&tb).Error
		if e != nil {
			return true, zservice.NewError(e).SetCode(zglobal.Code_ErrorBreakoff)
		}
		banTime = tb.Expires
		banCache = int(time.Until(banTime).Seconds())
		isBan = true
	}

	// 缓存
	_, e = Redis.Set(ctx, rKey, banTime.Format(time.RFC3339), time.Duration(banCache)*time.Second).Result()
	if e != nil {
		zservice.LogError(e)
	}
	return isBan, nil
}
