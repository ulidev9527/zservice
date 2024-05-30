package internal

import (
	"zservice/service/zauth/zauth_pb"
	"zservice/zservice"
	"zservice/zservice/zglobal"
)

// 组织创建，只需要 *name 和 perentOrgID
func Logic_OrgCreate(ctx *zservice.Context, in *zauth_pb.OrgInfo) *zauth_pb.OrgInfo_RES {

	if in.Name == "" {
		return &zauth_pb.OrgInfo_RES{Code: zglobal.Code_ParamsErr}
	}

	rootOrgID := uint32(0)
	parentOrgID := in.ParentOrgID

	if parentOrgID != 0 { // 非根组织
		// 验证组织是否存在
		if tab, e := GetOrgByID(ctx, uint(in.ParentOrgID)); e != nil {
			ctx.LogError(e)
			return &zauth_pb.OrgInfo_RES{Code: e.GetCode()}
		} else if tab == nil {
			return &zauth_pb.OrgInfo_RES{Code: zglobal.Code_Zauth_Org_NotFund}
		} else {
			// 配置根组织
			if tab.RootOrgID == 0 {
				rootOrgID = tab.OrgID
			} else {
				rootOrgID = tab.RootOrgID
			}
		}
	} else { // 根组织
		// 验证组织是否存在
		if tab, e := GetRootOrgByName(ctx, in.Name); e != nil {
			ctx.LogError(e)
			return &zauth_pb.OrgInfo_RES{Code: e.GetCode()}
		} else if tab != nil {
			return &zauth_pb.OrgInfo_RES{Code: zglobal.Code_Zauth_Org_AlreadyExist}
		}
	}

	// 获取一个未使用的组织 ID
	orgID, e := GetNewOrgID(ctx)
	if e != nil {
		ctx.LogError(e)
		return &zauth_pb.OrgInfo_RES{Code: e.GetCode()}
	}
	z := &ZauthOrgTable{
		Name:        in.Name,
		OrgID:       orgID,
		RootOrgID:   rootOrgID,
		ParentOrgID: parentOrgID,
		State:       in.State,
	}

	if e := z.Save(ctx); e != nil {
		ctx.LogError(e)
		return &zauth_pb.OrgInfo_RES{Code: e.GetCode()}
	}
	return &zauth_pb.OrgInfo_RES{Code: zglobal.Code_SUCC, Info: &zauth_pb.OrgInfo{
		OrgID:       uint32(z.OrgID),
		Name:        z.Name,
		ParentOrgID: uint32(z.ParentOrgID),
		RootOrgID:   uint32(z.RootOrgID),
		State:       uint32(z.State),
	}}

}
