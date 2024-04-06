package internal

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"blue/bsp"
	print2 "blue/common/print"
	"blue/common/timewheel"
	"blue/config"
	"blue/datastruct"
	"blue/datastruct/dict"
	str "blue/datastruct/string"
	"blue/log"

	"github.com/rosedblabs/rosedb/v2"
)

// DBConfig 存储数据库配置
type DBConfig struct {
	StorageOptions rosedb.Options         // 存储选项
	InitData       map[string]interface{} // 初始化数据
	DataDictSize   int                    // 数据字典大小
	Index          int                    // 数据库索引
}

// DbOption 定义数据库配置选项的函数类型
type DbOption func(*DBConfig)

var defaultDBConfig = DBConfig{
	StorageOptions: rosedb.DefaultOptions,
	DataDictSize:   1024,
	Index:          0,
}

// DB 表示一个数据库实例
type DB struct {
	index int

	data    *dict.ConcurrentDict
	storage *rosedb.DB
	rw      *sync.RWMutex
}

// NewDB 创建一个新的数据库实例
// opts: 一个或多个DbOption函数用于定制数据库配置
func NewDB(opts ...DbOption) *DB {
	// 应用配置选项
	dbConfig := defaultDBConfig
	for _, opt := range opts {
		opt(&dbConfig)
	}

	db := &DB{
		index: dbConfig.Index,
		data:  dict.MakeConcurrent(dbConfig.DataDictSize),
		rw:    &sync.RWMutex{},
	}

	// 如果提供了初始化数据，则加载到内存中
	if dbConfig.InitData != nil {
		for k, v := range dbConfig.InitData {
			db.data.Put(k, v)
		}
	}

	// 初始化存储
	initLen := db.InitStorage(dbConfig)
	log.Info(fmt.Sprintf("db{index[%d] initdata:[%v] initStorage[%v] }",
		db.index,
		len(dbConfig.InitData) != 0 || (initLen != 0),
		config.OpenStorage(strconv.Itoa(db.index))))
	return db
}

// InitStorage 初始化存储引擎
// dbConfig: 数据库配置
func (db *DB) InitStorage(dbConfig DBConfig) int {
	var l int
	if db.index != 0 && config.OpenStorage(strconv.Itoa(db.index)) {
		options := dbConfig.StorageOptions
		options.Sync = true

		// 检查并创建存储目录
		if _, err := os.Stat(dbConfig.StorageOptions.DirPath); errors.Is(err, os.ErrNotExist) {
			err = os.MkdirAll(dbConfig.StorageOptions.DirPath, os.ModePerm)
			if err != nil {
				panic(err)
			}
		}

		// 打开或创建存储
		storage, err := rosedb.Open(options)
		if err != nil {
			panic(err)
		}

		db.storage = storage

		// 从存储中加载数据到内存
		storage.Ascend(func(k []byte, v []byte) (bool, error) {
			l++
			db.data.Put(string(k), str.NewString(string(v)))
			return true, nil
		})

		print2.StorageInitSuccess(db.index)
	}

	return l
}

// ExecChain 执行链式操作
// ctx: 操作上下文
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

// ExecChainDB 处理数据库操作请求
// ctx: 操作上下文
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

// StoragePut 存储数据
// key: 键
// value: 值
func (db *DB) StoragePut(key []byte, value []byte) error {
	if db.storage == nil {
		return nil
	}

	return db.storage.Put(key, value)
}

// StorageDelete 删除存储的数据
// key: 键
func (db *DB) StorageDelete(key []byte) error {
	if db.storage == nil {
		return nil
	}

	return db.storage.Delete(key)
}

// RangeKV 获取所有键值对的字符串表示
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

// del 删除键值对
// ctx: 操作上下文
func (db *DB) del(ctx *bsp.BspProto) bsp.Reply {
	db.data.Remove(ctx.Key())
	err := db.StorageDelete(ctx.KeyBytes())
	if err != nil {
		return bsp.NewErr(bsp.ErrStorage)
	}

	return bsp.NewInfo(bsp.OK)
}

// expire 设置键的过期时间
// ctx: 操作上下文
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

// kvs 返回数据库中所有键值对的字符串表示
// ctx: 操作上下文
func (db *DB) kvs(ctx *bsp.BspProto) bsp.Reply {
	kv := db.RangeKV()

	if kv == "" {
		return bsp.NewInfo(bsp.NULL)
	}

	return bsp.NewStr(kv)
}

// dbsize 返回数据库中键值对的数量
// ctx: 操作上下文
func (db *DB) dbsize(ctx *bsp.BspProto) bsp.Reply {
	return bsp.NewNum(db.data.Len())
}

// type_ 返回键的数据类型
// ctx: 操作上下文
func (db *DB) type_(ctx *bsp.BspProto) bsp.Reply {
	val, ok := db.data.Get(ctx.Key())
	if !ok {
		return bsp.NewInfo(bsp.NULL)
	}

	return bsp.NewStr(val.(datastruct.Type).Type())
}
