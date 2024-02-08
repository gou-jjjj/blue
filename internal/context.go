package internal

import (
	"blue/bsp"
	"blue/common/rand"
	"context"
	"fmt"
	"net"
	"sync"
)

const (
	sessionLen = 10
)

type Context struct {
	context.Context
	conn     net.Conn
	db       uint8
	session  string
	request  *bsp.BspProto
	response bsp.Reply
	nextExec Exec
}

var bconnPool = sync.Pool{
	New: func() interface{} {
		return &Context{
			session: rand.RandString(sessionLen),
		}
	},
}

func NewContext(ctx context.Context, conn net.Conn) *Context {
	bconn, ok := bconnPool.Get().(*Context)
	if !ok {
		return &Context{
			Context: ctx,
			conn:    conn,
			session: rand.RandString(sessionLen),
		}
	}
	bconn.Context = ctx
	bconn.conn = conn
	fmt.Printf("%#v\n", bconn)
	return bconn
}

func (c *Context) SetNext(next Exec) {
	c.nextExec = next
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
	fmt.Printf("req:[%v]\n", c.response.String())
	return c.conn.Write(c.response.Bytes())
}

func (c *Context) Close() {
	fmt.Println("context close")
	if c.conn != nil {
		_ = c.conn.Close()
	}
	c.db = 0
	c.Context = nil
	c.nextExec = nil
	c.request = nil
	c.response = nil
	bconnPool.Put(c)
	return
}
