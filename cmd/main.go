package main

import (
	"blue/config"
	"blue/internal"
	"time"
)

func main() {
	handler := internal.NewHandler()

	server := internal.NewServer(
		func(c *internal.Config) {
			c.ClientLimit = config.BC.ClientLimit
			c.Timeout = time.Duration(config.BC.TimeOut) * time.Second
			c.HandlerFunc = handler
		},
	)

	server.Start()
}
