package cluster

import (
	"fmt"
	"net"
	"sync"
	"time"
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
	port        int32
	tryTimes    int
	dialTimeout time.Duration
}

func NewCluster(ip string, port int32, try int, dialTimeout time.Duration) *Cluster {
	clu := &Cluster{
		rw:          sync.RWMutex{},
		observers:   make([]string, 0),
		c:           NewConsistent(),
		ip:          ip,
		port:        port,
		tryTimes:    try,
		dialTimeout: dialTimeout,
	}

	clu.Listen()

	return clu
}

func (c *Cluster) Listen() {
	l, err := net.Listen(network, fmt.Sprintf("%s:%d", c.ip, c.port))
	if err != nil {
		panic(err)
	}
	c.listener = l

	go c.accept()
}

func (c *Cluster) accept() {
	for {
		conn, err := c.listener.Accept()
		if err != nil {
			panic(err)
		}

		mess := make([]byte, 1024)
		read, err := conn.Read(mess)
		if err != nil {
			conn.Close()
			continue
		}

		mess = mess[:read]
		if mess[0] == '+' {
			c.Register(string(mess[1:]))
		} else if mess[0] == '-' {
			c.Unregister(string(mess[1:]))
		}

		conn.Close()
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
