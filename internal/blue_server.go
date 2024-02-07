package internal

import (
	"blue/bsp"
	"blue/common/strbytes"
	"blue/common/timewheel"
	"context"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

type Exec interface {
	ExecChain(*Context) bool
}

type ServerInter interface {
	Exec
	Handle(context.Context, net.Conn)
	Close()
}

// BlueServer implements tcp.Handler and serves as a redis server
type BlueServer struct {
	activeConn sync.Map // *client -> placeholder
	db         []*DB
	closed     atomic.Int32
}

func NewBlueServer(dbs ...*DB) *BlueServer {
	b := &BlueServer{
		db:         make([]*DB, len(dbs)),
		activeConn: sync.Map{},
	}

	for i := 0; i < len(dbs); i++ {
		b.db[i] = dbs[i]
	}

	return b
}

func (svr *BlueServer) closeClient(client *Context) {
	if client == nil {
		return
	}
	client.Close()
	svr.activeConn.Delete(client)
}

// Handle receives and executes redis commands
func (svr *BlueServer) Handle(ctx context.Context, conn net.Conn) {
	if svr.isClose() {
		_ = conn.Close()
		return
	}

	client := NewContext(ctx, conn)
	svr.activeConn.Store(client, struct{}{})
	defer func() {
		svr.closeClient(client)
	}()

	canCtx, cancelFunc := context.WithCancel(*client)
	bch, errch := bsp.BspProtos(canCtx, conn)
	defer func() {
		cancelFunc()
		close(bch)
		close(errch)
	}()

	for {
		timewheel.Delay(1*time.Minute, client.session, func() {
			svr.closeClient(client)
		})

		select {
		case <-ctx.Done():
			return
		case req := <-bch:
			client.request = req
			client.response = bsp.Reply(nil)
			ok := svr.ExecChain(client)
			if !ok {
				svr.db[client.GetDB()].ExecChain(client)
			}

			client.Reply()

			bsp.BspPool.Put(req)
			continue
		case err := <-errch:
			client.response = err
			client.Reply()
			return
		}
	}
}

func (svr *BlueServer) ExecChain(ctx *Context) bool {
	switch ctx.request.Handle() {
	case bsp.VERSION:
		svr.version(ctx)
	case bsp.USE:
		svr.selectdb(ctx)
	default:
		return false
	}
	return true
}

func (svr *BlueServer) selected(ctx *Context) {
	ctx.response = bsp.NewNum(ctx.GetDB())
}

func (svr *BlueServer) selectdb(ctx *Context) {
	ctx.SetDB(strbytes.Bytes2Uint8(ctx.request.ValueBytes()))
	ctx.response = bsp.NewInfo(bsp.OK)
}

func (svr *BlueServer) version(ctx *Context) {
	ctx.response = bsp.NewStr([]byte("blue v0.1"))
}

func (svr *BlueServer) isClose() bool {
	return svr.closed.Load() == 1
}

// Close stops handler
func (svr *BlueServer) Close() {
	svr.closed.Add(1)

	svr.activeConn.Range(func(key interface{}, _ interface{}) bool {
		client := key.(*Context)
		client.Close()
		return true
	})
}
