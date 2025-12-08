package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"zserviceapps/packages/zservice"
	"zserviceapps/tools/fb_protoc/extcs"
	"zserviceapps/tools/fb_protoc/extgo"

	"github.com/spf13/viper"
)

var extHandlerMap = map[string]func(*viper.Viper, string){
	".go": extgo.Run,
	".cs": extcs.Run,
}
var conf *viper.Viper

var autoGenPrefix = "auto_gen_fb_ex_" // 文件输出前缀

func main() {

	if len(os.Args) < 3 {
		return
	}

	tmpDir := os.Args[1]
	targetDir := os.Args[2]

	info, e := os.Stat(tmpDir)
	if e != nil || !info.IsDir() {
		return
	}

	// 读取配置文件
	confPath := filepath.Join(targetDir, "conf.yaml")
	if b, e := os.ReadFile(confPath); e != nil {
		zservice.LogError(confPath, e)
		return
	} else {
		conf = viper.New()
		conf.SetConfigType("yaml")

		if e := conf.ReadConfig(bytes.NewBuffer(b)); e != nil {
			zservice.LogError(e)
			return
		}
	}

	filesList := []string{}

	// 处理tmpDir 文件
	if e = filepath.Walk(tmpDir, func(path string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !f.IsDir() {
			filesList = append(filesList, path)
			if fn, ok := extHandlerMap[filepath.Ext(path)]; ok {
				fn(conf, path)
			}
		}
		return nil
	}); e != nil {
		zservice.LogError(e)
		return
	}

	// 删除输入目录所有的 auto_gen_protoc_ex_开头的 文件
	rmAutoGenDirMap := map[string]int{}

	// 同步文件位置
	for _, inputFile := range filesList {
		ext := filepath.Ext(inputFile)
		syncPaths := conf.GetStringSlice("sync_dir" + ext) // sync_dir+ ".XXXX"
		if len(syncPaths) == 0 {
			continue
		}

		fileBody, _ := os.ReadFile(inputFile)
		fileBaseName := filepath.Base(inputFile)

		for _, path := range syncPaths {

			if !filepath.IsAbs(path) { // 绝对路径转到目标位置的相对路径
				path = filepath.Join(targetDir, path)
			}

			if _, ok := rmAutoGenDirMap[path]; !ok {
				rmAutoGenDirMap[path] = 1
				// 删除输出目录所有的 auto_gen_protoc_ex_开头的 文件
				if files, e := os.ReadDir(path); e != nil {
					zservice.LogErrorf("os.ReadDir() failed with \n %s", e)
				} else {
					for _, file := range files {
						if strings.HasPrefix(file.Name(), autoGenPrefix) {
							os.Remove(filepath.Join(path, file.Name()))
							zservice.LogInfo("rm", file.Name())
						}
					}
				}
			}

			if stat, err := os.Stat(path); err == nil && stat.IsDir() {
				output := filepath.Join(path, autoGenPrefix+fileBaseName)
				if e := os.WriteFile(output, fileBody, 0644); e != nil {
					zservice.LogError("write fail:", output, e)
				} else {
					zservice.LogInfo("cp", fileBaseName, output)
				}
			}
		}

	}

}
