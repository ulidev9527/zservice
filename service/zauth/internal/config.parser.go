package internal

import (
	"os"
	"zservice/zservice"
	"zservice/zservice/zglobal"
)

// 文件解析器方法接口
// 解析完成后会将解析后到数据存储到 redis
// 所有数据必须有 id 字段用于唯一标识
// 所有字段都会转换成小写格式处理
type FileParserFN func(file string) (map[string]string, *zservice.Error)

// 验证文件正确性
func parserFileVerify(fullpath string) *zservice.Error {
	fi, e := os.Stat(fullpath)
	if os.IsNotExist(e) {
		return zservice.NewError(e).SetCode(zglobal.Code_Zauth_config_FileNotExist)
	}
	if e != nil {
		return zservice.NewError(e).SetCode(zglobal.Code_Zauth_config_ParserFail)
	}
	if fi.IsDir() {
		return zservice.NewError("file is dir").SetCode(zglobal.Code_Zauth_config_PathIsDir)
	}

	return nil
}
