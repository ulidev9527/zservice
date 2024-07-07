package internal

import (
	"zservice/service/zauth/zauth_pb"
	"zservice/zservice"
)

func Logic_LoginCheck(ctx *zservice.Context, in *zauth_pb.LoginCheck_REQ) *zauth_pb.LoginCheck_RES {

	resultRES := &zauth_pb.LoginCheck_RES{Code: zservice.Code_Fail}

	if at, e := GetTokenInfo(ctx, in.Token); e != nil {
		ctx.LogError(e)
		resultRES.Code = e.GetCode()
	} else if at.UID == 0 { // uid 检查
		ctx.LogInfo("uid is zero", in.Token)
	} else if !at.HasLoginService(in.Service) { // 服务检查
		ctx.LogInfo("service is not in", in.Token, in.Service)
	} else if !at.TokenCheck(in.TokenSign) { // 签名检查
		ctx.LogInfo("sign check fail", in.Token, in.Service)
	} else { // 用户信息获取
		resultRES.Code = zservice.Code_SUCC
		if tab, e := GetUserByUID(ctx, at.UID); e != nil {
			ctx.LogError(e)
			resultRES.Code = e.GetCode()
		} else {
			resultRES.UserInfo = tab.ToUserInfo()
		}
	}

	return resultRES
}
