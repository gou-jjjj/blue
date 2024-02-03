package test

import (
	"fmt"
	"github.com/rosedblabs/rosedb/v2"
	"testing"
)

func TestQWQName(t *testing.T) {
	options := rosedb.DefaultOptions
	options.DirPath = "./data"
	open, err := rosedb.Open(options)
	if err != nil {
		panic(err)
		return
	}

	defer open.Close()

	//fmt.Printf("%v\n", open.Put([]byte("name"), []byte("rosedb")))
	//fmt.Printf("%v\n", open.Expire([]byte("name"), 50*time.Second))

	fmt.Println(open.TTL([]byte("name")))
	fmt.Println(open.TTL([]byte("aaa")))

}
