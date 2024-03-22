package addr

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
		if iface.Name == "en0" {
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
