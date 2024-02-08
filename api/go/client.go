package main

import (
	"blue/bsp"
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"net"
	"strconv"
)

var (
	ErrPip = errors.New("pipeline error")
)

type Option func(*Config)

type Config struct {
	ip   string
	port int
}

type Client struct {
	Config
	conn net.Conn
}

func WithDefaultOpt() Option {
	return func(c *Config) {
		c.ip = "127.0.0.1"
		c.port = 8080
	}
}

func NewClient(opts ...Option) *Client {
	c := &Config{}
	for _, opt := range opts {
		opt(c)
	}

	dial, err := net.Dial("tcp", c.ip+":"+strconv.Itoa(c.port))
	if err != nil {
		panic(err)
	}

	cli := &Client{
		Config: *c,
		conn:   dial,
	}
	return cli
}

func (c *Client) Addr() string {
	return c.conn.RemoteAddr().String()
}

func (c *Client) SetAddr(ip string, port int) {
	c.ip = ip
	c.port = port
}

func (c *Client) Version() (string, error) {
	build := bsp.NewRequestBuilder(bsp.VERSION).Build()

	return c.exec(build)
}

func (c *Client) Nset(k, num string) (string, error) {
	build := bsp.NewRequestBuilder(bsp.NSET).
		WithKey(k).
		WithValueNum(num).
		Build()

	return c.exec(build)
}

func (c *Client) Select(num ...string) (string, error) {
	build := bsp.NewRequestBuilder(bsp.SELECT)
	if len(num) != 0 {
		build.WithKey(num[0])
	}
	return c.exec(build.Build())
}

func (c *Client) exec(buf []byte) (string, error) {
	_, err := c.conn.Write(buf)
	if err != nil {
		return "", err
	}

	read := bufio.NewReader(c.conn)

	bys, err1 := read.ReadBytes(bsp.Done)
	if err1 != nil {
		return "", err1
	}
	fmt.Printf("%+b\n", bys)
	return bsp.NewReplyMessage(bys)
}

func (c *Client) execPipeline(buf [][]byte) (s []string, err error) {
	b := bytes.Buffer{}

	for _, v := range buf {
		b.Write(v)
	}

	_, err = c.conn.Write(b.Bytes())
	if err != nil {
		return nil, err
	}

	read := bufio.NewReader(c.conn)
	for {
		bys, err1 := read.ReadBytes(bsp.Done)
		if err1 != nil {
			if !errors.Is(err1, io.EOF) {
				err1 = errors.New("read error")
			}
			break
		}

		res, err := bsp.NewReplyMessage(bys)
		if err != nil {
			return nil, err
		}

		s = append(s, res)
	}

	if len(s) != len(buf) {
		return nil, ErrPip
	}

	return
}

func (c *Client) Close() {
	c.conn.Close()
}
