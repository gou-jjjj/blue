# Helper function to create request builder
from bsp_utils import str_to_bytes, append_split


class new_request_builder:
    def __init__(self, handle):
        self.data = bytes([handle])

    def with_key(self, key):
        self.data = self.data + str_to_bytes(key)
        return self

    def with_value_str(self, value):
        self.data = append_split(self.data)
        self.data = self.data + (str_to_bytes(value))
        return self

    def with_value_num(self, value):
        self.data.append(append_split(self.data))
        self.data.append(str_to_bytes(value))
        return self

    def with_values(self, *values):
        for value in values:
            self.data+(append_split(self.data))
            self.data+(str_to_bytes(value))
        return self

    def build(self):
        self.data = self.data + bytes([1])
        # print('build', self.data)
        return self.data
