package internal

import (
	"encoding/json"
	"fmt"
	"time"
	"zservice/service/zauth/zauth_pb"
	"zservice/zservice"
)

type AuthToken struct {
	UID           uint32    // 用户ID
	Token         string    // 令牌
	CreateAT      time.Time // 创建时间 毫秒
	ExpiresSecond uint32    // 过期时间 s
	Expires       time.Time // 过期时间
	Sign          string    // 签名，用于生成 token 和验证
	LoginServices []string  // 登陆的服务
}

// 创建一个 token
func CreateToken(ctx *zservice.Context, tokenSign string) (*AuthToken, *zservice.Error) {
	// 创建 token
	tk := &AuthToken{}

	tk.ExpiresSecond = uint32(zservice.Time_10m.Seconds())
	tk.CreateAT = time.Now()
	tk.Sign = tokenSign
	tk.Token = zservice.MD5String(fmt.Sprint(tk.Sign, zservice.RandomMD5(), zservice.RandomXID(), tk.CreateAT))

	if e := tk.Save(ctx); e != nil {
		return nil, e.AddCaller()
	}

	return tk, nil
}

// 获取 token
func GetTokenInfo(ctx *zservice.Context, tkStr string) (*AuthToken, *zservice.Error) {
	if tkStr == "" {
		return nil, zservice.NewError("token is empty string").SetCode(zservice.Code_NotFound)
	}

	rk := fmt.Sprintf(RK_TokenInfo, tkStr)

	if res, e := Redis.Get(rk).Result(); e != nil {
		if DBService.IsNotFoundErr(e) {
			return nil, zservice.NewError(e).SetCode(zservice.Code_NotFound)
		}
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

// token 登录
func TokenLogin(ctx *zservice.Context, in struct {
	Service string
	Expires uint32
}, at *AuthToken, user *UserTable) *zservice.Error {

	// 多点登录限制验证
	if res, e := GetOrCreateServiceKVTable(ctx, in.Service, KV_Service_Login_AllowMPOP); e != nil {
		return zservice.NewError(e)
	} else if !zservice.StringToBoolean(res.Value) {
		// 不允许多点登录，退出当前服务的其它授权 token
		rk := fmt.Sprintf(RK_UserLoginServices, user.UID, in.Service)
		if list, e := Redis.LRange(rk, 0, -1).Result(); e != nil { // 获取 token 所有登录服务
			ctx.LogError("redis get fail", rk, e)
		} else {
			remKeys := []string{}
			for _, tk := range list {
				if _at, e := GetTokenInfo(ctx, tk); e != nil { // 获取 token 信息
					if e.GetCode() == zservice.Code_NotFound {
						remKeys = append(remKeys, tk)
					} else {
						ctx.LogError("get tk fail", tk, e)
					}
				} else {
					// 退出相应 token 登录状态
					res := Logic_Logout(ctx, &zauth_pb.Logout_REQ{Token: _at.Token, TokenSign: _at.Sign})
					if res.Code != zservice.Code_SUCC {
						ctx.LogError("login out fail", tk)
					}
				}
			}
			if len(remKeys) > 0 { // 清理
				for _, tk := range remKeys {
					zservice.Go(func() {
						ctx.LogInfo("rem invalid token:", tk)
						if e := Redis.LRem(rk, 0, tk).Err(); e != nil {
							ctx.LogError(e)
						}
					})
				}
			}
		}
	}

	// 设置关联信息
	at.ExpiresSecond = in.Expires
	if at.ExpiresSecond == 0 {
		at.ExpiresSecond = uint32(zservice.Time_10Day.Seconds())
	}

	at.UID = user.UID
	at.AddLoginService(in.Service)

	if e := at.Save(ctx); e != nil {
		return e.AddCaller()
	}

	// 登录 token 更新
	if e := Redis.LPush(fmt.Sprintf(RK_UserLoginServices, at.UID, in.Service), at.Token).Err(); e != nil {
		zservice.Go(func() {
			at.Del(ctx)
		})
		return zservice.NewError(e)
	}
	return nil
}

// token 校验
func (l *AuthToken) TokenCheck(sign string) bool {
	return sign == l.Sign
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

	if l.ExpiresSecond < uint32(zservice.Time_10m.Seconds()) {
		l.ExpiresSecond = uint32(zservice.Time_10m.Seconds())
	}
	l.Expires = time.Now().Add(time.Second * time.Duration(l.ExpiresSecond))

	if e := Redis.SetEX(fmt.Sprintf(RK_TokenInfo, l.Token), zservice.JsonMustMarshalString(l), time.Until(l.Expires)).Err(); e != nil {
		return zservice.NewError(e)
	}
	return nil
}

// 删除 token
func (l *AuthToken) Del(ctx *zservice.Context) {

	rk_info := fmt.Sprintf(RK_TokenInfo, l.Token)

	// 删除 token 信息
	if e := Redis.Del(rk_info).Err(); e != nil {
		ctx.LogError(e)
	}
}
