package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"blue/config"
	"blue/internal"
	"blue/log"
)

var title = `
  _       _                
 | |__   | |  _   _    ___ 
 | '_ \  | | | | | |  / _ \
 | |_) | | | | |_| | | |__/
 |_.__/  |_|  \__,_|  \___|
                           `

var confPath = flag.String("c", "./blue-server.json", "config file path")

func init() {
	fmt.Printf("\033[34m%s\033[0m\n", title)
	fmt.Println(os.Getpid(), os.Getuid())
}

func main() {
	flag.Parse()

	configDB := config.InitConfig(*confPath)
	log.InitLog()

	dbs := make([]*internal.DB, config.SvrCfg.DBSum+1)
	dbs[0] = internal.NewDB(func(c *internal.DBConfig) {
		c.SetStorage = false
		c.DataDictSize = 1024
		c.Index = 0
		c.InitData = configDB
	})

	for i := 1; i <= config.SvrCfg.DBSum; i++ {
		dbs[i] = internal.NewDB(func(c *internal.DBConfig) {
			c.SetStorage = false
			c.DataDictSize = 1024
			c.Index = i
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
