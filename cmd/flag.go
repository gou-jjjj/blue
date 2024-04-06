package main

import (
	"blue/common/network"
	"blue/config"
	"flag"
)

// 定义默认集群地址
var defCluster = "127.0.0.1:13141"

// 定义默认配置文件路径
var defConfPath = "./blue-server.json"

// 通过命令行参数指定配置文件路径
var confPath = flag.String("c", defConfPath, "config file path")

// 通过命令行参数指定集群地址
var clusterPath = flag.String("p", defCluster, "cluster path")

/**
 * 获取集群地址
 * 若命令行指定了集群地址，则使用命令行指定的地址；
 * 否则，从配置文件中读取集群地址；
 * 若两者都未指定或指定的地址无效，则返回空字符串。
 * @return string 集群地址
 */
func clusterAddr() string {
	// 检查命令行指定的集群地址是否有效
	if *clusterPath != defCluster {
		if !network.ParseAddr(*clusterPath) {
			panic("cluster path is not valid")
		}
		return *clusterPath
	}

	// 检查配置文件中的集群地址是否有效
	if config.CluCfg.ClusterAddr != "" {
		if !network.ParseAddr(config.CluCfg.ClusterAddr) {
			panic("cluster path is not valid")
		}
		return config.CluCfg.ClusterAddr
	}

	return ""
}
