package cluster

import (
	"blue/internal"
	"sync"
)

type Cluster struct {
	rw        sync.RWMutex
	observers []string
	c         Consistent
	s         internal.Server
}

func NewCluster() *Cluster {
	clu := &Cluster{
		rw:        sync.RWMutex{},
		observers: nil,
		c:         Consistent{},
		s:         internal.Server{},
	}

	return clu
}

func (c *Cluster) Listen() {
	c.s.Start()
}

func (c *Cluster) Register(addr ...string) {
	c.rw.Lock()
	defer c.rw.Unlock()
	c.observers = append(c.observers, addr...)
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
}
