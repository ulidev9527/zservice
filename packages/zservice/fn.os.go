package zservice

import (
	"fmt"
	"net"
	"os"
	"sync"
	"time"

	"github.com/panjf2000/ants/v2"
)

const (
	GoPool_Max = 100000 // 最大协程数 10W
	GoPool_Min = 1000   // 最小协程数 1K
)

var (
	goPool, _             = ants.NewPool(GoPool_Min)
	goPool_mu             sync.Mutex
	goPool_maxErrOutputCD = 0
)

//	func Go(f func()) {
//		go f()
//	}
func Go(f func()) {
	goPool_mu.Lock()
	defer goPool_mu.Unlock()
	if e := goPool.Submit(f); e != nil {
		LogErrorCaller(2, "zserbice.GO :", e)
	}
	if goPool.Free() < 100 {

		if goPool.Cap() < GoPool_Max {
			goPool.Tune(goPool.Cap() + GoPool_Min)
			LogWarnCallerf(2, "zserbice.GO TUNE:%d", goPool.Cap())
		} else if goPool_maxErrOutputCD <= 0 {
			LogErrorCaller(2, "zserbice.GO TUNE MAX", goPool.Cap())
			goPool_maxErrOutputCD = GoPool_Max
		}
		goPool_maxErrOutputCD--
	}
}

// 写入文件到临时目录
func WriteFileToTempDir(name string, data []byte) *Error {
	if e := os.WriteFile(fmt.Sprintf("%s/%s", os.TempDir(), name), data, 0644); e != nil {
		return NewError(e).SetCode(Code_Fatal)
	}
	return nil
}

// 获取本机IP地址
func GetIp() ([]string, *Error) {

	addrs, e := net.InterfaceAddrs()
	if e != nil {
		return nil, NewError(e)
	}

	ips := make([]string, 0)
	for _, address := range addrs {

		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ips = append(ips, ipnet.IP.String())
			}
		}
	}
	return ips, nil
}

// 获取一个未使用的端口
func GetFreePort() int {
	loopCount := 10
	for {
		if loopCount < 10 {
			LogError("GetFreePort fail")
			os.Exit(1)
			return 0
		}
		loopCount--

		addr, e := net.ResolveTCPAddr("tcp", "127.0.0.1:0")
		if e != nil {
			LogError(e)
			continue
		}
		cli, e := net.ListenTCP("tcp", addr)
		if e != nil {
			LogError(e)
			continue
		}
		defer cli.Close()
		return cli.Addr().(*net.TCPAddr).Port
	}
}

// 检查断开是否能够打开

func IsPortOpen(host string, port int, timeout time.Duration) bool {
	target := net.JoinHostPort(host, fmt.Sprintf("%d", port))
	conn, err := net.DialTimeout("tcp", target, timeout)
	if err != nil {
		return false
	}
	defer conn.Close()
	return true
}

// 获取本机名
func GetHostname() string {
	if n, e := os.Hostname(); e != nil {
		LogError(e)
		return ""
	} else {
		return n
	}
}

// 是否运行在 docker 中
func IsRunInDocker() bool {
	// 方法1: 检查/.dockerenv
	if _, err := os.Stat("/.dockerenv"); err == nil {
		return true
	}
	return false
}
