package test

import (
	"fmt"
	"github.com/rosedblabs/rosedb/v2"
	"math/rand"
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

// BenchmarkName-24        1000000000               0.6006 ns/op
// BenchmarkName-24        1000000000               0.5981 ns/op
// BenchmarkName-24        478247127                2.584 ns/op
func BenchmarkName(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = rand.Uint64()
	}
}
