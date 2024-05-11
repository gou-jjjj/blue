//
// Created by calvin on 24-5-10.
//

#include <string.h>
#include <stdbool.h>

#define HEADER_TYPE_MASK 0b11100000

// Helper function to convert bytes to string
char *bytes_to_str(const char *data) {
    return strdup(data); // Assuming the data is already null-terminated
}

// Helper function to convert string to bytes
char *str_to_bytes(const char *s) {
    return strdup(s); // Assuming the string is null-terminated
}

// Helper function to extract header type
int header_type(int header) {
    return header & HEADER_TYPE_MASK;
}

// Helper function to check if string represents an integer
bool check_int(const char *s) {
    for (int i = 0; s[i] != '\0'; i++) {
        if (s[i] < '0' || s[i] > '9') {
            return false;
        }
    }
    return true;
}

// Helper function to append split to data
void append_split(char *data) {
    strcat(data, "\0");
}

// Helper function to append done to data
void append_done(char *data) {
    strcat(data, "\x01");
}
