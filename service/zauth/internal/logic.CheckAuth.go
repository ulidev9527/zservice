package internal

import (
	"strings"
	"time"
	"zservice/service/zauth/zauth_pb"
	"zservice/zservice"
	"zservice/zservice/ex/gormservice"
	"zservice/zservice/zglobal"
)

// 检查权限
func Logic_CheckAuth(ctx *zservice.Context, in *zauth_pb.CheckAuth_REQ) *zauth_pb.CheckAuth_RES {

	// 获取和检查 token
	resultRES := &zauth_pb.CheckAuth_RES{Code: zglobal.Code_Zauth_Fail}
	// 获取 token
	authToken := &AuthToken{}
	if at, e := GetToken(ctx, in.Token); e != nil {
		if e.GetCode() != zglobal.Code_NotFound { // 其它错误
			ctx.LogError(e)
			resultRES.Code = e.GetCode()
			return resultRES
		} else { // 未找到进行创建
			if at, e = CreateToken(ctx, in.TokenSign); e != nil {
				ctx.LogError(e)
				resultRES.Code = e.GetCode()
				return resultRES
			} else {
				resultRES.Token = at.Token
				in.Token = at.Token
				authToken = at
			}
		}
	} else {
		if at.TokenCheck(in.TokenSign) { // 检查 token 是否正确
			resultRES.Token = at.Token
			resultRES.Uid = at.UID
			authToken = at
		} else {
			ctx.LogError("token check fail", in.Token)
			return resultRES
		}
	}

	// 获取与指定参数最接近的权限对象
	permissionInfo := &PermissionTable{}
	if pInfo, e := func() (*PermissionTable, *zservice.Error) {
		var service = in.Service
		var action = in.Action
		var path = in.Path
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
		if e := Mysql.
			Order("LENGTH(action) DESC, LENGTH(path) DESC").
			Find(&tabs, "(service, action, path) IN ? AND state > 0", inArr).
			Error; e != nil {
			if !gormservice.IsNotFound(e) {
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
	}(); e != nil {
		ctx.LogError(e)
		resultRES.Code = e.GetCode()
		return resultRES
	} else {
		permissionInfo = pInfo
	}

	// 当前权限是否公开
	switch permissionInfo.State {
	case zglobal.E_PermissionState_IgnoreAll: // 拒绝所有访问
		ctx.LogWarn("permission is:", zglobal.E_PermissionState_IgnoreAll, in)
		resultRES.Code = zglobal.Code_Zauth_Fail
	case 1: // 公开访问
		resultRES.Code = zglobal.Code_SUCC
		return resultRES
	case 2: // 需要登录
		// 检查是否拥有该权限
		if authToken.UID == 0 { // 未登录, 不继续接下里用户判断流程
			ctx.LogInfo("not login", in.Service, in.Token, in.TokenSign)
			return resultRES
		}

		// 检查登陆服务是否正确
		if !authToken.HasLoginService(in.Service) {
			ctx.LogInfo("login service error", in.Service, in.Token, in.TokenSign)
			return resultRES
		}

		// 检查用户和权限组的绑定
		// 当前账号是否有权限配置
		if tab, e := GetPermissionBind(ctx, 2, authToken.UID, permissionInfo.PermissionID); e != nil {
			if e.GetCode() != zglobal.Code_NotFound { // 未找到进行父级查找，其它错误直接在这里返回
				ctx.LogError(e)
				resultRES.Code = e.GetCode()
				return resultRES
			}
		} else if !tab.IsExpired() { // 过期的检查权限表示无效，检查所在组织是否有权限
			if tab.State == 1 {
				resultRES.Code = zglobal.Code_SUCC
				return resultRES
			} else {
				resultRES.Code = zglobal.Code_Limit
				return resultRES
			}
		}

		// 查库
		bindCount := int64(0)
		if e := Mysql.Model(&UserOrgBindTable{}).Where( // 查找组中是否有当前账号的绑定信息
			"uid = ? AND org_id IN (?)",
			authToken.UID,
			Mysql.Model(&PermissionBindTable{}).Where( // 查找所有分配权限的组
				"permission_id = ? AND target_type = 1 AND state = 1 AND (expires IS NULL OR expires > ?)",
				permissionInfo.PermissionID,
				time.Now(),
			).Select("target_id"),
		).Count(&bindCount).Error; e != nil {
			ctx.LogError(e)
			return resultRES
		}

		if bindCount > 0 {
			resultRES.Code = zglobal.Code_SUCC
		}
	default:
		// 其它状态肯定是配置出错或者上面逻辑有问题
		ctx.LogError("permission config error", permissionInfo.State, in)
		return resultRES
	}

	return resultRES
}
