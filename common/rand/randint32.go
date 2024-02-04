package rand

import (
	_ "unsafe"
)

//go:nosplit
//go:linkname fastrand runtime.fastrand
func fastrand() uint32

//go:nosplit
//go:linkname fastrand64 runtime.fastrand64
func fastrand64() uint64

func Randu32() uint32 {
	return fastrand()
}

func Randu64() uint64 {
	return fastrand64()
}
