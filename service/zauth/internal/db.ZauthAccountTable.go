package internal

import (
	"fmt"
	"zservice/zglobal"
	"zservice/zservice"

	"gorm.io/gorm"
)

// 账号表
type ZauthAccountTable struct {
	gorm.Model
	AccountID     uint   `gorm:"unique"` // 用户唯一ID
	LoginName     string `gorm:"unique"` // 登陆账号
	LoginPassword string // 登陆密码
	Phone         string `gorm:"unique"`    // 手机号 含区号 +86******
	State         uint   `gorm:"default:1"` // 账号状态 0 禁用 1 启用
	PasswordToken string // 密码令牌
}

// 创建一个新的账号
func CreateAccount(ctx *zservice.Context) (*ZauthAccountTable, *zservice.Error) {
	accID, e := GetNewAccountID(ctx)
	if e != nil {
		return nil, e
	}
	return (&ZauthAccountTable{AccountID: accID}).Save(ctx)
}

// 获取一个新的账号ID
func GetNewAccountID(ctx *zservice.Context) (uint, *zservice.Error) {
	return GetNewID(ctx, func() uint {
		return uint(zservice.RandomIntRange(1000000, 999999999)) // 7-9位数
	}, HasAccountByID, func(e *zservice.Error) *zservice.Error {
		if e.GetCode() == zglobal.Code_Zauth_GenIDCountMaxErr {
			return e.SetCode(zglobal.Code_Zauth_AccountGenIDCountMaxErr)
		}
		return e
	})
}

// 是否存在这个账号
func HasAccountByID(ctx *zservice.Context, accountID uint) (bool, *zservice.Error) {
	return HasTableValue(ctx, &ZauthAccountTable{}, fmt.Sprintf(RK_AccountInfo, accountID), fmt.Sprintf("account_id = %v", accountID))
}

// 是否存在这个账号
func HasAccountByLoginName(ctx *zservice.Context, loginName string) (bool, *zservice.Error) {
	return HasTableValue(ctx, &ZauthAccountTable{}, fmt.Sprintf(RK_AccountLoginName, loginName), fmt.Sprintf("login_name = '%v'", loginName))
}

// 添加登陆名和密码
func (z *ZauthAccountTable) AddLoginNameAndPassword(ctx *zservice.Context, name, password string) *zservice.Error {

	rk := fmt.Sprintf(RK_AccountLoginName, name)
	// 锁
	un, e := Redis.Lock(rk)
	if e != nil {
		return e
	}
	defer un()

	// 验证重复
	if has, e := HasAccountByLoginName(ctx, name); e != nil {
		return e
	} else if has {
		return zservice.NewError("account already exist:", name).SetCode(zglobal.Code_Zauth_AccountAlreadyExist_LoginName)
	}

	z.LoginName = name
	z.PasswordToken = zservice.RandomMD5()
	z.LoginPassword = zservice.MD5String(fmt.Sprint(z.AccountID, z.PasswordToken, password))

	_, e = z.Save(ctx)
	return e
}

// 存储
func (z *ZauthAccountTable) Save(ctx *zservice.Context) (*ZauthAccountTable, *zservice.Error) {
	if z.AccountID == 0 {
		return nil, zservice.NewError("no account id").SetCode(zglobal.Code_ParamsErr)
	}

	rk_info := fmt.Sprintf(RK_AccountInfo, z.AccountID)

	// 上锁
	un, e := Redis.Lock(rk_info)
	if e != nil {
		return nil, e
	}
	defer un()

	if z.ID == 0 { // 创建
		if e := Mysql.Create(&z).Error; e != nil {
			return nil, zservice.NewError(e).SetCode(zglobal.Code_ErrorBreakoff)
		}
	} else { // 更新
		if e := Mysql.Save(&z).Error; e != nil {
			return nil, zservice.NewError(e).SetCode(zglobal.Code_ErrorBreakoff)
		}
	}

	// 存 redis
	if e := Redis.HMSet(rk_info, &z).Err(); e != nil {
		ctx.LogError(e)
	}

	return z, nil
}
