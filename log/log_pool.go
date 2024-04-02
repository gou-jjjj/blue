package log

import "sync"

var (
	poolMsg = sync.Pool{
		New: func() interface{} {
			return &Msg{}
		},
	}
)

func putMsg(m *Msg) {
	poolMsg.Put(m)
}

func getMsg() *Msg {
	return poolMsg.Get().(*Msg)
}
