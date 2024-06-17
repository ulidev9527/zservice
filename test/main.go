package main

import (
	"fmt"
	"net"
	"os"
	"zservice/zservice"
)

func init() {
	zservice.Init("test", "1.0.0")
}

type TT struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	NickName string `json:"nickName"`
	BT       []byte `json:"bt"`
}

func main() {
	addrs, err := net.InterfaceAddrs()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, address := range addrs {

		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				fmt.Println(ipnet.IP.String())
			}

		}
	}

}
