package internal

import (
	"encoding/json"
	"fmt"
	"time"
	"zservice/zservice"
	"zservice/zservice/zglobal"
)

type AuthToken struct {
	UID           uint      // 用户ID
	Token         string    // 令牌
	ExpiresSecond uint32    // 过期时间 单位: 秒
	Expires       time.Time // 过期时间 单位: 秒
	Sign          string    // 签名，用于生成 token 和验证
	TokenKey      string    // token key
	LoginTarget   string    // 登陆的目标平台
}

// 创建一个 token
func CreateToken(ctx *zservice.Context) (*AuthToken, *zservice.Error) {
	// 最小过期时间
	minExpiresSeconds := zservice.GetenvUInt32("TOKEN_MIN_EXPIRES")

	// 创建 token
	tk := &AuthToken{
		ExpiresSecond: uint32(minExpiresSeconds),
		Expires:       time.Now().Add(time.Second * time.Duration(minExpiresSeconds)),
		Sign:          zservice.MD5String(ctx.AuthSign),
		TokenKey:      zservice.RandomMD5(),
	}

	tk.Token = GenTokenSign(tk.Sign, tk.TokenKey)

	if e := tk.Save(); e != nil {
		return nil, e
	}

	return tk, nil
}

// 生成 token
func GenTokenSign(sign, key string) string {
	return zservice.MD5String(fmt.Sprint(sign, key))
}

// 获取 token
func GetToken(tkStr string) (*AuthToken, *zservice.Error) {
	if tkStr == "" {
		return nil, zservice.NewError("token is empty string").SetCode(zglobal.Code_Zauth_TokenIsNil)
	}

	rk := fmt.Sprintf(RK_TokenInfo, tkStr)
	if has, e := Redis.Exists(rk).Result(); e != nil {
		return nil, zservice.NewError(e).SetCode(zglobal.Code_ErrorBreakoff)
	} else if has == 0 {
		return nil, zservice.NewError("no token:", tkStr).SetCode(zglobal.Code_Zauth_TokenIsNil)
	}

	if res, e := Redis.Get(rk).Result(); e != nil {
		return nil, zservice.NewError(e).SetCode(zglobal.Code_ErrorBreakoff)
	} else {
		tk := &AuthToken{}
		if e := json.Unmarshal([]byte(res), &tk); e != nil {
			de := Redis.Del(rk).Err()
			if de != nil {
				return nil, zservice.NewError("convert token fail and del token fail:", res, e, de).SetCode(zglobal.Code_ErrorBreakoff)
			} else {
				return nil, zservice.NewError("convert token fail:", res, e).SetCode(zglobal.Code_ErrorBreakoff)
			}
		}
		return tk, nil
	}
}

// 校验 token
func (l *AuthToken) CheckToken(tk string, sign string) bool {
	if tk == "" || sign == "" {
		return false
	}
	s := GenTokenSign(zservice.MD5String(sign), l.TokenKey)
	return s == tk
}

// 保存
func (l *AuthToken) Save() *zservice.Error {
	rk := fmt.Sprintf(RK_TokenInfo, l.Token)

	l.Expires = time.Now().Add(time.Second * time.Duration(l.ExpiresSecond))

	if l.UID != 0 { // 登录 token 存储
		if e := Redis.SetEX(fmt.Sprintf(RK_AccountLoginToken, l.UID, l.Token), l.Token, time.Until(l.Expires)).Err(); e != nil {
			return zservice.NewError(e).SetCode(zglobal.Code_Zauth_TokenSaveFail)
		}
	}

	if e := Redis.SetEX(rk, zservice.JsonMustMarshalString(l), time.Until(l.Expires)).Err(); e != nil {
		return zservice.NewError(e).SetCode(zglobal.Code_Zauth_TokenSaveFail)
	}
	return nil
}

// 删除 token
func (l *AuthToken) Del() *zservice.Error {

	rk := fmt.Sprintf(RK_TokenInfo, l.Token)

	if e := Redis.Del(fmt.Sprintf(RK_AccountLoginToken, l.UID, l.Token)).Err(); e != nil {
		return zservice.NewError(e).SetCode(zglobal.Code_Zauth_TokenDelFail)
	}
	if e := Redis.Del(rk).Err(); e != nil {
		return zservice.NewError(e).SetCode(zglobal.Code_Zauth_TokenDelFail)
	}
	return nil
}
