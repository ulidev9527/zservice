package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"zserviceapps/packages/zservice"
)

// 生成常量
func Run_Gen_EConst(pbOutDir string) {

	zservice.LogInfo("start run gen EConst")

	outputStr_go := `package pb
`
	outputStr_cSharp := ""
	outputStr_cSharpTP := `namespace pb{
	public static class PBConst {
	$outputStr
	}
}`

	filePath := fmt.Sprintf("%s/econst.pb.go", pbOutDir)

	var lines []string
	if fileStr, e := os.ReadFile(filePath); e != nil {
		if !os.IsNotExist(e) {
			zservice.LogError(e)
		}
		return
	} else {
		lines = strings.Split(string(fileStr), "\n")
	}

	re_line := regexp.MustCompile(`EConst_.+ EConst =.+`) // 匹配对应行
	re_goTxt := regexp.MustCompile(`(= int32\(\d+)`)
	re_csTxt := regexp.MustCompile(`(= \d+)`)
	isStart := false
	for _, line := range lines {

		lineStr := strings.TrimSpace(line)
		// 忽略空行/注释行
		if lineStr == "" {
			continue
		}

		if isStart && lineStr == ")" {
			break
		}

		// 是否是匹配行
		if !re_line.MatchString(lineStr) {
			continue
		}

		isStart = true

		txt := strings.Replace(lineStr, "EConst_", "", 1)

		// go 输出
		goTxt := fmt.Sprintf("\nconst Const_%s",
			strings.Replace(txt, "EConst = ", "= int32(", 1))
		goTxt = re_goTxt.ReplaceAllString(goTxt, "$1)")

		outputStr_go += goTxt

		// c# 输出
		csTxt := fmt.Sprintf("\n        public const int %s",
			strings.Replace(txt, "EConst = ", "= ", 1))
		csTxt = re_csTxt.ReplaceAllString(csTxt, `$1;`)

		outputStr_cSharp += csTxt

	}

	os.WriteFile(fmt.Sprintf("%s/%s.go", pbOutDir, "pb.const"), []byte(outputStr_go), 0644)
	os.WriteFile(fmt.Sprintf("%s/%s.cs", pbOutDir, "pb.const"), []byte(strings.ReplaceAll(outputStr_cSharpTP, "$outputStr", outputStr_cSharp)), 0644)
	os.Remove(filePath)
	os.Remove(fmt.Sprintf("%s/Econst.cs", pbOutDir))

}
