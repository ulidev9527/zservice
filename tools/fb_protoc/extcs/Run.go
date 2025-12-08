package extcs

import (
	"os"
	"regexp"
	"strings"
	"zserviceapps/packages/zservice"

	"github.com/spf13/viper"
)

func Run(conf *viper.Viper, inputFile string) {

	// // 同步文件
	// syncPaths := conf.GetStringSlice("sync_dir.cs")
	// if len(syncPaths) == 0 {
	// 	return
	// }

	fileBody, _ := os.ReadFile(inputFile)
	// fileBaseName := filepath.Base(inputFile)
	inputBodyStr := string(fileBody)

	regex := regexp.MustCompile(`(this\..+ = null;)`)
	inputBodyStr = regex.ReplaceAllString(inputBodyStr, "// $1 ")

	regex = regexp.MustCompile(`(public List<.+ )\{ get; set; \}`)
	inputBodyStr = regex.ReplaceAllString(inputBodyStr, "$1 = new();")

	inputBodyStr = removeString(inputBodyStr)
	if e := os.WriteFile(inputFile, []byte(inputBodyStr), 0644); e != nil {
		zservice.LogError(inputFile, e)
	}
}

func removeString(str string) string {
	return strings.ReplaceAll(str, "public static void ValidateVersion()", "// public static void ValidateVersion()")
}

func parseStruct(str string) {

}
