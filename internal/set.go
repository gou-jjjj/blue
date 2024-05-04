package internal

import (
	"blue/bsp"
	"blue/datastruct/set"
)

func (db *DB) ExecChainSet(ctx *Context) {
	switch ctx.request.Handle() {
	case bsp.SADD:
		ctx.response = db.sadd(ctx.request)
	case bsp.SPOP:
		ctx.response = db.spop(ctx.request)
	case bsp.SIN:
		ctx.response = db.sin(ctx.request)
	case bsp.SDEL:
		ctx.response = db.sdel(ctx.request)
	case bsp.SGET:
		ctx.response = db.sget(ctx.request)

	default:
		ctx.response = bsp.NewErr(bsp.ErrCommand)
	}
}

func (db *DB) sadd(cmd *bsp.BspProto) bsp.Reply {
	val, ok := db.data.Get(cmd.Key())
	if !ok {
		db.dataCountIncr()
		db.data.Put(cmd.Key(), set.NewSet())
		val, _ = db.data.Get(cmd.Key())
	}

	st, ok := val.(set.Set)
	if !ok {
		return bsp.NewErr(bsp.ErrWrongType, cmd.Key())
	}

	st.Add(cmd.ValueStr())

	err := db.StorageSet(cmd.KeyBytes(), st.String())
	if err != nil {
		return bsp.NewErr(bsp.ErrStorage)
	}

	return bsp.NewInfo(bsp.OK)
}

func (db *DB) spop(cmd *bsp.BspProto) bsp.Reply {
	val, ok := db.data.Get(cmd.Key())
	if !ok {
		return bsp.NewInfo(bsp.NULL)
	}

	st, ok := val.(set.Set)
	if !ok {
		return bsp.NewErr(bsp.ErrWrongType, cmd.Key())
	}

	pop, ok := st.Pop()
	if !ok {
		return bsp.NewInfo(bsp.NULL)
	}

	err := db.StorageSet(cmd.KeyBytes(), st.String())
	if err != nil {
		return bsp.NewErr(bsp.ErrStorage)
	}

	return bsp.NewStr(pop)
}

func (db *DB) sin(cmd *bsp.BspProto) bsp.Reply {
	val, ok := db.data.Get(cmd.Key())
	if !ok {
		return bsp.NewInfo(bsp.False)
	}

	st, ok := val.(set.Set)
	if !ok {
		return bsp.NewErr(bsp.ErrWrongType, cmd.Key())
	}

	if st.ContainsOne(cmd.ValueStr()) {
		return bsp.NewInfo(bsp.True)
	} else {
		return bsp.NewInfo(bsp.False)
	}
}

func (db *DB) sdel(cmd *bsp.BspProto) bsp.Reply {
	val, ok := db.data.Get(cmd.Key())
	if !ok {
		return bsp.NewInfo(bsp.NULL)
	}

	st, ok := val.(set.Set)
	if !ok {
		return bsp.NewErr(bsp.ErrWrongType, cmd.Key())
	}

	st.Remove(cmd.ValueStr())

	err := db.StorageSet(cmd.KeyBytes(), st.String())
	if err != nil {
		return bsp.NewErr(bsp.ErrStorage)
	}

	return bsp.NewInfo(bsp.OK)
}

func (db *DB) sget(cmd *bsp.BspProto) bsp.Reply {
	val, ok := db.data.Get(cmd.Key())
	if !ok {
		return bsp.NewInfo(bsp.NULL)
	}

	st, ok := val.(set.Set)
	if !ok {
		return bsp.NewErr(bsp.ErrWrongType, cmd.Key())
	}

	return bsp.NewStr(st.String())
}
