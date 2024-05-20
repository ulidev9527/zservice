package zauth

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"zservice/service/zauth/zauth_pb"
	"zservice/zservice"
	"zservice/zservice/ex/grpcservice"
	"zservice/zservice/ex/redisservice"
	"zservice/zservice/zglobal"

	clientv3 "go.etcd.io/etcd/client/v3"
)

var grpcClient zauth_pb.ZauthClient
var fileConfigMap = &sync.Map{}

type ZAuthConfig struct {
	EtcdServiceName string
	Etcd            *clientv3.Client
	Redis           *redisservice.GoRedisEX
	NsqConsumerAddr string // nsq consumer addr
	IsNsqd          bool
}

func Init(c *ZAuthConfig) {
	func() {
		conn, e := grpcservice.NewGrpcClient(&grpcservice.GrpcClientConfig{
			EtcdServiceName: c.EtcdServiceName,
			EtcdServer:      c.Etcd,
		})
		if e != nil {
			zservice.LogPanic(e)
			return
		}

		grpcClient = zauth_pb.NewZauthClient(conn)
	}()

}

// 检查权限, 没返回错误表示检查成功
func CheckAuth(ctx *zservice.Context, req *zauth_pb.CheckAuth_REQ) *zservice.Error {
	if res, e := grpcClient.CheckAuth(context.WithValue(context.Background(), grpcservice.GRPC_contextEX_Middleware_Key, ctx.ContextS2S), req); e != nil {
		return zservice.NewError(e).SetCode(zglobal.Code_ErrorBreakoff)
	} else if res.Code != zglobal.Code_SUCC {
		return zservice.NewError("check auth fail").SetCode(res.Code)
	} else {
		return nil
	}
}

// 发送验证码
func SendVerifyCode(ctx *zservice.Context, phone string) *zservice.Error {
	if phone == "" || phone[0] != '+' { // 手机号初步验证
		return zservice.NewError("verify phone fail").SetCode(zglobal.Code_ParamsErr)
	}

	if res, e := grpcClient.SMSSendVerifyCode(context.WithValue(context.Background(), grpcservice.GRPC_contextEX_Middleware_Key, ctx.ContextS2S), &zauth_pb.SMSSendVerifyCode_REQ{
		Phone:  phone,
		Serive: zservice.GetServiceName(),
	}); e != nil {
		return zservice.NewError("send verify code fail").SetCode(zglobal.Code_ErrorBreakoff)
	} else if res.Code == zglobal.Code_SUCC {
		return nil
	} else {
		return zservice.NewError("send verify code fail").SetCode(res.Code)
	}

}

// 验证验证码
func VerifyCode(ctx *zservice.Context, phone string, verifyCode string) *zservice.Error {
	if phone == "" || phone[0] != '+' { // 手机号初步验证
		return zservice.NewError("verify phone fail").SetCode(zglobal.Code_ParamsErr)
	}
	if verifyCode == "" || len(verifyCode) != 6 {
		return zservice.NewError("verify code fail").SetCode(zglobal.Code_ParamsErr)
	}

	if res, e := grpcClient.SMSVerifyCode(context.WithValue(context.Background(), grpcservice.GRPC_contextEX_Middleware_Key, ctx.ContextS2S), &zauth_pb.SMSVerifyCode_REQ{
		Phone:      phone,
		VerifyCode: verifyCode,
		Serive:     zservice.GetServiceName(),
	}); e != nil {
		return zservice.NewError("verify code fail").SetCode(zglobal.Code_ErrorBreakoff)
	} else if res.Code == zglobal.Code_SUCC {
		return nil
	} else {
		return zservice.NewError("verify code fail").SetCode(res.Code)
	}
}

// 获取文件配置

// 获取指定文件的配置
// 不传 key 返回所有配置数组
// 一个 key 返回一个对象
// 多个 key 返回数组
func GetFileConfig(ctx *zservice.Context, fileName string, v any, keys ...string) *zservice.Error {
	// 是否有配置，没有拉取配置
	val, has := fileConfigMap.Load(fileName)

	if !has { // 配置中心拉取配置
		res, e := grpcClient.GetFileConfig(ctx, &zauth_pb.GetFileConfig_REQ{
			FileName: fileName,
		})
		if e != nil {
			return zservice.NewError(e)
		}

		if res.GetCode() != zglobal.Code_SUCC {
			return zservice.NewError("get file config fail").SetCode(res.Code)
		}

		var maps map[string]string
		ee := json.Unmarshal([]byte(res.Value), &maps)
		if ee != nil {
			return zservice.NewError(ee)
		}

		fileConfigMap.Store(fileName, maps)

		val = maps
	}

	// 开始解析
	maps := val.(map[string]string)
	keyCount := len(keys)

	if keyCount == 1 { // 解析一个配置
		str := maps[keys[0]]
		json.Unmarshal([]byte(str), v)
		return nil
	}

	useMap := map[string]string{}

	if keyCount == 0 {
		useMap = maps
	} else {
		for _, _v := range keys {
			str := maps[_v]
			if str == "" {
				return zservice.NewError("key not exist", _v).SetCode(zglobal.Code_Zconfig_GetConfigFail)
			}
			useMap[_v] = str
		}
	}

	// 数组模式解析
	jStr := ""
	for _, _v := range useMap {
		jStr += _v + ","
	}
	jStr = jStr[0 : len(jStr)-1]
	jStr = fmt.Sprintf("[ %s ]", jStr)
	zservice.LogInfo(jStr)
	if e := json.Unmarshal([]byte(jStr), v); e != nil {
		return zservice.NewError(e).SetCode(zglobal.Code_Zconfig_GetConfigFail)
	}
	return nil
}
