package internal

import "gorm.io/gorm"

// 账号表
type AccountTable struct {
	gorm.Model
	UID      uint64 // 用户唯一ID
	Account  string // 账号
	Password string // 密码
	Phone    string // 手机号 含区号 +86******
	State    uint   `gorm:"default:1"` // 账号状态 0 禁用 1 启用
}
