package internal

import (
	"errors"
	"fmt"
	"strings"
	"time"
	"zservice/service/zauth/zauth_pb"
	"zservice/zservice"
	"zservice/zservice/ex/gormservice"
	"zservice/zservice/zglobal"

	"gorm.io/gorm"
)

// 检查权限
func Logic_CheckAuth(ctx *zservice.Context, in *zauth_pb.CheckAuth_REQ) *zauth_pb.CheckAuth_RES {

	// 获取和检查 token
	// 获取 token
	at, e := GetToken(ctx, ctx.AuthToken)
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
		return &zauth_pb.CheckAuth_RES{Code: zglobal.Code_Zauth_TokenSignFail, IsTokenRefresh: isRefreshToken, Token: at.Token, Uid: at.UID}
	}

	// 权限相关参数列表
	authArr := zservice.JsonMustUnmarshal_StringArray([]byte(in.Auth))
	if len(authArr) != 3 {
		return &zauth_pb.CheckAuth_RES{Code: zglobal.Code_ParamsErr, IsTokenRefresh: isRefreshToken, Token: at.Token, Uid: at.UID}
	}

	authService := authArr[0]
	authAction := authArr[1]
	authPath := authArr[2]

	// 获取与指定参数最接近的权限对象
	permissionInfo, e := func() (*PermissionTable, *zservice.Error) {
		var service = authService
		var action = authAction
		var path = authPath
		actionArr := []string{action}
		if action != "" {
			actionArr = append(actionArr, "")
		}
		// 所有路径
		inArr := [][]string{}
		tmpPath := path
		for {
			lastIndex := strings.LastIndex(tmpPath, "/")
			if lastIndex == -1 {
				for _, v := range actionArr {
					inArr = append(inArr, []string{service, v, ""}) // 无 action
				}
				break // 已经到达路径根部，无需再查询
			}
			for _, v := range actionArr {
				inArr = append(inArr, []string{service, v, tmpPath}) // action
			}

			tmpPath = tmpPath[:lastIndex] // 获取父级路径
		}

		// 未找到 查表, 按权限最接近的查询
		tabs := []PermissionTable{}
		if e := Mysql.Model(&PermissionTable{}).Where("(service, action, path) IN ?", inArr).Order("LENGTH(action) DESC, LENGTH(path) DESC").Find(&tabs).Error; e != nil {
			if !errors.Is(e, gorm.ErrRecordNotFound) {
				return nil, zservice.NewError(e)
			}
		}

		if len(tabs) == 0 {
			return nil, zservice.NewError("not found").SetCode(zglobal.Code_NotFound)
		}

		for _, tab := range tabs {
			if tab.State != 3 { // 继承父级，向上查询
				return &tab, nil
			}
		}

		return &tabs[0], nil
	}()

	if e != nil {
		ctx.LogError(e)
		return &zauth_pb.CheckAuth_RES{Code: e.GetCode(), IsTokenRefresh: isRefreshToken, Token: at.Token, Uid: at.UID}
	}

	// 当前权限是否公开
	switch permissionInfo.State {
	case 0: // 权限禁用
	case 3: // 继承父级，父级未处理，权限配置有问题，查当前服务的顶级权限是否配置正确
		return &zauth_pb.CheckAuth_RES{Code: zglobal.Code_Zauth_Permission_ConfigErr, IsTokenRefresh: isRefreshToken, Token: at.Token, Uid: at.UID}
	case 1: // 公开访问
		return &zauth_pb.CheckAuth_RES{Code: zglobal.Code_SUCC, IsTokenRefresh: isRefreshToken, Token: at.Token, Uid: at.UID}
	}

	// 检查是否拥有该权限
	if at.UID == 0 { // 未登录, 不继续接下里用户判断流程
		return &zauth_pb.CheckAuth_RES{Code: zglobal.Code_Zauth_Fail, IsTokenRefresh: isRefreshToken, Token: at.Token, Uid: at.UID}
	}

	// 检查登陆服务是否正确
	if at.LoginService != authService {
		return &zauth_pb.CheckAuth_RES{Code: zglobal.Code_Zauth_Fail, IsTokenRefresh: isRefreshToken, Token: at.Token, Uid: at.UID}
	}

	// 服务登陆和token验证
	if s, e := Redis.Get(fmt.Sprintf(RK_UserLoginService, at.UID, authService)).Result(); e != nil {
		ctx.LogError(e)
		return &zauth_pb.CheckAuth_RES{Code: zglobal.Code_Zauth_Fail, IsTokenRefresh: isRefreshToken, Token: at.Token, Uid: at.UID}
	} else if s != at.Token { // token 不正确, 需要重新登陆
		return &zauth_pb.CheckAuth_RES{Code: zglobal.Code_LoginAgain, IsTokenRefresh: isRefreshToken, Token: at.Token, Uid: at.UID}
	}

	// 检查用户和权限组的绑定
	// 当前账号是否有权限配置
	if tab, e := GetPermissionBind(ctx, 2, at.UID, permissionInfo.PermissionID); e != nil && e.GetCode() != zglobal.Code_NotFound {
		ctx.LogError(e)
		return &zauth_pb.CheckAuth_RES{Code: e.GetCode(), IsTokenRefresh: isRefreshToken, Token: at.Token, Uid: at.UID}
	} else if tab != nil && !tab.IsExpired() { // 过期的检查权限表示无效，检查所在组织是否有权限
		if tab.State == 1 {
			return &zauth_pb.CheckAuth_RES{Code: zglobal.Code_SUCC, IsTokenRefresh: isRefreshToken, Token: at.Token, Uid: at.UID}
		} else {
			return &zauth_pb.CheckAuth_RES{Code: zglobal.Code_Fail, IsTokenRefresh: isRefreshToken, Token: at.Token, Uid: at.UID}
		}
	}

	// 查库
	bindCount := int64(0)
	if e := Mysql.Model(&UserOrgBindTable{}).Where( // 查找组中是否有当前账号的绑定信息
		"uid = ? AND org_id IN (?)",
		at.UID,
		Mysql.Model(&PermissionBindTable{}).Where( // 查找所有分配权限的组
			"permission_id = ? AND target_type = 1 AND state = 1 AND (expires = 0 OR expires > ?)",
			permissionInfo.PermissionID,
			time.Now().Unix(),
		).Select("target_id"),
	).Count(&bindCount).Error; e != nil {
		ctx.LogError(e)
		if gormservice.IsNotFound(e) {
			return &zauth_pb.CheckAuth_RES{Code: zglobal.Code_NotFound, IsTokenRefresh: isRefreshToken, Token: at.Token, Uid: at.UID}
		}
		return &zauth_pb.CheckAuth_RES{Code: zglobal.Code_Fail, IsTokenRefresh: isRefreshToken, Token: at.Token, Uid: at.UID}
	}

	if bindCount == 0 {
		return &zauth_pb.CheckAuth_RES{Code: zglobal.Code_NotFound, IsTokenRefresh: isRefreshToken, Token: at.Token, Uid: at.UID}
	} else {
		return &zauth_pb.CheckAuth_RES{Code: zglobal.Code_SUCC, IsTokenRefresh: isRefreshToken, Token: at.Token, Uid: at.UID}
	}
}
