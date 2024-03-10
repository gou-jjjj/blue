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
		c.Ip = BC.Ip
		c.Port = BC.Port
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
			ErrPrint(ErrRead(err.Error()))
			continue
		}

		split := TidyInput(input)

		if len(split) == 0 {
			continue
		}

	Resend:
		res, err := Exec(split)
		if err != nil {
			Connect()
			goto Resend
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

/*
Version() (string, error)
Select(...string) (string, error)
Del(string) (string, error)
Nset(string, string) (string, error)
Get(string) (string, error)
Set(string, string) (string, error)
Len(string) (string, error)
Kvs() (string, error)
Nget(string) (string, error)
*/
func Exec(s []string) (string, error) {
	if len(s) == 0 {
		return "", ErrCommand(s[0])
	}

	switch s[0] {
	case "set":
		if len(s) != 3 {
			return "", ErrArgu(s[0])
		}
		return conn.Set(s[1], s[2])
	case "del":
		if len(s) != 2 {
			return "", ErrArgu(s[0])
		}
		return conn.Del(s[1])
	case "version":
		if len(s) != 1 {
			return "", ErrArgu(s[0])
		}
		return conn.Version()
	case "get":
		if len(s) != 2 {
			return "", ErrArgu(s[0])
		}
		return conn.Get(s[1])

	default:
		return "", ErrCommand(s[0])
	}
}
