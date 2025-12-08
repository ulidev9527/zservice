package main

import (
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"zserviceapps/packages/zservice"
)

// proto message 信息
type ProtoMessageInfo struct {
	FileName string                   // 文件名称
	Name     string                   // 消息名称
	Fields   []*ProtoMessageFieldInfo // 字段信息
}

// proto 字段信息
type ProtoMessageFieldInfo struct {
	Name  string // 字段名称
	KType string // 字段类型 仅类名
	MType string // 字段类型 纯类型/含数组，不含类名
	Type  string // 字段类型 完整类型，类型+类名 []XX []*XX int
}

// 解析 proto 文件 message 信息
func getProtoMessageInfo(protoFiles []*ProtoFileInfo, pbOutDir string) []*ProtoMessageInfo {
	// 1. 收集所有消息类名，避免重复
	messageClassMap := make(map[string]struct{})
	for _, protoFile := range protoFiles {
		for _, messageClass := range protoFile.MessageClass {
			if _, exists := messageClassMap[messageClass]; exists {
				zservice.LogWarn("has duplicate messageClass =>", messageClass)
				continue
			}
			messageClassMap[messageClass] = struct{}{}
		}
	}

	// 2. 收集所有 .pb.go 文件
	goFiles := []string{}
	files, err := os.ReadDir(pbOutDir)
	if err != nil {
		zservice.LogError("os.ReadDir() failed:", err)
		return nil
	}
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".pb.go") {
			goFiles = append(goFiles, filepath.Join(pbOutDir, file.Name()))
		}
	}
	sort.Strings(goFiles) // 保证顺序一致

	// 3. 正则表达式
	reTypeName := regexp.MustCompile(`type (\w+) struct`)
	reMethodInfo := regexp.MustCompile(`func.+\*(\w+)\) Get(\w+)\(\) (.+) \{`)

	infoMap := make(map[string]*ProtoMessageInfo)

	// 4. 解析每个 go 文件
	for _, filePath := range goFiles {
		content, err := os.ReadFile(filePath)
		if err != nil {
			zservice.LogError("os.ReadFile() failed:", err)
			continue
		}
		lines := strings.Split(string(content), "\n")

		for _, line := range lines {
			lineStr := strings.TrimSpace(line)
			if lineStr == "" {
				continue
			}
			// 结构体类型
			if match := reTypeName.FindStringSubmatch(lineStr); len(match) > 1 {
				typeName := match[1]
				if typeName == "x" || typeName == "" {
					continue
				}
				if _, ok := infoMap[typeName]; ok {
					zservice.LogWarn("typeName already exists:", typeName)
					continue
				}
				if _, ok := messageClassMap[typeName]; !ok {
					zservice.LogWarn("Ignore typeName (not in proto):", typeName)
					continue
				}
				infoMap[typeName] = &ProtoMessageInfo{
					FileName: filepath.Base(filePath),
					Name:     typeName,
				}
				continue
			}
			// 字段方法
			if match := reMethodInfo.FindStringSubmatch(lineStr); len(match) == 4 {
				typeName := match[1]
				fieldName := match[2]
				fieldType := match[3]
				fieldMType, fieldKType := parseFieldType(fieldType)

				if _, ok := infoMap[typeName]; !ok {
					zservice.LogWarn("typeName not found for field:", typeName, fieldName)
					continue
				}
				infoMap[typeName].Fields = append(infoMap[typeName].Fields, &ProtoMessageFieldInfo{
					Name:  fieldName,
					MType: fieldMType,
					Type:  fieldType,
					KType: fieldKType,
				})
			}
		}
	}

	// 5. 转换为 slice 并排序
	protoMessageInfos := make([]*ProtoMessageInfo, 0, len(infoMap))
	for _, info := range infoMap {
		protoMessageInfos = append(protoMessageInfos, info)
	}
	sort.Slice(protoMessageInfos, func(i, j int) bool {
		return protoMessageInfos[i].Name < protoMessageInfos[j].Name
	})
	return protoMessageInfos
}

// 字段类型解析辅助
func parseFieldType(fieldType string) (mType, kType string) {
	switch {
	case strings.HasPrefix(fieldType, "[]*"):
		return "[]*", fieldType[3:]
	case strings.HasPrefix(fieldType, "[]"):
		return "[]", fieldType[2:]
	case strings.HasPrefix(fieldType, "*"):
		return "*", fieldType[1:]
	case strings.HasPrefix(fieldType, "E_"):
		return "E_", fieldType
	default:
		return fieldType, fieldType
	}
}
