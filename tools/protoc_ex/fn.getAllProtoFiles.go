package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"zserviceapps/packages/zservice"
)

// proto 文件信息
type ProtoFileInfo struct {
	Path         string   // 文件路径
	NewPath      string   // 新的文件路径
	Body         string   // 文件内容
	MessageClass []string // 消息类
}

// 获取所有 proto 文件信息
func getAllProtoFilesInfo(root string) []*ProtoFileInfo {
	var fileList []*ProtoFileInfo

	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			zservice.LogError("WalkDir error:", err)
			return nil
		}
		if d.IsDir() || !strings.HasSuffix(d.Name(), ".proto") {
			return nil
		}
		protoInfo, err := parseProtoFile(path)
		if err != nil {
			zservice.LogError("parseProtoFile error:", err)
			return nil
		}
		fileList = append(fileList, protoInfo)
		return nil
	})
	if err != nil {
		zservice.LogError("WalkDir failed:", err)
	}
	return fileList
}

// 解析单个 proto 文件
func parseProtoFile(filePath string) (*ProtoFileInfo, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("read file %s: %w", filePath, err)
	}
	lines := strings.Split(string(content), "\n")
	newLines := []string{}
	messageClasses := []string{}

	optionGoPkg := regexp.MustCompile(`^option\s+go_package\s*=`)
	optionCSharp := regexp.MustCompile(`^option\s+csharp_namespace\s*=`)

	for idx, line := range lines {
		lineStr := strings.TrimSpace(line)
		if lineStr == "" || strings.HasPrefix(lineStr, "//") {
			continue
		}
		words := strings.Fields(lineStr)
		if len(words) == 0 {
			continue
		}
		switch words[0] {
		case "import":
			importFile := strings.Trim(lineStr[len("import"):], " \";")
			// import 路径相对当前 proto 文件
			importPath := filepath.Join(filepath.Dir(filePath), importFile)
			absImportPath, _ := filepath.Abs(importPath)
			newLines = append(newLines, fmt.Sprintf("import \"%s\";", filepath.Base(importFile)))
			zservice.LogInfo("import =>", absImportPath)
		case "option":
			if optionGoPkg.MatchString(lineStr) || optionCSharp.MatchString(lineStr) {
				zservice.LogWarn("ignore option line: ", idx, filePath, "=>", lineStr)
				continue
			}
			newLines = append(newLines, lineStr)
		case "message":
			if len(words) > 1 {
				msg := strings.TrimRight(words[1], "{")
				messageClasses = append(messageClasses, msg)
			}
			newLines = append(newLines, lineStr)
		default:
			newLines = append(newLines, lineStr)
		}
	}
	// 添加 option
	newLines = append(newLines, "option go_package = \"../pb\";")
	newLines = append(newLines, "option csharp_namespace = \"pb\";")

	return &ProtoFileInfo{
		Path:         filePath,
		NewPath:      filepath.Base(filePath),
		Body:         strings.Join(newLines, "\n"),
		MessageClass: messageClasses,
	}, nil
}
