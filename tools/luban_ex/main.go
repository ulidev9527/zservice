package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
	"zserviceapps/packages/zservice"
)

// luban 功能扩展
func main() {

	// 获取启动目录
	projectDir, err := os.Getwd()
	if err != nil {
		zservice.LogPanicf("os.Getwd() failed with \n %s", err)
	}

	// 切换根目录
	if s, e := zservice.GetGomodDir(projectDir); e != nil {
		zservice.LogPanic(e)
	} else {
		projectDir = s
	}

	inputDir := fmt.Sprintf("%s/test/luban_ex_test", projectDir)
	// 获取传入参数
	if len(os.Args) >= 2 {
		inputDir = os.Args[1]
	}

	// inputDir 不是文件夹，切换到文件夹
	if s, e := os.Stat(inputDir); e != nil {
		zservice.LogPanic(e)
	} else if !s.IsDir() {
		inputDir = filepath.Dir(inputDir)
	}

	// 切换根目录
	if s, e := zservice.GetGomodDir(inputDir); e != nil {
		zservice.LogPanic(e)
	} else {
		projectDir = s
	}

	// 切换到项目目录
	if e := os.Chdir(projectDir); e != nil {
		zservice.LogPanic(e)
	}

	zservice.LogInfo("projectDir", projectDir)

	// 创建临时目录
	tmpDir := filepath.Join(projectDir, "~__luban_ex_temp_dir_"+time.Now().Format("20060102150405"))
	if err := os.MkdirAll(tmpDir, 0755); err != nil {
		zservice.LogPanicf("os.MkdirAll() failed with \n %s", err)
	}
	defer os.RemoveAll(tmpDir) // 删除临时目录

	lubanConfPath := fmt.Sprintf("%s/luban.conf", inputDir)
	if s, e := os.Stat(lubanConfPath); e != nil {
		zservice.LogPanic(e)
	} else if s.IsDir() {
		zservice.LogPanic("can`t find luban.conf file")
	}
	zservice.LogInfo("luban.conf =>", lubanConfPath)

	argsJson := []string{
		"dotnet",
		fmt.Sprintf("%s/tools/luban_ex/Tools/Luban/Luban.dll", projectDir),
		"-t all",
		"-d json",
		fmt.Sprintf("--conf %s", lubanConfPath),
		"-x outputDataDir=../static/luban_json",
	}

	argsGO := []string{
		"dotnet",
		fmt.Sprintf("%s/tools/luban_ex/Tools/Luban/Luban.dll", projectDir),
		"-t all",
		"-c go-json",
		fmt.Sprintf("--conf %s", lubanConfPath),
		fmt.Sprintf("-x outputCodeDir=%s/code_go", inputDir),
		"-x lubanGoModule=luban",
	}

	argsCSharp := []string{
		"dotnet",
		fmt.Sprintf("%s/tools/luban_ex/Tools/Luban/Luban.dll", projectDir),
		"-t all",
		"-c cs-newtonsoft-json",
		fmt.Sprintf("--conf %s", lubanConfPath),
		fmt.Sprintf("-x outputCodeDir=%s/Codecsharp", inputDir),
	}

	for _, args := range [][]string{argsJson, argsGO, argsCSharp} {
		// 执行命令

		var cmd_runtime *exec.Cmd

		switch runtime.GOOS {
		case "windows":
			cmd_runtime = exec.Command("cmd", "/c", strings.Join(args, " "))
		default:
			cmd_runtime = exec.Command("bash", "-c", strings.Join(args, " "))
		}

		// 执行命令
		cmd_runtime.Dir = inputDir
		cmd_runtime.Stdout = os.Stdout
		cmd_runtime.Stderr = os.Stderr

		zservice.LogInfo("cmd =>", cmd_runtime.String())

		if err := cmd_runtime.Run(); err != nil {
			zservice.LogPanicf("cmd run failed with %s", err)
		}

	}

}
