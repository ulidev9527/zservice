package zservice

import (
	"net"
)

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
