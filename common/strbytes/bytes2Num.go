package strbytes

func Uint642Bytes(n uint64) []byte {
	b := make([]byte, 8, 8)

	i := 7
	for n != 0 {
		b[i] = byte(n)
		n >>= 8
		i--
	}

	return b[i:]
}

func Bytes2Uint64(n []byte) uint64 {
	var b uint64
	for i := 0; i < len(n); i++ {
		b |= uint64(n[len(n)-i-1]) << (i * 8)
	}

	return b
}

func Bytes2Uint8(n []byte) uint8 {
	if len(n) != 1 {
		return 0
	}

	return n[0]
}
