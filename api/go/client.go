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
	Addr      string
	DB        int
	TryTimes  int
	token     string
	TimeOut   time.Duration
	DefaultDB bool
}

type Client struct {
	Config
	conn net.Conn
}

func WithDefaultOpt() Option {
	return func(c *Config) {
		c.Addr = "127.0.0.1:13140"
		c.TimeOut = 5 * time.Second
		c.TryTimes = 3
		c.DB = 1
	}
}

func WithCluster(addr string, token string) Option {
	return func(c *Config) {
		c.Addr = addr
		c.token = token
	}
}

func WithAddr(addr string) Option {
	return func(c *Config) {
		c.Addr = addr
	}
}

func WithToken(token string) Option {
	return func(c *Config) {
		c.token = token
	}
}

func NewClient(opts ...Option) (*Client, error) {
	c := &Config{}
	for _, opt := range opts {
		opt(c)
	}

	cli := &Client{
		Config: *c,
	}

	err := cli.connect()
	if err != nil {
		return nil, err
	}

	if c.DefaultDB {
		_, err = cli.Select(strconv.Itoa(c.DB))
		if err != nil {
			return nil, err
		}
	}

	return cli, nil
}

func (c *Client) connect() error {
	var err error
	for i := 0; i < c.TryTimes; i++ {
		c.conn, err = net.DialTimeout("tcp", c.Addr, c.TimeOut)
		if err == nil {
			return nil
		}
	}

	return err
}

func (c *Client) RemoteAddr() string {
	return c.conn.RemoteAddr().String()
}

func (c *Client) DirectExec(buf []byte) ([]byte, error) {
	_, err := c.conn.Write(buf)
	if err != nil {
		return nil, err
	}

	read := bufio.NewReader(c.conn)
	return read.ReadBytes(bsp.Done)
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
