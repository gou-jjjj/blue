package bsp

type ReplyType byte

const (
	ReplyInfo ReplyType = iota * (1 << 5)
	ReplyNumber
	ReplyString
	ReplyList
	ReplyError
)

var ReplyTypeMap = map[ReplyType]string{
	ReplyInfo:   "info",
	ReplyNumber: "number",
	ReplyString: "string",
	ReplyList:   "list",
	ReplyError:  "error",
}

// common -------------------------------------
const (
	OK ReplyType = iota + ReplyInfo
	NULL
	True
	False
)

// error -------------------------------------
const (
	// request
	ErrCommand ReplyType = iota + ReplyError
	ErrSyntax
	ErrWrongType
	ErrHeaderType
	ErrValueOutOfRange
	ErrNumberArguments

	// network
	ErrClient
	ErrConnection
	ErrTimeout
	ErrMaxClientsReached

	// user
	ErrPermissionDenied

	// server
	ErrReplication
	ErrConfiguration
	ErrOutOfMemory
	ErrStorage
)

var MessageMap = map[ReplyType]string{
	ReplyNumber: "number",
	ReplyString: "string",
	ReplyList:   "list",

	OK:                 "ok",
	NULL:               "null",
	ErrCommand:         "ERR unknown command",
	ErrSyntax:          "ERR syntax error",
	ErrWrongType:       "WRONGTYPE Operation against a key holding the wrong kind of value",
	ErrHeaderType:      "ERR header type error",
	ErrValueOutOfRange: "ERR value is out of range",
	ErrNumberArguments: "ERR wrong number of arguments",
}

var MessageMap2 = map[string]ReplyType{
	"number": ReplyNumber,
	"string": ReplyString,
	"list":   ReplyList,

	"ok":   OK,
	"null": NULL,
	"WRONGTYPE Operation against a key holding the wrong kind of value": ErrWrongType,
}
