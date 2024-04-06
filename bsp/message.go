package bsp

// ReplyType 定义了回复类型的枚举。
type ReplyType byte

// 定义了各种回复类型的常量。
const (
	ReplyInfo ReplyType = (iota + 1) * (1 << 5)
	ReplyNumber
	ReplyString
	ReplyList
	ReplyError
)

// ReplyTypeMap 提供了从ReplyType到其字符串表示的映射。
var ReplyTypeMap = map[ReplyType]string{
	ReplyInfo:   "info",
	ReplyNumber: "number",
	ReplyString: "string",
	ReplyList:   "list",
	ReplyError:  "error",
}

// common -------------------------------------
// 定义了一些通用的回复类型常量。
const (
	OK ReplyType = iota + ReplyInfo
	NULL
	True
	False
)

// error -------------------------------------
// 定义了各种错误类型的常量。
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

// MessageMap 提供了从错误类型到其错误信息字符串的映射。
var MessageMap = [...]string{
	OK:    "ok",
	NULL:  "null",
	True:  "true",
	False: "false",

	ErrCommand:          "ERR unknown command",
	ErrSyntax:           "ERR syntax error",
	ErrWrongType:        "err Operation against a key holding the wrong kind of value",
	ErrHeaderType:       "ERR header type error",
	ErrValueOutOfRange:  "ERR value is out of range",
	ErrNumberArguments:  "ERR wrong number of arguments",
	ErrReplication:      "ERR replication error",
	ErrRequestParameter: "err request parameter",
	ErrEnd:              "ERR end",
	ErrPermissionDenied: "ERR permission denied",
}
