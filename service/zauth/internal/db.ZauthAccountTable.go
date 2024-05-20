package internal

import (
	"fmt"
	"zservice/zservice"
	"zservice/zservice/zglobal"

	"gorm.io/gorm"
)

// 账号表
type ZauthAccountTable struct {
	gorm.Model
	AccountID     uint   `gorm:"unique"` // 用户唯一ID
	LoginName     string `gorm:"unique"` // 登陆账号
	LoginPass     string // 登陆密码
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
	z := &ZauthAccountTable{AccountID: accID}
	if e := z.Save(ctx); e != nil {
		return nil, e
	}
	return z, nil
}

// 获取一个新的账号ID
func GetNewAccountID(ctx *zservice.Context) (uint, *zservice.Error) {
	return dbhelper.GetNewTableID(ctx, func() uint {
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
	return dbhelper.HasTableValue(ctx, &ZauthAccountTable{}, fmt.Sprintf(RK_AccountInfo, accountID), fmt.Sprintf("account_id = %v", accountID))
}

// 是否存在这个账号
func HasAccountByLoginName(ctx *zservice.Context, loginName string) (bool, *zservice.Error) {
	return dbhelper.HasTableValue(ctx, &ZauthAccountTable{}, fmt.Sprintf(RK_AccountLoginName, loginName), fmt.Sprintf("login_name = '%v'", loginName))
}

// 账号密码签名
func AccountPassSign(z *ZauthAccountTable, password string) string {
	return zservice.MD5String(fmt.Sprint(z.AccountID, z.PasswordToken, password))
}

// 获取账号
func GetAccountByAccountID(ctx *zservice.Context, accountID uint) (*ZauthAccountTable, *zservice.Error) {
	tab := ZauthAccountTable{}

	if e := dbhelper.GetTableValue(ctx, &tab, fmt.Sprintf(RK_AccountInfo, accountID), fmt.Sprintf("account_id = %v", accountID)); e != nil {
		return nil, e
	}
	return &tab, nil
}

// 根据登陆名获取账号
func GetAccountByLoginName(ctx *zservice.Context, loginName string) (*ZauthAccountTable, *zservice.Error) {
	tab := ZauthAccountTable{}
	if e := dbhelper.GetTableValue(ctx, &tab, fmt.Sprintf(RK_AccountLoginName, loginName), fmt.Sprintf("login_name = '%v'", loginName)); e != nil {
		return nil, e
	}
	return &tab, nil
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
	z.LoginPass = AccountPassSign(z, password)

	return z.Save(ctx)
}

// 验证密码
func (z *ZauthAccountTable) VerifyPass(ctx *zservice.Context, password string) bool {
	return z.LoginPass == AccountPassSign(z, password)
}

// 存储
func (z *ZauthAccountTable) Save(ctx *zservice.Context) *zservice.Error {
	if z.AccountID == 0 {
		return zservice.NewError("no account id").SetCode(zglobal.Code_ParamsErr)
	}

	rk_info := fmt.Sprintf(RK_AccountInfo, z.AccountID)

	// 上锁
	un, e := Redis.Lock(rk_info)
	if e != nil {
		return e
	}
	defer un()

	if z.ID == 0 { // 创建
		if e := Mysql.Create(z).Error; e != nil {
			return zservice.NewError(e).SetCode(zglobal.Code_ErrorBreakoff)
		}
	} else { // 更新
		if e := Mysql.Save(z).Error; e != nil {
			return zservice.NewError(e).SetCode(zglobal.Code_ErrorBreakoff)
		}
	}

	// 删缓存
	if e := Redis.Del(rk_info).Err(); e != nil {
		zservice.LogError(zglobal.Code_Redis_DelFail, e)
	}

	return nil
}
