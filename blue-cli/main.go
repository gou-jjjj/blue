package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

var (
	remoteAddr net.Addr
)

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

func Connect() net.Conn {
	d, err := net.Dial("tcp", fmt.Sprintf("%s:%d", BC.Ip, BC.Port))
	if err != nil {
		ErrPrint(ErrConnect(err.Error()))
		os.Exit(1)
	}

	return d
}

// num get a
func main() {
	// 从标准输入创建一个新的 bufio.Reader
	reader := bufio.NewReader(os.Stdin)
	d := Connect()

	for {
		remoteAddr = d.RemoteAddr()
		fmt.Printf("%s> ", remoteAddr.String())
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

		exec, err := Exec(split)
		if err != nil {
			ErrPrint(err)
			continue
		}

		//for _, b := range exec {
		//	fmt.Printf("[%b]\n", b)
		//}
		//continue

		_, err = d.Write(exec)
		if err != nil {
			d = Connect()
			_, err = d.Write(exec)
			if err != nil {
				ErrPrint(ErrConnect(err.Error()))
			}
			continue
		}

		buf := make([]byte, 1024)
		n, err := d.Read(buf)
		if err != nil && err != io.EOF {
			ErrPrint(ErrConnect(err.Error()))
			continue
		}
		bytes := buf[:n]

		resp, err := NewReplyMessage(bytes)
		if err != nil {
			ErrPrint(err)
			continue
		}
		SuccessPrint(GreenMessage(resp))
	}
}
