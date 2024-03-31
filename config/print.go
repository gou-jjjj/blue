package config

import (
	"fmt"
	"log"
)

var title = `
  _       _                
 | |__   | |  _   _    ___ 
 | '_ \  | | | | | |  / _ \
 | |_) | | | | |_| | | |__/
 |_.__/  |_|  \__,_|  \___|
                           `

func PrintTitle() {
	fmt.Printf("\033[34m%s\033[0m\n", title)
}

func LogInitSuccess() {
	log.Println("log init success ...")
}

func ConfigInitSuccess() {
	log.Println("config init success ...")
}

func ClusterInitSuccess() {
	log.Println("cluster init success ...")
}

func StorageInitSuccess(dbidx int) {
	log.Printf("storage init success [%d]...", dbidx)
}

func ServerInitSuccess() {
	log.Println("Server init success ...")
}

func ErrPanic(err error, data ...string) {
	if err != nil {
		if len(data) > 0 {
			panic(fmt.Sprintf("%s: %s", err.Error(), data[0]))
		}
		panic(err)
	}
}
