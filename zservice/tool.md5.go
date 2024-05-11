package zservice

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

// 字符串 md5 编码
func MD5String(str string) string {
	hash := md5.New()
	io.WriteString(hash, str)
	return fmt.Sprintf("%x", hash.Sum(nil))
}

// 根据传入的文件路径对文件进行 MD5 计算
func Md5File(path string) (string, *Error) {
	hash := ""
	file, e := os.Open(path)

	if e != nil {
		return "", NewError(e)
	}
	defer file.Close()

	// 根据文件大小不同进行不同的处理
	stat, e := file.Stat()
	if e != nil {
		return "", NewError(e)
	}

	// 大于 100M 的文件，小内存读取
	if stat.Size() > 1024*1024*100 {
		md5hash := md5.New()
		if _, e := io.Copy(md5hash, file); e != nil {
			return "", NewError(e)
		}

		hash = hex.EncodeToString(md5hash.Sum(nil))

	} else {

		data, e := io.ReadAll(file)
		if e != nil {
			return "", NewError(e)
		}

		md5sum := md5.Sum(data)
		hash = hex.EncodeToString(md5sum[:])
	}

	return hash, nil
}
