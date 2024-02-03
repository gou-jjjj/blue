package internal

import (
	"blue/datastruct/dict"
	"bsp"
	"fmt"
	"github.com/rosedblabs/rosedb/v2"
	"os"
	"sync"
	"time"
)

const (
	dataDictSize = 1 << 16
	ttlDictSize  = 1 << 10
)

type DB struct {
	index int

	data    *dict.ConcurrentDict
	storage *rosedb.DB
	rw      *sync.RWMutex
}

type ExecFunc func(*DB, bsp.BspProto) bsp.Reply

func NewDB(index int) *DB {
	// 指定选项
	options := rosedb.DefaultOptions
	stoPath := fmt.Sprintf("./storage/data/%d", index)
	err := os.Mkdir(stoPath, 777)
	if err != nil {
		panic(err)
	}

	options.DirPath = stoPath
	storage, err := rosedb.Open(options)
	if err != nil {
		panic(err)
	}

	db := &DB{
		data:    dict.MakeConcurrent(dataDictSize),
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

func (db *DB) Exec(c *BConn, cmd *bsp.BspProto) error {
	var resp bsp.Reply
	switch cmd.Type() {
	case bsp.TypeNumber:
		resp = ExecNumberFunc(db, *cmd)
	case bsp.TypeString:
		//resp = ExecStringFunc(db, *cmd)
	default:
		resp = bsp.NewErr(bsp.ErrCommand)
	}

	_, err := c.Write(resp.Bytes())
	return err
}

func (db *DB) Put(key string, val []byte) error {
	db.data.Put(key, DataEntity{
		Val:    val,
		Expire: 0,
	})
	return db.storage.Put([]byte(key), val)
}

func (db *DB) PutWithExpire(key string, val []byte, expire uint64) error {
	db.data.Put(key, DataEntity{
		Val:    val,
		Expire: expire,
	})
	return db.storage.PutWithTTL([]byte(key), val, time.Duration(expire))
}

func (db *DB) Get(key string) ([]byte, error) {
	return nil, nil
}
