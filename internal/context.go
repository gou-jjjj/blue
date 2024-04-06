package internal

import (
	"context"
	"errors"
	"fmt"
	"net"
	"sync"
	"time"

	"blue/bsp"
	"blue/common/rand"
	"blue/config"
	"blue/log"
)

const (
	TokenLen = 10
)

type Context struct {
	context.Context
	nextExec  Exec
	conn      net.Conn
	db        uint8
	cliToken  string
	request   *bsp.BspProto
	response  bsp.Reply
	maxActive time.Duration
	isclose   int
}

var bconnPool = sync.Pool{
	New: func() interface{} {
		return &Context{
			db:       1,
			isclose:  0,
			cliToken: rand.RandString(TokenLen),
		}
	},
}

func (c *Context) String() string {
	return fmt.Sprintf("RemoteAddr:%s,cliToken:%s", c.conn.RemoteAddr().String(), c.cliToken)
}

func NewContext(ctx context.Context, conn net.Conn) *Context {
	bconn, ok := bconnPool.Get().(*Context)
	if !ok {
		nctx := &Context{
			Context:  ctx,
			conn:     conn,
			db:       1,
			isclose:  0,
			cliToken: rand.RandString(TokenLen),
			maxActive: time.Duration(config.CliCfg.ClientActive) *
				time.Second,
		}
		log.Info(bconn.String())
		return nctx
	}
	bconn.Context = ctx
	bconn.conn = conn
	bconn.db = 1
	bconn.isclose = 0
	bconn.cliToken = rand.RandString(TokenLen)
	bconn.maxActive = time.Duration(config.CliCfg.ClientActive) * time.Second
	log.Info(bconn.String())
	return bconn
}

func (c *Context) GetDB() uint8 {
	return c.db
}

func (c *Context) SetDB(index uint8) {
	c.db = index
}

func (c *Context) Reply() (int, error) {
	if c.response == nil {
		return c.conn.Write(bsp.NewErr(bsp.ErrReplication).Bytes())
	}
	log.Info(fmt.Sprintf("reply: %s", c.response.String()))
	if c.isClose() {
		return 0, errors.New("conn is close")
	}
	return c.conn.Write(c.response.Bytes())
}

func (c *Context) isClose() bool {
	return c.isclose == 1
}

func (c *Context) Close() {
	if c.conn != nil {
		_ = c.conn.Close()
	}
	c.db = 1
	c.Context = nil
	c.nextExec = nil
	c.request = nil
	c.response = nil
	c.maxActive = 0
	c.isclose = 1
	bconnPool.Put(c)
	return
}
