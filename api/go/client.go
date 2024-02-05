package main

import (
	"blue/bsp"
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"strconv"
	"sync"
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
	return c.Addr()
}

func (c *Client) exec() (err error) {
	sends := 0

	for i := range c.WriteBuf {
		req := c.WriteBuf[i]
		_, err := c.conn.Write(req)
		if err != nil {
			return err
		}
		sends++
	}

	if sends != len(c.WriteBuf) {
		err = ErrSend(sends)
	}

	wg := sync.WaitGroup{}
	read := bufio.NewReader(c.conn)
	for i := 0; i < sends; i++ {
		bytes, err := read.ReadBytes(bsp.Done)
		if err != nil {
			if !errors.Is(err, io.EOF) {
				err = errors.New("read error")
			}
			break
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			res := bsp.NewReplyMessage(bytes)
			c.ReadBuf = append(c.ReadBuf, []byte(res))
		}()
	}
	wg.Wait()
	return
}

func (c *Client) Close() {
	c.conn.Close()
	c.ReadBuf = nil
	c.WriteBuf = nil
}
