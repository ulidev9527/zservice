package internal

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var (
	Mysql *gorm.DB
	Redis *redis.Client
)
