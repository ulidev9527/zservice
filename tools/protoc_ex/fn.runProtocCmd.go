package main

import (
	"os"
	"os/exec"
	"zserviceapps/packages/zservice"
)

// 执行 protoc 命令
func runProtocCmd(tmpDir string, pbOutDir string, protoFilePaths []string) {

	// 调用 protoc 命令进行编译
	// 这里的问题在于 protoc 命令需要完整的文件路径作为参数
	// 而不是用空格分隔的字符串
	args := []string{
		"--go_out=.",
		"--go-grpc_out=.",
		"--csharp_out=.",
		"--proto_path=" + tmpDir,
	}
	// 将每个proto文件路径作为单独的参数
	args = append(args, protoFilePaths...)

	cmd_protoc := exec.Command("protoc", args...)
	cmd_protoc.Dir = pbOutDir
	cmd_protoc.Stdout = os.Stdout
	cmd_protoc.Stderr = os.Stderr

	zservice.LogInfo("cmd", cmd_protoc.String())
	if err := cmd_protoc.Run(); err != nil {
		zservice.LogError("cmd_protoc.Run() failed with \n", err)
		os.Exit(1)
	}
}
