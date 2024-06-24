package zauth

import "sync"

// 文件配置缓存
var fileConfigCache = &sync.Map{}
