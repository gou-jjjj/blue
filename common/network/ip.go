package network

import (
	"fmt"
	"net"
	"strings"
)

// LocalIpEn0 用于获取Mac电脑上en0或eth0接口的IPv4地址。
// 返回值为字符串类型的IPv4地址，如果无法获取则返回空字符串。
func LocalIpEn0() string {
	// 获取本机所有网络接口信息
	interfaces, err := net.Interfaces()
	if err != nil {
		panic(err)
	}

	// 遍历接口信息，查找名称为en0或eth0的接口
	for _, iface := range interfaces {
		if iface.Name == "en0" || iface.Name == "eth0" {
			// 获取接口的所有地址
			addrs, err := iface.Addrs()
			if err != nil {
				panic(err)
			}
			// 遍历地址，提取IPv4地址
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

// ParseAddr 用于解析地址字符串格式，判断其是否为有效的IP和端口组合。
// 参数addr为要解析的地址字符串。
// 返回值为布尔类型，表示地址是否有效。
func ParseAddr(addr string) bool {

	if addr == "" {
		return false
	}

	// 分割地址字符串为IP和端口
	addrs := strings.Split(addr, ":")

	// 判断是否包含IP和端口两部分
	if len(addrs) != 2 {
		return false
	}

	// 验证IP地址的有效性
	ip := net.ParseIP(addrs[0])
	if ip == nil {
		return false
	}

	// 验证端口的有效性
	_, err := net.ResolveTCPAddr("tcp", ":"+addrs[1])
	if err != nil {
		return false
	}
	return true
}

const space = "|"

// CombineAddr 用于将集群地址和客户端地址合并为一个字符串，中间以指定字符分隔。
// 参数cluAddr为集群地址，cliAddr为客户端地址。
// 返回值为合并后的地址字符串。
func CombineAddr(cluAddr, cliAddr string) string {
	return fmt.Sprintf("%s%s%s", cluAddr, space, cliAddr)
}

// SplitAddr 用于根据分隔符将地址字符串拆分为多个部分。
// 如果地址字符串以特定字符开头和结尾，会先去除这些字符，再进行拆分。
// 参数addr为要拆分的地址字符串。
// 返回值为拆分后的地址字符串切片。
func SplitAddr(addr string) []string {
	if addr[0] == '+' && addr[len(addr)-1] == '\n' {
		return strings.Split(addr[1:len(addr)-1], space)
	}

	return strings.Split(addr, space)
}
