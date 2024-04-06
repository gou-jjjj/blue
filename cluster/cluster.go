package cluster

//声明了用于集群通信的一系列结构和方法。

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"slices"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	g "blue/api/go"
	"blue/bsp"
	add "blue/common/network"
	"blue/log"
)

// 常量定义
const (
	network = "tcp" // 使用的网络协议
	sleep_  = 1     // 等待时间（单位未指定）
	done    = '\n'  // 消息结束标识

	addClu = '+' // 增加集群节点标识
	subClu = '-' // 删除集群节点标识
	getClu = '=' // 获取集群节点标识
)

// Subject 接口定义了注册、注销和通知的方法。
type Subject interface {
	Register(...string)
	Unregister(...string)
	Notify(string)
}

// Cluster 结构体包含了集群通信所需的全部属性。
type Cluster struct {
	ctx       context.Context    // 上下文，用于控制子协程的启动和停止
	cancel    context.CancelFunc // 取消函数，用于取消上下文
	isClose   *atomic.Bool       // 标记集群是否关闭
	rw        sync.RWMutex       // 读写锁，保护集群状态
	observers []string           // 观察者列表
	c         *Consistent        // 一致性哈希对象，用于节点分布
	listener  net.Listener       // 监听器，用于接收连接
	ip        string             // 服务器IP地址
	port      int                // 服务器端口号
	tryTimes  int                // 连接重试次数

	myAddr      string               // 本地地址
	cliAddr     string               // 客户端地址
	dialTimeout time.Duration        // 连接超时时间
	remoteCli   map[string]*g.Client // 远程客户端映射
	clu2cli     map[string]string    // 集群地址到客户端地址的映射
}

// NewCluster 创建并初始化一个新的Cluster实例。
//
// try: 连接重试次数。
// ip: 服务器IP地址。
// port: 服务器端口号。
// myAddr: 本地地址，用于标识当前节点。
// cliAddr: 客户端地址，用于与客户端建立连接。
// dialTimeout: 连接超时时间。
// 返回: 初始化后的Cluster实例。
func NewCluster(
	try int,
	ip string,
	port int,
	myAddr string,
	cliAddr string,
	dialTimeout time.Duration,
) *Cluster {

	ctx, cancelFunc := context.WithCancel(context.Background())
	isCls := atomic.Bool{}
	isCls.Store(false)

	clu := &Cluster{
		ctx:    ctx,
		cancel: cancelFunc,

		isClose:     &isCls,
		rw:          sync.RWMutex{},
		observers:   make([]string, 0),
		c:           NewConsistent(100),
		remoteCli:   make(map[string]*g.Client),
		myAddr:      myAddr,
		ip:          ip,
		port:        port,
		tryTimes:    try,
		cliAddr:     cliAddr,
		clu2cli:     make(map[string]string),
		dialTimeout: dialTimeout,
	}

	if ip == "" {
		en0 := add.LocalIpEn0()
		if en0 == "" {
			panic("en0 err")
		}
		clu.ip = en0
	}

	clu.addLocalAddr()
	clu.listen()

	return clu
}

// addLocalAddr 将本地地址添加到集群中。
func (c *Cluster) addLocalAddr() {
	log.Info(fmt.Sprintf("add local addr: [%s]", c.LocalAddr()))
	c.Register(add.CombineAddr(c.LocalAddr(), c.cliAddr))
}

// LocalAddr 返回本地地址。
//
// 返回: 格式化的本地地址字符串。
func (c *Cluster) LocalAddr() string {
	if c.myAddr != "" {
		return c.myAddr
	}

	return c.ip + ":" + strconv.Itoa(c.port)
}

// listen 启动监听器，等待客户端连接。
func (c *Cluster) listen() {
	listenAddr := fmt.Sprintf("%s:%d", c.ip, c.port)
	l, err := net.Listen(network, listenAddr)
	if err != nil {
		panic(err)
	}
	c.listener = l
	log.Info(fmt.Sprintf("cluster listen on %v ...", listenAddr))
	go c.accept()
}

// accept 在监听器上接收新的客户端连接，并处理这些连接。
func (c *Cluster) accept() {
	var conn net.Conn
	var err error

	for !c.isClose.Load() {
		conn, err = c.listener.Accept()
		if err != nil {
			if !errors.Is(err, net.ErrClosed) {
				log.Error(err.Error())
			}
			return
		}

		go c.handle(conn)
	}
}

