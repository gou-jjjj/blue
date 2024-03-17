package main

import (
	"bufio"
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

func init() {
	Connect()
}

func Connect() {
	conn = g.NewClient(func(c *g.Config) {
		c.Addr = BC.Addr
		c.TimeOut = time.Duration(BC.TimeOut) * time.Second
		c.TryTimes = BC.TryTimes
	})
}

// num get a
func main() {
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

			Connect()
			res, err = Exec(conn, split)
			if err != nil {
				ErrPrint(err.Error())
				os.Exit(0)
			}
		}

		SuccessPrint(res)
	}
}

func TidyInput(input string) []string {
	input = strings.TrimSpace(input)
	split := strings.Split(input, " ")
	newSplit := make([]string, 0, len(split))

	for i := range split {
		split[i] = strings.ToLower(split[i])
		if split[i] != "" {
			newSplit = append(newSplit, split[i])
		}
	}
	return newSplit
}
