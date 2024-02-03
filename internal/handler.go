package internal

import (
	"blue/common/timewheel"
	"bsp"
	"context"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

type HandlerInter interface {
	Handle(context.Context, net.Conn)
	Close() error
}

// BlueServer implements tcp.Handler and serves as a redis server
type BlueServer struct {
	activeConn sync.Map // *client -> placeholder
	db         []*DB
	closed     int32
}

func NewBlueServer(dbs int) *BlueServer {
	b := &BlueServer{
		db:         make([]*DB, dbs),
		activeConn: sync.Map{},
	}

	for i := 0; i < dbs; i++ {
		b.db[i] = NewDB(i)
	}

	return b
}

func (h *BlueServer) closeClient(client *BConn) {
	_ = client.Close()
	h.activeConn.Delete(client)
}

// Handle receives and executes redis commands
func (h *BlueServer) Handle(ctx context.Context, conn net.Conn) {
	if h.isClose() {
		_ = conn.Close()
		return
	}

	client := NewBConn(conn)
	h.activeConn.Store(client, struct{}{})

	canCtx, cancelFunc := context.WithCancel(ctx)
	bch, errch := bsp.BspProtos(canCtx, conn)
	for {
		timewheel.Delay(1*time.Minute, client.session, func() {
			h.closeClient(client)
		})
		select {
		case <-ctx.Done():
			h.closeClient(client)
			return
		case req := <-bch:
			if req.Handle() == bsp.SystemExit {
				client.Write(bsp.NewStr([]byte("good bye")).Bytes())
				h.closeClient(client)
				cancelFunc()
				close(bch)
				close(errch)
				return
			}
			errEsp := h.db.Exec(client, req)
			if errEsp != nil {
				h.closeClient(client)
			}

		case err := <-errch:
			client.Write(err.Bytes())
			h.closeClient(client)
			return
		}

		timewheel.Cancel(client.session)
	}
}

// Close stops handler
func (h *BlueServer) Close() error {
	atomic.AddInt32(&h.closed, 1)

	h.activeConn.Range(func(key interface{}, _ interface{}) bool {
		client := key.(*BConn)
		_ = client.Close()
		return true
	})
	return nil
}

func (db *BlueServer) selected() bsp.Reply {
	return bsp.NewNum(db.index)
}

func (db *BlueServer) selectdb(b *bsp.BspProto) bsp.Reply {
	return bsp.NewInfo(bsp.OK)
}

func (db *BlueServer) SystemExec(cmd *bsp.BspProto) (b bsp.Reply) {
	switch cmd.Handle() {
	case bsp.SystemVersion:
		b = db.version()
	default:
		return bsp.NewErr(bsp.ErrCommand)
	}
	return
}

func (db *BlueServer) version() bsp.Reply {
	return bsp.NewStr([]byte("blue v0.1"))
}

func (h *BlueServer) isClose() bool {
	return atomic.LoadInt32(&h.closed) == 1
}
