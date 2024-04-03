package network

import (
	"net"
	"strings"
)

func LocalIpEn0() string {
	interfaces, err := net.Interfaces()
	if err != nil {
		panic(err)
	}

	for _, iface := range interfaces {
		if iface.Name == "en0" || iface.Name == "eth0" {
			addrs, err := iface.Addrs()
			if err != nil {
				panic(err)
			}
			for _, addr := range addrs {
				ip4 := strings.Split(addr.String(), "/")[0]
				if net.ParseIP(ip4).To4() != nil {
					return ip4
				}
			}
		}
	}

	return ""
}

func ParseAddr(addr string) bool {

	if addr == "" {
		return false
	}

	addrs := strings.Split(addr, ":")

	if len(addrs) != 2 {
		return false
	}

	ip := net.ParseIP(addrs[0])
	if ip == nil {
		return false
	}

	_, err := net.ResolveTCPAddr("tcp", ":"+addrs[1])
	if err != nil {
		return false
	}
	return true
}
