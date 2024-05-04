package datastruct

const (
	Number = iota
	String
	List
	Set
	Json
)

const (
	NumberType = "number"
	StringType = "string"
	ListType   = "list"
	SetType    = "set"
	JsonType   = "json"
)

type BlueType uint8

var BlueTypes = map[uint8]string{
	Number: NumberType,
	String: StringType,
	List:   ListType,
	Set:    SetType,
	Json:   JsonType,
}

var BlueTypes_ = map[string]uint8{
	NumberType: Number,
	StringType: String,
	ListType:   List,
	SetType:    Set,
	JsonType:   Json,
}

const (
	TypeMask    = 0x0F
	SubTypeMask = 0xF0
)

type BlueObj struct {
	Type uint8
}

func (obj *BlueObj) GetType() string {
	return BlueTypes[obj.Type&TypeMask]
}

func (obj *BlueObj) GetSubType() string {
	return BlueTypes[obj.Type&SubTypeMask]
}

type Value interface {
	Value() string
}

type Type interface {
	Type() string
}
