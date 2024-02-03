package number

type Number interface {
	Add(int64) int64
	Sub(int64) int64
	Set(int64)
	Get() int64
	Incr() int64
	Decr() int64
}
