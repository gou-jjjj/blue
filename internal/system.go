package internal

import (
	"blue/bsp"
	"strconv"
)

var (
	pong = []byte("pong")
)

func (svr *BlueServer) ExecChain(ctx *Context) {
	switch ctx.request.Handle() {
	case bsp.VERSION:
		svr.version(ctx)
	case bsp.SELECT:
		if ctx.request.Key() != "" {
			svr.selectdb(ctx)
		} else {
			svr.selected(ctx)
		}
	case bsp.HELP:
		svr.help(ctx)
	case bsp.PING:
		svr.ping(ctx)
	case bsp.EXIT:
		svr.exit(ctx)
	default:
		svr.db[ctx.GetDB()].ExecChain(ctx)
	}
}

func (svr *BlueServer) selected(ctx *Context) {
	ctx.response = bsp.NewStr(ctx.GetDB())
}

func (svr *BlueServer) selectdb(ctx *Context) {
	dbIndex, err := strconv.Atoi(ctx.request.Key())
	if err != nil {
		ctx.response = bsp.NewErr(bsp.ErrRequestParameter)
		return
	}

	if dbIndex < 0 || dbIndex >= len(svr.db) {
		ctx.response = bsp.NewErr(bsp.ErrRequestParameter)
		return
	}

	ctx.SetDB(uint8(dbIndex))
	ctx.response = bsp.NewInfo(bsp.OK)
}

func (svr *BlueServer) version(ctx *Context) {
	ctx.response = bsp.NewStr([]byte(version_))
}

func (svr *BlueServer) help(ctx *Context) {
	v := ctx.request.HandleInfo().Summary
	ctx.response = bsp.NewStr([]byte(v))
}

func (svr *BlueServer) ping(ctx *Context) {
	ctx.response = bsp.NewStr(pong)
}

func (svr *BlueServer) exit(ctx *Context) {
	svr.closeClient(ctx)
}
