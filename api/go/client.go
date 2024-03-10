package blue

import (
	"blue/bsp"
	"bufio"
	"bytes"
	"errors"
	"io"
	"net"
	"strconv"
	"time"
)

var (
	ErrPip = errors.New("pipeline error")
)

type Option func(*Config)

type Config struct {
	Ip       string
	Port     int
	DB       int
	TryTimes int
	TimeOut  time.Duration
}

type Client struct {
	Config
	conn net.Conn
}

func WithDefaultOpt() Option {
	return func(c *Config) {
		c.Ip = "127.0.0.1"
		c.Port = 13140
		c.TimeOut = 5 * time.Second
		c.TryTimes = 3
		c.DB = 1
	}
}

func NewClient(opts ...Option) *Client {
	c := &Config{}
	for _, opt := range opts {
		opt(c)
	}

	cli := &Client{
		Config: *c,
	}

	cli.connect()

	_, err := cli.Select(strconv.Itoa(c.DB))
	if err != nil {
		panic(err)
	}

	return cli
}

func (c *Client) connect() {
	var err error
	for i := 0; i < c.TryTimes; i++ {
		c.conn, err = net.DialTimeout("tcp", c.Ip+":"+strconv.Itoa(c.Port), c.TimeOut)
		if err == nil {
			return
		}
	}

	panic(err)
}

func (c *Client) RemoteAddr() string {
	return c.conn.RemoteAddr().String()
}

func (c *Client) SetAddr(ip string, port int) {
	c.Ip = ip
	c.Port = port
}

func (c *Client) Version() (string, error) {
	build := bsp.NewRequestBuilder(bsp.VERSION).Build()

	return c.exec(build)
}

func (c *Client) Del(key string) (string, error) {
	build := bsp.NewRequestBuilder(bsp.DEL).WithKey(key).Build()

	return c.exec(build)
}

func (c *Client) Nset(k, num string) (string, error) {
	build := bsp.NewRequestBuilder(bsp.NSET).
		WithKey(k).
		WithValueNum(num).
		Build()

	return c.exec(build)
}

func (c *Client) Get(k string) (string, error) {
	build := bsp.NewRequestBuilder(bsp.GET).WithKey(k).Build()
	return c.exec(build)
}

func (c *Client) Set(k, v string) (string, error) {
	build := bsp.NewRequestBuilder(bsp.SET).
		WithKey(k).
		WithValueStr(v).
		Build()

	return c.exec(build)
}

func (c *Client) Len(k string) (string, error) {
	build := bsp.NewRequestBuilder(bsp.LEN).WithKey(k).Build()

	return c.exec(build)
}

func (c *Client) Kvs() (string, error) {
	build := bsp.NewRequestBuilder(bsp.KVS).Build()

	return c.exec(build)
}

func (c *Client) Nget(k string) (string, error) {
	build := bsp.NewRequestBuilder(bsp.NGET).
		WithKey(k).
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
