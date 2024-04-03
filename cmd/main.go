package main

import (
	"flag"
	"time"

	"blue/common/filename"
	"blue/common/network"
	print2 "blue/common/print"
	"blue/config"
	"blue/internal"
	"blue/log"
)

var defCluster = "127.0.0.1:13141"
var defConfPath = "./blue-server.json"

var confPath = flag.String("c", defConfPath, "config file path")
var clusterPath = flag.String("p", defCluster, "cluster path")

func init() {
	print2.PrintTitle()
}

func clusterAddr() string {
	if *clusterPath != defCluster {
		if !network.ParseAddr(*clusterPath) {
			panic("cluster path is not valid")
		}
		return *clusterPath
	}

	if config.CluCfg.ClusterAddr != "" {
		if !network.ParseAddr(config.CluCfg.ClusterAddr) {
			panic("cluster path is not valid")
		}
		return config.CluCfg.ClusterAddr
	}

	return ""
}

func main() {
	flag.Parse()
	// init config
	configDB := config.InitConfig(*confPath)

	// init log
	log.InitLog(config.LogCfg.LogLevel, config.LogCfg.LogOut)

	dbs := make([]*internal.DB, config.SvrCfg.DBSum+1)
	// init db0
	dbs[0] = internal.NewDB(func(c *internal.DBConfig) {
		c.DataDictSize = 1024
		c.Index = 0
		c.InitData = configDB
	},
	)

	// init db1-8
	for i := 1; i <= config.SvrCfg.DBSum; i++ {
		dbs[i] = internal.NewDB(func(c *internal.DBConfig) {
			c.DataDictSize = 1024
			c.Index = i
			c.StorageOptions.DirPath = filename.StorageName(
				config.StoCfg.StoragePath, i)
		})
	}

	// init handler
	handler := internal.NewBlueServer(dbs...)

	// init server
	internal.NewServer(
		func(c *internal.Config) {
			c.Ip = config.SvrCfg.Ip
			c.Port = config.SvrCfg.Port
			c.ClientLimit = config.CliCfg.ClientLimit
			c.Timeout = time.Duration(config.SvrCfg.TimeOut) * time.Second
			c.HandlerFunc = handler
		},
	).Start()
}
