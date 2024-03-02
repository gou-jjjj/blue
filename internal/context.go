package internal

import (
	"context"
	"fmt"
	"net"
	"sync"

	"blue/bsp"
	"blue/common/rand"
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
			db:      1,
			session: rand.RandString(sessionLen),
		}
	},
}

func NewContext(ctx context.Context, conn net.Conn) *Context {
	bconn, ok := bconnPool.Get().(*Context)
	if !ok {
		return &Context{
			db:      1,
			Context: ctx,
			conn:    conn,
			session: rand.RandString(sessionLen),
		}
	}
	bconn.Context = ctx
	bconn.conn = conn
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
	bconnPool.Put(c)
	return
}
