package cluster

import (
	"blue/config"
	"hash/crc32"
	"slices"
	"sort"
	"strconv"
	"sync"
)

var defalutReplicas = config.BC.Cluster.Replicas
var defalutHash = crc32.ChecksumIEEE

type hashfunc func(data []byte) uint32

type Node struct {
	hash uint32
	key  string
}

type Consistent struct {
	members  map[string]bool
	nodes    []Node
	replicas int
	count    int64
	rw       sync.RWMutex
	h        hashfunc
}

func NewConsistent() *Consistent {
	c := new(Consistent)
	c.replicas = defalutReplicas
	c.nodes = make([]Node, 0)
	c.members = make(map[string]bool)
	c.rw = sync.RWMutex{}
	c.h = defalutHash
	return c
}

func (c *Consistent) eltKey(elt string, idx int) string {
	return strconv.Itoa(idx) + elt
}

func (c *Consistent) Add(elt string) {
	c.rw.Lock()
	defer c.rw.Unlock()
	c.add(elt)
}

func (c *Consistent) add(elt string) {
	for i := 0; i < c.replicas; i++ {
		c.nodes = append(c.nodes, Node{
			hash: c.hashKey(c.eltKey(elt, i)),
			key:  elt,
		})
	}
	c.members[elt] = true

	sort.Slice(c.nodes, func(i, j int) bool {
		return c.nodes[i].hash < c.nodes[j].hash
	})

	c.count++
}

func (c *Consistent) Remove(elt string) {
	c.rw.Lock()
	defer c.rw.Unlock()
	c.remove(elt)
}

func (c *Consistent) remove(elt string) {
	slices.DeleteFunc(c.nodes, func(n Node) bool {
		return n.key == elt
	})
	delete(c.members, elt)
	//c.updateSortedHashes()
	c.count--
}

func (c *Consistent) Members() []string {
	c.rw.RLock()
	defer c.rw.RUnlock()

	m := make([]string, 0, len(c.members))
	for k := range c.members {
		m = append(m, k)
	}
	return m
}

func (c *Consistent) Get(name string) string {
	c.rw.RLock()
	defer c.rw.RUnlock()

	key := c.hashKey(name)
	i := c.search(key)

	return c.nodes[i%len(c.nodes)].key
}

func (c *Consistent) search(key uint32) int {
	s, _ := slices.BinarySearchFunc(c.nodes, key, func(n Node, i uint32) int {
		if n.hash >= i {
			return 1
		}
		return -1
	})

	return s
}

func (c *Consistent) hashKey(key string) uint32 {
	if len(key) < 64 {
		var scratch [64]byte
		copy(scratch[:], key)
		return c.h(scratch[:len(key)])
	}
	return c.h([]byte(key))
}
