package zservice

var (
	__OpenFlowLog       = false
	__IsWatchFlowLogEnv = false
)

// 流程日志
func FlowLog(name string, fn func(log func(vals ...any))) {

	if !__IsWatchFlowLogEnv {
		WatchEnvChange("zservice_flowlog_open", func(key, newVal, oldVal string) {
			__OpenFlowLog = StringToBoolean(newVal)
		})
	}

	if !__OpenFlowLog {
		return
	}

	fn(func(vals ...any) {

		vals = append([]any{"---FLowStart-", name, "---"}, vals...)
		vals = append(vals, "---FlowEnd---")
		LogWarnCaller(2, vals...)

	})
}
