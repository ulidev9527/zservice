package zservice

// 内置状态码
const (
	Code_Zero         = 0  // 如果出现 0 表示业务未处理状态码
	Code_Succ         = 1  // 成功
	Code_Fail         = 2  // 失败
	Code_Limit        = 3  // 限制/上限
	Code_Auth         = 4  // 鉴权失败/无权
	Code_NotImplement = 5  // 未实现
	Code_Params       = 6  // 参数错误
	Code_NotFound     = 7  // 资源不存在/没找到/数据未查询到
	Code_Again        = 8  // 等待，重试
	Code_Repetition   = 9  // 数据重复, 数据已存在，数据相同
	Code_Reject       = 10 // 拒绝
	Code_Fatal        = 11 // 代码执行阻断执行错误，严重错误，服务断开
)
