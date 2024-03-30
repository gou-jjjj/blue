package internal

import (
	"blue/bsp"
	"blue/datastruct/list"
	"strconv"
	"strings"
)

var (
	pong = []byte("pong")
)

func (svr *BlueServer) ExecChain(ctx *Context) {
	if !svr.authGuest(ctx) {
		if ctx.request.Handle() == bsp.AUTH {
			svr.auth(ctx)
			return
		}
		ctx.response = bsp.NewErr(bsp.ErrPermissionDenied)
		return
	}

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
	case bsp.AUTH:
		svr.auth(ctx)

	default:
		svr.db[ctx.GetDB()].ExecChain(ctx)
	}
}

func (svr *BlueServer) selected(ctx *Context) {
	ctx.response = bsp.NewStr(ctx.GetDB())
}

func (svr *BlueServer) authGuest(ctx *Context) bool {
	conf := svr.db[0].data
	val, ok := conf.Get("GuestToken")
	if !ok {
		return true
	}
	l := val.(*list.QuickList)
	if l.Len() == 0 {
		return true
	}
	if l.Contains(func(a interface{}) bool {
		return a.(string) == ctx.cliToken
	}) {
		return true
	}
	val, ok = conf.Get("RootToken")
	if !ok {
		return false
	}
	l = val.(*list.QuickList)
	if l.Len() == 0 {
		return false
	}
	return l.Contains(func(a interface{}) bool {
		return a.(string) == ctx.cliToken
	})
}

func (svr *BlueServer) authRoot(ctx *Context) bool {
	conf := svr.db[0].data

	val, ok := conf.Get("RootToken")
	if !ok {
		return false
	}

	l := val.(*list.QuickList)
	if l.Len() == 0 {
		return true
	}

	return l.Contains(func(a interface{}) bool {
		return a.(string) == ctx.cliToken
	})
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

	if dbIndex == 0 && !svr.authRoot(ctx) {
		ctx.response = bsp.NewErr(bsp.ErrPermissionDenied)
		return
	}

	ctx.SetDB(uint8(dbIndex))
	ctx.response = bsp.NewInfo(bsp.OK)
}

func (svr *BlueServer) version(ctx *Context) {
	ctx.response = bsp.NewStr([]byte(version_))
}

func (svr *BlueServer) help(ctx *Context) {
	k := ctx.request.Key()
	upk := strings.ToUpper(k)
	if handleId, ok := bsp.HandleMap2[upk]; !ok {
		ctx.response = bsp.NewErr(bsp.ErrRequestParameter, k)
	} else {
		summary := bsp.CommandsMap[handleId].Summary
		ctx.response = bsp.NewStr(summary)
	}
}

func (svr *BlueServer) ping(ctx *Context) {
	ctx.response = bsp.NewStr(pong)
}

func (svr *BlueServer) exit(ctx *Context) {
	svr.closeClient(ctx)
}
