package cluster

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

const (
	sleep_ = 1
	done   = '\n'

	addClu = '+'
	subClu = '-'
	getClu = '='
)

type Subject interface {
	Register(...string)
	Unregister(...string)
	Notify(string)
}

const network = "tcp"

type Cluster struct {
	ctx     context.Context
	cancel  context.CancelFunc
	isClose *atomic.Bool

	rw        sync.RWMutex
	observers []string
	c         *Consistent
	listener  net.Listener
	ip        string
	port      int
	tryTimes  int

	myAddr      string
	cliAddr     string
	dialTimeout time.Duration
	remoteCli   map[string]*g.Client
	clu2cli     map[string]string
}

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

func (c *Cluster) addLocalAddr() {
	log.Info(fmt.Sprintf("add local addr: [%s]", c.LocalAddr()))
	c.Register(add.CombineAddr(c.LocalAddr(), c.cliAddr))
}

func (c *Cluster) LocalAddr() string {
	if c.myAddr != "" {
		return c.myAddr
	}

	return c.ip + ":" + strconv.Itoa(c.port)
}

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

func (c *Cluster) Dial(ctx *bsp.BspProto) ([]byte, bool) {
	if ctx.Key() == "" {
		return nil, false
	}
	fmt.Println(1)
W:
	remoteAddr := c.c.Get(ctx.Key())
	if c.LocalAddr() == remoteAddr {
		return nil, false
	}
	fmt.Println(2)
	cli, err := c.getClient(remoteAddr)
	if err != nil {
		c.Offline(remoteAddr)
		goto W
	}
	fmt.Println(3)
	exec, err := cli.DirectExec(ctx.Buf())
	if err != nil {
		c.Offline(remoteAddr)
		goto W
	}
	fmt.Println(4)
	if exec == nil {
		c.Unregister(remoteAddr)
		goto W
	}
	fmt.Println(5)
	return exec, true
}

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

func (c *Cluster) InitClusterAddr(addr ...string) {
	c.Register(addr...)
}

func (c *Cluster) handle(conn net.Conn) {
	defer func() {
		_ = conn.Close()
	}()

	reader := bufio.NewReader(conn)
	for {
		select {
		case <-c.ctx.Done():
			return
		default:
			addr, err := reader.ReadString(done)
			if err != nil {
				if err != io.EOF {
					log.Error(err.Error())
				}
				return
			}

			if len(addr) == 0 || addr[0] != '+' && addr[0] != '-' && addr[0] != '=' {
				continue
			}

			fmt.Printf("[%+v]\n", addr)

			switch addr[0] {
			case addClu:
				c.Register(addr[1 : len(addr)-1])

			case subClu:
				c.Unregister(addr[1 : len(addr)-1])

			case getClu:
				bys := c.GetClu2cli()
				bys = append(bys, done)

				_, err = conn.Write(bys)
				if err != nil {
					log.Error(err.Error())
					return
				}
			}

		}
	}
}

func (c *Cluster) Register(addr ...string) {
	c.rw.Lock()
	defer c.rw.Unlock()

	cluaddrs := make([]string, 0)
	cliaddrs := make([]string, 0)

	for i := range addr {
		addrs := add.SplitAddr(addr[i])

		if add.ParseAddr(addrs[0]) && add.ParseAddr(addrs[1]) {
			cluaddrs = append(cluaddrs, addrs[0])
			cliaddrs = append(cliaddrs, addrs[1])
		}
	}

	for i, s := range cluaddrs {
		if !slices.Contains(c.observers, s) {
			if cliaddrs[i] != c.cliAddr {
				client, err := c.setClient(cliaddrs[i])
				if err != nil {
					log.Warn(fmt.Sprintf("set client err: %v", err))
					continue
				}

				c.remoteCli[s] = client
			}

			c.clu2cli[s] = cliaddrs[i]

			c.observers = append(c.observers, s)
			c.c.Add(s)
			log.Info(fmt.Sprintf("register [%v]", s))
		}
	}
}

func (c *Cluster) Unregister(addr ...string) {
	c.rw.Lock()
	defer c.rw.Unlock()

	idx := 0
	for _, a := range addr {
		for i, obs := range c.observers {
			if obs == a {
				log.Info(fmt.Sprintf("unregister %+v", obs))
				c.delClient(obs)
				delete(c.clu2cli, obs)
				c.observers[idx], c.observers[i] = c.observers[i], c.observers[idx]
				idx++
				break
			}
		}
	}

	clear(c.observers[:idx])
	c.observers = c.observers[idx:]
	for i := range addr {
		c.c.Remove(addr[i])
	}
}

func (c *Cluster) Notify(addr string) {
	splitAddr := add.SplitAddr(addr)
	cluAddr := splitAddr[0]

	for _, observer := range c.observers {
		if observer != cluAddr {
			c.notify(addr, observer)
		}
	}
}

func (c *Cluster) notify(addr string, observer string) {
	conn, err := net.DialTimeout(network, observer, c.dialTimeout)
	if err != nil {
		for i := 0; i < c.tryTimes; i++ {
			conn, err = net.DialTimeout(network, observer, c.dialTimeout)
			if err != nil {
				time.Sleep(time.Duration(sleep_) * time.Second)
			}
		}
	}

	if err != nil {
		log.Error(fmt.Sprintf("notify: %v", err))
		return
	}

	_, _ = conn.Write([]byte(addr))
	_ = conn.Close()
}

func (c *Cluster) Online(addr string) {
	c.Register(addr)
	c.Notify(fmt.Sprintf("+%s\n", addr))
}

func (c *Cluster) Offline(addr string) {
	c.Unregister(addr)
	c.Notify(fmt.Sprintf("-%s\n", addr))
}

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

func (c *Cluster) getClient(addr string) (*g.Client, error) {
	var err error

	conn, ok := c.remoteCli[addr]
	if !ok {

		cliAddr, ok := c.clu2cli[addr]
		if !ok {
			errMsg := fmt.Sprintf("cliaddr not found: %s", cliAddr)
			log.Error(errMsg)
			return nil, errors.New(errMsg)
		}

		conn, err = c.setClient(cliAddr)
		if err != nil {
			log.Error(err.Error())
			return nil, err
		}
	}

	return conn, nil
}

func (c *Cluster) delClient(addr string) {
	if _, ok := c.remoteCli[addr]; ok {
		delete(c.remoteCli, addr)
	}
}

func (c *Cluster) GetObservers() []string {
	return c.observers
}

func (c *Cluster) GetConsistent() *Consistent {
	return c.c
}

func (c *Cluster) GetClu2cli() []byte {
	res := bytes.Buffer{}
	for k, v := range c.clu2cli {
		res.WriteString(add.CombineAddr(k, v))
		res.WriteString(" ")
	}

	return res.Bytes()
}

func (c *Cluster) Close() {
	addr := c.LocalAddr()

	c.isClose.Store(true)
	c.Offline(addr)
	c.cancel()
	_ = c.listener.Close()
}
