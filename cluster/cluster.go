package cluster

import (
	"bufio"
	"bytes"
	"io"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	g "blue/api/go"
	"blue/bsp"
	add "blue/common/network"
)

const (
	sleep_ = 1
	done   = '\n'
)

type Subject interface {
	Register(...string)
	Unregister(...string)
	Notify(string)
}

const network = "tcp"

type Cluster struct {
	rw          sync.RWMutex
	observers   []string
	c           *Consistent
	listener    net.Listener
	ip          string
	port        int
	token       string
	tryTimes    int
	configAddr  string
	dialTimeout time.Duration
	remoteCli   map[string]*g.Client
}

func NewCluster(try int, port int, token string, dialTimeout time.Duration) *Cluster {
	clu := &Cluster{
		rw:          sync.RWMutex{},
		observers:   make([]string, 0),
		c:           NewConsistent(100),
		port:        port,
		token:       token,
		tryTimes:    try,
		dialTimeout: dialTimeout,
	}

	en0 := add.LocalIpEn0()
	if en0 == "" {
		panic("en0 err")
	}
	clu.ip = en0

	clu.addLocalAddr()
	clu.listen()

	return clu
}

func (c *Cluster) addLocalAddr() {
	c.Register(c.LocalAddr())
}

func (c *Cluster) LocalAddr() string {
	return c.ip + ":" + strconv.Itoa(c.port)
}

func (c *Cluster) listen() {
	l, err := net.Listen(network, c.LocalAddr())
	if err != nil {
		panic(err)
	}
	c.listener = l

	go c.accept()
}

func (c *Cluster) accept() {
	var conn net.Conn
	var err error

	for {
		conn, err = c.listener.Accept()
		if err != nil {
			panic(err)
		}

		go c.handle(conn)
	}
}

func (c *Cluster) Dial(ctx *bsp.BspProto) ([]byte, bool) {
W:
	if ctx.Key() == "" {
		return nil, false
	}

	remoteAddr := c.c.Get(ctx.Key())
	if c.LocalAddr() == remoteAddr {
		return nil, false
	}

	cli, err := c.getClient(remoteAddr)
	if err != nil {
		c.Unregister(remoteAddr)
		goto W
	}

	exec, err := cli.DirectExec(ctx.Buf())
	if err != nil {
		c.Unregister(remoteAddr)
		goto W
	}

	if exec == nil {
		c.Unregister(remoteAddr)
		goto W
	}

	return exec, true
}

func (c *Cluster) GetClusterAddrs(addr string) []string {
	conn, err := net.Dial(network, addr)
	if err != nil {
		println(err.Error())
		return []string{}
	}
	defer func() {
		_ = conn.Close()
	}()

	_, err = conn.Write([]byte("=\n"))
	if err != nil {
		println(err.Error())
		return []string{}
	}

	time.Sleep(1 * time.Second)
	red := bufio.NewReader(conn)
	byt, err := red.ReadBytes(done)
	if err != nil || len(byt) == 0 {

		return []string{}
	}

	res := string(byt)
	return strings.Fields(res)
}

func (c *Cluster) handle(conn net.Conn) {
	defer func() {
		_ = conn.Close()
	}()

	reader := bufio.NewReader(conn)
	for {
		addr, err := reader.ReadString(done)
		if err != nil && err != io.EOF {
			return
		}

		if len(addr) == 0 || addr[0] != '+' && addr[0] != '-' && addr[0] != '=' {
			continue
		}

		if addr[0] == '+' {
			c.Register(addr[1 : len(addr)-1])
		} else if addr[0] == '-' {
			c.Unregister(addr[1 : len(addr)-1])
		} else {
			res := bytes.Buffer{}
			for _, addr = range c.observers {
				res.WriteString(addr)
				res.WriteString(" ")
			}
			res.WriteByte(done)

			_, err = conn.Write(res.Bytes())
			if err != nil {
				return
			}
		}
	}
}

func (c *Cluster) Register(addr ...string) {
	c.rw.Lock()
	defer c.rw.Unlock()
	c.observers = append(c.observers, addr...)
	for i := range addr {
		c.c.Add(addr[i])
	}
}

func (c *Cluster) Unregister(addr ...string) {
	c.rw.Lock()
	defer c.rw.Unlock()
	for _, a := range addr {
		for i, obs := range c.observers {
			if obs == a {
				c.observers = append(c.observers[:i], c.observers[i+1:]...)
				break
			}
		}
	}

	for i := range addr {
		c.c.Remove(addr[i])
	}
}

func (c *Cluster) Notify(addr string) {
	for _, observer := range c.observers {
		if observer == addr {
			continue
		}

		go func(addr string, observer string) {
			conn, err := net.DialTimeout(network, observer, c.dialTimeout)
			if err != nil {
				for i := 0; i < c.tryTimes; i++ {
					conn, err = net.DialTimeout(network, observer, c.dialTimeout)
					if err == nil {
						break
					} else {
						time.Sleep(time.Duration(sleep_) * time.Second)
					}
				}
			}

			if err != nil {
				return
			}

			_, _ = conn.Write([]byte(addr))
			_ = conn.Close()
		}(addr, observer)
	}
}

func (c *Cluster) online(addr string) {
	c.Notify("+" + addr)
}

func (c *Cluster) offline(addr string) {
	c.Notify("-" + addr)
}

func (c *Cluster) setClient(addr string) (*g.Client, error) {
	cli, err := g.NewClient(g.WithDefaultOpt(), func(c *g.Config) {
		c.Addr = addr
	})

	if err != nil {
		return nil, err
	}

	c.remoteCli[addr] = cli
	return cli, nil
}

func (c *Cluster) getClient(addr string) (*g.Client, error) {
	if conn, ok := c.remoteCli[addr]; ok {
		return conn, nil
	} else {
		nconn, err := c.setClient(addr)
		if err != nil {
			return nil, err
		}

		return nconn, nil
	}
}

func (c *Cluster) Close() {
	addr := c.LocalAddr()

	c.offline(addr)
	_ = c.listener.Close()
}
