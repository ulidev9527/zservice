package internal

import (
	"fmt"
	"os"
	"zservice/zglobal"
	"zservice/zservice"
)

// 文件解析器方法接口
// 解析完成后会将解析后到数据存储到 redis
// 所有数据必须有 id 字段用于唯一标识
// 所有字段都会转换成小写格式处理
type FileParserFN func(file string) *zservice.Error

var fileParserMap = map[uint32]FileParserFN{}

func init() {
	// 注册解析器
	fileParserMap[zglobal.E_ZConfig_Parser_Excel] = ParserExcel
}

// 验证文件正确性
func ParserFileVerify(file string) *zservice.Error {
	fi, e := os.Stat(file)
	if os.IsNotExist(e) {
		return zservice.NewError(e).SetCode(zglobal.Code_Zconfig_FileNotExist)
	}
	if e != nil {
		return zservice.NewError(e).SetCode(zglobal.Code_Zconfig_ParserFail)
	}
	if fi.IsDir() {
		return zservice.NewError("file is dir").SetCode(zglobal.Code_Zconfig_PathIsDir)
	}

	return nil
}

// 获取文件的 md5
func GetMd5(fullPath string) (md5Str string, e *zservice.Error) {
	if e := ParserFileVerify(fullPath); e != nil {
		if e.GetCode() != zglobal.Code_SUCC {
			return "", e
		}
	}

	// md5信息文件
	md5FileFullPath := fmt.Sprintf("%s.md5", fullPath)
	e = ParserFileVerify(md5FileFullPath)
	if e != nil {
		// 是否需要创建文件
		if e.GetCode() == zglobal.Code_Zconfig_FileNotExist {
			md5Str, e := zservice.Md5File(fullPath)
			if e != nil {
				return "", e.SetCode(zglobal.Code_Zconfig_GetFileMd5Fail)
			}
			// 写入md5信息到文件
			ee := os.WriteFile(md5FileFullPath, []byte(md5Str), 0644)
			if ee != nil {
				return "", zservice.NewError(ee).SetCode(zglobal.Code_Zconfig_GetFileMd5Fail)
			}
		} else {
			return "", e
		}
	}

	// 读取 md5 信息文件
	data, ee := os.ReadFile(md5FileFullPath)
	if ee != nil {
		return "", zservice.NewError(ee).SetCode(zglobal.Code_Zconfig_GetFileMd5Fail)
	}
	if len(data) == 0 {
		return "", zservice.NewError("file md5 is empty").SetCode(zglobal.Code_Zconfig_GetFileMd5Fail)
	}
	return string(data), nil
}

// 解析文件
func ParserFile(fileName string, parserType uint32) *zservice.Error {
	// 解析器
	parserFN, ok := fileParserMap[parserType]
	if !ok {
		return zservice.NewError("parser not found").SetCode(zglobal.Code_Zconfig_ParserNotExist)
	}
	if e := parserFN(fileName); e != nil && e.GetCode() != zglobal.Code_SUCC {
		return e
	}

	if e := Nsq.Publish(NSQ_FileConfig_Change, []byte(fileName)); e != nil {
		return zservice.NewError(e).SetCode(zglobal.Code_ErrorBreakoff)
	}
	return nil
}
