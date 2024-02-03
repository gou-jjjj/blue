package internal

import (
	"blue/datastruct/number"
	"bsp"
)

func ExecNumberFunc(db *DB, cmd bsp.BspProto) bsp.Reply {
	switch cmd.Handle() {
	case bsp.NumSet:
		return Set(db, cmd)
	case bsp.NumGet:
		return Get(db, cmd)
	default:
		return bsp.NewErr(bsp.ErrCommand)
	}

}

func Set(db *DB, cmd bsp.BspProto) bsp.Reply {
	db.data.RemoveWithLock(cmd.ValueStr())

	newNumber, err := number.NewNumber(cmd.ValueBytes())
	if err != nil {
		return bsp.NewErr(bsp.ErrWrongType, cmd.ValueStr())
	}

	db.data.Put(cmd.Key(), newNumber)
	err = db.storage.Put(cmd.KeyBytes(), cmd.ValueBytes())
	if err != nil {
		return bsp.NewErr(bsp.ErrStorage)
	}

	return bsp.NewInfo(bsp.OK)
}

func Get(db *DB, cmd bsp.BspProto) bsp.Reply {
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
