package zauth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"zservice/service/zauth/internal"
	"zservice/service/zauth/zauth_pb"
	"zservice/zservice"
	"zservice/zservice/ex/ginservice"
	"zservice/zservice/ex/grpcservice"
	"zservice/zservice/ex/nsqservice"
	"zservice/zservice/ex/redisservice"
	"zservice/zservice/zglobal"

	"github.com/gin-gonic/gin"
	"github.com/nsqio/go-nsq"
	clientv3 "go.etcd.io/etcd/client/v3"
)

var grpcClient zauth_pb.ZauthClient
var fileConfigMap = &sync.Map{}

type ZAuthInitConfig struct {
	ZauthServiceName string // 权限服务名称
	Etcd             *clientv3.Client
	Redis            *redisservice.GoRedisEX
	NsqConsumerAddrs string // nsq consumer addr
	IsNsqdAddr       bool
}

func Init(c *ZAuthInitConfig) {
	func() {
		conn, e := grpcservice.NewGrpcClient(&grpcservice.GrpcClientConfig{
			ZauthServiceName: c.ZauthServiceName,
			EtcdServer:       c.Etcd,
		})
		if e != nil {
			zservice.LogPanic(e)
			return
		}

		grpcClient = zauth_pb.NewZauthClient(conn)
	}()

	if c.ZauthServiceName == "" {
		zservice.LogPanic("ZauthServiceName is nil")
	}

	nsqservice.NewNsqConsumer(&nsqservice.NsqConsumerConfig{
		Addrs:      c.NsqConsumerAddrs,
		IsNsqdAddr: c.IsNsqdAddr,
		Topic:      internal.NSQ_FileConfig_Change,
		Channel:    fmt.Sprintf("%s-%s", zservice.GetServiceName(), zservice.RandomXID()),
		OnMessage: func(m *nsq.Message) error {
			fileName := string(m.Body)
			zservice.LogInfo("Update config ", fileName)
			fileConfigMap.Delete(fileName)
			return nil
		},
	})

}

// 授权检查
func GinCheckAuthMiddleware(zs *zservice.ZService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		zctx := ginservice.GetCtxEX(ctx)
		zctx.AuthSign = zservice.MD5String(ctx.Request.UserAgent()) // 生成签名

		// 授权查询
		if e := CheckAuth(zctx, &zauth_pb.CheckAuth_REQ{
			Auth: string(zservice.JsonMustMarshal([]string{zservice.GetServiceName(), strings.ToLower(ctx.Request.Method), ctx.Request.URL.Path})),
		}); e != nil {

			ctx.JSON(http.StatusOK, &zglobal.Default_RES{
				Code: e.GetCode(),
				Msg:  zctx.TraceID,
			})
			ctx.Abort()
			return
		}

		ctx.Next()
	}
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
