#include "bsp.h"

// 实现 NewBspProto 函数
BspProto* NewBspProto() {
    BspProto *b = (BspProto*)malloc(sizeof(BspProto));
    if (b == NULL) {
        fprintf(stderr, "Memory allocation failed\n");
        exit(EXIT_FAILURE);
    }
    b->key = NULL;
    b->value = NULL;
    b->buf = NULL;
    return b;
}

// 实现 PutBspProto 函数
void PutBspProto(BspProto *b) {
    if (b == NULL) return;
    if (b->key != NULL) free(b->key);
    if (b->value != NULL) {
        for (int i = 0; b->value[i] != NULL; i++) {
            free(b->value[i]);
        }
        free(b->value);
    }
    if (b->buf != NULL) free(b->buf);
    free(b);
}

// 实现 SetHeader 函数
void SetHeader(BspProto *b, Header h) {
    // 根据你的需求实现
}

// 实现 Key 函数
char* Key(BspProto *b) {
    return b->key;
}

// 实现 SetKey 函数
void SetKey(BspProto *b, char *key) {
    if (b->key != NULL) free(b->key);
    b->key = strdup(key);
}

// 实现 ValueBytes 函数
char* ValueBytes(BspProto *b) {
    return b->value[0];
}

// 实现 ValueStr 函数
char* ValueStr(BspProto *b) {
    return b->value[0];
}

// 实现 SetValue 函数
void SetValue(BspProto *b, char *value) {
    if (b->value != NULL) {
        free(b->value[0]);
        free(b->value);
    }
    b->value = (char**)malloc(2 * sizeof(char*));
    b->value[0] = strdup(value);
    b->value[1] = NULL;
}

// 实现 SetValues 函数
void SetValues(BspProto *b, char **values) {
    if (b->value != NULL) {
        for (int i = 0; b->value[i] != NULL; i++) {
            free(b->value[i]);
        }
        free(b->value);
    }
    int numValues = 0;
    while (values[numValues] != NULL) {
        numValues++;
    }
    b->value = (char**)malloc((numValues + 1) * sizeof(char*));
    for (int i = 0; i < numValues; i++) {
        b->value[i] = strdup(values[i]);
    }
    b->value[numValues] = NULL;
}

// 实现 Buf 函数
char* Buf(BspProto *b) {
    return b->buf;
}

// 实现 SetBuf 函数
void SetBuf(BspProto *b, char *buf) {
    if (b->buf != NULL) free(b->buf);
    b->buf = strdup(buf);
}
