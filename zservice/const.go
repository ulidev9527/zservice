package zservice

const (
	S_S2S = "ZSERVICE-S2S" // 链路记录的 KEY
	S_C2S = "ZSERVICE-C2S" // 客户端请求的 KEY

	ENV_ZSERVICE_NAME            = "ZSERVICE_NAME"            // 服务名称
	ENV_ZSERVICE_VERSION         = "ZSERVICE_VERSION"         // 服务版本
	ENV_ZSERVICE_REMOTE_ENV_ADDR = "ZSERVICE_REMOTE_ENV_ADDR" // 远程环境变量地址
	ENV_ZSERVICE_REMOTE_ENV_AUTH = "ZSERVICE_REMOTE_ENV_AUTH" // 远程环境变量认证
	ENV_ZSERVICE_FILES_ENV       = "ZSERVICE_FILES_ENV"       // 环境变量文件
)
