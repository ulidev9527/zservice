package internal

import (
	"bufio"
	"bytes"
	"zservice/service/zauth/zauth_pb"
	"zservice/zservice"

	"github.com/derekparker/trie"
)

type ___ZZZZString struct {
	RuneMap map[rune]int
	Trie    *trie.Trie
}

var ZZZZString = &___ZZZZString{}

func InitZZZZString() {

	ZZZZString.RuneMap = make(map[rune]int)
	ZZZZString.Trie = trie.New()
	ctx := zservice.NewContext()
	if e := ZZZZString.Reload(ctx); e != nil {
		ctx.LogWarn(e.Error())
	}
}

// 重新加载 zzzz 字符串
func (s *___ZZZZString) Reload(ctx *zservice.Context) *zservice.Error {
	kvRES := &zauth_pb.GetServiceKV_RES{}
	if res := Logic_GetServiceKV(ctx, &zauth_pb.GetServiceKV_REQ{
		Key:     KV_ZZZZString,
		Service: zservice.GetServiceName(),
	}); res.Code != zservice.Code_SUCC {
		return zservice.NewError("error get KV: "+KV_ZZZZString, res)
	} else {
		kvRES = res
	}

	if res := Logic_DownloadAsset(ctx, &zauth_pb.DownloadAsset_REQ{AssetID: kvRES.Value}); res.Code != zservice.Code_SUCC {
		return zservice.NewError("error download asset: "+kvRES.Value, res)
	} else {
		scanner := bufio.NewScanner(bytes.NewReader(res.Info.Data))
		size := 0
		for scanner.Scan() {
			for _, v := range scanner.Text() {
				size += 1
				s.RuneMap[v] = 1
			}
		}
		ctx.LogInfo("reload zzzz string size: ", len(res.Info.Data))
	}
	return nil
}

// 是否有 zzzz 字符串
func (s *___ZZZZString) Has(ctx *zservice.Context, str string) bool {
	_, has := s.Trie.Find(str)
	return has
}
