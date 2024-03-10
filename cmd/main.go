package main

import (
	"blue/config"
	"blue/internal"
	"blue/log"
	"fmt"
	"time"
)

var title = `  _       _                
 | |__   | |  _   _    ___ 
 | '_ \  | | | | | |  / _ \
 | |_) | | | | |_| | |  __/
 |_.__/  |_|  \__,_|  \___|
                           `

func main() {
	fmt.Printf("\033[34m%s\033[0m\n", title)
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
