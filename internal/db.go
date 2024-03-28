package internal

import (
	"blue/config"
	"blue/log"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"blue/bsp"
	"blue/datastruct"
	"blue/datastruct/dict"

	"github.com/rosedblabs/rosedb/v2"
)

type DBConfig struct {
	StorageOptions rosedb.Options

	InitData     map[string]interface{}
	SetStorage   bool
	DataDictSize int
	Index        int
}

type DbOption func(*DBConfig)

var defaultDBConfig = DBConfig{
	StorageOptions: rosedb.DefaultOptions,
	DataDictSize:   1024,
	Index:          0,
}

type DB struct {
	index int

	data    *dict.ConcurrentDict
	storage *rosedb.DB
	rw      *sync.RWMutex
}

func NewDB(opts ...DbOption) *DB {
	// 指定选项
	dbConfig := defaultDBConfig
	for _, opt := range opts {
		opt(&dbConfig)
	}

	db := &DB{
		index: dbConfig.Index,
		data:  dict.MakeConcurrent(dbConfig.DataDictSize),
		rw:    &sync.RWMutex{},
	}

	if dbConfig.InitData != nil {
		for k, v := range dbConfig.InitData {
			db.data.Put(k, v)
		}
	}

	db.InitStorage(dbConfig)
	log.Info(fmt.Sprintf("db{index[%d] initdata:[%v] initStorage[%v] }",
		db.index, len(dbConfig.InitData) != 0, dbConfig.SetStorage))
	return db
}

func (db *DB) InitStorage(dbConfig DBConfig) {
	if dbConfig.SetStorage {
		options := dbConfig.StorageOptions

		if _, err := os.Stat(dbConfig.StorageOptions.DirPath); err != nil {
			err = os.Mkdir(dbConfig.StorageOptions.DirPath, 777)
			if err != nil {
				panic(err)
			}
		}

		storage, err := rosedb.Open(options)
		if err != nil {
			panic(err)
		}

		db.storage = storage

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

		config.StorageInitSuccess()
	}
}

func (db *DB) ExecChain(ctx *Context) {
	switch ctx.request.Type() {
	case bsp.TypeDB:
		db.ExecChainDB(ctx)
	case bsp.TypeNumber:
		db.ExecChainNumber(ctx)
	case bsp.TypeString:
		db.ExecChainString(ctx)
	case bsp.TypeList:
		db.ExecChainList(ctx)
	case bsp.TypeSet:
		db.ExecChainSet(ctx)
	case bsp.TypeJson:

	default:
		fmt.Printf("db:[%+b]", ctx.request.Type())
		ctx.response = bsp.NewErr(bsp.ErrCommand)
	}
}

func (db *DB) ExecChainDB(ctx *Context) {
	switch ctx.request.Handle() {
	case bsp.DEL:
		ctx.response = db.del(ctx.request)
	case bsp.EXPIRE:
		ctx.response = db.expire(ctx.request)
	case bsp.KVS:
		ctx.response = db.kvs(ctx.request)

	default:
		ctx.response = bsp.NewErr(bsp.ErrCommand)
	}
}

func (db *DB) StoragePut(key []byte, value []byte) error {
	if db.storage == nil {
		return nil
	}

	return db.storage.Put(key, value)
}

func (db *DB) StorageDelete(key []byte) error {
	if db.storage == nil {
		return nil
	}

	return db.storage.Delete(key)
}

func (db *DB) RangeKV() string {
	if db.data.Len() == 0 {
		return ""
	}

	builder := strings.Builder{}
	db.data.ForEach(func(key string, val interface{}) bool {
		builder.WriteString(fmt.Sprintf("%s: %s\n", key, val.(datastruct.Value).Value()))
		return true
	})
	return builder.String()[:builder.Len()-1]
}

func (db *DB) del(ctx *bsp.BspProto) bsp.Reply {
	key := ctx.Values()
	for i := range key {
		db.data.Remove(string(key[i]))
		err := db.StorageDelete(key[i])
		if err != nil {
			return bsp.NewErr(bsp.ErrStorage)
		}
	}
	return bsp.NewInfo(bsp.OK)
}

func (db *DB) expire(ctx *bsp.BspProto) bsp.Reply {
	return bsp.NewInfo(bsp.OK)
}

func (db *DB) kvs(ctx *bsp.BspProto) bsp.Reply {
	kv := db.RangeKV()

	if kv == "" {
		return bsp.NewInfo(bsp.NULL)
	}

	return bsp.NewStr(kv)
}
