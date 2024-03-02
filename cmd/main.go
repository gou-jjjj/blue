package main

import (
	"blue/config"
	"blue/internal"
	"blue/log"
	"time"
)

func main() {
	configDB := config.InitConfig()
	log.InitSyncLog()

	dbs := make([]*internal.DB, config.BC.ServerConfig.DBSum+1)
	dbs[0] = configDB

	for i := 1; i <= config.BC.ServerConfig.DBSum; i++ {
		dbs[i] = internal.NewDB(func(c *internal.DBConfig) {
			c.SetStorage = false
			c.DataDictSize = 1024
			c.Index = i
		})
	}

	handler := internal.NewBlueServer(dbs...)

	server := internal.NewServer(
		func(c *internal.Config) {
			c.Ip = config.BC.ServerConfig.Ip
			c.Port = config.BC.ServerConfig.Port
			c.ClientLimit = config.BC.ClientConfig.ClientLimit
			c.Timeout = time.Duration(config.BC.ServerConfig.TimeOut) * time.Second
			c.HandlerFunc = handler
		},
	)

	server.Start()
}
