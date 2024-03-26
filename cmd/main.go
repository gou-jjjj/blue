package main

import (
	"fmt"
	"os"
	"time"

	"blue/config"
	"blue/internal"
	"blue/log"
)

var title = `  _       _                
 | |__   | |  _   _    ___ 
 | '_ \  | | | | | |  / _ \
 | |_) | | | | |_| | | |__/
 |_.__/  |_|  \__,_|  \___|
                           `

func init() {
	fmt.Printf("\033[34m%s\033[0m\n", title)
	fmt.Println(os.Getpid(), os.Getuid())
}

func main() {
	configDB := config.InitConfig()
	log.InitLog()

	dbs := make([]*internal.DB, config.BC.ServerConfig.DBSum+1)
	dbs[0] = internal.NewDB(func(c *internal.DBConfig) {
		c.SetStorage = false
		c.DataDictSize = 1024
		c.Index = 0
		c.InitData = configDB
	})

	for i := 1; i <= config.BC.ServerConfig.DBSum; i++ {
		dbs[i] = internal.NewDB(func(c *internal.DBConfig) {
			c.SetStorage = false
			c.DataDictSize = 1024
			c.Index = i
		})
	}

	handler := internal.NewBlueServer(dbs...)

	internal.NewServer(
		func(c *internal.Config) {
			c.Ip = config.BC.ServerConfig.Ip
			c.Port = config.BC.ServerConfig.Port
			c.ClientLimit = config.BC.ClientConfig.ClientLimit
			c.Timeout = time.Duration(config.BC.ServerConfig.TimeOut) * time.Second
			c.HandlerFunc = handler
		},
	).Start()
}
