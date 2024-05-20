package internal

import (
	"zservice/service/zauth/zauth_pb"
	"zservice/zservice"
	"zservice/zservice/zglobal"
)

// 检查权限
func CheckAuth(ctx *zservice.Context, in *zauth_pb.CheckAuth_REQ) *zauth_pb.CheckAuth_RES {

	// 获取 token
	tk, e := GetToken(ctx.AuthToken)
	if e != nil {
		if e.GetCode() != zglobal.Code_Zauth_TokenIsNil {

			ctx.LogError(e)
			return &zauth_pb.CheckAuth_RES{
				Code: e.GetCode(),
			}
		} else {
			tk, e = CreateToken(ctx)
			if e != nil {
				ctx.LogError(e)
				return &zauth_pb.CheckAuth_RES{
					Code: e.GetCode(),
				}
			}
		}
	}

	// 检查 token 合法性
	if tk.Sign != ctx.AuthSign {
		ctx.LogError(zservice.NewError("no token:", ctx.AuthToken).SetCode(zglobal.Code_Zauth_TokenIsNil))
		return &zauth_pb.CheckAuth_RES{
			Code: zglobal.Code_Zauth_TokenSignFail,
		}
	}

	// 检查权限
	// 查询用户所有的组
	if e := Mysql.Raw(`
	WITH RECURSIVE cte(id) AS (
		SELECT g_id FROM account_group_bind_tables WHERE uid=? 
		UNION ALL SELECT 
		agt.g_id FROM cte JOIN account_group_tables agt ON cte.id = agt.id
	) SELECT DISTINCT id FROM cte WHERE id > 0;
	`, 1001).Find(&[]struct{}{}).Error; e != nil {
		zservice.LogError(e)
	}

	return &zauth_pb.CheckAuth_RES{
		Code: zglobal.Code_SUCC,
	}
}