// Dial 尝试与指定的集群地址建立连接，并执行通信。
//
// ctx: 传递给客户端的协议数据。
// 返回: 执行结果，如果执行成功则为true。
func (c *Cluster) Dial(ctx *bsp.BspProto) ([]byte, bool) {
	if ctx.Key() == "" {
		return nil, false
	}

	for {
		remoteAddr := c.c.Get(ctx.Key())
		if c.LocalAddr() == remoteAddr {
			return nil, false
		}

		cli, err := c.getClient(remoteAddr)
		if err != nil {
			c.Offline(remoteAddr)
			continue
		}

		exec, err := cli.DirectExec(ctx.Buf())
		if err != nil {
			c.Offline(remoteAddr)
			continue
		}

		if exec == nil {
			c.Unregister(remoteAddr)
			continue
		}
		return exec, true
	}
}

// GetClusterAddr 通过指定的集群地址获取整个集群的地址列表。
//
// addr: 指定的集群地址。
// 返回: 集群中其他节点的地址列表。
func (c *Cluster) GetClusterAddr(addr string) []string {
	conn, err := net.Dial(network, addr)
	if err != nil {
		log.Error(err.Error())
		return []string{}
	}
	defer func() {
		_ = conn.Close()
	}()

	_, err = conn.Write([]byte("=\n"))
	if err != nil {
		log.Error(err.Error())
		return []string{}
	}

	time.Sleep(1 * time.Second)
	red := bufio.NewReader(conn)
	byt, err := red.ReadBytes(done)
	if err != nil || len(byt) == 0 {
		return []string{}
	}

	res := string(byt)
	addrs := strings.Fields(res)
	log.Info(fmt.Sprintf("clster addrs %+v", addrs))
	return addrs
}

// InitClusterAddr 初始化集群地址。
//
// addr: 要初始化的集群地址列表。
func (c *Cluster) InitClusterAddr(addr ...string) {
	c.Register(addr...)
}

// handle 处理网络连接请求。
// 对于每个连接，读取指令并根据指令类型执行相应的注册、注销或获取集群状态操作。
func (c *Cluster) handle(conn net.Conn) {
	defer func() {
		_ = conn.Close() // 确保连接关闭
	}()

	reader := bufio.NewReader(conn) // 使用缓冲读取连接数据
	for {
		select {
		case <-c.ctx.Done(): // 检查上下文是否取消
			return
		default:
			addr, err := reader.ReadString(done) // 读取地址信息
			if err != nil {
				if err != io.EOF { // 排除正常结束的情况
					log.Error(err.Error())
				}
				return
			}

			if len(addr) == 0 || addr[0] != '+' && addr[0] != '-' && addr[0] != '=' { // 检查地址格式
				continue
			}

			switch addr[0] { // 根据地址的第一个字符执行相应操作
			case addClu:
				c.Register(addr[1 : len(addr)-1])

			case subClu:
				c.Unregister(addr[1 : len(addr)-1])

			case getClu:
				bys := c.GetClu2cli() // 获取集群到客户端的映射
				bys = append(bys, done)

				_, err = conn.Write(bys) // 将映射写回连接
				if err != nil {
					log.Error(err.Error())
					return
				}
			}

		}
	}
}

// Register 注册集群地址。
// 遍历提供的地址列表，将有效的集群和客户端地址注册到集群中。
func (c *Cluster) Register(addr ...string) {
	c.rw.Lock()
	defer c.rw.Unlock()

	cluaddrs := make([]string, 0)
	cliaddrs := make([]string, 0)

	for i := range addr {
		addrs := add.SplitAddr(addr[i])

		if add.ParseAddr(addrs[0]) && add.ParseAddr(addrs[1]) { // 解析并验证地址
			cluaddrs = append(cluaddrs, addrs[0])
			cliaddrs = append(cliaddrs, addrs[1])
		}
	}

	for i, s := range cluaddrs {
		if !slices.Contains(c.observers, s) { // 检查是否已注册
			if cliaddrs[i] != c.cliAddr {
				client, err := c.setClient(cliaddrs[i]) // 设置客户端连接
				if err != nil {
					log.Warn(fmt.Sprintf("set client err: %v", err))
					continue
				}

				c.remoteCli[s] = client
			}

			c.clu2cli[s] = cliaddrs[i]

			c.observers = append(c.observers, s) // 更新观察者列表
			c.c.Add(s)                           // 在一致性哈希中添加
			log.Info(fmt.Sprintf("register [%v]", s))
		}
	}
}

