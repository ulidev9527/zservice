package internal

const (
	RK_FileConfig = "zconfig:fileConfig:%s" // 文件配置缓存
	RK_FileMD5    = "zconfig:fileMD5:%s"    // 文件配置的 md5, 用于标识是否需要重新解析

	NSQ_FileConfig_Change = "zconfig_fileConfig_change" // 文件配置变更通知
)
