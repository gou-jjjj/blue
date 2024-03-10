package internal

import (
	"context"
	"fmt"
	"net"
	"testing"
	"time"
)

type Hello struct {
}

func (h *Hello) ExecChain(c *Context) bool {
	return true
}

func (h *Hello) Close() {

}

func (h *Hello) Handle(ctx context.Context, conn net.Conn) {
	req := make([]byte, 1024)
	read, err := conn.Read(req)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(req[:read]))
	_, err = conn.Write(append([]byte("hello: "), req[:read]...))
	if err != nil {
		panic(err)
	}
}

func HelloClient(ip string, port int) {
	time.Sleep(1 * time.Second)
	dial, err := net.Dial("tcp", fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		panic(err)
	}
	defer dial.Close()

	_, err = dial.Write([]byte("blue"))
	if err != nil {
		panic(err)
	}

	resp := make([]byte, 1024)
	read, err := dial.Read(resp)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s\n", resp[:read])
}

func TestName(t *testing.T) {
	ip := "127.0.0.1"
	port := 8080

	HelloServer := NewServer(func(c *Config) {
		c.Ip = ip
		c.Port = port

		c.HandlerFunc = &Hello{}
	})

	go HelloClient(ip, port)

	HelloServer.Start()
}
