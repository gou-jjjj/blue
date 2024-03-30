package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	g "blue/api/go"
)

var (
	conn *g.Client
)

func Connect() (*g.Client, error) {
	return g.NewClient(func(c *g.Config) {
		c.Addr = BC.Addr
		c.TimeOut = time.Duration(BC.TimeOut) * time.Second
		c.TryTimes = BC.TryTimes
		c.DB = BC.DB
	})
}

var defaultConf = `./blue-cli.conf`
var defaultAddr = `127.0.0.1:13140`
var addr = flag.String("a", defaultAddr, "server address")
var conf = flag.String("c", defaultConf, "config file path")

func main() {
	flag.Parse()

	if defaultAddr != *addr {
		BC.Addr = *addr
	}
	var err error
	conn, err = Connect()
	if err != nil {
		panic(err)
	}
	// 从标准输入创建一个新的 bufio.Reader
	reader := bufio.NewReader(os.Stdin)

	for {
		remoteAddr := conn.RemoteAddr()
		fmt.Printf("%s> %s>", remoteAddr, BlueMessage(strconv.Itoa(BC.DB)))
		// 读取直到遇到换行符
		input, err := reader.ReadString('\n')
		if err != nil {
			ErrPrint(ErrRead(err.Error()).Error())
			continue
		}

		split := TidyInput(input)

		if len(split) == 0 {
			continue
		}

		res, err := Exec(conn, split)
		if err != nil {
			if !strings.Contains(err.Error(), "broken pipe") {
				ErrPrint(err.Error())
				continue
			}

			conn, err = Connect()
			if err != nil {
				panic(err)
			}
			res, err = Exec(conn, split)
			if err != nil {
				ErrPrint(err.Error())
				os.Exit(0)
			}
		}

		if split[0] == "select" && len(split) == 2 && res == "ok" {
			BC.DB, _ = strconv.Atoi(split[1])
		}
		SuccessPrint(res)
	}
}

func TidyInput(input string) []string {
	input = strings.TrimSpace(input)
	split := strings.Fields(input)
	newSplit := make([]string, 0, len(split))

	for i := range split {
		split[i] = strings.ToLower(split[i])
		if split[i] != "" {
			newSplit = append(newSplit, split[i])
		}
	}
	return newSplit
}
