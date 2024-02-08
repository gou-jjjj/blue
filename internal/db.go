package internal

import (
	"blue/bsp"
	"blue/datastruct/dict"
	"blue/datastruct/number"
	"fmt"
	"github.com/rosedblabs/rosedb/v2"
	"os"
	"sync"
	"time"
)

type DBConfig struct {
	rosedb.Options

	DataDictSize int
	index        int
}

type DbOption func(*DBConfig)

var defaultDBConfig = DBConfig{
	Options:      rosedb.DefaultOptions,
	DataDictSize: 1024,
	index:        0,
}

type DB struct {
	index int

	data    *dict.ConcurrentDict
	storage *rosedb.DB
	rw      *sync.RWMutex
}

func NewDB(opts ...DbOption) *DB {
	// 指定选项
	config := defaultDBConfig
	for _, opt := range opts {
		opt(&config)
	}

	options := config.Options

	if _, err := os.Stat(config.DirPath); err != nil {
		err = os.Mkdir(config.DirPath, 777)
		if err != nil {
			panic(err)
		}
	}

	storage, err := rosedb.Open(options)
	if err != nil {
		panic(err)
	}

	db := &DB{
		index:   config.index,
		data:    dict.MakeConcurrent(config.DataDictSize),
		storage: storage,
		rw:      &sync.RWMutex{},
	}

	storage.Ascend(func(k []byte, v []byte) (bool, error) {
		ttl, err := storage.TTL(k)
		if err != nil {
			db.data.Put(string(k), DataEntity{
				Val:    v,
				Expire: 0,
			})
		}

		db.data.Put(string(k), DataEntity{
			Val:    v,
			Expire: uint64(time.Now().Second()) + uint64(ttl),
		})

		return true, nil
	})

	return db
}

func (db *DB) ExecChain(ctx *Context) bool {
	switch ctx.request.Type() {
	case bsp.TypeDB:
		db.ExecChainDB(ctx)
	case bsp.TypeNumber:
		db.ExecChainNumber(ctx)
	case bsp.TypeString:
		db.ExecChainString(ctx)
	case bsp.TypeList:

	case bsp.TypeSet:

	case bsp.TypeJson:

	default:
		fmt.Printf("db:[%+b]", ctx.request.Type())
		ctx.response = bsp.NewErr(bsp.ErrCommand)
	}

	return true
}

func (db *DB) ExecChainDB(ctx *Context) {
	switch ctx.request.Handle() {
	case bsp.DEL:
		ctx.response = db.del(ctx.request.Values())
	default:
		ctx.response = bsp.NewErr(bsp.ErrCommand)
	}
}

func (db *DB) ExecChainNumber(ctx *Context) {
	switch ctx.request.Handle() {
	case bsp.NSET:
		ctx.response = db.nset(ctx.request)
	case bsp.NGET:
		ctx.response = db.nget(ctx.request)
	default:
		ctx.response = bsp.NewErr(bsp.ErrCommand)
	}

}

func (db *DB) ExecChainString(ctx *Context) {}

func (db *DB) del(key [][]byte) bsp.Reply {
	for i := range key {
		db.data.Remove(string(key[i]))
		err := db.storage.Delete(key[i])
		if err != nil {
			return bsp.NewErr(bsp.ErrStorage)
		}
	}
	return bsp.NewInfo(bsp.OK)
}

func (db *DB) nset(cmd *bsp.BspProto) bsp.Reply {
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
