package cluster

import (
	"hash/crc32"
	"slices"
	"sort"
	"strconv"
	"sync"
)

// 默认副本数
var defalutReplicas = 100

// 默认哈希函数
var defalutHash = crc32.ChecksumIEEE

// 定义哈希函数类型
type hashfunc func(data []byte) uint32

// 节点结构体
type Node struct {
	hash uint32
	key  string
}

// 一致性哈希结构体
type Consistent struct {
	members  map[string]bool // 成员映射
	nodes    []Node          // 节点列表
	replicas int             // 副本数
	count    int64           // 节点计数
	rw       sync.RWMutex    // 读写锁
	h        hashfunc        // 哈希函数
}

// 新建一致性哈希对象
// replicas: 副本数量
func NewConsistent(replicas int) *Consistent {
	if replicas >= defalutReplicas {
		replicas = defalutReplicas
	}

	c := new(Consistent)
	c.replicas = defalutReplicas
	c.nodes = make([]Node, 0)
	c.members = make(map[string]bool)
	c.rw = sync.RWMutex{}
	c.h = defalutHash
	return c
}

// 为一致性哈希添加元素
// elt: 需要添加的元素
func (c *Consistent) Add(elt string) {
	c.rw.Lock()
	defer c.rw.Unlock()

	if _, ok := c.members[elt]; !ok {
		c.add(elt)
	}
}

func (c *Consistent) eltKey(elt string, idx int) string {
	return strconv.Itoa(idx) + elt
}

// 实际添加元素操作
// elt: 需要添加的元素
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

// 从一致性哈希中移除元素
// elt: 需要移除的元素
func (c *Consistent) Remove(elt string) {
	c.rw.Lock()
	defer c.rw.Unlock()
	c.remove(elt)
}

// 实际移除元素操作
// elt: 需要移除的元素
func (c *Consistent) remove(elt string) {
	slices.DeleteFunc(c.nodes, func(n Node) bool {
		return n.key == elt
	})
	delete(c.members, elt)
	//c.updateSortedHashes()
	c.count--
}

// 获取一致性哈希的所有成员
func (c *Consistent) Members() []string {
	c.rw.RLock()
	defer c.rw.RUnlock()

	m := make([]string, 0, len(c.members))
	for k := range c.members {
		m = append(m, k)
	}
	return m
}

// 根据名称获取对应的服务节点
// name: 服务名称
// 返回: 对应的服务节点名称
func (c *Consistent) Get(name string) string {
	c.rw.RLock()
	defer c.rw.RUnlock()

	key := c.hashKey(name)
	i := c.search(key)

	return c.nodes[i%len(c.nodes)].key
}

// 查找节点的位置
// key: 需要查找的哈希值
// 返回: 查找到的节点索引
func (c *Consistent) search(key uint32) int {
	s, _ := slices.BinarySearchFunc(c.nodes, key, func(n Node, i uint32) int {
		if n.hash >= i {
			return 1
		}
		return -1
	})

	return s
}

// 计算键的哈希值
// key: 需要计算哈希值的键
// 返回: 键的哈希值
func (c *Consistent) hashKey(key string) uint32 {
	if len(key) < 64 {
		var scratch [64]byte
		copy(scratch[:], key)
		return c.h(scratch[:len(key)])
	}
	return c.h([]byte(key))
}
