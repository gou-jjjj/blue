package strbytes

import "strconv"

//var numPool = &sync.Pool{
//	New: func() interface{} {
//		return &big.Int{}
//	},
//}
//
//func Str2Bytes(s string) (by []byte) {
//	b := numPool.Get().(*big.Int)
//	b.SetString(s, 10)
//	by = b.Bytes()
//	numPool.Put(b)
//	return
//}
//
//func Bytes2Str(s []byte) (ss string) {
//	b := numPool.Get().(*big.Int)
//	ss = b.SetBytes(s).String()
//	numPool.Put(b)
//	return
//}
//
//func CheckInt(str string) (b2 bool) {
//	b := numPool.Get().(*big.Int)
//	_, b2 = b.SetString(str, 10)
//	numPool.Put(b)
//	return
//}

func Str2Bytes(s string) (by []byte) {
	by = []byte(s)
	return
}

func Bytes2Str(s []byte) (ss string) {
	ss = string(s)
	return
}

func CheckInt(str string) bool {
	_, err := strconv.Atoi(str)
	if err != nil {
		return false
	}
	return true
}
