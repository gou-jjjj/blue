package cluster

import "sync"

type Cluster struct {
	rw        sync.RWMutex
	observers []string
}

func NewCluster() {

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

func (c *Cluster) Online(addr ...string) {
	c.Register(addr...)

	for i := range c.observers {

		_ = i
	}
}

func (c *Cluster) Offline(addr ...string) {
	c.Unregister(addr...)
	for i := range c.observers {
		_ = i
	}
}
