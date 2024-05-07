from bsp_const import MESSAGE_MAP, Done, Error, ReplyType
from bsp_utils import bytes_to_str


def reply_message(reply):
    if reply is None or len(reply) < 2 or reply[-1] != Done:
        return "", MESSAGE_MAP[Error.REPLICATION]

    reply = reply[:-1]

    if reply[0] == ReplyType.NUMBER:
        print("num", reply)
        return bytes_to_str(reply[1:]), None
    elif reply[0] == ReplyType.STRING:
        return reply[1:].decode("utf-8"), None
    elif reply[0] == ReplyType.LIST:
        print("list", reply)
        return str(reply[1:]), None

    return MESSAGE_MAP.get(reply[0], None), None


class new_response:
    def __init__(self, reply):
        self._reply, self._err = reply_message(reply)

    def get_reply(self):
        if self._err is not None:
            raise ValueError(self._err)
        return self._reply

    def get_err(self):
        return self._err
