package internal

import (
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
	config := defaultDBConfig
	for _, opt := range opts {
		opt(&config)
	}

	db := &DB{
		index: config.Index,
		data:  dict.MakeConcurrent(config.DataDictSize),
		rw:    &sync.RWMutex{},
	}

	if config.SetStorage {
		options := config.StorageOptions

		if _, err := os.Stat(config.StorageOptions.DirPath); err != nil {
			err = os.Mkdir(config.StorageOptions.DirPath, 777)
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
	}

	if config.InitData != nil {
		for k, v := range config.InitData {
			db.data.Put(k, v)
		}
	}

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
		db.ExecChainList(ctx)
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

func (db *DB) del(key [][]byte) bsp.Reply {
	for i := range key {
		db.data.Remove(string(key[i]))
		err := db.StorageDelete(key[i])
		if err != nil {
			return bsp.NewErr(bsp.ErrStorage)
		}
	}
	return bsp.NewInfo(bsp.OK)
}
