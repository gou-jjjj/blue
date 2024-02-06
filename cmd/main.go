package main

import (
	"blue/config"
	"blue/internal"
	"blue/log"
	"time"
)

func main() {
	config.InitConfig()
	log.InitSyncLog()

	db := internal.NewDB(func(c *internal.DBConfig) {
		c.DirPath = config.BC.Storage.Path
	})

	handler := internal.NewBlueServer(db)

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
