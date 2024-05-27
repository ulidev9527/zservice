package internal

import (
	"fmt"
	"strings"
	"zservice/service/zauth/zauth_pb"
	"zservice/zservice"
	"zservice/zservice/zglobal"
)

// 检查权限
func Logic_CheckAuth(ctx *zservice.Context, in *zauth_pb.CheckAuth_REQ) *zauth_pb.CheckAuth_RES {

	// 获取和检查 token
	// 获取 token
	at, e := GetToken(ctx.AuthToken)
	isRefreshToken := false
	if e != nil {
		if e.GetCode() != zglobal.Code_Zauth_TokenIsNil {
			ctx.LogError(e)
			return &zauth_pb.CheckAuth_RES{Code: e.GetCode()}
		} else {

			// 没有 token 创建
			at, e = CreateToken(ctx)
			if e != nil {
				ctx.LogError(e)
				return &zauth_pb.CheckAuth_RES{Code: e.GetCode()}
			}
			ctx.AuthToken = at.Token
			isRefreshToken = true
		}
	}

	// 检查 token 正确性
	if !at.CheckToken(ctx.AuthToken, ctx.AuthSign) {
		ctx.LogError("token check fail", ctx.AuthToken)
		return &zauth_pb.CheckAuth_RES{Code: zglobal.Code_Zauth_TokenSignFail, IsTokenRefresh: isRefreshToken, Token: at.Token}
	}

	// 权限相关参数列表
	authArr := zservice.JsonMustUnmarshal_StringArray([]byte(in.Auth))
	if len(authArr) == 0 {
		return &zauth_pb.CheckAuth_RES{Code: zglobal.Code_ParamsErr, IsTokenRefresh: isRefreshToken, Token: at.Token}
	}

	authService := authArr[0]
	authAction := authArr[1]
	authPath := authArr[2]

	// 获取权限信息
	permissionInfo, e := func() (*ZauthPermissionTable, *zservice.Error) {
		var service = authService
		var action = authAction
		var path = authPath
		var tab *ZauthPermissionTable
		var e *zservice.Error

		// 尝试直接查询指定路径、动作和服务的权限
		tab, e = GetPermissionBySAP(ctx, service, action, path)
		if e != nil && e.GetCode() != zglobal.Code_DB_NotFound {
			return nil, e
		}
		if tab != nil && tab.State != 3 { // 3 需要查询父级
			return tab, nil
		}

		// 如果没有匹配到指定路径、动作和服务的权限，则尝试查询动作为空字符串的权限
		tab, e = GetPermissionBySAP(ctx, service, "", path)
		if e != nil && e.GetCode() != zglobal.Code_DB_NotFound {
			return nil, e
		}
		if tab != nil && tab.State != 3 { // 3 需要查询父级
			return tab, nil
		}

		// 如果没有匹配到指定路径和服务的权限，则尝试查询父级路径的权限
		for {
			lastIndex := strings.LastIndex(path, "/")
			if lastIndex == -1 {
				break // 已经到达路径根部，无需再查询
			}
			path = path[:lastIndex] // 获取父级路径

			tab, e = GetPermissionBySAP(ctx, service, "", path)
			if e != nil && e.GetCode() != zglobal.Code_DB_NotFound {
				return nil, e
			}
			if tab != nil && tab.State != 3 { // 3 需要查询父级
				return tab, nil
			}
		}
		return tab, nil
	}()

	if e != nil {
		ctx.LogError(e)
		return &zauth_pb.CheckAuth_RES{Code: e.GetCode(), IsTokenRefresh: isRefreshToken, Token: at.Token}
	}

	// 当前权限是否公开
	switch permissionInfo.State {
	case 0: // 权限禁用
	case 3: // 继承父级，父级未处理，权限配置有问题
		return &zauth_pb.CheckAuth_RES{Code: zglobal.Code_Zauth_Permission_ConfigErr, IsTokenRefresh: isRefreshToken, Token: at.Token}
	case 1: // 公开访问
		return &zauth_pb.CheckAuth_RES{Code: zglobal.Code_SUCC, IsTokenRefresh: isRefreshToken, Token: at.Token}
	}

	// 检查是否拥有该权限
	if at.UID == 0 { // 未登录, 不继续接下里用户判断流程
		return &zauth_pb.CheckAuth_RES{Code: zglobal.Code_AuthFail, IsTokenRefresh: isRefreshToken, Token: at.Token}
	}

	// 服务登陆和token验证
	if s, e := Redis.Get(fmt.Sprintf(RK_AccountLoginService, at.UID, authService)).Result(); e != nil {
		ctx.LogError(e)
		return &zauth_pb.CheckAuth_RES{Code: zglobal.Code_AuthFail, IsTokenRefresh: isRefreshToken, Token: at.Token}
	} else if s != at.Token { // token 不正确, 需要重新登陆
		return &zauth_pb.CheckAuth_RES{Code: zglobal.Code_LoginAgain, IsTokenRefresh: isRefreshToken, Token: at.Token}
	}

	// 检查是否有权限
	isAllow, e := func() (bool, *zservice.Error) {
		// 当前账号是否有权限配置
		if tab, e := GetPermissionBind(ctx, 2, at.UID, permissionInfo.PermissionID); e != nil && e.GetCode() != zglobal.Code_DB_NotFound {
			return false, e
		} else if tab != nil && tab.IsExpired() { // 过期的检查权限表示无效，检查所在组织是否有权限
			return tab.Allow, nil
		}

		// 查找所有分配权限的组
		bindInfo := &ZauthAccountOrgBindTable{}

		if e := Mysql.Model(&ZauthAccountOrgBindTable{}).Where(
			"account_id = ? AND org_id IN (?)",
			at.UID,
			Mysql.Model(&ZauthPermissionBindTable{}).Where(
				"permission_id = ? AND target_type = 1 AND state = 1 AND (is_expired IS NULL OR is_expired > NOW())",
				permissionInfo.PermissionID,
			).Select("target_id"),
		).First(bindInfo).Error; e != nil {
			return false, zservice.NewError(e).SetCode(zglobal.Code_ErrorBreakoff)
		}
		return bindInfo.ID > 0, nil
	}()

	if e != nil {
		ctx.LogError(e)
		return &zauth_pb.CheckAuth_RES{Code: e.GetCode(), IsTokenRefresh: isRefreshToken, Token: at.Token}
	}
	if isAllow {
		return &zauth_pb.CheckAuth_RES{Code: zglobal.Code_SUCC, IsTokenRefresh: isRefreshToken, Token: at.Token}
	} else {
		return &zauth_pb.CheckAuth_RES{Code: zglobal.Code_AuthFail, IsTokenRefresh: isRefreshToken, Token: at.Token}
	}

	// // 检查权限
	// // 查询用户所有的组
	// if e := Mysql.Raw(`
	// WITH RECURSIVE cte(id) AS (
	// 	SELECT g_id FROM account_group_bind_tables WHERE uid=?
	// 	UNION ALL SELECT
	// 	agt.g_id FROM cte JOIN account_group_tables agt ON cte.id = agt.id
	// ) SELECT DISTINCT id FROM cte WHERE id > 0;
	// `, 1001).Find(&[]struct{}{}).Error; e != nil {
	// 	zservice.LogError(e)
	// }

	// if e := at.Save(); e != nil {
	// 	ctx.LogError(e)
	// 	return &zauth_pb.CheckAuth_RES{Code: e.GetCode()}
	// }

	// return &zauth_pb.CheckAuth_RES{
	// 	Code:           zglobal.Code_SUCC,
	// 	IsTokenRefresh: ctx.AuthToken != at.Token,
	// 	Token:          at.Token,
	// }
}
