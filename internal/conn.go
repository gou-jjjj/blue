package internal

import (
	"net"
	"sync"
)

type BConn struct {
	conn net.Conn
	db   uint8
}

var bconnPool = sync.Pool{
	New: func() any {
		return &BConn{}
	},
}

func NewBConn(conn net.Conn) *BConn {
	bconn, ok := bconnPool.Get().(*BConn)
	if !ok {
		return &BConn{conn: conn}
	}
	bconn.conn = conn
	return bconn
}

func (c *BConn) Close() error {
	_ = c.conn.Close()
	c.db = 0
	bconnPool.Put(c)
	return nil
}

func (c *BConn) DBIndex() uint8 {
	return c.db
}

func (c *BConn) SelectDB(index uint8) {
	c.db = index
}

func (c *BConn) Write(b []byte) (int, error) {
	if len(b) == 0 {
		return 0, nil
	}

	return c.conn.Write(b)
}
