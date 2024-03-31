package main

import (
	"flag"
	"time"

	"blue/common/filename"
	"blue/config"
	"blue/internal"
	"blue/log"
)

var confPath = flag.String("c", "./blue-server.json", "config file path")

func init() {
	config.PrintTitle()
}

func main() {
	flag.Parse()

	configDB := config.InitConfig(*confPath)
	log.InitLog()

	dbs := make([]*internal.DB, config.SvrCfg.DBSum+1)
	dbs[0] = internal.NewDB(func(c *internal.DBConfig) {
		c.DataDictSize = 1024
		c.Index = 0
		c.InitData = configDB
	},
	)

	for i := 1; i <= config.SvrCfg.DBSum; i++ {
		dbs[i] = internal.NewDB(func(c *internal.DBConfig) {
			c.DataDictSize = 1024
			c.Index = i
			c.StorageOptions.DirPath = filename.StorageName(
				config.StoCfg.StoragePath, i)
		})
	}

	handler := internal.NewBlueServer(dbs...)

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
