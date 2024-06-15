package internal

import (
	"bufio"
	"os"
	"zservice/zservice"
	"zservice/zservice/zglobal"

	"github.com/derekparker/trie"
)

type ___ZZZZString struct {
	RuneMap map[rune]int
	Trie    *trie.Trie
}

var ZZZZString = &___ZZZZString{}

func init() {

	ZZZZString.RuneMap = make(map[rune]int)
	ZZZZString.Trie = trie.New()

}

// 重新加载 zzzz 字符串
func (s *___ZZZZString) Reload(ctx *zservice.Context, path string) *zservice.Error {
	if file, e := os.Open(path); e != nil {
		return zservice.NewError("error open " + path).SetCode(zglobal.Code_ErrorBreakoff)
	} else {
		defer file.Close()
		scanner := bufio.NewScanner(file)
		size := 0
		for scanner.Scan() {
			for _, v := range scanner.Text() {
				size += 1
				s.RuneMap[v] = 1
			}
		}
		if e := scanner.Err(); e != nil {
			return zservice.NewError(e).SetCode(zglobal.Code_Fail)
		} else {
			return nil
		}
	}
}

// 是否有 zzzz 字符串
func (s *___ZZZZString) Has(ctx *zservice.Context, str string) bool {
	_, has := s.Trie.Find(str)
	return has
}