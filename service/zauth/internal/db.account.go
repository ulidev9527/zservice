package internal

import "gorm.io/gorm"

// 账号表
type AccountTable struct {
	gorm.Model
	UID      uint64 // 用户唯一ID
	Phone    string // 手机号 含区号 +86******
	Account  string // 账号
	Password string // 密码
	Status   int    `gorm:"default:0"` // 状态 0为初始状态,无法使用
}
