package strbytes

import (
	"math/big"
	"sync"
)

var numPool = &sync.Pool{
	New: func() interface{} {
		return &big.Int{}
	},
}

func Str2Bytes(s string) (by []byte) {
	b := numPool.Get().(*big.Int)
	b.SetString(s, 10)
	by = b.Bytes()
	numPool.Put(b)
	return
}

func Bytes2Str(s []byte) (ss string) {
	b := numPool.Get().(*big.Int)
	ss = b.SetBytes(s).String()
	numPool.Put(b)
	return
}
