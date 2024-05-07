import string
import socket
import time

from bsp_const import COMMANDS
from bsp_request import new_request_builder
from bsp_response import new_response


# Define client class
class Client:
    def __init__(self, addr="127.0.0.1:13140", db=1, try_times=3, token="", timeout=10):
        self.addr = addr
        self.db = db
        self.try_times = try_times
        self.token = token
        self.timeout = timeout
        self.conn = None
        self.connect()

    def connect(self):
        host, port = self.addr.split(":")
        for _ in range(self.try_times):
            try:
                self.conn = socket.create_connection((host, int(port)), timeout=self.timeout)
                return
            except Exception as e:
                time.sleep(1)
        raise e

    def send_command(self, command):
        try:
            self.conn.sendall(command)
            return self.conn.recv(4096)
        except Exception as e:
            raise e

    def exec(self, buf):
        if self.conn is None:
            self.connect()

        reply = self.send_command(buf)
        return new_response(reply)

    def close(self):
        if self.conn:
            self.conn.close()

    def exec_pipeline(self, buf):
        self.connect()
        self.conn.sendall(b"".join(buf))
        replies = []
        while True:
            try:
                reply = self.conn.recv(4096)
                if not reply:
                    break
                replies.append(reply)
            except Exception as e:
                if isinstance(e, EOFError):
                    break
                raise e
        if len(replies) != len(buf):
            raise ValueError("Pipeline error")
        return replies

    def __del__(self):
        self.close()

    ###---------------- 指令-------------------------------------
    ###---------------- 指令-------------------------------------
    ###---------------- 指令-------------------------------------

    def version(self):
        return self.exec(
            new_request_builder(COMMANDS["VERSION"]).build()
        )

    def get(self, key):
        return self.exec(
            new_request_builder(COMMANDS["GET"]).with_key(key).build()
        )

    def set(self, key, value):
        return self.exec(
            new_request_builder(COMMANDS["SET"]).with_key(key).with_value_str(value).build()
        )

    def delete(self, key):
        return self.exec(
            new_request_builder(COMMANDS["DEL"]).with_key(key).build()
        )

    def dbsize(self):
        return self.exec(
            new_request_builder(COMMANDS["DBSIZE"]).build()
        )

    def auth(self, credentials):
        return self.exec(
            new_request_builder(COMMANDS["AUTH"]).with_value_str(credentials).build()
        )

    def exit(self):
        return self.exec(
            new_request_builder(COMMANDS["EXIT"]).build()
        )

    def expire(self, key, seconds):
        return self.exec(
            new_request_builder(COMMANDS["EXPIRE"]).with_key(key).with_value_str(seconds).build()
        )

    def help(self, key):
        return self.exec(
            new_request_builder(COMMANDS["HELP"]).with_value_str(key).build()
        )

    def incr(self, key):
        return self.exec(
            new_request_builder(COMMANDS["INCR"]).with_key(key).build()
        )

    def kvs(self):
        return self.exec(
            new_request_builder(COMMANDS["KVS"]).build()
        )

    def len(self, key):
        return self.exec(
            new_request_builder(COMMANDS["LEN"]).with_key(key).build()
        )

    def lget(self, key):
        return self.exec(
            new_request_builder(COMMANDS["LGET"]).with_key(key).build()
        )

    def llen(self, key):
        return self.exec(
            new_request_builder(COMMANDS["LLEN"]).with_key(key).build()
        )

    def lpop(self, key):
        return self.exec(
            new_request_builder(COMMANDS["LPOP"]).with_key(key).build()
        )

    def lpush(self, key, *values):
        return self.exec(
            new_request_builder(COMMANDS["LPUSH"]).with_key(key).with_values(*values).build()
        )

    def nget(self, key):
        return self.exec(
            new_request_builder(COMMANDS["NGET"]).with_value_str(key).build()
        )

    def nset(self, key, value):
        return self.exec(
            new_request_builder(COMMANDS["NSET"]).with_key(key).with_value_str(value).build()
        )

    def ping(self):
        return self.exec(
            new_request_builder(COMMANDS["PING"]).build()
        )

    def rpop(self, key):
        return self.exec(
            new_request_builder(COMMANDS["RPOP"]).with_key(key).build()
        )

    def rpush(self, key, *values):
        return self.exec(
            new_request_builder(COMMANDS["RPUSH"]).with_key(key).with_values(*values).build()
        )

    def sadd(self, key, *values):
        return self.exec(
            new_request_builder(COMMANDS["SADD"]).with_key(key).with_values(*values).build()
        )

    def sdel(self, key, value):
        return self.exec(
            new_request_builder(COMMANDS["SDEL"]).with_key(key).with_value_str(value).build()
        )

    def sget(self, key):
        return self.exec(
            new_request_builder(COMMANDS["SGET"]).with_key(key).build()
        )

    def sin(self, key, value):
        return self.exec(
            new_request_builder(COMMANDS["SIN"]).with_key(key).with_value_str(value).build()
        )

    def spop(self, key):
        return self.exec(
            new_request_builder(COMMANDS["SPOP"]).with_key(key).build()
        )

    def select(self, db_index):
        return self.exec(
            new_request_builder(COMMANDS["SELECT"]).with_value_str(db_index).build()
        )

    def type(self, key):
        return self.exec(
            new_request_builder(COMMANDS["TYPE"]).with_key(key).build()
        )


if __name__ == '__main__':
    c = Client()
    rep = c.set("a", "bdasd251752asdsad")
    print(rep.get_reply())

    rep = c.dbsize()
    print(rep.get_reply())

    rep = c.get("a")
    print(rep.get_reply())

    rep=c.type("a")
    print(rep.get_reply())
