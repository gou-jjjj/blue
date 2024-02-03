package number

import (
	iter "blue/datastruct"
	"bsp"
	"bsp/common"
	"strconv"
	"sync/atomic"
)

const (
	MaxNorNumber = 1<<63 - 1
	MinNorNumber = -1 << 63
)

type NorNum struct {
	iter.BlueObj
	v int64
}

func NewNumber(number ...any) (*NorNum, error) {
	n := &NorNum{
		BlueObj: iter.BlueObj{
			Type: iter.Number,
		}}

	if len(number) > 0 {
		switch number[0].(type) {
		case []byte:
			parseInt := common.BytesToUint64(number[0].([]byte))
			n.Set(int64(parseInt))
		case string:
			parseInt, err := strconv.ParseInt(number[0].(string), 10, 64)
			if err != nil {
				return nil, err
			}
			n.Set(parseInt)
		default:
			return nil, bsp.NewErr(bsp.ErrWrongType)
		}
	}

	return n, nil
}

func (n *NorNum) Value() int64 {
	return n.v
}

func (n *NorNum) Add(val int64) int64 {
	if MaxNorNumber-n.v < val {
		panic("overflow")
	}
	return atomic.AddInt64(&n.v, val)
}

func (n *NorNum) Sub(val int64) int64 {
	if n.v < val+MinNorNumber {
		panic("overflow")
	}
	return atomic.AddInt64(&n.v, -val)
}

func (n *NorNum) Set(i int64) {
	atomic.StoreInt64(&n.v, i)
}

func (n *NorNum) Get() int64 {
	return atomic.LoadInt64(&n.v)
}

func (n *NorNum) Incr() int64 {
	return n.Add(1)
}

func (n *NorNum) Decr() int64 {
	return n.Sub(1)
}
