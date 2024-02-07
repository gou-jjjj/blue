package bsp

const (
	Split byte = 0b0
	Done  byte = 0b1

	Expire byte = 0b110
	From   byte = 0b111
	End    byte = 0b1000
)

func AppendSplit(d []byte) []byte {
	return append(d, Split)
}

func AppendDone(d []byte) []byte {
	return append(d, Done)
}
