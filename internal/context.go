package internal

import (
	"blue/config"
	"context"
	"fmt"
	"net"
	"sync"
	"time"

	"blue/bsp"
	"blue/common/rand"
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
}

var bconnPool = sync.Pool{
	New: func() interface{} {
		return &Context{
			db:       1,
			cliToken: rand.RandString(TokenLen),
		}
	},
}

func NewContext(ctx context.Context, conn net.Conn) *Context {
	bconn, ok := bconnPool.Get().(*Context)
	if !ok {
		return &Context{
			Context:  ctx,
			conn:     conn,
			db:       1,
			cliToken: rand.RandString(TokenLen),
			maxActive: time.Duration(config.BC.ClientConfig.ClientActive) *
				time.Second,
		}
	}
	bconn.Context = ctx
	bconn.conn = conn
	bconn.db = 1
	bconn.cliToken = rand.RandString(TokenLen)
	bconn.maxActive = time.Duration(config.BC.ClientConfig.ClientActive) * time.Second

	return bconn
}

func (c *Context) SetNext(next Exec) {
	c.nextExec = next
}

func (c *Context) GetDB() uint8 {
	fmt.Println("db: ", c.db)
	return c.db
}

func (c *Context) SetDB(index uint8) {
	c.db = index
}

func (c *Context) Reply() (int, error) {
	if c.response == nil {
		return c.conn.Write(bsp.NewErr(bsp.ErrReplication).Bytes())
	}
	fmt.Printf("reply:[%v]\n", c.response.String())
	return c.conn.Write(c.response.Bytes())
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
	bconnPool.Put(c)
	return
}