// Unregister 注销集群地址。
// 从集群中移除指定的地址。
func (c *Cluster) Unregister(addr ...string) {
	c.rw.Lock()
	defer c.rw.Unlock()

	idx := 0
	for _, a := range addr {
		for i, obs := range c.observers {
			if obs == a {
				log.Info(fmt.Sprintf("unregister %+v", obs))
				c.delClient(obs)                                                    // 删除客户端连接
				delete(c.clu2cli, obs)                                              // 从集群到客户端映射中移除
				c.observers[idx], c.observers[i] = c.observers[i], c.observers[idx] // 交换并更新观察者列表
				idx++
				break
			}
		}
	}

	clear(c.observers[:idx]) // 清理已移除的观察者
	c.observers = c.observers[idx:]
	for i := range addr {
		c.c.Remove(addr[i]) // 在一致性哈希中移除
	}
}

// Notify 通知集群中除自身外的所有观察者关于地址的变化。
func (c *Cluster) Notify(addr string) {
	splitAddr := add.SplitAddr(addr)
	cluAddr := splitAddr[0]

	for _, observer := range c.observers {
		if observer != cluAddr {
			c.notify(addr, observer) // 发送通知
		}
	}
}

// notify 向指定的观察者发送地址变化通知。
func (c *Cluster) notify(addr string, observer string) {
	conn, err := net.DialTimeout(network, observer, c.dialTimeout) // 建立连接
	if err != nil {
		for i := 0; i < c.tryTimes; i++ {
			conn, err = net.DialTimeout(network, observer, c.dialTimeout)
			if err != nil {
				time.Sleep(time.Duration(sleep_) * time.Second) // 重试间隔
			}
		}
	}

	if err != nil {
		log.Error(fmt.Sprintf("notify: %v", err))
		return
	}

	_, _ = conn.Write([]byte(addr)) // 发送地址信息
	_ = conn.Close()                // 关闭连接
}

// Online 将地址注册到集群中并通知所有观察者。
func (c *Cluster) Online(addr string) {
	c.Register(addr)
	c.Notify(fmt.Sprintf("+%s\n", addr))
}

// Offline 从集群中注销地址并通知所有观察者。
func (c *Cluster) Offline(addr string) {
	c.Unregister(addr)
	c.Notify(fmt.Sprintf("-%s\n", addr))
}

// setClient 创建并配置与集群地址对应的客户端。
func (c *Cluster) setClient(addr string) (*g.Client, error) {
	cli, err := g.NewClient(g.WithDefaultOpt(), func(c *g.Config) {
		c.Addr = addr
		c.DB = 1
		c.DefaultDB = false
		c.TimeOut = 1 * time.Minute
	})

	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	return cli, nil
}

// getClient 获取与集群地址对应的客户端。
// 如果客户端不存在，则尝试创建并缓存。
func (c *Cluster) getClient(addr string) (*g.Client, error) {
	c.rw.RLock()
	cli, ok := c.remoteCli[addr]
	c.rw.RUnlock()
	if ok {
		return cli, nil
	}

	c.rw.Lock()
	defer c.rw.Unlock()

	cliAddr, ok := c.clu2cli[addr]
	if !ok {
		return nil, fmt.Errorf("client address not found for: %s", addr)
	}

	cli, err := c.setClient(cliAddr)
	if err != nil {
		return nil, err
	}

	c.remoteCli[addr] = cli
	return cli, nil
}

// delClient 删除与集群地址对应的客户端。
func (c *Cluster) delClient(addr string) {
	if _, ok := c.remoteCli[addr]; ok {
		delete(c.remoteCli, addr)
	}
}

// GetObservers 返回当前所有观察者的地址列表。
func (c *Cluster) GetObservers() []string {
	return c.observers
}

// GetConsistent 返回集群使用的一致性哈希对象。
func (c *Cluster) GetConsistent() *Consistent {
	return c.c
}

// GetClu2cli 返回集群到客户端的地址映射。
func (c *Cluster) GetClu2cli() []byte {
	res := bytes.Buffer{}
	for k, v := range c.clu2cli {
		res.WriteString(add.CombineAddr(k, v))
		res.WriteString(" ")
	}

	return res.Bytes()
}

// Close 关闭集群连接，注销本地地址并通知所有观察者。
func (c *Cluster) Close() {
	addr := c.LocalAddr()

	c.isClose.Store(true)
	c.Offline(addr)
	c.cancel()
	_ = c.listener.Close()
}
