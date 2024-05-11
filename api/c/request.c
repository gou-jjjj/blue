//
// Created by calvin on 24-5-10.
//

#include <stdlib.h>
#include <string.h>
#include <stdio.h>
#include <stdarg.h>

#include "utils.c"

typedef struct {
    char *data;
} new_request_builder;

// Helper function to create new request builder
new_request_builder *new_request_builder_create(char handle) {
    new_request_builder *builder = malloc(sizeof(new_request_builder));
    if (builder == NULL) {
        return NULL; // Memory allocation failed
    }
    builder->data = (char *)malloc(1);
    if (builder->data == NULL) {
        free(builder);
        return NULL; // Memory allocation failed
    }
    builder->data[0] = handle;
    return builder;
}

// Method to add key to request builder
void new_request_builder_with_key(new_request_builder *builder, const char *key) {
    size_t key_len = strlen(key);
    builder->data = realloc(builder->data, strlen(builder->data) + key_len + 1);
    strcat(builder->data, key);
}

// Method to add string value to request builder
void new_request_builder_with_value_str(new_request_builder *builder, const char *value) {
    append_split(builder->data);
    strcat(builder->data, value);
}

// Method to add numeric value to request builder
void new_request_builder_with_value_num(new_request_builder *builder, int value) {
    append_split(builder->data);
    char num_str[20]; // Assuming int won't exceed 20 digits
    sprintf(num_str, "%d", value);
    strcat(builder->data, num_str);
}

// Method to add multiple values to request builder
void new_request_builder_with_values(new_request_builder *builder, int count, ...) {
    va_list args;
    va_start(args, count);
    for (int i = 0; i < count; i++) {
        const char *value = va_arg(args, const char *);
        append_split(builder->data);
        strcat(builder->data, value);
    }
    va_end(args);
}

// Method to build request
char *new_request_builder_build(new_request_builder *builder) {
    append_done(builder->data);
    return builder->data;
}

// Method to destroy request builder
void new_request_builder_destroy(new_request_builder *builder) {
    free(builder->data);
    free(builder);
}
