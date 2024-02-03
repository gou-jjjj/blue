package bsp

type HeaderInter interface {
	Type() Header
	Handle() Header
	Len() uint8

	TypeStr() string
	HandleStr() string

	ValueLenMax() uint8
	Bytes() []byte
}

const (
	typeMask   = 0b111_00000_00000000
	handleMask = 0b000_11111_00000000
	lenMask    = 0b000_00000_11111111
)

const (
	TypeSystem Header = iota * (1 << 5)
	TypeDB
	TypeNumber
	TypeString
	TypeList
	TypeSet
	TypeJson
)

var TypeMap = map[Header]string{
	TypeSystem: "system",
	TypeDB:     "db",
	TypeNumber: "num",
	TypeString: "str",
	TypeList:   "list",
	TypeSet:    "set",
	TypeJson:   "json",
}

var TypeMap2 = map[string]Header{
	"system": TypeSystem,
	"db":     TypeDB,
	"num":    TypeNumber,
	"str":    TypeString,
	"list":   TypeList,
	"set":    TypeSet,
	"json":   TypeJson,
}

type Header uint16

func NewHeader(header Header, Len int8) Header {
	return (header << 8) | Header(Len)
}

func (h Header) Type() Header {
	return h & typeMask >> 8
}

func (h Header) Handle() Header {
	return h >> 8
}

func (h Header) Len() uint8 {
	return uint8(h & lenMask)
}

func (h Header) Bytes() []byte {
	return []byte{byte(h >> 8), byte(h)}
}

func (h Header) TypeStr() string {
	return TypeMap[h.Type()]
}

func (h Header) HandleStr() string {
	return HandleMap[h.Handle()]
}

func (h Header) ValueLenMax() uint8 {
	return HandlePara[h>>8]
}
