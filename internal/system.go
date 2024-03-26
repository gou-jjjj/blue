package internal

import (
	"blue/bsp"
	"strconv"
)

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

func (svr *BlueServer) kvs(ctx *Context) {
	kv := svr.db[ctx.GetDB()].RangeKV()

	if kv == "" {
		ctx.response = bsp.NewInfo(bsp.NULL)
	} else {
		ctx.response = bsp.NewStr(kv)
	}
}
