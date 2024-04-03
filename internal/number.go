package internal

import (
	"blue/bsp"
	"blue/datastruct/number"
	"strconv"
)

func (db *DB) ExecChainNumber(ctx *Context) {
	switch ctx.request.Handle() {
	case bsp.NSET:
		ctx.response = db.nset(ctx.request)
	case bsp.NGET:
		ctx.response = db.nget(ctx.request)
	case bsp.INCR:
		ctx.response = db.incr(ctx.request)
	default:
		ctx.response = bsp.NewErr(bsp.ErrCommand)
	}

}

func (db *DB) nset(cmd *bsp.BspProto) bsp.Reply {

	newNumber, err := number.NewNumber(cmd.ValueBytes())
	if err != nil {
		return bsp.NewErr(bsp.ErrWrongType, cmd.ValueStr())
	}

	db.data.Put(cmd.Key(), newNumber)
	err = db.StoragePut(cmd.KeyBytes(), cmd.ValueBytes())
	if err != nil {
		return bsp.NewErr(bsp.ErrStorage)
	}

	return bsp.NewInfo(bsp.OK)
}

func (db *DB) nget(cmd *bsp.BspProto) bsp.Reply {
	v, ok := db.data.Get(cmd.Key())
	if !ok {
		return bsp.NewInfo(bsp.NULL)
	}

	n, ok := v.(number.Number)
	if !ok {
		return bsp.NewErr(bsp.ErrWrongType, cmd.Key())
	}

	return bsp.NewNum(n.Get())
}

func (db *DB) incr(cmd *bsp.BspProto) bsp.Reply {
	v, ok := db.data.Get(cmd.Key())
	if !ok {
		return bsp.NewInfo(bsp.NULL)
	}

	n, ok := v.(number.Number)
	if !ok {
		return bsp.NewErr(bsp.ErrWrongType, cmd.Key())
	}

	incr := n.Incr()

	err := db.StoragePut(cmd.KeyBytes(), []byte(strconv.FormatInt(incr, 10)))
	if err != nil {
		return bsp.NewErr(bsp.ErrStorage)
	}

	return bsp.NewNum(incr)
}
