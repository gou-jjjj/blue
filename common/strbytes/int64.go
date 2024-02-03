package strbytes

func Uint64ToBytes(n uint64) []byte {
	b := make([]byte, 0, 8)
	for n != 0 {
		b = append(b, byte(n))
		n >>= 8
	}

	return b
}

func BytesToUint64(n []byte) uint64 {
	var b uint64
	for i := 0; i < len(n); i++ {
		b |= uint64(n[i]) << (i * 8)
	}

	return b
}
