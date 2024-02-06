package main

import (
	"blue/bsp"
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"strconv"
)

var (
	ErrSend = func(i int) error {
		return errors.New(fmt.Sprintf("Error, just send: %d ", i))
	}
)

type Option func(*Config)

type Config struct {
	ip   string
	port int
}

type Client struct {
	Config
	conn net.Conn

	ReadBuf  [][]byte
	WriteBuf [][]byte
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

func (c *Client) Version() string {
	c.WriteBuf = append(c.WriteBuf, bsp.NewReq(bsp.VERSION))

	return string(c.ReadBuf[0])
}

func (c *Client) exec() (err error) {
	sends := 0

	for i := range c.WriteBuf {
		req := c.WriteBuf[i]
		_, err = c.conn.Write(req)
		if err != nil {
			return
		}
		sends++
	}

	if sends != len(c.WriteBuf) {
		err = ErrSend(sends)
	}

	read := bufio.NewReader(c.conn)
	for i := 0; i < sends; i++ {
		bytes, err1 := read.ReadBytes(bsp.Done)
		if err1 != nil {
			if !errors.Is(err1, io.EOF) {
				err1 = errors.New("read error")
			}
			break
		}

		res, err1 := bsp.NewReplyMessage(bytes)
		if err1 != nil {
			return err1
		}
		c.ReadBuf = append(c.ReadBuf, []byte(res))
	}

	return
}

func (c *Client) Close() {
	c.conn.Close()
	c.ReadBuf = nil
	c.WriteBuf = nil
}
