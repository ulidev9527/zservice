package zservice

const (
	NSQ_Topic_zlog_AddKV = "zlog_add_kv" // 添加 kv 日志

	EV_Config_ServiceFileConfigChange = "ConfigService_%s_FileConfigChange" // 服务文件配置变更通知

	FI_ServiceFileConfig = "static/fileConfig/%s/%s" // 服务配置文件路径
	FI_ServiceEnvFile    = "static/envConfig/%s.env" // 服务环境文件路径
	FI_UploadDir         = "static/upload/%s"        // 上传文件路径
)
