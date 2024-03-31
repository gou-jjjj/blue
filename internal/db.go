package internal

import (
	str "blue/datastruct/string"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"blue/bsp"
	"blue/common/timewheel"
	"blue/config"
	"blue/datastruct"
	"blue/datastruct/dict"
	"blue/log"

	"github.com/rosedblabs/rosedb/v2"
)

type DBConfig struct {
	StorageOptions rosedb.Options
	InitData       map[string]interface{}
	DataDictSize   int
	Index          int
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

	initLen := db.InitStorage(dbConfig)
	log.Info(fmt.Sprintf("db{index[%d] initdata:[%v] initStorage[%v] }",
		db.index,
		len(dbConfig.InitData) != 0 || (initLen != 0),
		config.StoCfg.OpenStorage(strconv.Itoa(db.index))))
	return db
}

func (db *DB) InitStorage(dbConfig DBConfig) int {
	var l int
	if db.index != 0 && config.StoCfg.OpenStorage(strconv.Itoa(db.index)) {
		options := dbConfig.StorageOptions

		if _, err := os.Stat(dbConfig.StorageOptions.DirPath); errors.Is(err, os.ErrNotExist) {
			err = os.MkdirAll(dbConfig.StorageOptions.DirPath, os.ModePerm)
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
			l++
			db.data.Put(string(k), str.NewString(string(v)))
			return true, nil
		})

		config.StorageInitSuccess(db.index)
	}

	return l
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
	case bsp.DBSIZE:
		ctx.response = db.dbsize(ctx.request)
	case bsp.TYPE:
		ctx.response = db.type_(ctx.request)

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
	db.data.Remove(ctx.Key())
	err := db.StorageDelete(ctx.KeyBytes())
	if err != nil {
		return bsp.NewErr(bsp.ErrStorage)
	}

	return bsp.NewInfo(bsp.OK)
}

func (db *DB) expire(ctx *bsp.BspProto) bsp.Reply {
	key := ctx.Key()
	ttl := ctx.ValueStr()

	if _, ok := db.data.Get(key); !ok {
		return bsp.NewInfo(bsp.NULL)
	}

	ttlInt, err := strconv.Atoi(ttl)
	if err != nil {
		return bsp.NewErr(bsp.ErrRequestParameter, ttl)
	}
	timewheel.Delay(time.Duration(ttlInt)*time.Second, key, func() {
		if db != nil {
			db.data.Remove(key)
		}
	})

	return bsp.NewInfo(bsp.OK)
}

func (db *DB) kvs(ctx *bsp.BspProto) bsp.Reply {
	kv := db.RangeKV()

	if kv == "" {
		return bsp.NewInfo(bsp.NULL)
	}

	return bsp.NewStr(kv)
}

func (db *DB) dbsize(ctx *bsp.BspProto) bsp.Reply {
	return bsp.NewNum(db.data.Len())
}

func (db *DB) type_(ctx *bsp.BspProto) bsp.Reply {
	val, ok := db.data.Get(ctx.Key())
	if !ok {
		return bsp.NewInfo(bsp.NULL)
	}

	return bsp.NewStr(val.(datastruct.Type).Type())
}
