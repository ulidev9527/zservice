package internal

import (
	"bufio"
	"os"
	"time"
	"zservice/zservice"

	"github.com/derekparker/trie"
)

type ___ZZZZString struct {
	RuneMap map[rune]int
	Trie    *trie.Trie
}

var ZZZZString = &___ZZZZString{}

func InitZZZZ() {

	ZZZZString.RuneMap = make(map[rune]int)
	ZZZZString.Trie = trie.New()
	zservice.Go(func() {
		time.Sleep(time.Second)
		if e := ZZZZString.Reload(zservice.NewContext()); e != nil {
			zservice.LogError(e)
		}
	})
}

// 重新加载 zzzz 字符串
func (s *___ZZZZString) Reload(ctx *zservice.Context) *zservice.Error {

	if file, e := os.Open(FI_ZZZZStringFile); e != nil {
		return zservice.NewError("error open "+FI_ZZZZStringFile, e)
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
			return zservice.NewError(e)
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
