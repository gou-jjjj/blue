# Helper function to convert bytes to string
import string

from bsp_const import HEADER_TYPE_MASK


def bytes_to_str(data):
    return data.decode("utf-8")


# Helper function to convert string to bytes
def str_to_bytes(s):
    return s.encode('utf-8')


# Helper function to extract header type
def header_type(header):
    return header & HEADER_TYPE_MASK


# Helper function to check if string represents an integer
def check_int(s):
    return all(c in string.digits for c in s)


# Helper function to append split to data
def append_split(data):
    return data + (bytes([0]))


# Helper function to append done to data
def append_done(data):
    return data + (bytes([1]))