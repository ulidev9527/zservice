package zservice

import (
	"os"
	"strings"
)

func Getenv(key string) string {
	return os.Getenv(key)
}

func GetenvBool(key string) bool {
	return Convert_StringToBoolean(Getenv(key))
}

// json
func GetenvStringSplit(key string, split ...string) []string {
	str := Getenv(key)
	if str == "" {
		return []string{}
	}

	if len(split) > 0 {
		return strings.Split(str, split[0])
	}

	return strings.Split(str, ",")
}
