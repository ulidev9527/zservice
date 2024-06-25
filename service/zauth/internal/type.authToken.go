package internal

import (
	"encoding/json"
	"fmt"
	"time"
	"zservice/zservice"
	"zservice/zservice/zglobal"
)

type AuthToken struct {
	UID           uint32    // 用户ID
	Token         string    // 令牌
	ExpiresSecond uint32    // 过期时间 单位: 秒
	Expires       time.Time // 过期时间 单位: 秒
	Sign          string    // 签名，用于生成 token 和验证
	TokenKey      string    // token key
	LoginServices []string  // 登陆的服务
}

// 创建一个 token
func CreateToken(ctx *zservice.Context, tokenSign string) (*AuthToken, *zservice.Error) {
	// 创建 token
	tk := &AuthToken{
		Sign:     zservice.MD5String(tokenSign),
		TokenKey: zservice.RandomMD5(),
	}

	tk.Token = genTokenSign(tk.Sign, tk.TokenKey)

	if e := tk.Save(ctx); e != nil {
		return nil, e
	}

	return tk, nil
}

// 生成 token
func genTokenSign(sign, key string) string {
	return zservice.MD5String(fmt.Sprint(sign, key))
}

// 获取 token
func GetToken(ctx *zservice.Context, tkStr string) (*AuthToken, *zservice.Error) {
	if tkStr == "" {
		return nil, zservice.NewError("token is empty string").SetCode(zglobal.Code_NotFound)
	}

	rk := fmt.Sprintf(RK_TokenInfo, tkStr)
	if has, e := Redis.Exists(rk).Result(); e != nil {
		return nil, zservice.NewError(e)
	} else if has == 0 {
		return nil, zservice.NewError("no token:", tkStr).SetCode(zglobal.Code_NotFound)
	}

	if res, e := Redis.Get(rk).Result(); e != nil {
		return nil, zservice.NewError(e)
	} else {
		tk := &AuthToken{}
		if e := json.Unmarshal([]byte(res), &tk); e != nil {
			de := Redis.Del(rk).Err()
			if de != nil {
				return nil, zservice.NewError("convert token fail and del token fail:", res, e, de)
			} else {
				return nil, zservice.NewError("convert token fail:", res, e)
			}
		}
		if e := tk.Save(ctx); e != nil {
			ctx.LogError(e)
		}
		return tk, nil
	}
}

// token 校验
func (l *AuthToken) TokenCheck(tk string, sign string) bool {
	if tk == "" || sign == "" {
		return false
	}
	s := genTokenSign(zservice.MD5String(sign), l.TokenKey)
	return s == tk
}

// 是否有登陆服务
func (l *AuthToken) HasLoginService(service string) bool {
	for _, v := range l.LoginServices {
		if v == service {
			return true
		}
	}
	return false
}

// 添加登陆服务
func (l *AuthToken) AddLoginService(service string) {
	if l.HasLoginService(service) {
		return
	}
	l.LoginServices = append(l.LoginServices, service)
}

// 保存
func (l *AuthToken) Save(ctx *zservice.Context) *zservice.Error {

	if l.ExpiresSecond < 600 {
		l.ExpiresSecond = 600
	}
	l.Expires = time.Now().Add(time.Second * time.Duration(l.ExpiresSecond))

	if l.UID != 0 { // 登录 token 存储
		for _, service := range l.LoginServices {
			// 更新 service
			if e := Redis.SetEX(fmt.Sprintf(RK_UserLoginService, l.UID, service), l.Token, time.Until(l.Expires)).Err(); e != nil {
				return zservice.NewError(e).SetCode(zglobal.Code_Fail)
			}
		}
	}

	if e := Redis.SetEX(fmt.Sprintf(RK_TokenInfo, l.Token), zservice.JsonMustMarshalString(l), time.Until(l.Expires)).Err(); e != nil {
		return zservice.NewError(e).SetCode(zglobal.Code_Fail)
	}
	return nil
}

// 删除 token
func (l *AuthToken) Del(ctx *zservice.Context) {

	rk_info := fmt.Sprintf(RK_TokenInfo, l.Token)

	// 用户登陆信息删除
	zservice.Go(func() {
		if len(l.LoginServices) > 0 {
			delKeys := []string{}
			for _, service := range l.LoginServices {
				rk := fmt.Sprintf(RK_UserLoginService, l.UID, service)
				if s, e := Redis.Get(rk).Result(); e != nil {
					ctx.LogError(e)
				} else if s == l.Token {
					delKeys = append(delKeys, rk)
				}
			}
			if e := Redis.Del(delKeys...).Err(); e != nil {
				ctx.LogError(e)
			}
		}
	})

	// 删除 token 信息
	if e := Redis.Del(rk_info).Err(); e != nil {
		ctx.LogError(e)
	}
}
