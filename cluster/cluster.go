package cluster

import (
	g "blue/api/go"
	"bufio"
	"fmt"
	"io"
	"net"
	"sync"
	"time"
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
	addr        string
	token       string
	tryTimes    int
	configAddr  string
	dialTimeout time.Duration
}

func NewCluster(try int, confAddr string, token string, dialTimeout time.Duration) *Cluster {
	clu := &Cluster{
		rw:          sync.RWMutex{},
		observers:   make([]string, 0),
		c:           NewConsistent(100),
		addr:        confAddr,
		token:       token,
		tryTimes:    try,
		dialTimeout: dialTimeout,
	}

	clu.Listen()

	return clu
}

func (c *Cluster) Init() {
	once := sync.Once{}
	once.Do(func() {
		g.NewClient(g.WithCluster(c.addr, c.token))
	})
}

func (c *Cluster) Listen() {
	l, err := net.Listen(network, c.addr)
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

func (c *Cluster) handle(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	for {
		addr, err := reader.ReadString(done)
		if err != nil && err != io.EOF {
			return
		}

		if len(addr) == 0 || addr[0] != '+' && addr[0] != '-' {
			continue
		}

		fmt.Println("addr:", addr)
		if addr[0] == '+' {
			c.Register(addr[1 : len(addr)-1])
		} else {
			c.Unregister(addr[1 : len(addr)-1])
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

			conn.Write([]byte(addr))
			conn.Close()
		}(addr, observer)
	}
}

func (c *Cluster) Online(addr string) {
	c.Notify("+" + addr)
}

func (c *Cluster) Offline(addr string) {
	c.Notify("-" + addr)
}

func (c *Cluster) Close() {
	addr := c.listener.Addr().String()

	c.Offline(addr)
	c.listener.Close()
}
