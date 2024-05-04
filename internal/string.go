package internal

import (
	"blue/bsp"
	str "blue/datastruct/string"
)

func (db *DB) ExecChainString(ctx *Context) {
	switch ctx.request.Handle() {
	case bsp.SET:
		ctx.response = db.set(ctx.request)
	case bsp.GET:
		ctx.response = db.get(ctx.request)
	case bsp.LEN:
		ctx.response = db.len(ctx.request)
	default:
		ctx.response = bsp.NewErr(bsp.ErrCommand)
	}
}

func (db *DB) len(cmd *bsp.BspProto) bsp.Reply {
	v, ok := db.data.Get(cmd.Key())
	if !ok {
		return bsp.NewInfo(bsp.NULL)
	}

	s, ok := v.(str.String)
	if !ok {
		return bsp.NewErr(bsp.ErrWrongType, cmd.Key())
	}

	return bsp.NewNum(int64(s.Len()))
}

func (db *DB) set(cmd *bsp.BspProto) bsp.Reply {
	newString := str.NewString(cmd.ValueStr())
	db.data.Put(cmd.Key(), newString)
	err := db.StorageStr(cmd.KeyBytes(), newString.Get())
	if err != nil {
		return bsp.NewErr(bsp.ErrStorage)
	}

	db.dataCountIncr()
	return bsp.NewInfo(bsp.OK)
}

func (db *DB) get(cmd *bsp.BspProto) bsp.Reply {
	v, ok := db.data.Get(cmd.Key())
	if !ok {
		return bsp.NewInfo(bsp.NULL)
	}

	s, ok := v.(str.String)
	if !ok {
		return bsp.NewErr(bsp.ErrWrongType, cmd.Key())
	}

	return bsp.NewStr(s.Get())
}
