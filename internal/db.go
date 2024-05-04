package internal

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"blue/bsp"
	print2 "blue/common/print"
	"blue/common/timewheel"
	"blue/config"
	"blue/datastruct"
	"blue/datastruct/dict"
	"blue/datastruct/list"
	"blue/datastruct/number"
	"blue/datastruct/set"
	str "blue/datastruct/string"
	"blue/log"

	"github.com/rosedblabs/rosedb/v2"
)

const (
	StorageNum  = ":1"
	StorageStr  = ":2"
	StorageList = ":3"
	StorageSet  = ":4"
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

	data      *dict.ConcurrentDict
	storage   *rosedb.DB
	rw        *sync.RWMutex
	dataCount atomic.Int64
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
		index:     dbConfig.Index,
		data:      dict.MakeConcurrent(dbConfig.DataDictSize),
		rw:        &sync.RWMutex{},
		dataCount: atomic.Int64{},
	}

	// 如果提供了初始化数据，则加载到内存中
	if dbConfig.InitData != nil {
		for k, v := range dbConfig.InitData {
			db.data.Put(k, v)
		}
	}

	// 初始化存储
	initLen := db.InitStorage(dbConfig)
	log.Info(fmt.Sprintf("db{index[%d] initdata:[%v] initStorage[%v]{%d} }",
		db.index,
		len(dbConfig.InitData) != 0 || (initLen != 0),
		config.OpenStorage(strconv.Itoa(db.index)), initLen))
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
			db.dataCountIncr()
			db.loadMem(k, v)
			return true, nil
		})

		print2.StorageInitSuccess(db.index)
	}

	return l
}

func (db *DB) loadMem(k []byte, v []byte) {
	slt := len(k) - 2
	ty := k[slt:]
	k1 := string(k[:slt])

	switch string(ty) {
	case StorageNum:
		newNumber, _ := number.NewNumber(v)
		db.data.Put(k1, newNumber)

	case StorageList:
		newList := list.NewQuickList()
		by := bytes.Split(v, []byte{' '})
		for i := range by {
			newList.Add(string(by[i]))
		}

		db.data.Put(k1, newList)

	case StorageStr:
		newStr := str.NewString(string(v))
		db.data.Put(k1, newStr)

	case StorageSet:
		newSet := set.NewSet()

		by := bytes.Split(v, []byte{' '})
		for i := range by {
			newSet.Add(string(by[i]))
		}
		db.data.Put(k1, newSet)

	default:
		panic("init storage unknown type")
	}
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
	v, ok := db.data.Get(ctx.Key())
	if !ok {
		return bsp.NewInfo(bsp.NULL)
	}

	objType := v.(datastruct.Type)
	err := db.StorageDelete(objType.Type(), ctx.KeyBytes())
	if err != nil {
		return bsp.NewErr(bsp.ErrStorage)
	}

	db.data.Remove(ctx.Key())
	db.dataCountDecr()
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
			db.dataCountDecr()
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
	count := db.dataCountLoad()
	s := strconv.Itoa(count)
	return bsp.NewStr(s)
}

// type_ 返回键的数据类型
// ctx: 操作上下文
func (db *DB) type_(ctx *bsp.BspProto) bsp.Reply {
	val, ok := db.data.Get(ctx.Key())
	if !ok {
		return bsp.NewInfo(bsp.NULL)
	}
	obj, ok := val.(datastruct.BlueObj)
	if !ok {
		return bsp.NewErr(bsp.ErrWrongType)
	}

	return bsp.NewStr(obj.GetType())
}

func (db *DB) dataCountIncr() {
	db.dataCount.Add(1)
}

func (db *DB) dataCountDecr() {
	db.dataCount.Add(-1)
}

func (db *DB) dataCountLoad() int {
	return int(db.dataCount.Load())
}

// storagePut 存储数据
// key: 键
// value: 值
func (db *DB) storagePut(key []byte, value []byte) error {

	return db.storage.Put(key, value)
}

// StorageDelete 删除存储的数据
// key: 键
func (db *DB) StorageDelete(ty string, key []byte) error {
	if db.storage == nil {
		return nil
	}

	switch ty {
	case datastruct.StringType:
		key = append(key, StorageStr...)
	case datastruct.NumberType:
		key = append(key, StorageNum...)
	case datastruct.ListType:
		key = append(key, StorageList...)
	case datastruct.SetType:
		key = append(key, StorageSet...)

	default:

	}

	err := db.storage.Delete(key)
	if err != nil {
		return err
	}

	return db.storage.Sync()
}

// StorageNum 存储数字
func (db *DB) StorageNum(key []byte, value int64) error {
	if db.storage == nil {
		return nil
	}
	key = append(key, StorageNum...)

	return db.storagePut(key, []byte(strconv.FormatInt(value, 10)))
}

// StorageStr 存储字符串
func (db *DB) StorageStr(key []byte, value string) error {
	if db.storage == nil {
		return nil
	}
	key = append(key, StorageStr...)

	return db.storagePut(key, []byte(value))
}

// StorageList 存储列表
func (db *DB) StorageList(key []byte, value string) error {
	if db.storage == nil {
		return nil
	}
	key = append(key, StorageList...)

	return db.storagePut(key, []byte(value))
}

// StorageSet 存储集合
func (db *DB) StorageSet(key []byte, value string) error {
	if db.storage == nil {
		return nil
	}
	key = append(key, StorageSet...)

	return db.storagePut(key, []byte(value))
}
