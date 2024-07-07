package zservice

const (
	E_ConfigAsset_Parser_Excel = 1 // excel 解析器

	E_PermissionState_IgnoreAll  = 0 // 忽略全部访问
	E_PermissionState_AllowAll   = 1 // 允许所有访问
	E_PermissionState_AllowLogin = 2 // 允许登录访问
	E_PermissionState_Parent     = 3 // 继承父级状态

)
