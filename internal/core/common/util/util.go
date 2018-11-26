package util

import (
	"fmt"
	"net"
)

func LookupHostIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println("Error:", err)
	}

	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}

	return ""
}
