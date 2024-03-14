package internal

import "blue/bsp"

func (db *DB) ExecChainList(ctx *Context) {
	switch ctx.request.Handle() {
	case bsp.LSET:
		db.lset(ctx.request)
	case bsp.LGET:
		db.lget(ctx.request)
	case bsp.LLEN:
		db.llen(ctx.request)
	case bsp.LPUSH:
		db.lpush(ctx.request)
	case bsp.LPOP:
		db.lpop(ctx.request)
	case bsp.RPUSH:
		db.rpush(ctx.request)
	case bsp.RPOP:
		db.rpop(ctx.request)
	default:
		ctx.response = bsp.NewErr(bsp.ErrCommand)
	}
}

func (db *DB) lset(cmd *bsp.BspProto) {

}

func (db *DB) lget(cmd *bsp.BspProto) {}

func (db *DB) llen(cmd *bsp.BspProto) {}

func (db *DB) lpush(cmd *bsp.BspProto) {}

func (db *DB) lpop(cmd *bsp.BspProto) {}

func (db *DB) rpush(cmd *bsp.BspProto) {}

func (db *DB) rpop(cmd *bsp.BspProto) {}
