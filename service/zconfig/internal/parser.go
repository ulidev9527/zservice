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
type FileParserFN func(file string) (map[string]string, *zservice.Error)

var fileParserMap = map[uint32]FileParserFN{}

func init() {
	// 注册解析器
	fileParserMap[zglobal.E_ZConfig_Parser_Excel] = ParserExcel
}

// 验证文件正确性
func parserFileVerify(fullpath string) *zservice.Error {
	fi, e := os.Stat(fullpath)
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
func getFileMD5(fullPath string) (md5Str string, e *zservice.Error) {
	// md5 信息文件路径
	md5FileFullPath := fmt.Sprintf("%s.md5", fullPath)
	// 验证文件正确性，文件不存在则创建
	if e := parserFileVerify(md5FileFullPath); e != nil {
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
			return md5Str, nil
		} else {
			return "", e
		}
	}

	// 文件存在，直接读取 md5 文件
	// 读取 md5 信息文件
	data, ee := os.ReadFile(md5FileFullPath)
	if ee != nil {
		return "", zservice.NewError(ee).SetCode(zglobal.Code_Zconfig_GetFileMd5Fail)
	}
	if len(data) == 0 {
		// 空文件 删除
		ee := os.Remove(md5FileFullPath)
		if ee != nil {
			return "", zservice.NewError("file md5 is empty, del fail", ee).SetCode(zglobal.Code_Zconfig_GetFileMd5Fail)
		}
		return "", zservice.NewError("file md5 is empty").SetCode(zglobal.Code_Zconfig_GetFileMd5Fail)
	}
	return string(data), nil
}

// 解析文件
func ParserFile(fileName string, parserType uint32) *zservice.Error {

	// 上锁
	RK_lock := fmt.Sprintf(RK_FileConfigLcok, fileName)
	un, e := Redis.Lock(RK_lock)
	if e != nil {
		return e
	}
	defer un()

	// 全路径
	fullPath := fmt.Sprintf("%s/%s", FI_StaticRoot, fileName)

	// 解析器获取
	parserFN, ok := fileParserMap[parserType]
	if !ok {
		return zservice.NewError("parser not found").SetCode(zglobal.Code_Zconfig_ParserNotExist)
	}

	// md5 检查
	fileMD5, e := getFileMD5(fullPath)
	if e != nil {
		return e
	}

	// md5 匹配, 没有数据或者无变化，返回 nil 进行解析
	rKeyMd5 := fmt.Sprintf(RK_FileMD5, fileName)
	if e := func() *zservice.Error {
		has, e := Redis.Exists(rKeyMd5).Result()
		if e != nil {
			return zservice.NewError(e).SetCode(zglobal.Code_Zconfig_ParserFail)
		}
		if has == 0 { // 不存在 需要更新
			return nil
		}
		str, e := Redis.Get(rKeyMd5).Result()
		if e != nil {
			return zservice.NewError(e).SetCode(zglobal.Code_Zconfig_ParserFail)
		}
		if str == fileMD5 {
			return zservice.NewError("file md5 not change:", fileName).SetCode(zglobal.Code_Zconfig_FileMd5NotChange)
		}
		return nil // 有变化
	}(); e != nil {
		return e
	}

	// 解析
	maps, e := parserFN(fileName)
	if e != nil {
		return e
	}

	// 存储到 redis
	rKeyFile := fmt.Sprintf(RK_FileConfig, fileName)
	if e := func() *zservice.Error {

		if e := Redis.Del(rKeyMd5).Err(); e != nil {
			return zservice.NewError(e).SetCode(zglobal.Code_ErrorBreakoff)
		}

		if e := Redis.HMSet(rKeyFile, maps).Err(); e != nil {
			return zservice.NewError(e).SetCode(zglobal.Code_ErrorBreakoff)
		}
		return nil
	}(); e != nil {
		return e
	}

	// 通知文件变更
	if e := NsqFileConfigChange(fileName); e != nil {
		return e
	}

	// 保存 md5
	if e := Redis.Set(rKeyMd5, fileMD5, 0).Err(); e != nil {
		zservice.LogError(e)
	}

	return nil
}
