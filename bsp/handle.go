package bsp

const (
	SystemExit = iota + TypeSystem
	SystemVersion
)

const (
	DBSelected Header = iota + TypeDB
	DBSelectdb
)

const (
	NumGet Header = iota + TypeNumber
	NumSet
	NumIncr
	NumDecr
)

const (
	StrGet Header = iota + TypeString
	StrSet
)

var HandlePara = [256]uint8{
	SystemExit:    0,
	SystemVersion: 0,
	DBSelected:    1,
	DBSelectdb:    1,
	NumGet:        1,
	NumSet:        2,
	NumIncr:       0,
	NumDecr:       0,
	StrGet:        1,
	StrSet:        2,
}

var HandleMap = [uint(256)]string{
	SystemExit:    "sys exit",
	SystemVersion: "sys ver",
	NumGet:        "num get",
	NumSet:        "num set",
	NumIncr:       "num incr",
	NumDecr:       "num decr",
	StrGet:        "str get",
	StrSet:        "str set",
}

var HandleMap2 = map[string]Header{
	"sys exit":    SystemExit,
	"sys ver":     SystemVersion,
	"num get":     NumGet,
	"num set":     NumSet,
	"num incr":    NumIncr,
	"num decr":    NumDecr,
	"str get":     StrGet,
	"str set":     StrSet,
}
