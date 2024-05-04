package internal

import (
	"blue/bsp"
	"blue/datastruct/list"
)

func (db *DB) ExecChainList(ctx *Context) {
	switch ctx.request.Handle() {
	case bsp.LSET:
		ctx.response = db.lset(ctx.request)
	case bsp.LGET:
		ctx.response = db.lget(ctx.request)
	case bsp.LLEN:
		ctx.response = db.llen(ctx.request)
	case bsp.LPUSH:
		ctx.response = db.lpush(ctx.request)
	case bsp.LPOP:
		ctx.response = db.lpop(ctx.request)
	case bsp.RPUSH:
		ctx.response = db.rpush(ctx.request)
	case bsp.RPOP:
		ctx.response = db.rpop(ctx.request)
	default:
		ctx.response = bsp.NewErr(bsp.ErrCommand)
	}
}

func (db *DB) lset(cmd *bsp.BspProto) bsp.Reply {
	newlist := list.NewQuickList()
	db.data.Put(cmd.Key(), newlist)

	db.dataCountIncr()

	//err := db.StorageList(cmd.KeyBytes(), "")
	//if err != nil {
	//	return bsp.NewErr(bsp.ErrStorage)
	//}

	return bsp.NewInfo(bsp.OK)
}

func (db *DB) lget(cmd *bsp.BspProto) bsp.Reply {
	quickList, ok := db.data.Get(cmd.Key())
	if !ok {
		return bsp.NewInfo(bsp.NULL)
	}
	l, ok := quickList.(*list.QuickList)
	if !ok {
		return bsp.NewErr(bsp.ErrWrongType, cmd.Key())
	}

	if l.Len() == 0 {
		return bsp.NewInfo(bsp.NULL)
	}

	return bsp.NewStr(l.String())
}

func (db *DB) llen(cmd *bsp.BspProto) bsp.Reply {
	quickList, ok := db.data.Get(cmd.Key())
	if !ok {
		return bsp.NewInfo(bsp.NULL)

	}
	l, ok := quickList.(*list.QuickList)
	if !ok {
		return bsp.NewErr(bsp.ErrWrongType, cmd.Key())
	}
	return bsp.NewNum(l.Len())
}

func (db *DB) lpush(cmd *bsp.BspProto) bsp.Reply {
	quickList, ok := db.data.Get(cmd.Key())
	if !ok {
		db.lset(cmd)
		quickList, _ = db.data.Get(cmd.Key())
	}
	l, ok := quickList.(*list.QuickList)
	if !ok {
		return bsp.NewErr(bsp.ErrWrongType, cmd.Key())
	}
	for i := range cmd.Values() {
		l.Insert(0, string(cmd.Values()[i]))
	}

	err := db.StorageList(cmd.KeyBytes(), l.String())
	if err != nil {
		return bsp.NewErr(bsp.ErrStorage)
	}

	return bsp.NewInfo(bsp.OK)
}

func (db *DB) lpop(cmd *bsp.BspProto) bsp.Reply {
	quickList, ok := db.data.Get(cmd.Key())
	if !ok {
		return bsp.NewInfo(bsp.NULL)
	}
	l, ok := quickList.(*list.QuickList)
	if !ok {
		return bsp.NewErr(bsp.ErrWrongType, cmd.Key())
	}

	v := l.Remove(0)

	err := db.StorageList(cmd.KeyBytes(), l.String())
	if err != nil {
		return bsp.NewErr(bsp.ErrStorage)
	}

	return bsp.NewStr(v)
}

func (db *DB) rpush(cmd *bsp.BspProto) bsp.Reply {
	quickList, ok := db.data.Get(cmd.Key())
	if !ok {
		db.lset(cmd)
		quickList, _ = db.data.Get(cmd.Key())
	}
	l, ok := quickList.(*list.QuickList)
	if !ok {
		return bsp.NewErr(bsp.ErrWrongType, cmd.Key())
	}
	for i := range cmd.Values() {
		l.Add(string(cmd.Values()[i]))
	}

	err := db.StorageList(cmd.KeyBytes(), l.String())
	if err != nil {
		return bsp.NewErr(bsp.ErrStorage)
	}

	return bsp.NewInfo(bsp.OK)
}

func (db *DB) rpop(cmd *bsp.BspProto) bsp.Reply {
	quickList, ok := db.data.Get(cmd.Key())
	if !ok {
		return bsp.NewInfo(bsp.NULL)
	}
	l, ok := quickList.(*list.QuickList)
	if !ok {
		return bsp.NewErr(bsp.ErrWrongType, cmd.Key())
	}

	v := l.RemoveLast()

	err := db.StorageList(cmd.KeyBytes(), l.String())
	if err != nil {
		return bsp.NewErr(bsp.ErrStorage)
	}

	return bsp.NewStr(v)
}
