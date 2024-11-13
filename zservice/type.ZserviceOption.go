package zservice

// 初始化配置
type ZserviceOption struct {
	Name    string // 显示在日志中的名称
	Version string // 版本号
	Debug   bool   // 是否开启调试模式

	// 服务启动回调
	OnStart func(*ZService)
}
