package internal

import (
	"fmt"
	"zservice/zservice"
	"zservice/zservice/ex/gormservice"
	"zservice/zservice/zglobal"
)

// 账号表
type UserTable struct {
	gormservice.Model
	UID            uint32 `gorm:"unique"` // 用户唯一ID
	LoginName      string `gorm:"unique"` // 登陆账号
	LoginPass      string // 登陆密码
	LoginPassToken string // 密码令牌
	Phone          string `gorm:"unique"`    // 手机号 含区号 +86******
	State          uint32 `gorm:"default:1"` // 账号状态 0 禁用 1 启用
}

// 创建一个新的账号
func CreateUser(ctx *zservice.Context) (*UserTable, *zservice.Error) {
	accID, e := GetNewUID(ctx)
	if e != nil {
		return nil, e
	}
	z := &UserTable{UID: accID}
	if e := z.Save(ctx); e != nil {
		return nil, e
	}
	return z, nil
}

// 获取一个新的账号ID
func GetNewUID(ctx *zservice.Context) (uint32, *zservice.Error) {
	return dbhelper.GetNewTableID(ctx, func() uint32 {
		return zservice.RandomUInt32Range(1000000, 999999999) // 7-9位数
	}, HasUserByID)
}

// 是否存在这个账号
func HasUserByID(ctx *zservice.Context, id uint32) (bool, *zservice.Error) {
	return dbhelper.HasTableValue(ctx, &UserTable{}, fmt.Sprintf(RK_UserInfo, id), fmt.Sprintf("uid = %v", id))
}

// 是否存在这个账号
func HasUserByLoginName(ctx *zservice.Context, loginName string) (bool, *zservice.Error) {
	return dbhelper.HasTableValue(ctx, &UserTable{}, fmt.Sprintf(RK_UserLoginName, loginName), fmt.Sprintf("login_name = '%v'", loginName))
}

// 账号密码签名
func UserGenPassSign(z *UserTable, password string) string {
	return zservice.MD5String(fmt.Sprint(z.UID, z.LoginPassToken, password))
}

// 获取账号
func GetUserByID(ctx *zservice.Context, id uint) (*UserTable, *zservice.Error) {
	tab := UserTable{}

	if e := dbhelper.GetTableValue(ctx, &tab, fmt.Sprintf(RK_UserInfo, id), fmt.Sprintf("uid = %v", id)); e != nil {
		return nil, e
	}
	return &tab, nil
}

// 根据登陆名获取账号
func GetUserByLoginName(ctx *zservice.Context, loginName string) (*UserTable, *zservice.Error) {

	rk := fmt.Sprintf(RK_UserLoginName, loginName)
	if has, e := Redis.Exists(rk).Result(); e != nil { // 是否有缓存
		return nil, zservice.NewError(e).SetCode(zglobal.Code_ErrorBreakoff)
	} else if has > 0 {
		if s, e := Redis.Get(rk).Result(); e != nil { // 是否有数据
			return nil, zservice.NewError(e).SetCode(zglobal.Code_ErrorBreakoff)
		} else {
			if tab, e := GetUserByID(ctx, zservice.StringToUint(s)); e != nil {
				return nil, e
			} else {
				return tab, nil
			}
		}
	}

	// 未找到 查表
	tab := UserTable{}

	// 验证数据库中是否存在
	if e := Mysql.Model(&tab).Where(fmt.Sprintf("login_name = '%v'", loginName)).First(&tab).Error; e != nil {
		if gormservice.IsNotFound(e) {
			return nil, zservice.NewError(e).SetCode(zglobal.Code_NotFound)
		}
		return nil, zservice.NewError(e).SetCode(zglobal.Code_ErrorBreakoff)
	}

	// 更新缓存
	if e := Redis.Set(rk, zservice.Uint32ToString(tab.UID)).Err(); e != nil {
		ctx.LogError(e)
	}
	if e := Redis.Set(fmt.Sprintf(RK_UserInfo, tab.UID), zservice.JsonMustMarshalString(tab)).Err(); e != nil {
		ctx.LogError(e)
	}

	return &tab, nil
}

// 根据手机号获取账号
func GetUserByPhone(ctx *zservice.Context, phone string) (*UserTable, *zservice.Error) {
	rk := fmt.Sprintf(RK_UserLoginPhone, phone)
	if has, e := Redis.Exists(rk).Result(); e != nil { // 是否有缓存
		return nil, zservice.NewError(e).SetCode(zglobal.Code_ErrorBreakoff)
	} else if has >= 0 {
		if s, e := Redis.Get(rk).Result(); e != nil { // 是否有数据
			return nil, zservice.NewError(e).SetCode(zglobal.Code_ErrorBreakoff)
		} else {
			if tab, e := GetUserByID(ctx, zservice.StringToUint(s)); e != nil {
				return nil, e
			} else {
				return tab, nil
			}
		}
	}

	// 未找到 查表
	tab := UserTable{}

	// 验证数据库中是否存在
	if e := Mysql.Model(&tab).Where(fmt.Sprintf("phone = '%v'", phone)).First(&tab).Error; e != nil {
		return nil, zservice.NewError(e).SetCode(zglobal.Code_ErrorBreakoff)
	}

	// 更新缓存
	if e := Redis.Set(rk, zservice.Uint32ToString(tab.UID)).Err(); e != nil {
		ctx.LogError(e)
	}
	if e := Redis.Set(fmt.Sprintf(RK_UserInfo, tab.UID), zservice.JsonMustMarshalString(tab)).Err(); e != nil {
		ctx.LogError(e)
	}

	return &tab, nil
}

// 添加登陆名和密码
func (z *UserTable) AddLoginNameAndPassword(ctx *zservice.Context, name, password string) *zservice.Error {

	rk := fmt.Sprintf(RK_UserLoginName, name)
	// 锁
	un, e := Redis.Lock(rk)
	if e != nil {
		return e
	}
	defer un()

	// 验证重复
	if has, e := HasUserByLoginName(ctx, name); e != nil {
		return e
	} else if has {
		return zservice.NewError("user already exist:", name).SetCode(zglobal.Code_Zauth_UserAlreadyExist_LoginName)
	}

	z.LoginName = name
	z.LoginPassToken = zservice.RandomMD5()
	z.LoginPass = UserGenPassSign(z, password)

	return z.Save(ctx)
}

// 验证密码
func (z *UserTable) VerifyPass(ctx *zservice.Context, password string) bool {
	return z.LoginPass == UserGenPassSign(z, password)
}

// 存储
func (z *UserTable) Save(ctx *zservice.Context) *zservice.Error {

	rk_info := fmt.Sprintf(RK_UserInfo, z.UID)

	// 上锁
	un, e := Redis.Lock(rk_info)
	if e != nil {
		return e
	}
	defer un()

	if e := Mysql.Save(z).Error; e != nil {
		return zservice.NewError(e).SetCode(zglobal.Code_ErrorBreakoff)
	}

	// 删缓存
	zservice.Go(func() {
		if e := Redis.Del(rk_info).Err(); e != nil {
			ctx.LogError(zglobal.Code_Redis_DelFail, e)
		}
	})

	return nil
}
