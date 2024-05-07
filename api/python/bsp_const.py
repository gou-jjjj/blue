Done = 0x01
# Define constants for header types
HEADER_TYPE_MASK = 0b11100000
TYPE_SYSTEM = 1 << 5
TYPE_DB = 2 << 5
TYPE_NUMBER = 3 << 5
TYPE_STRING = 4 << 5
TYPE_LIST = 5 << 5
TYPE_SET = 6 << 5
TYPE_JSON = 7 << 5

# Define header handle error value
HANDLE_ERROR = 255


# Define common reply types
class ReplyType:
    INFO = 32
    NUMBER = INFO * 2
    STRING = INFO * 3
    LIST = INFO * 4
    ERROR = INFO * 5


# Common reply messages
class Common:
    OK = ReplyType.INFO
    NULL = OK + 1
    TRUE = NULL + 1
    FALSE = TRUE + 1


# Error messages
class Error:
    COMMAND = ReplyType.ERROR
    SYNTAX = COMMAND + 1
    WRONG_TYPE = SYNTAX + 1
    HEADER_TYPE = WRONG_TYPE + 1
    VALUE_OUT_OF_RANGE = HEADER_TYPE + 1
    NUMBER_ARGUMENTS = VALUE_OUT_OF_RANGE + 1
    REQUEST_PARAMETER = NUMBER_ARGUMENTS + 1
    END = REQUEST_PARAMETER + 1
    CLIENT = END + 1
    CONNECTION = CLIENT + 1
    TIMEOUT = CONNECTION + 1
    MAX_CLIENTS_REACHED = TIMEOUT + 1
    PERMISSION_DENIED = MAX_CLIENTS_REACHED + 1
    REPLICATION = PERMISSION_DENIED + 1
    CONFIGURATION = REPLICATION + 1
    OUT_OF_MEMORY = CONFIGURATION + 1
    STORAGE = OUT_OF_MEMORY + 1


# Map reply types to strings
REPLY_TYPE_MAP = {
    ReplyType.INFO: "info",
    ReplyType.NUMBER: "number",
    ReplyType.STRING: "string",
    ReplyType.LIST: "list",
    ReplyType.ERROR: "error",
}

# Map error codes to messages
MESSAGE_MAP = {
    Common.OK: "ok",
    Common.NULL: "null",
    Common.TRUE: "true",
    Common.FALSE: "false",
    Error.COMMAND: "ERR unknown command",
    Error.SYNTAX: "ERR syntax error",
    Error.WRONG_TYPE: "ERR Operation against a key holding the wrong kind of value",
    Error.HEADER_TYPE: "ERR header type error",
    Error.VALUE_OUT_OF_RANGE: "ERR value is out of range",
    Error.NUMBER_ARGUMENTS: "ERR wrong number of arguments",
    Error.REPLICATION: "ERR replication error",
    Error.REQUEST_PARAMETER: "ERR request parameter",
    Error.END: "ERR end",
    Error.PERMISSION_DENIED: "ERR permission denied",
}

# Command handles
COMMAND_HANDLES = {
    TYPE_SYSTEM + 1: "AUTH", TYPE_DB + 1: "DBSIZE", TYPE_DB + 2: "DEL", TYPE_SYSTEM + 2: "EXIT",
    TYPE_DB + 3: "EXPIRE", TYPE_STRING + 1: "GET", TYPE_SYSTEM + 3: "HELP", TYPE_NUMBER + 1: "INCR",
    TYPE_DB + 4: "KVS", TYPE_STRING + 2: "LEN", TYPE_LIST + 1: "LGET", TYPE_LIST + 2: "LLEN",
    TYPE_LIST + 3: "LPOP", TYPE_LIST + 4: "LPUSH", TYPE_NUMBER + 2: "NGET", TYPE_NUMBER + 3: "NSET",
    TYPE_SYSTEM + 4: "PING", TYPE_LIST + 5: "RPOP", TYPE_LIST + 6: "RPUSH", TYPE_SET + 1: "SADD",
    TYPE_SET + 2: "SDEL", TYPE_SET + 3: "SGET", TYPE_SET + 4: "SIN", TYPE_SET + 5: "SPOP",
    TYPE_SYSTEM + 5: "SELECT", TYPE_STRING + 3: "SET", TYPE_DB + 5: "TYPE", TYPE_SYSTEM + 6: "VERSION"
}

COMMANDS = {v: k for k, v in COMMAND_HANDLES.items()}


# Define command class
class Command:
    def __init__(self, name, summary, group, arity, key, value, arguments):
        self.name = name
        self.summary = summary
        self.group = group
        self.arity = arity
        self.key = key
        self.value = value
        self.arguments = arguments


# Command mappings
COMMANDS_MAP = {
    handle: Command(
        name=cmd_name, summary="", group="", arity=0, key="", value="", arguments=[]
    ) for handle, cmd_name in COMMAND_HANDLES.items()
}
