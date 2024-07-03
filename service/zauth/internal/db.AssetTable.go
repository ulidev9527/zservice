package internal

import (
	"encoding/json"
	"fmt"
	"time"
	"zservice/service/zauth/zauth_pb"
	"zservice/zservice"
	"zservice/zservice/zglobal"

	"gorm.io/gorm"
)

// 资源
type AssetTable struct {
	gorm.Model
	Name    string         // 名称
	MD5     string         // md5
	Token   string         // 资源token
	Expires zservice.Ztime // 过期时间
	Size    uint64         // 文件大小
}

// 创建资源
func AssetCreate(ctx *zservice.Context, in *zauth_pb.AssetInfo) (*AssetTable, *zservice.Error) {

	if in.Md5 == "" {
		return nil, zservice.NewError("invalid md5").SetCode(zglobal.Code_ParamsErr)
	}

	// 准备写入数据
	tab := &AssetTable{
		Name:    in.Name,
		MD5:     in.Md5,
		Token:   zservice.MD5String(fmt.Sprintf("%s_%d_%d", in.Md5, in.Expires, time.Now().UnixMicro())),
		Expires: zservice.ZtimeUnixMilli(in.Expires),
		Size:    in.Size,
	}

	if e := tab.Save(ctx); e != nil {
		return nil, e
	}
	return tab, nil
}

// 获取资源
func AssetGetByToken(ctx *zservice.Context, token string) (*AssetTable, *zservice.Error) {

	rk_token := fmt.Sprintf(RK_AssetToken, token)

	tab := &AssetTable{}

	// 读取缓存
	if s, e := Redis.Get(rk_token).Result(); e != nil {
		if !DBService.IsNotFoundErr(e) {
			return nil, zservice.NewError(e)
		}
	} else {
		if e := json.Unmarshal([]byte(s), tab); e != nil {
			ctx.LogError(e)
		} else {
			return tab, nil
		}
	}

	// 读取数据库
	if e := Gorm.Where("token = ?", token).First(tab).Error; e != nil {
		if DBService.IsNotFoundErr(e) {
			return nil, zservice.NewError(e).SetCode(zglobal.Code_NotFound)
		}
		return nil, zservice.NewError(e)
	}
	return tab, nil
}

func (tab *AssetTable) Save(ctx *zservice.Context) *zservice.Error {

	if tab.Token == "" || tab.MD5 == "" {
		return zservice.NewError("invalid token or md5").SetCode(zglobal.Code_ParamsErr)
	}

	RK_token := fmt.Sprintf(RK_AssetToken, tab.Token)
	RK_md5 := fmt.Sprintf(RK_AssetMd5, tab.MD5)

	un, e := Redis.Lock(RK_md5)
	if e != nil {
		return e
	}
	defer un()

	if e := Gorm.Save(tab).Error; e != nil {
		return zservice.NewError("save asset error:", e.Error())
	}

	// 缓存
	zservice.Go(func() {

		if e := Redis.Del(RK_token).Err(); e != nil {
			ctx.LogError(e)
		}

		if e := Redis.Del(RK_md5).Err(); e != nil {
			ctx.LogError(e)
		}

	})

	return nil
}
