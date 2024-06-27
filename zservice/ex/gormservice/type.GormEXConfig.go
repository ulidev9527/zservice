package gormservice

import "fmt"

type GormEXConfig struct {
	DBName string
	Addr   string // IP:PORT
	User   string
	Passwd string
	Debug  bool
}

func (conf *GormEXConfig) GetUri(isUsePasswd ...bool) string {

	if isUsePasswd != nil && isUsePasswd[0] {
		return fmt.Sprintf("%s:%s@tcp(%v)/%s?parseTime=true", conf.User, conf.Passwd, conf.Addr, conf.DBName)
	} else {
		return fmt.Sprintf("%s:%s@tcp(%v)/%s?parseTime=true", conf.User, "******", conf.Addr, conf.DBName)
	}
}
