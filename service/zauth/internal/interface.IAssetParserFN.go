package internal

import (
	"zservice/zservice"
)

// 文件解析器方法接口
// 解析完成后会将解析后到数据存储到 redis
// 所有数据必须有 id 字段用于唯一标识
// 所有字段都会转换成小写格式处理
type IAssetConfigParserFN func(file string) (map[string]string, *zservice.Error)
