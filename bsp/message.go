package bsp

type ReplyType byte

const (
	ReplyInfo ReplyType = (iota + 1) * (1 << 5)
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
	ErrRequestParameter
	ErrEnd

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

var MessageMap = [...]string{
	OK:    "ok",
	NULL:  "null",
	True:  "true",
	False: "false",

	ErrCommand:          "ERR unknown command",
	ErrSyntax:           "ERR syntax error",
	ErrWrongType:        "Err Operation against a key holding the wrong kind of value",
	ErrHeaderType:       "ERR header type error",
	ErrValueOutOfRange:  "ERR value is out of range",
	ErrNumberArguments:  "ERR wrong number of arguments",
	ErrReplication:      "ERR replication error",
	ErrRequestParameter: "Err request parameter",
	ErrEnd:              "ERR end",
}
