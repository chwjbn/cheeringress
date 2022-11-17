package cheerlib

import (
	"net"
	"os"
	"os/user"
	"runtime"
	"strconv"
)

func OsProcessNo() string {
	if os.Getpid() > 0 {
		return strconv.Itoa(os.Getpid())
	}
	return ""
}

func OsHostName() string {

	xHostName, xErr := os.Hostname()

	if xErr == nil {
		return xHostName
	}
	return "unknown"
}

func OSName() string {
	return runtime.GOOS
}

func OSUserName() string {

	xUser, xErr := user.Current()
	if xErr == nil {
		return xUser.Username
	}
	return "unknown"
}

func OsAllIPV4() (ipv4s []string) {
	adders, err := net.InterfaceAddrs()
	if err != nil {
		return
	}

	for _, addr := range adders {
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				ipv4 := ipNet.IP.String()
				if ipv4 == "127.0.0.1" || ipv4 == "localhost" {
					continue
				}
				ipv4s = append(ipv4s, ipv4)
			}
		}
	}
	return
}

func OsIPV4() string {
	ipv4s := OsAllIPV4()
	if len(ipv4s) > 0 {
		return ipv4s[0]
	}
	return "unknown"
}
