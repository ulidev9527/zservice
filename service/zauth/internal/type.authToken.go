package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	"zservice/zglobal"
	"zservice/zservice"
)

type AuthToken struct {
	UID              uint64    // 用户ID
	Token            string    // 令牌
	ExpiresSecond    uint32    // 过期时间 单位: 秒
	Expires          time.Time // 过期时间 单位: 秒
	RefreshTokenTime time.Time // 下次刷新token时间, 用于自动刷新
	NewToken         string    // 新令牌
	Sign             string    // 签名，用于生成 token 和验证
}

// 创建一个 token
func CreateToken(ctx *zservice.Context) (*AuthToken, *zservice.Error) {
	// 最小过期时间
	minExpiresSeconds := zservice.GetenvUInt32("TOKEN_MIN_EXPIRES")

	// 创建 token
	tk := &AuthToken{
		ExpiresSecond: uint32(minExpiresSeconds),
		Expires:       time.Now().Add(time.Second * time.Duration(minExpiresSeconds)),
		Sign:          zservice.MD5String(fmt.Sprint(0, ctx.AuthSign)),
	}

	if e := tk.GenToken(); e != nil {
		return nil, e
	}

	return tk, nil
}

// 获取 token
func GetToken(tkStr string) (*AuthToken, *zservice.Error) {
	if tkStr == "" {
		return nil, zservice.NewError("no token:", tkStr).SetCode(zglobal.Code_Zauth_TokenIsNil)
	}

	rk := fmt.Sprintf(RK_Token, tkStr)
	if has, e := Redis.Exists(context.TODO(), rk).Result(); e != nil {
		return nil, zservice.NewError(e).SetCode(zglobal.Code_ErrorBreakoff)
	} else if has == 0 {
		return nil, zservice.NewError("no token:", tkStr).SetCode(zglobal.Code_Zauth_TokenIsNil)
	}

	if res, e := Redis.Get(context.TODO(), rk).Result(); e != nil {
		return nil, zservice.NewError(e).SetCode(zglobal.Code_ErrorBreakoff)
	} else {
		tk := &AuthToken{}
		if e := json.Unmarshal([]byte(res), &tk); e != nil {
			de := Redis.Del(context.TODO(), rk).Err()
			if de != nil {
				return nil, zservice.NewError("convert token fail and del token fail:", res, e, de).SetCode(zglobal.Code_ErrorBreakoff)
			} else {
				return nil, zservice.NewError("convert token fail:", res, e).SetCode(zglobal.Code_ErrorBreakoff)
			}
		}
		return tk, nil
	}
}

// 刷新 token
func (l *AuthToken) GenToken() *zservice.Error {

	l.Token = zservice.MD5String(fmt.Sprint(l.Sign, zservice.RandomXID()))

	return l.Save()
}

// 保存
func (l *AuthToken) Save() *zservice.Error {
	rk := fmt.Sprintf(RK_Token, l.Token)
	_, e := Redis.Set(context.TODO(), rk, l.Token, time.Second*time.Duration(l.ExpiresSecond)).Result()
	if e != nil {
		return zservice.NewError(e).SetCode(zglobal.Code_Zauth_TokenSaveFail)
	}
	return nil
}
